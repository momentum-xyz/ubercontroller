package db

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/database"
)

var _ database.DB = (*DB)(nil)

type DB struct {
	conn   *pgxpool.Pool
	common database.CommonDB
	database.NodesDB
	database.WorldsDB
	database.SpacesDB
	database.UsersDB
	database.Assets2dDB
	database.Assets3dDB
	database.SpaceTypesDB
}

func NewDB(
	conn *pgxpool.Pool,
	common database.CommonDB,
	nodes database.NodesDB,
	worlds database.WorldsDB,
	spaces database.SpacesDB,
	users database.UsersDB,
	assets2d database.Assets2dDB,
	assets3d database.Assets3dDB,
	spaceTypes database.SpaceTypesDB,
) *DB {
	return &DB{
		conn:         conn,
		common:       common,
		NodesDB:      nodes,
		WorldsDB:     worlds,
		SpacesDB:     spaces,
		UsersDB:      users,
		Assets2dDB:   assets2d,
		Assets3dDB:   assets3d,
		SpaceTypesDB: spaceTypes,
	}
}
