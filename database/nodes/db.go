package nodes

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getNodeQuery = `SELECT * FROM space WHERE space_id = parent_id;`
)

var _ database.NodesDB = (*DB)(nil)

type DB struct {
	conn   *pgxpool.Pool
	common database.CommonDB
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB) *DB {
	return &DB{
		conn:   conn,
		common: commonDB,
	}
}

func (db *DB) GetNode(ctx context.Context) (*entry.Node, error) {
	var node entry.Node
	if err := pgxscan.Get(ctx, db.conn, &node, getNodeQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &node, nil
}
