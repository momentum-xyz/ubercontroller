package node

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

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
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/space"
	"github.com/momentum-xyz/ubercontroller/universe/streamchat"
	"github.com/momentum-xyz/ubercontroller/universe/user"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Node = (*Node)(nil)

type Node struct {
	*space.Space
	name                string
	cfg                 *config.Config
	ctx                 context.Context
	log                 *zap.SugaredLogger
	db                  database.DB
	router              *gin.Engine
	httpServer          *http.Server
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
	chatService         *streamchat.StreamChat
}

func NewNode(
	id uuid.UUID,
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
		Space:            space.NewSpace(id, db, nil),
		db:               db,
		worlds:           worlds,
		assets2d:         assets2D,
		assets3d:         assets3D,
		plugins:          plugins,
		spaceTypes:       spaceTypes,
		userTypes:        userTypes,
		attributeTypes:   attributeTypes,
		nodeAttributes:   generic.NewSyncMap[entry.AttributeID, *entry.AttributePayload](0),
		spaceIDToWorldID: generic.NewSyncMap[uuid.UUID, uuid.UUID](0),
	}
}

func (n *Node) GetName() string {
	return n.name
}

func (n *Node) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.Errorf("failed to get config from context: %T", ctx.Value(types.ConfigContextKey))
	}

	n.ctx = ctx
	n.cfg = cfg
	n.log = log

	consoleWriter := zapcore.Lock(os.Stdout)
	gin.DefaultWriter = consoleWriter

	//TODO: hash salt once it is present in the DB
	utils.SetAnonymizer(n.GetID(), uuid.Nil)

	r := gin.New()
	r.Use(gin.LoggerWithWriter(consoleWriter))
	r.Use(gin.RecoveryWithWriter(consoleWriter))

	n.router = r
	n.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", n.cfg.Settings.Address, n.cfg.Settings.Port),
		Handler: n.router,
	}
	n.pluginController = mplugin.NewPluginController(n.GetID())
	n.chatService = streamchat.NewStreamChat()
	if err := n.chatService.Initialize(ctx); err != nil {
		return err
	}

	return n.Space.Initialize(ctx)
}

func (n *Node) CreateSpace(spaceID uuid.UUID) (universe.Space, error) {
	return nil, errors.Errorf("not permitted for node")
}

func (n *Node) SetParent(parent universe.Space, updateDB bool) error {
	return errors.Errorf("not permitted for node")
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
	spaces := map[uuid.UUID]universe.Space{
		n.GetID(): n,
	}

	for _, world := range n.GetWorlds().GetWorlds() {
		for spaceID, space := range world.GetAllSpaces() {
			spaces[spaceID] = space
		}
	}

	return spaces
}

func (n *Node) FilterAllSpaces(predicateFn universe.SpacesFilterPredicateFn) map[uuid.UUID]universe.Space {
	spaces := make(map[uuid.UUID]universe.Space)
	if predicateFn(n.GetID(), n) {
		spaces[n.GetID()] = n
	}

	for _, world := range n.GetWorlds().GetWorlds() {
		for spaceID, space := range world.FilterAllSpaces(predicateFn) {
			spaces[spaceID] = space
		}
	}

	return spaces
}

func (n *Node) GetSpaceFromAllSpaces(spaceID uuid.UUID) (universe.Space, bool) {
	if spaceID == n.GetID() {
		return n, true
	}

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
	if space.GetID() == n.GetID() {
		return errors.Errorf("not permitted for node")
	}

	n.spaceIDToWorldID.Store(space.GetID(), space.GetWorld().GetID())
	return nil
}

func (n *Node) RemoveSpaceFromAllSpaces(space universe.Space) (bool, error) {
	if space.GetID() == n.GetID() {
		return false, errors.Errorf("not permitted for node")
	}

	n.spaceIDToWorldID.Mu.Lock()
	defer n.spaceIDToWorldID.Mu.Unlock()

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
	n.SetEnabled(true)

	// in goroutine for graceful shutdown
	go func() {
		if err := n.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			n.log.Fatal(errors.WithMessage(err, "Node: Run: failed to run http server"))
		}
	}()

	<-n.ctx.Done()
	gracePeriod := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()

	if err := n.httpServer.Shutdown(ctx); err != nil {
		return errors.WithMessage(err, "failed to shutdown http server")
	}

	return nil
}

func (n *Node) Stop() error {
	if err := n.worlds.Stop(); err != nil {
		return errors.WithMessage(err, "failed to stop worlds")
	}
	n.SetEnabled(false)

	return nil
}

func (n *Node) Load() error {
	n.log.Infof("Loading node %s...", n.GetID())

	group, _ := errgroup.WithContext(n.ctx)
	// main loading thread
	group.Go(func() error {
		// first stage
		group, _ := errgroup.WithContext(n.ctx)
		group.Go(n.assets2d.Load)
		group.Go(n.assets3d.Load)
		group.Go(n.userTypes.Load)
		group.Go(n.attributeTypes.Load)
		if err := group.Wait(); err != nil {
			return errors.WithMessage(err, "failed to load basic data")
		}

		// second stage
		group, _ = errgroup.WithContext(n.ctx)
		group.Go(n.spaceTypes.Load)
		group.Go(n.plugins.Load)
		if err := group.Wait(); err != nil {
			return errors.WithMessage(err, "failed to load additional data")
		}

		// third stage
		group, _ = errgroup.WithContext(n.ctx)
		group.Go(
			func() error {
				nodeEntry, err := n.db.NodesGetNode(n.ctx)
				if err != nil {
					return errors.WithMessage(err, "failed to get node")
				}
				if err := n.LoadFromEntry(nodeEntry.Space, false); err != nil {
					return errors.WithMessage(err, "failed to load node from entry")
				}

				if err := n.loadNodeAttributes(); err != nil {
					return errors.WithMessage(err, "failed to load node attributes")
				}

				return nil
			},
		)
		group.Go(n.worlds.Load)
		if err := group.Wait(); err != nil {
			return errors.WithMessage(err, "failed to load space tree")
		}

		return nil
	})
	// background loading thread
	group.Go(func() error {
		if err := n.chatService.Load(); err != nil {
			return errors.WithMessage(err, "failed to load chat service")
		}

		return nil
	})
	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to load universe")
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
