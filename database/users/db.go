package users

import (
	"context"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"

	"github.com/momentum-xyz/ubercontroller/database"
)

const (
	getUserByIDQuery     = `SELECT * FROM "user" WHERE user_id = $1;`
	getUsersByIDsQuery   = `SELECT * FROM "user" WHERE user_id = ANY($1);`
	getUserByWalletQuery = `SELECT * FROM "user"
         						WHERE user_id = (SELECT user_id FROM user_attribute
         						                    /* Kusama plugin umid */
         						                	WHERE plugin_id = '86DC3AE7-9F3D-42CB-85A3-A71ABC3C3CB8'
         						                    AND attribute_name = 'wallet'
         						                    AND value->'wallet' ? $1
         						                );`
	getWalletByUserID = `SELECT value -> 'wallet' ->> 0 AS wallet
						FROM user_attribute
						WHERE user_id = $1
						  AND plugin_id = '86DC3AE7-9F3D-42CB-85A3-A71ABC3C3CB8'
						  AND attribute_name = 'wallet'`
	checkIsUserExistsByNameQuery   = `SELECT EXISTS(SELECT 1 FROM "user" WHERE profile->>'name' = $1);`
	checkIsUserExistsByWalletQuery = `SELECT EXISTS(SELECT 1 FROM "user" WHERE user_id = (SELECT user_id FROM user_attribute
         						                    /* Kusama plugin umid */
         						                	WHERE plugin_id = '86DC3AE7-9F3D-42CB-85A3-A71ABC3C3CB8'
         						                    AND attribute_name = 'wallet'
         						                    AND value->'wallet' ? $1
         						                );`
	getUserProfileByUserIDQuery = `SELECT profile FROM "user" WHERE user_id = $1;`

	upsertUserQuery = `INSERT INTO "user"
    						(user_id, user_type_id, profile, options, created_at, updated_at)
						VALUES
						    ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
						ON CONFLICT (user_id)
						DO UPDATE SET
							user_type_id = $2, profile = $3, options = $4, updated_at = CURRENT_TIMESTAMP;`

	updateUserUserTypeIDQuery = `UPDATE "user" SET user_type_id = $2, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1;`
	updateUserOptionsQuery    = `UPDATE "user" SET options = $2, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1;`
	updateUserProfileQuery    = `UPDATE "user" SET profile = $2, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1;`

	removeUserByIDQuery   = `DELETE FROM "user" WHERE user_id = $1;`
	removeUsersByIDsQuery = `DELETE FROM "user" WHERE user_id = ANY($1);`
)

var _ database.UsersDB = (*DB)(nil)

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

func (db *DB) GetUserByID(ctx context.Context, userID umid.UMID) (*entry.User, error) {
	var user entry.User
	if err := pgxscan.Get(ctx, db.conn, &user, getUserByIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &user, nil
}

func (db *DB) GetUsersByIDs(ctx context.Context, userIDs []umid.UMID) ([]*entry.User, error) {
	var users []*entry.User
	if err := pgxscan.Select(ctx, db.conn, &users, getUsersByIDsQuery, userIDs); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return users, nil
}

func (db *DB) GetUserByWallet(ctx context.Context, wallet string) (*entry.User, error) {
	var user entry.User
	if err := pgxscan.Get(ctx, db.conn, &user, getUserByWalletQuery, wallet); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &user, nil
}

func (db *DB) GetUserWalletByUserID(ctx context.Context, userID umid.UMID) (*string, error) {
	var wallet string
	if err := db.conn.QueryRow(ctx, getWalletByUserID, userID).
		Scan(&wallet); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}

	return &wallet, nil
}

func (db *DB) CheckIsUserExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	if err := pgxscan.Get(ctx, db.conn, &exists, checkIsUserExistsByNameQuery, name); err != nil {
		return false, errors.WithMessage(err, "failed to query db")
	}
	return exists, nil
}

func (db *DB) CheckIsUserExistsByWallet(ctx context.Context, wallet string) (bool, error) {
	var exists bool
	if err := pgxscan.Get(ctx, db.conn, &exists, checkIsUserExistsByWalletQuery, wallet); err != nil {
		return false, errors.WithMessage(err, "failed to query db")
	}
	return exists, nil
}

func (db *DB) GetUserProfileByUserID(ctx context.Context, userID umid.UMID) (*entry.UserProfile, error) {
	var profile entry.UserProfile
	err := db.conn.QueryRow(
		ctx,
		getUserProfileByUserIDQuery, userID,
	).Scan(&profile)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &profile, nil
}

func (db *DB) UpsertUser(ctx context.Context, user *entry.User) error {
	if _, err := db.conn.Exec(
		ctx, upsertUserQuery,
		user.UserID, user.UserTypeID, user.Profile, user.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil

}

func (db *DB) UpsertUsers(ctx context.Context, users []*entry.User) error {
	batch := &pgx.Batch{}
	for _, user := range users {
		batch.Queue(
			upsertUserQuery, user.UserID, user.UserTypeID, user.Profile, user.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to exec db for: %s", users[i].UserID))
		}
	}

	return errs.ErrorOrNil()

}

func (db *DB) RemoveUsersByIDs(ctx context.Context, userIDs []umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removeUsersByIDsQuery, userIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveUserByID(ctx context.Context, userID umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removeUserByIDQuery, userID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil

}
func (db *DB) UpdateUserUserTypeID(ctx context.Context, userID umid.UMID, userTypeID umid.UMID) error {
	if _, err := db.conn.Exec(ctx, updateUserUserTypeIDQuery, userID, userTypeID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateUserOptions(ctx context.Context, userID umid.UMID, options *entry.UserOptions) error {
	if _, err := db.conn.Exec(ctx, updateUserOptionsQuery, userID, options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateUserProfile(ctx context.Context, userID umid.UMID, profile *entry.UserProfile) error {
	if _, err := db.conn.Exec(ctx, updateUserProfileQuery, userID, profile); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
