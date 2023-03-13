package object_types

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
	getObjectTypesQuery = `SELECT * FROM object_type;`

	upsertObjectTypeQuery = `INSERT INTO object_type
								(object_type_id, asset_2d_id, asset_3d_id, object_type_name,
								category_name, description, options, created_at, updated_at)
							VALUES
								($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
							ON CONFLICT (object_type_id)
							DO UPDATE SET
								asset_2d_id = $2, asset_3d_id = $3, object_type_name = $4, category_name = $5,
								description = $6, options = $7, updated_at = CURRENT_TIMESTAMP;`

	updateObjectTypeNameQuery         = `UPDATE object_type SET object_type_name = $2, updated_at = CURRENT_TIMESTAMP WHERE object_type_id = $1;`
	updateObjectTypeCategoryNameQuery = `UPDATE object_type SET category_name = $2, updated_at = CURRENT_TIMESTAMP WHERE object_type_id = $1;`
	updateObjectTypeDescriptionQuery  = `UPDATE object_type SET description = $2, updated_at = CURRENT_TIMESTAMP WHERE object_type_id = $1;`
	updateObjectTypeOptionsQuery      = `UPDATE object_type SET options = $2, updated_at = CURRENT_TIMESTAMP WHERE object_type_id = $1;`

	removeObjectTypeByIDQuery   = `DELETE FROM object_type WHERE object_type_id = $1;`
	removeObjectTypesByIDsQuery = `DELETE FROM object_type WHERE object_type_id = ANY($1);`
)

var _ database.ObjectTypesDB = (*DB)(nil)

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

func (db *DB) GetObjectTypes(ctx context.Context) ([]*entry.ObjectType, error) {
	var objectTypes []*entry.ObjectType
	if err := pgxscan.Select(ctx, db.conn, &objectTypes, getObjectTypesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return objectTypes, nil
}

func (db *DB) UpsertObjectType(ctx context.Context, objectType *entry.ObjectType) error {
	if _, err := db.conn.Exec(
		ctx, upsertObjectTypeQuery,
		objectType.ObjectTypeID, objectType.Asset2dID, objectType.Asset3dID, objectType.ObjectTypeName,
		objectType.CategoryName, objectType.Description, objectType.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertObjectTypes(ctx context.Context, objectTypes []*entry.ObjectType) error {
	batch := &pgx.Batch{}
	for _, objectType := range objectTypes {
		batch.Queue(
			upsertObjectTypeQuery, objectType.ObjectTypeID, objectType.Asset2dID, objectType.Asset3dID,
			objectType.ObjectTypeName, objectType.CategoryName, objectType.Description, objectType.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %+v", objectTypes[i].ObjectTypeID),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) RemoveObjectTypeByID(ctx context.Context, objectTypeID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeObjectTypeByIDQuery, objectTypeID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveObjectTypesByIDs(ctx context.Context, objectTypeIDs []uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeObjectTypesByIDsQuery, objectTypeIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectTypeName(ctx context.Context, objectTypeID uuid.UUID, name string) error {
	if _, err := db.conn.Exec(ctx, updateObjectTypeNameQuery, objectTypeID, name); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectTypeCategoryName(ctx context.Context, objectTypeID uuid.UUID, categoryName string) error {
	if _, err := db.conn.Exec(ctx, updateObjectTypeCategoryNameQuery, objectTypeID, categoryName); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectTypeDescription(ctx context.Context, objectTypeID uuid.UUID, description *string) error {
	if _, err := db.conn.Exec(ctx, updateObjectTypeDescriptionQuery, objectTypeID, description); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateObjectTypeOptions(ctx context.Context, objectTypeID uuid.UUID, options *entry.ObjectOptions) error {
	if _, err := db.conn.Exec(ctx, updateObjectTypeOptionsQuery, objectTypeID, options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
