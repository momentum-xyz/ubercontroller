package worlds

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getWorldIDsQuery = `SELECT space_id FROM space WHERE space_id != parent_id
						AND parent_id = (SELECT space_id FROM space WHERE space_id = parent_id);`
	getWorldsQuery = `SELECT * FROM space WHERE space_id != parent_id
						AND parent_id = (SELECT space_id FROM space WHERE space_id = parent_id);`
)

var _ database.WorldsDB = (*DB)(nil)

type DB struct {
	conn   *pgxpool.Pool
	common database.CommonDB
	spaces database.ObjectsDB
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB, spacesDB database.ObjectsDB) *DB {
	return &DB{
		conn:   conn,
		common: commonDB,
		spaces: spacesDB,
	}
}

func (db *DB) GetWorldIDs(ctx context.Context) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := pgxscan.Select(ctx, db.conn, &ids, getWorldIDsQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return ids, nil
}

func (db *DB) GetWorlds(ctx context.Context) ([]*entry.Object, error) {
	var worlds []*entry.Object
	if err := pgxscan.Select(ctx, db.conn, &worlds, getWorldsQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return worlds, nil
}
