package assets_2d

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getAssetsQuery = `SELECT * FROM asset_2d;`

	removeAssetByIDQuery   = `DELETE FROM asset_2d WHERE asset_2d_id = $1;`
	removeAssetsByIDsQuery = `DELETE FROM asset_2d WHERE asset_2d_id IN ($1);`

	updateAssetMetaQuery    = `UPDATE asset_2d SET meta = $2 WHERE asset_2d_id = $1;`
	updateAssetOptionsQuery = `UPDATE asset_2d SET options = $2 WHERE asset_2d_id = $1;`

	upsertAssetQuery = `INSERT INTO asset_2d
							(asset_2d_id, meta, options, created_at, updated_at)
						VALUES
							($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
						ON CONFLICT (asset_2d_id)
						DO UPDATE SET
							meta = $2, options = $3, updated_at = CURRENT_TIMESTAMP;`
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

func (db *DB) Assets2dUpsertAsset(ctx context.Context, asset2d *entry.Asset2d) error {
	if _, err := db.conn.Exec(ctx, upsertAssetQuery, asset2d.Asset2dID, asset2d.Meta, asset2d.Options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) Assets2dUpsertAssets(ctx context.Context, assets2d []*entry.Asset2d) error {
	batch := &pgx.Batch{}
	for _, asset := range assets2d {
		batch.Queue(upsertAssetQuery, asset.Asset2dID, asset.Meta, asset.Options)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %s", assets2d[i].Asset2dID),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) Assets2dRemoveAssetByID(ctx context.Context, asset2dID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeAssetByIDQuery, asset2dID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) Assets2dRemoveAssetsByIDs(ctx context.Context, asset2dIDs []uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeAssetsByIDsQuery, asset2dIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) Assets2dUpdateAssetMeta(ctx context.Context, asset2dID uuid.UUID, meta *entry.Meta) error {
	if _, err := db.conn.Exec(ctx, updateAssetMetaQuery, asset2dID, meta); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) Assets2dUpdateAssetOptions(ctx context.Context, asset2dID uuid.UUID, asset2dOptions *entry.Asset2dOptions) error {
	if _, err := db.conn.Exec(ctx, updateAssetOptionsQuery, asset2dID, asset2dOptions); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
