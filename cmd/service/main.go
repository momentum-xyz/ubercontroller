package main

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	assets2dDB "github.com/momentum-xyz/ubercontroller/database/assets_2d"
	assets3dDB "github.com/momentum-xyz/ubercontroller/database/assets_3d"
	attributesDB "github.com/momentum-xyz/ubercontroller/database/attribute"
	commonDB "github.com/momentum-xyz/ubercontroller/database/common"
	nodeAttributesDB "github.com/momentum-xyz/ubercontroller/database/node_attributes"
	nodesDB "github.com/momentum-xyz/ubercontroller/database/nodes"
	pluginsDB "github.com/momentum-xyz/ubercontroller/database/plugin"
	spaceAttributesDB "github.com/momentum-xyz/ubercontroller/database/space_attributes"
	spaceTypesDB "github.com/momentum-xyz/ubercontroller/database/space_types"
	spaceUserAttributesDB "github.com/momentum-xyz/ubercontroller/database/space_user_attributes"
	spacesDB "github.com/momentum-xyz/ubercontroller/database/spaces"
	userAttributesDB "github.com/momentum-xyz/ubercontroller/database/user_attributes"
	userTypesDB "github.com/momentum-xyz/ubercontroller/database/user_types"
	userUserAttributesDB "github.com/momentum-xyz/ubercontroller/database/user_user_attributes"
	usersDB "github.com/momentum-xyz/ubercontroller/database/users"
	worldsDB "github.com/momentum-xyz/ubercontroller/database/worlds"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/database/db"
	"github.com/momentum-xyz/ubercontroller/database/migrations"
	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/universe/assets_2d"
	"github.com/momentum-xyz/ubercontroller/universe/assets_3d"
	"github.com/momentum-xyz/ubercontroller/universe/attributes"
	"github.com/momentum-xyz/ubercontroller/universe/node"
	"github.com/momentum-xyz/ubercontroller/universe/plugins"
	"github.com/momentum-xyz/ubercontroller/universe/space_types"
	"github.com/momentum-xyz/ubercontroller/universe/worlds"
)

var log = logger.L()

func main() {
	if err := run(); err != nil {
		log.Fatal(errors.WithMessage(err, "failed to run service"))
	}
}

func run() error {
	cfg := config.GetConfig()

	ctx := context.WithValue(context.Background(), types.ContextLoggerKey, log)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := api.Initialize(ctx, cfg); err != nil {
		return errors.WithMessage(err, "failed to initialize api")
	}

	pool, err := createDBConnection(ctx, &cfg.Postgres)
	if err != nil {
		return errors.WithMessage(err, "failed to create db connection")
	}
	defer pool.Close()

	db, err := createDB(pool)
	if err != nil {
		return errors.WithMessage(err, "failed to create db")
	}

	node, err := createNode(ctx, cfg, db)
	if err != nil {
		return errors.WithMessage(err, "failed to create node")
	}

	if err := node.Load(); err != nil {
		return errors.WithMessagef(err, "failed to load node: %s", node.GetID())
	}

	defer func() {
		if err := node.Stop(); err != nil {
			log.Error(errors.WithMessagef(err, "failed to stop node: %s", node.GetID()))
		}
	}()

	if err := node.Run(); err != nil {
		return errors.WithMessagef(err, "failed to run node: %s", node.GetID())
	}

	return nil
}

func createNode(ctx context.Context, cfg *config.Config, db database.DB) (universe.Node, error) {
	worlds := worlds.NewWorlds(db)
	assets2d := assets_2d.NewAssets2d(db)
	assets3d := assets_3d.NewAssets3d(db)
	spaceTypes := space_types.NewSpaceTypes(db)
	attributes := attributes.NewAttributes(db)
	plugins := plugins.NewPlugins(db)

	nodeEntry, err := db.NodesGetNode(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get node")
	}

	node := node.NewNode(*nodeEntry.SpaceID, cfg, db, worlds, assets2d, assets3d, spaceTypes, attributes, plugins)
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
	if err := attributes.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize attributes")
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
	config, err := cfg.GenConfig()
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
		spaceTypesDB.NewDB(conn, common),
		userTypesDB.NewDB(conn, common),
		attributesDB.NewDB(conn, common),
		pluginsDB.NewDB(conn, common),
		spaceAttributesDB.NewDB(conn, common),
		spaceUserAttributesDB.NewDB(conn, common),
		userAttributesDB.NewDB(conn, common),
		userUserAttributesDB.NewDB(conn, common),
		nodeAttributesDB.NewDB(conn, common),
	), nil
}
