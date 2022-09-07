package assets3d

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getAssetsQuery          = `SELECT * FROM asset_3d;`
	updateAssetNameQuery    = `UPDATE asset_3d SET asset_3d_name = $2 WHERE asset_3d_id = $1;`
	updateAssetOptionsQuery = `UPDATE asset_3d SET options = $2 WHERE asset_3d_id = $1;`
	removeAssetByIDQuery    = `DELETE FROM asset_3d WHERE asset_3d_id = $1;`
	removeAssetsByIDsQuery  = `DELETE FROM asset_3d WHERE asset_3d_id IN ($1);`
	upsertAssetQuery        = `INSERT INTO asset_3d
									(asset_3d_id, asset_3d_name, options)
								VALUES
									($1, $2, $3)
								ON CONFLICT (asset_2d_id)
								DO UPDATE SET
									asset_2d_name = $2,
									options = $3;`
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

func (db *DB) Assets3dGetAssets(ctx context.Context) ([]*entry.Asset3d, error) {
	var assets []*entry.Asset3d
	if err := pgxscan.Select(ctx, db.conn, &assets, getAssetsQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return assets, nil
}

func (db *DB) Assets3dUpsertAsset(ctx context.Context, asset3d *entry.Asset3d) error {
	if _, err := db.conn.Exec(ctx, upsertAssetQuery, asset3d.Asset3dID, asset3d.Name, asset3d.Options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) Assets3dUpsertAssets(ctx context.Context, assets3d []*entry.Asset3d) error {
	batch := &pgx.Batch{}
	for _, asset := range assets3d {
		batch.Queue(upsertAssetQuery, asset.Asset3dID, asset.Name, asset.Options)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			return errors.WithMessage(err, "failed to exec db batch")
		}
	}

	return nil
}

func (db *DB) Assets3dRemoveAssetByID(ctx context.Context, asset3dID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeAssetByIDQuery, asset3dID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) Assets3dRemoveAssetsByIDs(ctx context.Context, asset3dIDs []uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeAssetsByIDsQuery, asset3dIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) Assets3dUpdateAssetName(ctx context.Context, asset3dID uuid.UUID, name string) error {
	if _, err := db.conn.Exec(ctx, updateAssetNameQuery, asset3dID, name); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) Assets3dUpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, asset3dOptions *entry.Asset3dOptions) error {
	if _, err := db.conn.Exec(ctx, updateAssetOptionsQuery, asset3dID, asset3dOptions); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
