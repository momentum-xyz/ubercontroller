package assets_3d

import (
	"context"

	"github.com/momentum-xyz/ubercontroller/universe"
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
	getAssetsQuery     = `SELECT * FROM asset_3d;`
	getUserAssetsQuery = `SELECT * FROM asset_3d_user;`

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

	updateAssetMetaQuery     = `UPDATE asset_3d SET meta = $2, updated_at = CURRENT_TIMESTAMP WHERE asset_3d_id = $1;`
	updateUserAssetMetaQuery = `UPDATE asset_3d_user SET meta = $3, updated_at = CURRENT_TIMESTAMP WHERE asset_3d_id = $1 AND user_id = $2;`

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

func (db *DB) GetUserAssets(ctx context.Context) ([]*entry.UserAsset3d, error) {
	var assets []*entry.UserAsset3d
	if err := pgxscan.Select(ctx, db.conn, &assets, getUserAssetsQuery); err != nil {
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

func (db *DB) UpsertUserAsset(ctx context.Context, userAsset3d *entry.UserAsset3d) error {
	if _, err := db.conn.Exec(ctx, upsertUserAssetQuery, userAsset3d.Asset3dID, userAsset3d.UserID, userAsset3d.Meta, userAsset3d.Private); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}

	return nil
}

func (db *DB) UpsertAssets(ctx context.Context, assets3d []*entry.Asset3d, userAssets3d []*entry.UserAsset3d) error {
	batch := &pgx.Batch{}
	for _, asset := range assets3d {
		batch.Queue(upsertAssetQuery, asset.Asset3dID, asset.Meta, asset.Options)
	}
	for _, userAsset := range userAssets3d {
		batch.Queue(upsertUserAssetQuery, userAsset.Asset3dID, userAsset.UserID, userAsset.Meta, userAsset.Private)
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

func (db *DB) RemoveUserAssetByID(ctx context.Context, assetUserID universe.AssetUserIDPair) error {
	if _, err := db.conn.Exec(ctx, removeAssetByIDQuery, assetUserID.AssetID, assetUserID.UserID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

// func (db *DB) RemoveAssetsByIDs(ctx context.Context, asset3dIDs []umid.UMID) error {
// 	if _, err := db.conn.Exec(ctx, removeAssetsByIDsQuery, asset3dIDs); err != nil {
// 		return errors.WithMessage(err, "failed to exec db")
// 	}
// 	return nil
// }

func (db *DB) UpdateAssetMeta(ctx context.Context, assetID umid.UMID, meta *entry.Asset3dMeta) error {
	if _, err := db.conn.Exec(ctx, updateAssetMetaQuery, assetID, meta); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateUserAssetMeta(ctx context.Context, assetUserID universe.AssetUserIDPair, meta *entry.Asset3dMeta) error {
	if _, err := db.conn.Exec(ctx, updateUserAssetMetaQuery, assetUserID.AssetID, assetUserID.UserID, meta); err != nil {
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
