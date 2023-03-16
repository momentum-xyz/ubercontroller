package user_objects

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"sync"

	"github.com/momentum-xyz/ubercontroller/utils/modify"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getUserObjectsQuery           = `SELECT * FROM user_object;`
	getUserObjectByIDQuery        = `SELECT * FROM user_object WHERE user_id = $1 AND object_id = $2;`
	getUserObjectsByUserIDQuery   = `SELECT * FROM user_object WHERE user_id = $1;`
	getUserObjectsByObjectIDQuery = `SELECT * FROM user_object WHERE object_id = $1;`
	getUserObjectValueByIDQuery   = `SELECT value FROM user_object WHERE user_id = $1 AND object_id = $2;`

	getObjectIndirectAdminsQuery  = `SELECT GetIndirectObjectAdmins($1);`
	checkIsIndirectAdminByIDQuery = `WITH object_admins AS (
										SELECT GetIndirectObjectAdmins($2) AS user_id
									)
									SELECT EXISTS(
										SELECT 1
										FROM object_admins
										WHERE user_id = $1
									) AS user_is_admin;`

	upsertUserObjectQuery = `INSERT INTO user_object
    							(user_id, object_id, value, created_at, updated_at)
							VALUES
							    ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
							ON CONFLICT (user_id, object_id)
							DO UPDATE SET
							    value = $3;`

	updateUserObjectValueByIDQuery = `UPDATE user_object SET value = $3, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1 AND object_id = $2;`

	removeUserObjectByIDQuery = `DELETE FROM user_object WHERE user_id = $1 AND object_id = $2;`
)

var _ database.UserObjectsDB = (*DB)(nil)

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

func (db *DB) GetUserObjectsByUserID(ctx context.Context, userID umid.UMID) ([]*entry.UserObject, error) {
	var userObjects []*entry.UserObject
	if err := pgxscan.Select(ctx, db.conn, &userObjects, getUserObjectsByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userObjects, nil
}

func (db *DB) GetUserObjectsByObjectID(ctx context.Context, objectID umid.UMID) ([]*entry.UserObject, error) {
	var userObjects []*entry.UserObject
	if err := pgxscan.Select(ctx, db.conn, &userObjects, getUserObjectsByObjectIDQuery, objectID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userObjects, nil
}

func (db *DB) GetUserObjectByID(ctx context.Context, userObjectID entry.UserObjectID) (*entry.UserObject, error) {
	var userObject *entry.UserObject
	if err := pgxscan.Select(
		ctx, db.conn, &userObject, getUserObjectByIDQuery, userObjectID.UserID, userObjectID.ObjectID,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userObject, nil
}

func (db *DB) GetUserObjectValueByID(
	ctx context.Context, userObjectID entry.UserObjectID,
) (*entry.UserObjectValue, error) {
	var value entry.UserObjectValue
	if err := db.conn.QueryRow(
		ctx, getUserObjectValueByIDQuery,
		userObjectID.UserID, userObjectID.ObjectID,
	).
		Scan(&value); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &value, nil
}

func (db *DB) GetUserObjects(ctx context.Context) ([]*entry.UserObject, error) {
	var userObjects []*entry.UserObject
	if err := pgxscan.Select(ctx, db.conn, &userObjects, getUserObjectsQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userObjects, nil
}

func (db *DB) GetObjectIndirectAdmins(ctx context.Context, objectID umid.UMID) ([]*umid.UMID, error) {
	var userIDs []*umid.UMID
	if err := pgxscan.Select(ctx, db.conn, &userIDs, getObjectIndirectAdminsQuery, objectID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userIDs, nil
}

func (db *DB) CheckIsIndirectAdminByID(ctx context.Context, userObjectID entry.UserObjectID) (bool, error) {
	var isIndirectAdmin bool
	if err := db.conn.QueryRow(ctx, checkIsIndirectAdminByIDQuery, userObjectID.UserID, userObjectID.ObjectID).
		Scan(&isIndirectAdmin); err != nil {
		return false, errors.WithMessage(err, "failed to query db")
	}
	return isIndirectAdmin, nil
}

func (db *DB) UpdateUserObjectValue(
	ctx context.Context, userObjectID entry.UserObjectID, modifyFn modify.Fn[entry.UserObjectValue],
) (*entry.UserObjectValue, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	value, err := db.GetUserObjectValueByID(ctx, userObjectID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get value by umid")
	}

	value, err = modifyFn(value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify value")
	}

	if _, err := db.conn.Exec(
		ctx, updateUserObjectValueByIDQuery,
		userObjectID.UserID, userObjectID.ObjectID,
		value,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return value, nil
}

func (db *DB) UpsertUserObject(
	ctx context.Context, userObjectID entry.UserObjectID, modifyFn modify.Fn[entry.UserObjectValue],
) (*entry.UserObjectValue, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	value, err := db.GetUserObjectValueByID(ctx, userObjectID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.WithMessage(err, "failed to get user object value by umid")
		}
	}

	value, err = modifyFn(value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify user object value")
	}

	if _, err := db.conn.Exec(
		ctx, upsertUserObjectQuery,
		userObjectID.UserID, userObjectID.ObjectID,
		value,
	); err != nil {
		return nil, errors.WithMessage(err, "failed to exec db")
	}

	return value, nil
}

func (db *DB) RemoveUserObjectByID(ctx context.Context, userObjectID entry.UserObjectID) error {
	res, err := db.conn.Exec(ctx, removeUserObjectByIDQuery, userObjectID.UserID, userObjectID.ObjectID)
	if err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (db *DB) RemoveUserObjectsByIDs(ctx context.Context, userObjectIDs []entry.UserObjectID) error {
	batch := &pgx.Batch{}
	for _, userObjectID := range userObjectIDs {
		batch.Queue(removeUserObjectByIDQuery, userObjectID.UserID, userObjectID.ObjectID)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %+v", userObjectIDs[i]),
			)
		}
	}

	return errs.ErrorOrNil()
}
