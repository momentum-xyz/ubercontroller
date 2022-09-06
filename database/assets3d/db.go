package assets3d

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

var _ database.Assets3dDB = (*DB)(nil)

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
func (db *DB) Assets3dGetAssets(ctx context.Context) ([]*entry.Asset3d, error) {
	return nil, nil
}

// TODO: implement
func (db *DB) Assets3dUpsetAsset(ctx context.Context, asset3d *entry.Asset3d) error {
	return nil
}

// TODO: implement
func (db *DB) Assets3dUpsetAssets(ctx context.Context, assets3d []*entry.Asset3d) error {
	return nil
}

// TODO: implement
func (db *DB) Assets3dRemoveAssetByID(ctx context.Context, asset3dID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) Assets3dRemoveAssetsByIDs(ctx context.Context, asset3dIDs []uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) Assets3dUpdateAssetName(ctx context.Context, asset3dID uuid.UUID, name string) error {
	return nil
}

// TODO: implement
func (db *DB) Assets3dUpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, asset3dOptions *entry.Asset3dOptions) error {
	return nil
}
