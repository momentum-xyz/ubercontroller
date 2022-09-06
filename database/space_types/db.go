package space_types

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
	getSpaceTypesQuery = `SELECT * FROM space_type;`
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

// TODO: implement
func (db *DB) SpaceTypesUpsetSpaceType(ctx context.Context, spaceType *entry.SpaceType) error {
	return nil
}

// TODO: implement
func (db *DB) SpaceTypesUpsetSpaceTypes(ctx context.Context, spaceTypes []*entry.SpaceType) error {
	return nil
}

// TODO: implement
func (db *DB) SpaceTypesRemoveSpaceTypeByID(ctx context.Context, spaceTypeID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpaceTypesRemoveSpaceTypesByIDs(ctx context.Context, spaceTypeIDs []uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpaceTypesUpdateSpaceTypeName(ctx context.Context, spaceTypeID uuid.UUID, name string) error {
	return nil
}

// TODO: implement
func (db *DB) SpaceTypesUpdateSpaceTypeCategoryName(ctx context.Context, spaceTypeID uuid.UUID, categoryName string) error {
	return nil
}

// TODO: impelement
func (db *DB) SpaceTypesUpdateSpaceTypeDescription(ctx context.Context, spaceTypeID uuid.UUID, description *string) error {
	return nil
}

// TODO: implement
func (db *DB) SpaceTypesUpdateSpaceTypeOptions(ctx context.Context, spaceTypeID uuid.UUID, options *entry.SpaceOptions) error {
	return nil
}
