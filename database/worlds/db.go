package worlds

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getAllWorldIDsQuery = `SELECT object_id FROM object
                 			WHERE parent_id = (SELECT object_id FROM object WHERE object_id = parent_id)
                 			AND object_id != parent_id;`
	getWorldIDsQuery = `SELECT object_id FROM object
                 			WHERE parent_id = (SELECT object_id FROM object WHERE object_id = parent_id)
                 			AND object_id != parent_id 
                 			ORDER BY created_at `
	getWorldsQuery = `SELECT * FROM object
         					WHERE parent_id = (SELECT object_id FROM object WHERE object_id = parent_id)
         					AND object_id != parent_id;`
	getRecentWorldIDsQuery = `SELECT object_id FROM object
         					WHERE parent_id = (SELECT object_id FROM object WHERE object_id = parent_id)
         					AND object_id != parent_id
         					ORDER BY created_at DESC
							LIMIT 6;`
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

func (db *DB) GetAllWorldIDs(ctx context.Context) ([]umid.UMID, error) {
	var ids []umid.UMID
	if err := pgxscan.Select(ctx, db.conn, &ids, getAllWorldIDsQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return ids, nil
}

func (db *DB) GetWorldIDs(ctx context.Context, sortType universe.SortType, limit string) ([]umid.UMID, error) {
	limitQuery := " LIMIT " + limit + ";"
	var ids []umid.UMID
	if err := pgxscan.Select(ctx, db.conn, &ids, getWorldIDsQuery+string(sortType)+limitQuery); err != nil {
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

func (db *DB) GetRecentWorldIDs(ctx context.Context) ([]umid.UMID, error) {
	var worldIDs []umid.UMID
	if err := pgxscan.Select(ctx, db.conn, &worldIDs, getRecentWorldIDsQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return worldIDs, nil
}
