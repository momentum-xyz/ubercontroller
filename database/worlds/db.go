package worlds

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
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

func (db *DB) WorldsGetWorldIDs(ctx context.Context) ([]uuid.UUID, error) {
	return db.spaces.SpacesGetSpaceIDsByParentID(ctx, uuid.Nil)
}

func (db *DB) WorldsGetWorlds(ctx context.Context) ([]*entry.Space, error) {
	return db.spaces.SpacesGetSpacesByParentID(ctx, uuid.Nil)
}
