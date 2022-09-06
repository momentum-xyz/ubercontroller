package assets2d

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
	getAssetsQuery = `SELECT * FROM asset_2d;`
)

var _ database.Assets2dDB = (*DB)(nil)

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

func (db *DB) Assets2dGetAssets(ctx context.Context) ([]*entry.Asset2d, error) {
	var assets []*entry.Asset2d
	if err := pgxscan.Select(ctx, db.conn, &assets, getAssetsQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

// TODO: implement
func (db *DB) Assets2dUpsetAsset(ctx context.Context, asset2d *entry.Asset2d) error {
	return nil
}

// TODO: implement
func (db *DB) Assets2dUpsetAssets(ctx context.Context, assets2d []*entry.Asset2d) error {
	return nil
}

// TODO: implement
func (db *DB) Assets2dRemoveAssetByID(ctx context.Context, asset2dID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) Assets2dRemoveAssetsByIDs(ctx context.Context, asset2dIDs []uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) Assets2dUpdateAssetName(ctx context.Context, asset2dID uuid.UUID, name string) error {
	return nil
}

// TODO: implement
func (db *DB) Assets2dUpdateAssetOptions(ctx context.Context, asset2dID uuid.UUID, asset2dOptions *entry.Asset2dOptions) error {
	return nil
}
