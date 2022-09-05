package spaces

import (
	"context"
	"github.com/pkg/errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getSpaceByID = `SELECT * FROM space WHERE space_id = $1;`
)

var _ database.SpacesDB = (*DB)(nil)

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
func (db *DB) SpacesGetSpaceByID(ctx context.Context, spaceID uuid.UUID) (*entry.Space, error) {
	var space entry.Space
	if err := pgxscan.Get(ctx, db.conn, &space, getSpaceByID, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &space, nil
}

// TODO: implement
func (db *DB) SpacesGetSpaceIDsByParentID(ctx context.Context, parentID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

// TODO: implement
func (db *DB) SpacesGetSpacesByParentID(ctx context.Context, parentID uuid.UUID) ([]*entry.Space, error) {
	return nil, nil
}

// TODO: implement
func (db *DB) SpacesRemoveSpaceByID(ctx context.Context, spaceID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesRemoveSpacesByIDs(ctx context.Context, spaceIDs []uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceParentID(ctx context.Context, spaceID uuid.UUID, parentID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpacePosition(ctx context.Context, spaceID uuid.UUID, position *cmath.Vec3) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceOwnerID(ctx context.Context, spaceID, ownerID uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceAsset2dID(ctx context.Context, spaceID uuid.UUID, asset2dID *uuid.UUID) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpdateSpaceAsset3dID(ctx context.Context, spaceID uuid.UUID, asset3dID *uuid.UUID) error {
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

// TODO: implement
func (db *DB) SpacesUpsertSpace(ctx context.Context, space *entry.Space) error {
	return nil
}

// TODO: implement
func (db *DB) SpacesUpsertSpaces(ctx context.Context, spaces []*entry.Space) error {
	return nil
}
