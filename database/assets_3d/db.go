package assets_3d

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
	getAssetsQuery = `SELECT * FROM asset_3d;`

	removeAssetByIDQuery   = `DELETE FROM asset_3d WHERE asset_3d_id = $1;`
	removeAssetsByIDsQuery = `DELETE FROM asset_3d WHERE asset_3d_id = ANY($1);`

	updateAssetMetaQuery    = `UPDATE asset_3d SET meta = $2 WHERE asset_3d_id = $1;`
	updateAssetOptionsQuery = `UPDATE asset_3d SET options = $2 WHERE asset_3d_id = $1;`

	upsertAssetQuery = `INSERT INTO asset_3d
							(asset_3d_id, meta, options, created_at, updated_at)
						VALUES
							($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
						ON CONFLICT (asset_3d_id)
						DO UPDATE SET
							meta = $2, options = $3, updated_at = CURRENT_TIMESTAMP;`
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

func (db *DB) GetAssets(ctx context.Context) ([]*entry.Asset3d, error) {
	var assets []*entry.Asset3d
	if err := pgxscan.Select(ctx, db.conn, &assets, getAssetsQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) UpsertAsset(ctx context.Context, asset3d *entry.Asset3d) error {
	if _, err := db.conn.Exec(ctx, upsertAssetQuery, asset3d.Asset3dID, asset3d.Meta, asset3d.Options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertAssets(ctx context.Context, assets3d []*entry.Asset3d) error {
	batch := &pgx.Batch{}
	for _, asset := range assets3d {
		batch.Queue(upsertAssetQuery, asset.Asset3dID, asset.Meta, asset.Options)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %s", assets3d[i].Asset3dID),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) RemoveAssetByID(ctx context.Context, asset3dID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeAssetByIDQuery, asset3dID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveAssetsByIDs(ctx context.Context, asset3dIDs []uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeAssetsByIDsQuery, asset3dIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateAssetMeta(ctx context.Context, asset3dID uuid.UUID, meta *entry.Asset3dMeta) error {
	if _, err := db.conn.Exec(ctx, updateAssetMetaQuery, asset3dID, meta); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, asset3dOptions *entry.Asset3dOptions) error {
	if _, err := db.conn.Exec(ctx, updateAssetOptionsQuery, asset3dID, asset3dOptions); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
