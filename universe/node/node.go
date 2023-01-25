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
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/object"
	"github.com/momentum-xyz/ubercontroller/universe/streamchat"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Node = (*Node)(nil)

type Node struct {
	*object.Object
	cfg        *config.Config
	ctx        context.Context
	log        *zap.SugaredLogger
	db         database.DB
	router     *gin.Engine
	httpServer *http.Server

	worlds         universe.Worlds
	assets2d       universe.Assets2d
	assets3d       universe.Assets3d
	objectTypes    universe.ObjectTypes
	userTypes      universe.UserTypes
	attributeTypes universe.AttributeTypes
	plugins        universe.Plugins

	userObjects *userObjects

	nodeAttributes       *nodeAttributes // WARNING: the Node is sharing the same mutex ("Mu") with it
	userAttributes       *userAttributes
	userUserAttributes   *userUserAttributes
	objectUserAttributes *objectUserAttributes

	objectIDToWorld *generic.SyncMap[uuid.UUID, universe.World] // TODO: introduce GC for lost Worlds and Objects

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
	objectTypes universe.ObjectTypes,
	userTypes universe.UserTypes,
	attributeTypes universe.AttributeTypes,
) *Node {
	node := &Node{
		Object:          object.NewObject(id, db, nil),
		db:              db,
		worlds:          worlds,
		assets2d:        assets2D,
		assets3d:        assets3D,
		plugins:         plugins,
		objectTypes:     objectTypes,
		userTypes:       userTypes,
		attributeTypes:  attributeTypes,
		objectIDToWorld: generic.NewSyncMap[uuid.UUID, universe.World](0),
	}
	node.userObjects = newUserObjects(node)
	node.nodeAttributes = newNodeAttributes(node)
	node.userAttributes = newUserAttributes(node)
	node.userUserAttributes = newUserUserAttributes(node)
	node.objectUserAttributes = newObjectUserAttributes(node)

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

	return n.Object.Initialize(ctx)
}

func (n *Node) ToObject() universe.Object {
	return n.Object
}

func (n *Node) GetUserObjects() universe.UserObjects {
	return n.userObjects
}

func (n *Node) GetNodeAttributes() universe.NodeAttributes {
	return n.nodeAttributes
}

func (n *Node) GetUserAttributes() universe.UserAttributes {
	return n.userAttributes
}

func (n *Node) GetUserUserAttributes() universe.UserUserAttributes {
	return n.userUserAttributes
}

func (n *Node) GetObjectUserAttributes() universe.ObjectUserAttributes {
	return n.objectUserAttributes
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

func (n *Node) GetObjectTypes() universe.ObjectTypes {
	return n.objectTypes
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
		group.Go(n.objectTypes.Load)
		group.Go(n.plugins.Load)
		if err := group.Wait(); err != nil {
			return errors.WithMessage(err, "failed to load additional data")
		}

		// third stage
		group, _ = errgroup.WithContext(n.ctx)
		group.Go(
			func() error {
				nodeEntry, err := n.db.GetNodesDB().GetNode(n.ctx)
				if err != nil {
					return errors.WithMessage(err, "failed to get node")
				}
				if err := n.LoadFromEntry(nodeEntry.Object, false); err != nil {
					return errors.WithMessage(err, "failed to load node from entry")
				}

				if err := n.GetNodeAttributes().Load(); err != nil {
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

		if err := n.attributeTypes.Save(); err != nil {
			errsMu.Lock()
			defer errsMu.Unlock()
			errs = multierror.Append(errs, errors.WithMessage(err, "failed to save attribute types"))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := n.plugins.Save(); err != nil {
			errsMu.Lock()
			defer errsMu.Unlock()
			errs = multierror.Append(errs, errors.WithMessage(err, "failed to save plugins"))
		}
	}()

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

		if err := n.objectTypes.Save(); err != nil {
			errsMu.Lock()
			defer errsMu.Unlock()
			errs = multierror.Append(errs, errors.WithMessage(err, "failed to save object types"))
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
