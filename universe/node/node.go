package node

import (
	"context"
	"fmt"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	influx_api "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/mplugin"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/user"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Node = (*Node)(nil)

type Node struct {
	id                  uuid.UUID
	name                string
	cfg                 *config.Config
	ctx                 context.Context
	log                 *zap.SugaredLogger
	db                  database.DB
	router              *gin.Engine
	worlds              universe.Worlds
	assets2d            universe.Assets2d
	assets3d            universe.Assets3d
	spaceTypes          universe.SpaceTypes
	userTypes           universe.UserTypes
	attributeTypes      universe.AttributeTypes
	plugins             universe.Plugins
	nodeAttributes      *generic.SyncMap[entry.AttributeID, *entry.AttributePayload]
	spaceIDToWorldID    *generic.SyncMap[uuid.UUID, uuid.UUID]
	influx              influx_api.WriteAPIBlocking
	pluginController    *mplugin.PluginController
	corePluginInterface *mplugin.PluginInterface
}

func NewNode(
	id uuid.UUID,
	cfg *config.Config,
	db database.DB,
	worlds universe.Worlds,
	assets2D universe.Assets2d,
	assets3D universe.Assets3d,
	plugins universe.Plugins,
	spaceTypes universe.SpaceTypes,
	userTypes universe.UserTypes,
	attributeTypes universe.AttributeTypes,
) *Node {
	return &Node{
		id:               id,
		cfg:              cfg,
		db:               db,
		worlds:           worlds,
		assets2d:         assets2D,
		assets3d:         assets3D,
		plugins:          plugins,
		spaceTypes:       spaceTypes,
		userTypes:        userTypes,
		attributeTypes:   attributeTypes,
		nodeAttributes:   generic.NewSyncMap[entry.AttributeID, *entry.AttributePayload](),
		spaceIDToWorldID: generic.NewSyncMap[uuid.UUID, uuid.UUID](),
	}
}

func (n *Node) GetID() uuid.UUID {
	return n.id
}

func (n *Node) GetName() string {
	return n.name
}

func (n *Node) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	n.ctx = ctx
	n.log = log

	consoleWriter := zapcore.Lock(os.Stdout)
	gin.DefaultWriter = consoleWriter

	//TODO: hash salt once it is present in the DB
	utils.SetAnonymizer(n.id, uuid.Nil)

	r := gin.New()
	r.Use(gin.LoggerWithWriter(consoleWriter))
	r.Use(gin.RecoveryWithWriter(consoleWriter))

	n.router = r
	n.pluginController = mplugin.NewPluginController(n.id)

	return nil
}

func (n *Node) GetWorlds() universe.Worlds {
	return n.worlds
}

func (n *Node) GetAssets2d() universe.Assets2d {
	return n.assets2d
}

func (n *Node) GetAssets3d() universe.Assets3d {
	return n.assets3d
}

func (n *Node) GetPlugins() universe.Plugins {
	return n.plugins
}

func (n *Node) GetAttributeTypes() universe.AttributeTypes {
	return n.attributeTypes
}

func (n *Node) GetSpaceTypes() universe.SpaceTypes {
	return n.spaceTypes
}

func (n *Node) GetUserTypes() universe.UserTypes {
	return n.userTypes
}

func (n *Node) GetAllSpaces() map[uuid.UUID]universe.Space {
	spaces := make(map[uuid.UUID]universe.Space)

	for _, world := range n.GetWorlds().GetWorlds() {
		for spaceID, space := range world.GetAllSpaces() {
			spaces[spaceID] = space
		}
	}

	return spaces
}

func (n *Node) FilterAllSpaces(predicateFn universe.SpaceFilterPredicateFn) map[uuid.UUID]universe.Space {
	spaces := make(map[uuid.UUID]universe.Space)
	for _, world := range n.GetWorlds().GetWorlds() {
		for spaceID, space := range world.FilterAllSpaces(predicateFn) {
			spaces[spaceID] = space
		}
	}
	// TODO: check self
	return spaces
}

func (n *Node) GetSpaceFromAllSpaces(spaceID uuid.UUID) (universe.Space, bool) {
	worldID, ok := n.spaceIDToWorldID.Load(spaceID)
	if !ok {
		return nil, false
	}
	world, ok := n.GetWorlds().GetWorld(worldID)
	if !ok {
		return nil, false
	}
	return world.GetSpaceFromAllSpaces(spaceID)
}

func (n *Node) AddSpaceToAllSpaces(space universe.Space) error {
	n.spaceIDToWorldID.Store(space.GetID(), space.GetWorld().GetID())
	return nil
}

func (n *Node) RemoveSpaceFromAllSpaces(space universe.Space) (bool, error) {
	n.spaceIDToWorldID.Mu.RLock()
	defer n.spaceIDToWorldID.Mu.RUnlock()

	if _, ok := n.spaceIDToWorldID.Data[space.GetID()]; ok {
		delete(n.spaceIDToWorldID.Data, space.GetID())

		return true, nil
	}

	return false, nil
}

func (n *Node) AddAPIRegister(register universe.APIRegister) {
	register.RegisterAPI(n.router)
}

func (n *Node) Run() error {
	if err := n.worlds.Run(); err != nil {
		return errors.WithMessage(err, "failed to run worlds")
	}

	return n.router.Run(fmt.Sprintf("%s:%d", n.cfg.Settings.Address, n.cfg.Settings.Port))
}

func (n *Node) Stop() error {
	return n.worlds.Stop()
}

func (n *Node) Load() error {
	n.log.Infof("Loading node %s...", n.GetID())

	group, _ := errgroup.WithContext(n.ctx)
	group.Go(
		func() error {
			return n.assets2d.Load()
		},
	)
	group.Go(
		func() error {
			return n.assets3d.Load()
		},
	)
	group.Go(
		func() error {
			return n.userTypes.Load()
		},
	)
	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to load assets")
	}

	group, _ = errgroup.WithContext(n.ctx)
	group.Go(
		func() error {
			return n.attributeTypes.Load()
		},
	)
	group.Go(
		func() error {
			return n.spaceTypes.Load()
		},
	)
	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to load additional data")
	}

	if err := n.plugins.Load(); err != nil {
		return errors.WithMessage(err, "failed to load space types")
	}

	if err := n.loadSelfData(); err != nil {
		return errors.WithMessage(err, "failed to load self data")
	}

	if err := n.worlds.Load(); err != nil {
		return errors.WithMessage(err, "failed to load worlds")
	}

	n.AddAPIRegister(n)

	n.log.Infof("Node loaded: %s", n.GetID())

	return nil
}

func (n *Node) Save() error {
	n.log.Infof("Saving node: %s...", n.GetID())

	var wg sync.WaitGroup
	var errs *multierror.Error
	var errsMu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := n.assets2d.Save(); err != nil {
			errsMu.Lock()
			defer errsMu.Unlock()
			errs = multierror.Append(errs, errors.WithMessage(err, "failed to save assets 2d"))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := n.assets3d.Save(); err != nil {
			errsMu.Lock()
			defer errsMu.Unlock()
			errs = multierror.Append(errs, errors.WithMessage(err, "failed to save assets 3d"))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := n.spaceTypes.Save(); err != nil {
			errsMu.Lock()
			defer errsMu.Unlock()
			errs = multierror.Append(errs, errors.WithMessage(err, "failed to save space types"))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := n.worlds.Save(); err != nil {
			errsMu.Lock()
			defer errsMu.Unlock()
			errs = multierror.Append(errs, errors.WithMessage(err, "failed to save worlds"))
		}
	}()

	wg.Wait()

	n.log.Infof("Node saved: %s", n.GetID())

	return errs.ErrorOrNil()
}

func (n *Node) loadSelfData() error {
	go func() {
		if err := n.loadNodeAttributes(); err != nil {
			n.log.Error(errors.WithMessage(err, "Node: loadSelfData: failed to load node attributes"))
		}
	}()

	return nil
}

func (n *Node) detectSpawnWorld(userId uuid.UUID) universe.World {
	// TODO: implement. Temporary, just first world from the list
	wid := uuid.MustParse("d83670c7-a120-47a4-892d-f9ec75604f74")
	if world, ok := n.worlds.GetWorld(wid); ok != false {
		return world

	}
	return nil
}

func (n *Node) LoadUser(userID uuid.UUID) (universe.User, error) {
	user := user.NewUser(userID, n.db)
	if err := user.Initialize(n.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize user: %s", userID)
	}

	if err := user.Load(); err != nil {
		return nil, errors.WithMessagef(err, "failed to load user: %s", userID)
	}

	fmt.Printf("%+v\n", user.GetPosition())
	user.SetPosition(cmath.Vec3{X: 50, Y: 50, Z: 150})
	fmt.Printf("%+v\n", user.GetPosition())
	return user, nil
}
