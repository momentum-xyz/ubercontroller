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
	getUserSpaceIndirectAdmins           = `SELECT GetIndirectSpaceAdmins($1);`

	checkIsIndirectAdminQuery = `WITH space_admins AS (
									SELECT GetIndirectSpaceAdmins($2) AS user_id
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

var _ database.UserObjectDB = (*DB)(nil)

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

func (db *DB) GetUserObjectsByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserObject, error) {
	var userSpaces []*entry.UserObject
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpacesByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) GetUserObjectsByObjectID(ctx context.Context, spaceID uuid.UUID) ([]*entry.UserObject, error) {
	var userSpaces []*entry.UserObject
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpacesBySpaceIDQuery, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) GetUserObjectByID(ctx context.Context, userSpaceID entry.UserObjectID) (*entry.UserObject, error) {
	var userSpace *entry.UserObject
	if err := pgxscan.Select(
		ctx, db.conn, &userSpace, getUserSpacesByUserIDAndSpaceIDQuery, userSpaceID.UserID, userSpaceID.ObjectID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpace, nil
}

func (db *DB) UserSpaceGetUserSpaceValueByID(
	ctx context.Context, userSpaceID entry.UserObjectID,
) (*entry.UserObjectValue, error) {
	var value entry.UserObjectValue
	err := db.conn.QueryRow(ctx,
		getUserSpaceValueByIDQuery,
		userSpaceID.UserID,
		userSpaceID.ObjectID).Scan(&value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &value, nil
}

func (db *DB) GetUserObjects(ctx context.Context) ([]*entry.UserObject, error) {
	var userSpaces []*entry.UserObject
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpacesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) GetObjectIndirectAdmins(ctx context.Context, spaceID uuid.UUID) ([]*uuid.UUID, error) {
	var userIDs []*uuid.UUID
	if err := pgxscan.Select(ctx, db.conn, &userIDs, getUserSpaceIndirectAdmins, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userIDs, nil
}

func (db *DB) CheckIsUserIndirectObjectAdmin(ctx context.Context, userID, spaceID uuid.UUID) (bool, error) {
	var isIndirectAdmin bool
	if err := db.conn.QueryRow(ctx, checkIsIndirectAdminQuery, userID, spaceID).
		Scan(&isIndirectAdmin); err != nil {
		return false, errors.WithMessage(err, "failed to query db")
	}
	return isIndirectAdmin, nil
}

func (db *DB) GetValueByID(ctx context.Context, userSpaceID entry.UserObjectID) (*entry.UserObjectValue, error) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) UpdateValueByID(ctx context.Context, userSpaceID entry.UserObjectID, modifyFn modify.Fn[entry.UserObjectValue]) (*entry.UserObjectValue, error) {
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
		userSpaceID.UserID, userSpaceID.ObjectID,
		value,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return value, nil
}

func (db *DB) UpsertUserObject(ctx context.Context, userSpace *entry.UserObject) error {
	if _, err := db.conn.Exec(
		ctx, upsertUserSpaceQuery, userSpace.ObjectID, userSpace.UserID,
		userSpace.Value,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertUserObjects(ctx context.Context, userSpaces []*entry.UserObject) error {
	batch := &pgx.Batch{}
	for _, userSpace := range userSpaces {
		batch.Queue(
			upsertUserSpaceQuery, userSpace.UserID, userSpace.ObjectID,
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
					userSpaces[i].UserID, userSpaces[i].ObjectID,
				),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) RemoveUserObject(ctx context.Context, userSpaces *entry.UserObject) error {
	if _, err := db.conn.Exec(ctx, removeUserSpaceByIDQuery, userSpaces.UserID, userSpaces.ObjectID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}

	return nil
}

func (db *DB) RemoveUserObjects(ctx context.Context, userSpaces []*entry.UserObject) error {
	batch := &pgx.Batch{}
	for _, userSpace := range userSpaces {
		batch.Queue(removeUserSpaceByIDQuery, userSpace.UserID, userSpace.ObjectID)
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
					userSpaces[i].UserID, userSpaces[i].ObjectID,
				),
			)
		}
	}

	return errs.ErrorOrNil()
}
