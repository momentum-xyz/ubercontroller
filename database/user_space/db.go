package user_spaces

import (
	"context"
	"sync"

	"github.com/momentum-xyz/ubercontroller/utils/modify"

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
	getUserSpacesQuery                   = `SELECT * FROM user_space;`
	getUserSpacesBySpaceIDQuery          = `SELECT * FROM user_space WHERE space_id = $1;`
	getUserSpacesByUserIDQuery           = `SELECT * FROM user_space WHERE user_id = $1;`
	getUserSpacesByUserIDAndSpaceIDQuery = `SELECT * FROM user_space WHERE user_id = $1 AND space_id = $2;`
	getUserSpaceValueByIDQuery           = `SELECT value FROM user_space WHERE user_id = $1 AND space_id = $2;`
	getUserSpaceIndirectAdmins           = `SELECT getinderectspaceadmins($1);`

	checkIsIndirectAdminQuery = `WITH space_admins AS (
									SELECT getinderectspaceadmins($2) AS user_id
								)
								SELECT EXISTS(
									SELECT 1
									FROM space_admins
    								WHERE user_id = $1
    							) AS user_is_admin;`

	updateUserSpacesValueQuery = `UPDATE user_space SET value = $3 WHERE user_id = $1 AND space_id = $2;`

	removeUserSpaceByIDQuery = `DELETE FROM user_space WHERE user_id = $1 AND space_id = $2;`

	upsertUserSpaceQuery = `INSERT INTO user_space
											(user_id, space_id, value, created_at, updated_at)
										VALUES
											($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
										ON CONFLICT (user_id, space_id)
										DO UPDATE SET
											value = $3;`
)

var _ database.UserSpaceDB = (*DB)(nil)

type DB struct {
	conn   *pgxpool.Pool
	common database.CommonDB
	mu     sync.Mutex
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB) *DB {
	return &DB{
		conn:   conn,
		common: commonDB,
	}
}

func (db *DB) GetUserSpacesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserSpace, error) {
	var userSpaces []*entry.UserSpace
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpacesByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) GetUserSpacesBySpaceID(ctx context.Context, spaceID uuid.UUID) ([]*entry.UserSpace, error) {
	var userSpaces []*entry.UserSpace
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpacesBySpaceIDQuery, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) GetUserSpaceByID(ctx context.Context, userSpaceID entry.UserSpaceID) (*entry.UserSpace, error) {
	var userSpace *entry.UserSpace
	if err := pgxscan.Select(
		ctx, db.conn, &userSpace, getUserSpacesByUserIDAndSpaceIDQuery, userSpaceID.UserID, userSpaceID.SpaceID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpace, nil
}

func (db *DB) UserSpaceGetUserSpaceValueByID(
	ctx context.Context, userSpaceID entry.UserSpaceID,
) (*entry.UserSpaceValue, error) {
	var value entry.UserSpaceValue
	err := db.conn.QueryRow(ctx,
		getUserSpaceValueByIDQuery,
		userSpaceID.UserID,
		userSpaceID.SpaceID).Scan(&value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &value, nil
}

func (db *DB) GetUserSpaces(ctx context.Context) ([]*entry.UserSpace, error) {
	var userSpaces []*entry.UserSpace
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpacesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) GetSpaceIndirectAdmins(ctx context.Context, spaceID uuid.UUID) ([]*uuid.UUID, error) {
	var userIDs []*uuid.UUID
	if err := pgxscan.Select(ctx, db.conn, &userIDs, getUserSpaceIndirectAdmins, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userIDs, nil
}

func (db *DB) CheckIsUserIndirectSpaceAdmin(ctx context.Context, userID, spaceID uuid.UUID) (bool, error) {
	var isIndirectAdmin bool
	if err := db.conn.QueryRow(ctx, checkIsIndirectAdminQuery, userID, spaceID).
		Scan(&isIndirectAdmin); err != nil {
		return false, errors.WithMessage(err, "failed to query db")
	}
	return isIndirectAdmin, nil
}

func (db *DB) GetValueByID(ctx context.Context, userSpaceID entry.UserSpaceID) (*entry.UserSpaceValue, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) UpdateValueByID(ctx context.Context, userSpaceID entry.UserSpaceID, modifyFn modify.Fn[entry.UserSpaceValue]) (*entry.UserSpaceValue, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	value, err := db.UserSpaceGetUserSpaceValueByID(ctx, userSpaceID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get value by id")
	}

	value, err = modifyFn(value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify value")
	}

	if _, err := db.conn.Exec(
		ctx, updateUserSpacesValueQuery,
		userSpaceID.UserID, userSpaceID.SpaceID,
		value,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return value, nil
}

func (db *DB) UpsertUserSpace(ctx context.Context, userSpace *entry.UserSpace) error {
	if _, err := db.conn.Exec(
		ctx, upsertUserSpaceQuery, userSpace.SpaceID, userSpace.UserID,
		userSpace.Value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertUserSpaces(ctx context.Context, userSpaces []*entry.UserSpace) error {
	batch := &pgx.Batch{}
	for _, userSpace := range userSpaces {
		batch.Queue(
			upsertUserSpaceQuery, userSpace.UserID, userSpace.SpaceID,
			userSpace.Value,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs,
				errors.WithMessagef(
					err,
					"failed to exec db for user id: %s, space id: %s",
					userSpaces[i].UserID, userSpaces[i].SpaceID,
				),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) RemoveUserSpace(ctx context.Context, userSpaces *entry.UserSpace) error {
	if _, err := db.conn.Exec(ctx, removeUserSpaceByIDQuery, userSpaces.UserID, userSpaces.SpaceID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}

	return nil
}

func (db *DB) RemoveUserSpaces(ctx context.Context, userSpaces []*entry.UserSpace) error {
	batch := &pgx.Batch{}
	for _, userSpace := range userSpaces {
		batch.Queue(removeUserSpaceByIDQuery, userSpace.UserID, userSpace.SpaceID)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs,
				errors.WithMessagef(
					err,
					"failed to exec db for user id: %s, space id: %s",
					userSpaces[i].UserID, userSpaces[i].SpaceID,
				),
			)
		}
	}

	return errs.ErrorOrNil()
}
