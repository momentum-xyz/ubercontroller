package main

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/universe/common"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/zakaria-chahboun/cute"
	"go.uber.org/zap"

	assets2dDB "github.com/momentum-xyz/ubercontroller/database/assets_2d"
	assets3dDB "github.com/momentum-xyz/ubercontroller/database/assets_3d"
	attributesTypeDB "github.com/momentum-xyz/ubercontroller/database/attribute_types"
	commonDB "github.com/momentum-xyz/ubercontroller/database/common"
	nodeAttributesDB "github.com/momentum-xyz/ubercontroller/database/node_attributes"
	nodesDB "github.com/momentum-xyz/ubercontroller/database/nodes"
	pluginsDB "github.com/momentum-xyz/ubercontroller/database/plugins"
	spaceAttributesDB "github.com/momentum-xyz/ubercontroller/database/space_attributes"
	spaceTypesDB "github.com/momentum-xyz/ubercontroller/database/space_types"
	spaceUserAttributesDB "github.com/momentum-xyz/ubercontroller/database/space_user_attributes"
	spacesDB "github.com/momentum-xyz/ubercontroller/database/spaces"
	userAttributesDB "github.com/momentum-xyz/ubercontroller/database/user_attributes"
	userSpaceDB "github.com/momentum-xyz/ubercontroller/database/user_space"
	userTypesDB "github.com/momentum-xyz/ubercontroller/database/user_types"
	userUserAttributesDB "github.com/momentum-xyz/ubercontroller/database/user_user_attributes"
	usersDB "github.com/momentum-xyz/ubercontroller/database/users"
	worldsDB "github.com/momentum-xyz/ubercontroller/database/worlds"
	"github.com/momentum-xyz/ubercontroller/utils"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/database/db"
	"github.com/momentum-xyz/ubercontroller/database/migrations"
	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/assets_2d"
	"github.com/momentum-xyz/ubercontroller/universe/assets_3d"
	"github.com/momentum-xyz/ubercontroller/universe/attribute_types"
	"github.com/momentum-xyz/ubercontroller/universe/node"
	"github.com/momentum-xyz/ubercontroller/universe/plugins"
	"github.com/momentum-xyz/ubercontroller/universe/space_types"
	"github.com/momentum-xyz/ubercontroller/universe/user_types"
	"github.com/momentum-xyz/ubercontroller/universe/worlds"
)

var log = logger.L()

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(errors.WithMessage(err, "failed to run service"))
	}
}

func run(ctx context.Context) error {
	cfg := config.GetConfig()

	ctx = context.WithValue(ctx, types.LoggerContextKey, log)
	ctx = context.WithValue(ctx, types.ConfigContextKey, cfg)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//todo: change to pool
	message.InitBuilder(20, 1024*32)
	if err := universe.InitializeIDs(
		uuid.MustParse("f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0"),
		uuid.MustParse("86DC3AE7-9F3D-42CB-85A3-A71ABC3C3CB8"),
	); err != nil {
		return errors.WithMessage(err, "failed to initialize universe")
	}
	if err := common.Initialize(ctx); err != nil {
		return errors.WithMessage(err, "failed to initialize common package")
	}

	pool, err := createDBConnection(ctx, &cfg.Postgres)
	if err != nil {
		return errors.WithMessage(err, "failed to create db connection")
	}
	defer pool.Close()

	tm1 := time.Now()
	db, err := createDB(pool)
	if err != nil {
		return errors.WithMessage(err, "failed to create db")
	}

	node, err := createNode(ctx, db)
	if err != nil {
		return errors.WithMessage(err, "failed to create node")
	}

	if err := node.Load(); err != nil {
		return errors.WithMessagef(err, "failed to load node: %s", node.GetID())
	}
	tm2 := time.Now()

	cute.SetTitleColor(cute.BrightGreen)
	cute.SetMessageColor(cute.BrightBlue)
	cute.Println("Node loaded", "Loading time:", tm2.Sub(tm1))

	defer func() {
		if err := node.Stop(); err != nil {
			log.Error(errors.WithMessagef(err, "failed to stop node: %s", node.GetID()))
		}
	}()

	if err := node.Run(); err != nil {
		return errors.WithMessagef(err, "failed to run node: %s", node.GetID())
	}

	cute.SetTitleColor(cute.BrightPurple)
	cute.SetMessageColor(cute.BrightBlue)
	cute.Println("Node stopped", "That's all folks!")

	return nil
}

func createNode(ctx context.Context, db database.DB) (universe.Node, error) {
	worlds := worlds.NewWorlds(db)
	assets2d := assets_2d.NewAssets2d(db)
	assets3d := assets_3d.NewAssets3d(db)
	plugins := plugins.NewPlugins(db)
	spaceTypes := space_types.NewSpaceTypes(db)
	userTypes := user_types.NewUserTypes(db)
	attributeTypes := attribute_types.NewAttributeTypes(db)

	nodeEntry, err := db.NodesGetNode(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get node")
	}

	node := node.NewNode(
		nodeEntry.SpaceID,
		db,
		worlds,
		assets2d,
		assets3d,
		plugins,
		spaceTypes,
		userTypes,
		attributeTypes,
	)
	universe.InitializeNode(node)

	if err := worlds.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize worlds")
	}
	if err := assets2d.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize assets 2d")
	}
	if err := assets3d.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize assets 3d")
	}
	if err := spaceTypes.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize space types")
	}
	if err := userTypes.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize user types")
	}
	if err := attributeTypes.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize attribute types")
	}
	if err := plugins.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize plugins")
	}
	if err := node.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize node")
	}

	return node, nil
}

func createDBConnection(ctx context.Context, cfg *config.Postgres) (*pgxpool.Pool, error) {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	config, err := cfg.GenConfig(log.Desugar())
	if err != nil {
		return nil, errors.WithMessage(err, "failed to gen postgres config")
	}

	if err := migrations.MigrateDatabase(ctx, cfg); err != nil {
		return nil, errors.WithMessage(err, "failed to migrate database")
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create db pool")
	}
	return pool, nil
}

func createDB(conn *pgxpool.Pool) (database.DB, error) {
	common := commonDB.NewDB(conn)
	spaces := spacesDB.NewDB(conn, common)
	return db.NewDB(
		conn,
		common,
		nodesDB.NewDB(conn, common),
		worldsDB.NewDB(conn, common, spaces),
		spaces,
		usersDB.NewDB(conn, common),
		assets2dDB.NewDB(conn, common),
		assets3dDB.NewDB(conn, common),
		pluginsDB.NewDB(conn, common),
		userSpaceDB.NewDB(conn, common),
		spaceTypesDB.NewDB(conn, common),
		userTypesDB.NewDB(conn, common),
		attributesTypeDB.NewDB(conn, common),
		nodeAttributesDB.NewDB(conn, common),
		spaceAttributesDB.NewDB(conn, common),
		spaceUserAttributesDB.NewDB(conn, common),
		userAttributesDB.NewDB(conn, common),
		userUserAttributesDB.NewDB(conn, common),
	), nil
}
