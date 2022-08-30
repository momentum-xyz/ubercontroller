package spaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

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

// TODO: implement
func (db *DB) SpacesGetSpacesByParentID(ctx context.Context, parentID uuid.UUID) ([]*entry.Space, error) {
	return nil, nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceParentID(ctx context.Context, spaceID, parentID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpacePosition(ctx context.Context, spaceID uuid.UUID, position cmath.Vec3) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceOwnerID(ctx context.Context, spaceID, ownerID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceAsset2dID(ctx context.Context, spaceID, asset2dID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceAsset3dID(ctx context.Context, spaceID, asset3dID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceSpaceTypeID(ctx context.Context, spaceID, spaceTypeID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceOptions(ctx context.Context, spaceID uuid.UUID, options *entry.SpaceOptions) error {
	return nil
}
