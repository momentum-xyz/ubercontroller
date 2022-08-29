package space_types

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/database"
)

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
