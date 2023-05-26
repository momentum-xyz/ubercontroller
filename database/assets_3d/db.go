package assets_3d

import (
	"context"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getAssetsQuery = `SELECT au.asset_3d_id, au.user_id, au.meta || a.meta as meta, a.options, au.is_private, a.created_at, a.updated_at
										FROM asset_3d_user as au
										JOIN asset_3d as a USING (asset_3d_id);`

	upsertAssetQuery = `INSERT INTO asset_3d
							(asset_3d_id, meta, options, created_at, updated_at)
						VALUES
							($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
						ON CONFLICT (asset_3d_id)
						DO NOTHING;`
	upsertUserAssetQuery = `INSERT INTO asset_3d_user
							(asset_3d_id, user_id, meta, is_private, created_at, updated_at)
						VALUES
							($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
						ON CONFLICT (asset_3d_id, user_id)
						DO UPDATE SET
							meta = $3, is_private = $4, updated_at = CURRENT_TIMESTAMP;`

	updateAssetMetaQuery = `UPDATE asset_3d_user SET meta = $2, updated_at = CURRENT_TIMESTAMP WHERE asset_3d_id = $1 AND user_id = $2;`

	updateAssetOptionsQuery = `UPDATE asset_3d SET options = $2, updated_at = CURRENT_TIMESTAMP WHERE asset_3d_id = $1;`

	removeAssetByIDQuery   = `DELETE FROM asset_3d_user WHERE asset_3d_id = $1 AND user_id = $2;`
	removeAssetsByIDsQuery = `DELETE FROM asset_3d_user WHERE asset_3d_id = ANY($1) AND user_id = $2;`
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
	// TODO pick only 'category' and 'type' fields from meta map
	if _, err := db.conn.Exec(ctx, upsertAssetQuery, asset3d.Asset3dID, asset3d.Meta, asset3d.Options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}

	if _, err := db.conn.Exec(ctx, upsertUserAssetQuery, asset3d.Asset3dID, asset3d.UserID, asset3d.Meta, asset3d.Private); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}

	return nil
}

func (db *DB) UpsertAssets(ctx context.Context, assets3d []*entry.Asset3d) error {
	batch := &pgx.Batch{}
	for _, asset := range assets3d {
		// TODO pick only 'category' and 'type' fields from meta map
		batch.Queue(upsertAssetQuery, asset.Asset3dID, asset.Meta, asset.Options)

		batch.Queue(upsertUserAssetQuery, asset.Asset3dID, asset.UserID, asset.Meta, asset.Private)
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

// TODO add userID
func (db *DB) RemoveAssetByID(ctx context.Context, asset3dID umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removeAssetByIDQuery, asset3dID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

// TODO add userID
func (db *DB) RemoveAssetsByIDs(ctx context.Context, asset3dIDs []umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removeAssetsByIDsQuery, asset3dIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

// TODO add userID
func (db *DB) UpdateAssetMeta(ctx context.Context, asset3dID umid.UMID, meta *entry.Asset3dMeta) error {
	if _, err := db.conn.Exec(ctx, updateAssetMetaQuery, asset3dID, meta); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateAssetOptions(ctx context.Context, asset3dID umid.UMID, asset3dOptions *entry.Asset3dOptions) error {
	if _, err := db.conn.Exec(ctx, updateAssetOptionsQuery, asset3dID, asset3dOptions); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
