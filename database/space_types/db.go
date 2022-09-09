package space_types

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
	getSpaceTypesQuery               = `SELECT * FROM space_type;`
	updateSpaceTypeNameQuery         = `UPDATE space_type SET space_type_name = $2 WHERE space_type_id = $1;`
	updateSpaceTypeCategoryNameQuery = `UPDATE space_type SET category_name = $2 WHERE space_type_id = $1;`
	updateSpaceTypeDescriptionQuery  = `UPDATE space_type SET description = $2 WHERE space_type_id = $1;`
	updateSpaceTypeOptionsQuery      = `UPDATE space_type SET options = $2 WHERE space_type_id = $1;`
	removeSpaceTypeByIDQuery         = `DELETE FROM space_type WHERE space_type_id = $1;`
	removeSpaceTypesByIDsQuery       = `DELETE FROM space_type WHERE space_type_id IN ($1);`
	upsertSpaceTypeQuery             = `INSERT INTO space_type
											(space_type_id, asset_2d_id, asset_3d_id, space_type_name,
											 category_name, description, options, created_at, updated_at)
										VALUES
											($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
										ON CONFLICT (space_type_id)
										DO UPDATE SET
											asset_2d_id = $2, asset_3d_id = $3, space_type_name = $4, category_name = $5,
											description = $6, options = $7;`
)

var _ database.SpaceTypesDB = (*DB)(nil)

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

func (db *DB) SpaceTypesGetSpaceTypes(ctx context.Context) ([]*entry.SpaceType, error) {
	var spaceTypes []*entry.SpaceType
	if err := pgxscan.Select(ctx, db.conn, &spaceTypes, getSpaceTypesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return spaceTypes, nil
}

func (db *DB) SpaceTypesUpsertSpaceType(ctx context.Context, spaceType *entry.SpaceType) error {
	if _, err := db.conn.Exec(ctx, upsertSpaceTypeQuery, spaceType.SpaceTypeID, spaceType.Asset2dID, spaceType.Asset3dID,
		spaceType.SpaceTypeName, spaceType.CategoryName, spaceType.Description, spaceType.Options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceTypesUpsertSpaceTypes(ctx context.Context, spaceTypes []*entry.SpaceType) error {
	batch := &pgx.Batch{}
	for _, spaceType := range spaceTypes {
		batch.Queue(upsertSpaceTypeQuery, spaceType.SpaceTypeID, spaceType.Asset2dID, spaceType.Asset3dID,
			spaceType.SpaceTypeName, spaceType.CategoryName, spaceType.Description, spaceType.Options)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) SpaceTypesRemoveSpaceTypeByID(ctx context.Context, spaceTypeID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeSpaceTypeByIDQuery, spaceTypeID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceTypesRemoveSpaceTypesByIDs(ctx context.Context, spaceTypeIDs []uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeSpaceTypesByIDsQuery, spaceTypeIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceTypesUpdateSpaceTypeName(ctx context.Context, spaceTypeID uuid.UUID, name string) error {
	if _, err := db.conn.Exec(ctx, updateSpaceTypeNameQuery, spaceTypeID, name); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceTypesUpdateSpaceTypeCategoryName(ctx context.Context, spaceTypeID uuid.UUID, categoryName string) error {
	if _, err := db.conn.Exec(ctx, updateSpaceTypeCategoryNameQuery, spaceTypeID, categoryName); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceTypesUpdateSpaceTypeDescription(ctx context.Context, spaceTypeID uuid.UUID, description *string) error {
	if _, err := db.conn.Exec(ctx, updateSpaceTypeDescriptionQuery, spaceTypeID, description); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpaceTypesUpdateSpaceTypeOptions(ctx context.Context, spaceTypeID uuid.UUID, options *entry.SpaceOptions) error {
	if _, err := db.conn.Exec(ctx, updateSpaceTypeOptionsQuery, spaceTypeID, options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
