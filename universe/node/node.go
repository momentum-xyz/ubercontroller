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
	"github.com/sasha-s/go-deadlock"
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
	"github.com/momentum-xyz/ubercontroller/universe/space"
	"github.com/momentum-xyz/ubercontroller/universe/streamchat"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Node = (*Node)(nil)

type Node struct {
	*space.Space
	cfg        *config.Config
	ctx        context.Context
	log        *zap.SugaredLogger
	db         database.DB
	router     *gin.Engine
	httpServer *http.Server

	//mu             sync.RWMutex
	mu             deadlock.RWMutex
	nodeAttributes *nodeAttributes // WARNING: the Node is sharing the same mutex ("mu") with it

	worlds         universe.Worlds
	assets2d       universe.Assets2d
	assets3d       universe.Assets3d
	spaceTypes     universe.SpaceTypes
	userTypes      universe.UserTypes
	attributeTypes universe.AttributeTypes
	plugins        universe.Plugins

	spaceIDToWorld *generic.SyncMap[uuid.UUID, universe.World] // TODO: introduce GC for lost Worlds and Spaces

	chatService *streamchat.StreamChat

	pluginController    *mplugin.PluginController
	corePluginInterface *mplugin.PluginInterface

	influx influx_api.WriteAPIBlocking
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
	node := &Node{
		Space:          space.NewSpace(id, db, nil),
		db:             db,
		worlds:         worlds,
		assets2d:       assets2D,
		assets3d:       assets3D,
		plugins:        plugins,
		spaceTypes:     spaceTypes,
		userTypes:      userTypes,
		attributeTypes: attributeTypes,
		spaceIDToWorld: generic.NewSyncMap[uuid.UUID, universe.World](0),
	}
	node.nodeAttributes = newNodeAttributes(node)

	return node
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

func (n *Node) ToSpace() universe.Space {
	return n.Space
}

func (n *Node) GetNodeAttributes() universe.Attributes[entry.AttributeID] {
	return n.nodeAttributes
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

// TODO: investigate how to load with limited database connection pool
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
			return errors.WithMessage(err, "failed to load universe tree")
		}

		return nil
	})
	// background loading thread
	group.Go(n.chatService.Load)
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
