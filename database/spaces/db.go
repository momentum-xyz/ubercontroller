package spaces

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
	getSpaceByIDQuery           = `SELECT * FROM space WHERE space_id = $1;`
	getSpaceIDsByParentIDQuery  = `SELECT space_id FROM space WHERE parent_id = $1;`
	getSpacesByParentIDQuery    = `SELECT * FROM space WHERE parent_id = $1;`
	removeSpaceByIDQuery        = `DELETE FROM space WHERE space_id = $1;`
	removeSpacesByIDsQuery      = `DELETE FROM space WHERE space_id IN ($1);`
	updateSpaceParentIDQuery    = `UPDATE space SET parent_id = $2 WHERE space_id = $1;`
	updateSpacePositionQuery    = `UPDATE space SET position = $2 WHERE space_id = $1;`
	updateSpaceOwnerIDQuery     = `UPDATE space SET owner_id = $2 WHERE space_id = $1;`
	updateSpaceAsset2dIDQuery   = `UPDATE space SET asset_2d_id = $2 WHERE space_id = $1;`
	updateSpaceAsset3dIDQuery   = `UPDATE space SET asset_3d_id = $2 WHERE space_id = $1;`
	updateSpaceSpaceTypeIDQuery = `UPDATE space SET space_type_id = $2 WHERE space_id = $1;`
	updateSpaceOptionsQuery     = `UPDATE space SET options = $2 WHERE space_id = $1;`
	upsertSpaceQuery            = `INSERT INTO space
    									(space_id, space_type_id, owner_id, parent_id, asset_2d_id,
    									 asset_3d_id, options, position, created_at, updated_at)
									VALUES
									    ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
									ON CONFLICT (space_id)
									DO UPDATE SET
										space_type_id = $2, owner_id = $3, parent_id = $4, asset_2d_id = $5,
									    asset_3d_id = $6, options = $7, position = $8, updated_at = CURRENT_TIMESTAMP;`
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

func (db *DB) SpacesUpdateSpaceSpaceTypeID(ctx context.Context, spaceID, spaceTypeID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateSpaceSpaceTypeIDQuery, spaceID, spaceTypeID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpacesUpdateSpaceOptions(ctx context.Context, spaceID uuid.UUID, options *entry.SpaceOptions) error {
	if _, err := db.conn.Exec(ctx, updateSpaceOptionsQuery, spaceID, options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpacesUpsertSpace(ctx context.Context, space *entry.Space) error {
	if _, err := db.conn.Exec(ctx, upsertSpaceQuery,
		space.SpaceID, space.SpaceTypeID, space.OwnerID, space.ParentID, space.Asset2dID, space.Asset3dID,
		space.Options, space.Position); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) SpacesUpsertSpaces(ctx context.Context, spaces []*entry.Space) error {
	batch := &pgx.Batch{}
	for _, space := range spaces {
		batch.Queue(upsertSpaceQuery, space.SpaceID, space.SpaceTypeID, space.OwnerID,
			space.ParentID, space.Asset2dID, space.Asset3dID, space.Options, space.Position)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to exec db for: %s", spaces[i].SpaceID))
		}
	}

	return errs.ErrorOrNil()
}
