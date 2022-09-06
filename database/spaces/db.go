package spaces

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getSpaceByIDQuery          = `SELECT * FROM space WHERE space_id = $1;`
	getSpaceIDsByParentIDQuery = `SELECT space_id FROM space WHERE parent_id = $1;`
	getSpacesByParentIDQuery   = `SELECT * FROM space WHERE parent_id = $1;`
	removeSpaceByIDQuery       = `DELETE FROM space WHERE space_id = $1;`
	removeSpacesByIDsQuery     = `DELETE FROM space WHERE space_id IN ($1);`
	updateSpaceParentIDQuery   = `UPDATE space SET parent_id = $2 WHERE space_id = $1;`
	updateSpacePositionQuery   = `UPDATE space SET position = $2 WHERE space_id = $1;`
	updateSpaceOwnerIDQuery    = `UPDATE space SET owner_id = $2 WHERE space_id = $1;`
	updateSpaceAsset2dIDQuery  = `UPDATE space SET asset_2d_id = $2 WHERE space_id = $1;`
	updateSpaceAsset3dIDQuery  = `UPDATE space SET asset_3d_id = $2 WHERE space_id = $1;`
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

func (db *DB) SpacesGetSpaceByID(ctx context.Context, spaceID uuid.UUID) (*entry.Space, error) {
	var space entry.Space
	if err := pgxscan.Get(ctx, db.conn, &space, getSpaceByIDQuery, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &space, nil
}

func (db *DB) SpacesGetSpaceIDsByParentID(ctx context.Context, parentID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := pgxscan.Select(ctx, db.conn, &ids, getSpaceIDsByParentIDQuery, parentID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return ids, nil
}

func (db *DB) SpacesGetSpacesByParentID(ctx context.Context, parentID uuid.UUID) ([]*entry.Space, error) {
	var spaces []*entry.Space
	if err := pgxscan.Select(ctx, db.conn, &spaces, getSpacesByParentIDQuery, parentID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return spaces, nil
}

func (db *DB) SpacesRemoveSpaceByID(ctx context.Context, spaceID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeSpaceByIDQuery, spaceID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpacesRemoveSpacesByIDs(ctx context.Context, spaceIDs []uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeSpacesByIDsQuery, spaceIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpacesUpdateSpaceParentID(ctx context.Context, spaceID uuid.UUID, parentID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateSpaceParentIDQuery, spaceID, parentID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpacesUpdateSpacePosition(ctx context.Context, spaceID uuid.UUID, position *cmath.Vec3) error {
	if _, err := db.conn.Exec(ctx, updateSpacePositionQuery, spaceID, position); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpacesUpdateSpaceOwnerID(ctx context.Context, spaceID, ownerID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateSpaceOwnerIDQuery, spaceID, ownerID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpacesUpdateSpaceAsset2dID(ctx context.Context, spaceID uuid.UUID, asset2dID *uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateSpaceAsset2dIDQuery, spaceID, asset2dID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpacesUpdateSpaceAsset3dID(ctx context.Context, spaceID uuid.UUID, asset3dID *uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateSpaceAsset3dIDQuery, spaceID, asset3dID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
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
