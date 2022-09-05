package common

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/database"
)

var _ database.CommonDB = (*DB)(nil)

type DB struct {
	conn *pgxpool.Pool
}

func NewDB(conn *pgxpool.Pool) *DB {
	return &DB{
		conn: conn,
	}
}
