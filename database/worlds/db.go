package worlds

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getWorldIDsQuery = `SELECT object_id FROM object
                 			WHERE parent_id = (SELECT object_id FROM object WHERE object_id = parent_id)
                 			AND object_id != parent_id;`
	getWorldsQuery = `SELECT * FROM object
         					WHERE parent_id = (SELECT object_id FROM object WHERE object_id = parent_id)
         					AND object_id != parent_id;`
)

var _ database.WorldsDB = (*DB)(nil)

type DB struct {
	conn    *pgxpool.Pool
	common  database.CommonDB
	objects database.ObjectsDB
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB, objectsDB database.ObjectsDB) *DB {
	return &DB{
		conn:    conn,
		common:  commonDB,
		objects: objectsDB,
	}
}

func (db *DB) GetWorldIDs(ctx context.Context) ([]umid.UMID, error) {
	var ids []umid.UMID
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
