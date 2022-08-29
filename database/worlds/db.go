package worlds

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/universe"
)

type DB struct {
	conn   *pgxpool.Pool
	common database.CommonDB
	spaces database.SpacesDB
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB, spacesDB database.SpacesDB) *DB {
	return &DB{
		conn:   conn,
		common: commonDB,
		spaces: spacesDB,
	}
}

func (db *DB) WorldsGetWorlds(ctx context.Context) ([]universe.SpaceEntry, error) {
	return db.spaces.SpacesGetSpacesByParentID(ctx, uuid.Nil)
}
