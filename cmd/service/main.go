package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	assets2dDB "github.com/momentum-xyz/ubercontroller/database/assets2d"
	assets3dDB "github.com/momentum-xyz/ubercontroller/database/assets3d"
	commonDB "github.com/momentum-xyz/ubercontroller/database/common"
	"github.com/momentum-xyz/ubercontroller/database/db"
	nodesDB "github.com/momentum-xyz/ubercontroller/database/nodes"
	spaceTypesDB "github.com/momentum-xyz/ubercontroller/database/space_types"
	spacesDB "github.com/momentum-xyz/ubercontroller/database/spaces"
	usersDB "github.com/momentum-xyz/ubercontroller/database/users"
	worldsDB "github.com/momentum-xyz/ubercontroller/database/worlds"
	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/assets2d"
	"github.com/momentum-xyz/ubercontroller/universe/assets3d"
	"github.com/momentum-xyz/ubercontroller/universe/node"
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

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), types.ContextLoggerKey, log))
	defer cancel()

	pool, err := createDBConnection(ctx, &cfg.Postgres)
	if err != nil {
		return errors.WithMessage(err, "failed to create db connection")
	}
	defer pool.Close()

	db, err := createDB(pool)
	if err != nil {
		return errors.WithMessage(err, "failed to create db")
	}

	node, err := createNode(ctx, db)
	if err != nil {
		return errors.WithMessage(err, "failed to create node")
	}

	if err := node.Load(ctx); err != nil {
		return errors.WithMessage(err, "failed to load node")
	}

	defer func() {
		if err := node.Stop(); err != nil {
			log.Error(errors.WithMessage(err, "failed to stop node"))
		}
	}()

	return node.Run()
}

func createNode(ctx context.Context, db database.DB) (universe.Node, error) {
	worlds := worlds.NewWorlds(db)
	if err := worlds.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize worlds")
	}

	assets2d := assets2d.NewAssets2D(db)
	if err := assets2d.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize assets 2d")
	}

	assets3d := assets3d.NewAssets3D(db)
	if err := assets3d.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize assets 3d")
	}

	spaceTypes := space_types.NewSpaceTypes(db)
	if err := spaceTypes.Initialize(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to initialize space types")
	}

	node := node.NewNode(uuid.Nil, db, worlds, assets2d, assets3d, spaceTypes)
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
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create db pool")
	}
	return pool, nil
}

func createDB(conn *pgxpool.Pool) (database.DB, error) {
	common := commonDB.NewDB(conn)
	return db.NewDB(
		conn,
		common,
		nodesDB.NewDB(conn, common),
		worldsDB.NewDB(conn, common),
		spacesDB.NewDB(conn, common),
		usersDB.NewDB(conn, common),
		assets2dDB.NewDB(conn, common),
		assets3dDB.NewDB(conn, common),
		spaceTypesDB.NewDB(conn, common),
	), nil
}
