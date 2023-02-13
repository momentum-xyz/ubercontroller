package objects

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getObjectByIDQuery          = `SELECT * FROM object WHERE object_id = $1;`
	getObjectIDsByParentIDQuery = `SELECT object_id FROM object WHERE parent_id = $1;`
	getObjectsByParentIDQuery   = `SELECT * FROM object WHERE parent_id = $1;`

	upsertObjectQuery = `INSERT INTO object
    						(object_id, object_type_id, owner_id, parent_id, asset_2d_id,
    						asset_3d_id, options, position, created_at, updated_at)
						VALUES
							($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
						ON CONFLICT (object_id)
						DO UPDATE SET
							object_type_id = $2, owner_id = $3, parent_id = $4, asset_2d_id = $5,
							asset_3d_id = $6, options = $7, position = $8, updated_at = CURRENT_TIMESTAMP;`

	updateObjectParentIDQuery     = `UPDATE object SET parent_id = $2, updated_at = CURRENT_TIMESTAMP WHERE object_id = $1;`
	updateObjectPositionQuery     = `UPDATE object SET position = $2, updated_at = CURRENT_TIMESTAMP WHERE object_id = $1;`
	updateObjectOwnerIDQuery      = `UPDATE object SET owner_id = $2, updated_at = CURRENT_TIMESTAMP WHERE object_id = $1;`
	updateObjectAsset2dIDQuery    = `UPDATE object SET asset_2d_id = $2, updated_at = CURRENT_TIMESTAMP WHERE object_id = $1;`
	updateObjectAsset3dIDQuery    = `UPDATE object SET asset_3d_id = $2, updated_at = CURRENT_TIMESTAMP WHERE object_id = $1;`
	updateObjectObjectTypeIDQuery = `UPDATE object SET object_type_id = $2, updated_at = CURRENT_TIMESTAMP WHERE object_id = $1;`
	updateObjectOptionsQuery      = `UPDATE object SET options = $2, updated_at = CURRENT_TIMESTAMP WHERE object_id = $1;`

	removeObjectByIDQuery   = `DELETE FROM object WHERE object_id = $1;`
	removeObjectsByIDsQuery = `DELETE FROM object WHERE object_id = ANY($1);`
)

var _ database.ObjectsDB = (*DB)(nil)

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

func (db *DB) GetObjectByID(ctx context.Context, objectID uuid.UUID) (*entry.Object, error) {
	var object entry.Object
	if err := pgxscan.Get(ctx, db.conn, &object, getObjectByIDQuery, objectID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &object, nil
}

func (db *DB) GetObjectIDsByParentID(ctx context.Context, parentID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := pgxscan.Select(ctx, db.conn, &ids, getObjectIDsByParentIDQuery, parentID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return ids, nil
}

func (db *DB) GetObjectsByParentID(ctx context.Context, parentID uuid.UUID) ([]*entry.Object, error) {
	var objects []*entry.Object
	if err := pgxscan.Select(ctx, db.conn, &objects, getObjectsByParentIDQuery, parentID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return objects, nil
}

func (db *DB) RemoveObjectByID(ctx context.Context, objectID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeObjectByIDQuery, objectID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectsByIDs(ctx context.Context, objectIDs []uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeObjectsByIDsQuery, objectIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectParentID(ctx context.Context, objectID uuid.UUID, parentID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateObjectParentIDQuery, objectID, parentID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectPosition(ctx context.Context, objectID uuid.UUID, position *cmath.ObjectPosition) error {
	if _, err := db.conn.Exec(ctx, updateObjectPositionQuery, objectID, position); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectOwnerID(ctx context.Context, objectID, ownerID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateObjectOwnerIDQuery, objectID, ownerID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectAsset2dID(ctx context.Context, objectID uuid.UUID, asset2dID *uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateObjectAsset2dIDQuery, objectID, asset2dID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectAsset3dID(ctx context.Context, objectID uuid.UUID, asset3dID *uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateObjectAsset3dIDQuery, objectID, asset3dID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectObjectTypeID(ctx context.Context, objectID, objectTypeID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateObjectObjectTypeIDQuery, objectID, objectTypeID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectOptions(ctx context.Context, objectID uuid.UUID, options *entry.ObjectOptions) error {
	if _, err := db.conn.Exec(ctx, updateObjectOptionsQuery, objectID, options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertObject(ctx context.Context, object *entry.Object) error {
	if _, err := db.conn.Exec(
		ctx, upsertObjectQuery,
		object.ObjectID, object.ObjectTypeID, object.OwnerID, object.ParentID, object.Asset2dID, object.Asset3dID,
		object.Options, object.Position,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertObjects(ctx context.Context, objects []*entry.Object) error {
	batch := &pgx.Batch{}
	for _, object := range objects {
		batch.Queue(
			upsertObjectQuery, object.ObjectID, object.ObjectTypeID, object.OwnerID,
			object.ParentID, object.Asset2dID, object.Asset3dID, object.Options, object.Position,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %s", objects[i].ObjectID),
			)
		}
	}

	return errs.ErrorOrNil()
}
