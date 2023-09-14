package node

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	influx_api "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/arbitrum_nova_adapter"
	"github.com/momentum-xyz/ubercontroller/mplugin"
	"github.com/momentum-xyz/ubercontroller/pkg/media"
	"github.com/momentum-xyz/ubercontroller/seed"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/universe/object"
	"github.com/momentum-xyz/ubercontroller/universe/streamchat"
	"github.com/momentum-xyz/ubercontroller/universe/user"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.Node = (*Node)(nil)

type Node struct {
	*object.Object
	ctx        types.NodeContext
	cfg        *config.Config
	log        *zap.SugaredLogger
	db         database.DB
	router     *gin.Engine
	httpServer *http.Server

	worlds         universe.Worlds
	assets2d       universe.Assets2d
	assets3d       universe.Assets3d
	activities     universe.Activities
	media          *media.Media
	objectTypes    universe.ObjectTypes
	userTypes      universe.UserTypes
	attributeTypes universe.AttributeTypes
	plugins        universe.Plugins

	userObjects *userObjects

	nodeAttributes       *nodeAttributes // WARNING: the Node is sharing the same mutex ("Mu") with it
	userAttributes       *userAttributes
	userUserAttributes   *userUserAttributes
	objectUserAttributes *objectUserAttributes

	objectIDToWorld *generic.SyncMap[umid.UMID, universe.World] // TODO: introduce GC for lost Worlds and Objects

	chatService *streamchat.StreamChat

	pluginController    *mplugin.PluginController
	corePluginInterface *mplugin.PluginInterface

	influx influx_api.WriteAPIBlocking
}

func NewNode(
	id umid.UMID,
	db database.DB,
	worlds universe.Worlds,
	assets2D universe.Assets2d,
	assets3D universe.Assets3d,
	activities universe.Activities,
	media *media.Media,
	plugins universe.Plugins,
	objectTypes universe.ObjectTypes,
	userTypes universe.UserTypes,
	attributeTypes universe.AttributeTypes,
) *Node {
	node := &Node{
		Object:          object.NewObject(id, db, nil, media),
		db:              db,
		worlds:          worlds,
		assets2d:        assets2D,
		assets3d:        assets3D,
		activities:      activities,
		media:           media,
		plugins:         plugins,
		objectTypes:     objectTypes,
		userTypes:       userTypes,
		attributeTypes:  attributeTypes,
		objectIDToWorld: generic.NewSyncMap[umid.UMID, universe.World](0),
	}
	node.userObjects = newUserObjects(node)
	node.nodeAttributes = newNodeAttributes(node)
	node.userAttributes = newUserAttributes(node)
	node.userUserAttributes = newUserUserAttributes(node)
	node.objectUserAttributes = newObjectUserAttributes(node)

	node.chatService = streamchat.NewStreamChat()
	node.pluginController = mplugin.NewPluginController(id)

	return node
}

func (n *Node) Initialize(ctx types.NodeContext) error {
	n.ctx = ctx
	n.log = ctx.Logger()
	n.cfg = ctx.Config()

	consoleWriter := zapcore.Lock(os.Stdout)
	gin.DefaultWriter = consoleWriter

	//TODO: hash salt once it is present in the DB
	utils.SetAnonymizer(n.GetID(), umid.Nil)

	r := gin.New()
	r.Use(gin.LoggerWithWriter(consoleWriter, "/health"))
	r.Use(gin.RecoveryWithWriter(consoleWriter))

	n.router = r
	n.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", n.cfg.Settings.Address, n.cfg.Settings.Port),
		Handler: n.router,
	}

	if n.cfg.Common.AllowCORS {
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		}))
	}

	if err := n.chatService.Initialize(ctx); err != nil {
		return errors.WithMessage(err, "failed to initialize chat service")
	}

	return n.ToObject().Initialize(ctx)
}

func (n *Node) GetConfig() *config.Config {
	return n.cfg
}

func (n *Node) GetMedia() *media.Media {
	return n.media
}

func (n *Node) GetLogger() *zap.SugaredLogger {
	return n.log
}

func (n *Node) LoadUser(userID umid.UMID) (universe.User, error) {
	newUser := user.NewUser(userID, n.db)
	if err := newUser.Initialize(n.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize user: %s", userID)
	}

	if err := newUser.Load(); err != nil {
		return nil, errors.WithMessagef(err, "failed to load user: %s", userID)
	}

	return newUser, nil
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

func (n *Node) GetWorldsByOwnerID(userID umid.UMID) map[umid.UMID]universe.World {
	n.Children.Mu.RLock()
	defer n.Children.Mu.RUnlock()

	worlds := make(map[umid.UMID]universe.World, len(n.Children.Data))
	for id, world := range n.worlds.GetWorlds() {

		if world.GetOwnerID() == userID {
			worlds[id] = world
		}
	}

	return worlds
}

func (n *Node) GetAssets2d() universe.Assets2d {
	return n.assets2d
}

func (n *Node) GetAssets3d() universe.Assets3d {
	return n.assets3d
}

func (n *Node) GetActivities() universe.Activities {
	return n.activities
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
	temporaryUserTypeID, err := common.GetGuestUserTypeID()
	if err != nil {
		return errors.WithMessage(err, "failed to get guest user type id")
	}

	users, err := n.db.GetUsersDB().GetUsersByUserType(n.ctx, temporaryUserTypeID)
	for _, user := range users {
		loadedUser, _ := n.LoadUser(user.UserID)

		isTemporaryUser, err := loadedUser.IsTemporaryUser()
		if err != nil {
			return errors.WithMessagef(err, "failed to assess if user is temporary user: %s", loadedUser.GetID())
		}

		if isTemporaryUser {
			ok, err := loadedUser.SetOfflineTimer()
			if !ok || err != nil {
				return errors.WithMessage(err, "failed to set offline timer")
			}
		}
	}

	if err := n.worlds.Run(); err != nil {
		return errors.WithMessage(err, "failed to run worlds")
	}
	n.SetEnabled(true)

	//harvester.Initialise(ctx, log, cfg, pool)
	//if cfg.Arbitrum.ArbitrumMOMTokenAddress != "" {
	//	arbitrumAdapter := arbitrum_nova_adapter.NewArbitrumNovaAdapter(cfg)
	//	arbitrumAdapter.Run()
	//	if err := harvester.GetInstance().RegisterAdapter(arbitrumAdapter); err != nil {
	//		return errors.WithMessage(err, "failed to register arbitrum adapter")
	//	}
	//}
	//err = harvester.SubscribeAllWallets(ctx, harvester.GetInstance(), cfg, pool)
	//if err != nil {
	//	log.Error(err)
	//}

	/**
	Simplified version of harvester
	*/
	if n.cfg.Arbitrum.MOMTokenAddress != "" {
		logger := n.GetLogger()
		adapter := arbitrum_nova_adapter.NewArbitrumNovaAdapter(n.cfg, logger)
		adapter.Run()

		pgConfig, err := n.cfg.Postgres.GenConfig(logger.Desugar())
		pool, err := pgxpool.ConnectConfig(context.Background(), pgConfig)
		if err != nil {
			n.log.Fatal("failed to create db pool")
		}
		defer pool.Close()

		t := harvester.NewTable(pool, adapter, n.Listener)
		t.Run()
	}

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
	n.log.Infof("Loading node: %s...", n.GetID())

	// main loading thread
	group, ctx := errgroup.WithContext(n.ctx)
	group.Go(
		func() error {
			// first stage
			group, _ := errgroup.WithContext(ctx)
			group.Go(n.GetPlugins().Load)
			group.Go(n.GetAssets2d().Load)
			group.Go(n.GetAssets3d().Load)
			group.Go(n.GetUserTypes().Load)
			group.Go(n.GetAttributeTypes().Load)
			if err := group.Wait(); err != nil {
				return errors.WithMessage(err, "failed to load basic data")
			}

			// second stage
			group, _ = errgroup.WithContext(ctx)
			group.Go(n.GetObjectAttributes().Load)
			group.Go(n.GetObjectTypes().Load)
			group.Go(n.GetActivities().Load)
			if err := group.Wait(); err != nil {
				return errors.WithMessage(err, "failed to load additional data")
			}

			// third stage
			group, _ = errgroup.WithContext(ctx)
			group.Go(n.load)
			group.Go(n.GetWorlds().Load)
			if err := group.Wait(); err != nil {
				return errors.WithMessage(err, "failed to load universe tree")
			}

			return nil
		},
	)
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

	if err := n.GetPlugins().Save(); err != nil {
		return errors.WithMessage(err, "failed to save plugins")
	}
	if err := n.GetAttributeTypes().Save(); err != nil {
		return errors.WithMessage(err, "failed to save AttributeTypes")
	}
	if err := n.GetNodeAttributes().Save(); err != nil {
		return errors.WithMessage(err, "failed to save NodeAttributes")
	}
	if err := n.GetAssets2d().Save(); err != nil {
		return errors.WithMessage(err, "failed to save assets 2d")
	}
	if err := n.GetAssets3d().Save(); err != nil {
		return errors.WithMessage(err, "failed to save assets 3d")
	}
	if err := n.GetUserTypes().Save(); err != nil {
		return errors.WithMessage(err, "failed to save user types")
	}
	if err := n.GetObjectTypes().Save(); err != nil {
		return errors.WithMessage(err, "failed to save object types")
	}

	objectType, ok := n.GetObjectTypes().GetObjectType(umid.MustParse(seed.NodeObjectTypeID))
	if !ok {
		return errors.New("failed to get node object_type by UMID")
	}

	if err := n.SetObjectType(objectType, false); err != nil {
		return errors.WithMessage(err, "failed to set object_type to node")
	}

	//if err := n.SetParent(n, false); err != nil {
	//	return errors.WithMessage(err, "failed to set parent for node")
	//}

	if err := n.save(); err != nil {
		return errors.WithMessage(err, "failed to save node")
	}

	return nil
	//
	//var errs *multierror.Error
	//var errsMu sync.Mutex
	//addError := func(err error, msg string) {
	//	errsMu.Lock()
	//	defer errsMu.Unlock()
	//
	//	errs = multierror.Append(errs, errors.WithMessage(err, msg))
	//}
	//
	//addWGTask := func(wg *sync.WaitGroup, task func() error, errMsg string) {
	//	wg.Add(1)
	//
	//	go func() {
	//		defer wg.Done()
	//
	//		if err := task(); err != nil {
	//			addError(err, errMsg)
	//		}
	//	}()
	//}
	//
	//// first stage
	//wg := &sync.WaitGroup{}
	//addWGTask(wg, n.GetPlugins().Save, "failed to save plugins")
	//addWGTask(wg, n.GetAssets2d().Save, "failed to save assets 2d")
	//addWGTask(wg, n.GetAssets3d().Save, "failed to save assets 3d")
	//addWGTask(wg, n.GetUserTypes().Save, "failed to save user types")
	//addWGTask(wg, n.GetAttributeTypes().Save, "failed to save attribute types")
	//wg.Wait()
	//
	//// second stage
	//wg = &sync.WaitGroup{}
	//addWGTask(wg, n.GetObjectTypes().Save, "failed to save object types")
	//wg.Wait()
	//
	//// third stage
	//wg = &sync.WaitGroup{}
	//addWGTask(wg, n.save, "failed to save node data")
	//wg.Wait()
	//
	//wg = &sync.WaitGroup{}
	//addWGTask(wg, n.GetWorlds().Save, "failed to save worlds")
	//wg.Wait()
	//
	//n.log.Infof("Node saved: %s", n.GetID())
	//
	//return errs.ErrorOrNil()
}

func (n *Node) save() error {
	if err := n.ToObject().Save(); err != nil {
		return errors.WithMessage(err, "failed to save node object")
	}
	if err := n.GetNodeAttributes().Save(); err != nil {
		return errors.WithMessage(err, "failed to save node attributes")
	}
	return nil
}

func (n *Node) load() error {
	n.log.Infof("Loading node data: %s...", n.GetID())

	group, ctx := errgroup.WithContext(n.ctx)
	group.Go(n.GetNodeAttributes().Load)
	group.Go(
		func() error {
			nodeEntry, err := n.db.GetNodesDB().GetNode(ctx)
			if err != nil {
				return errors.WithMessage(err, "failed to get node")
			}
			if err := n.LoadFromEntry(nodeEntry.Object, false); err != nil {
				return errors.WithMessage(err, "failed to load node from entry")
			}
			return nil
		},
	)
	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to load node data")
	}

	n.log.Infof("Node data loaded: %s", n.GetID())

	return nil
}
