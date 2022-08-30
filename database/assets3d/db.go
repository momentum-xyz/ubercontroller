package assets3d

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
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB) *DB {
	return &DB{
		conn:   conn,
		common: commonDB,
	}
}

// TODO: implement
func (db *DB) Assets3dUpdateAssetName(ctx context.Context, asset3dID uuid.UUID, name string) error {
	return nil
}

// TODO: implement
func (db *DB) Assets3dUpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, asset3dOptions *entry.Asset3dOptions) error {
	return nil
}
