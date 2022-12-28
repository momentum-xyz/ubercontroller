package users

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
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
         						                                WHERE plugin_id = '86DC3AE7-9F3D-42CB-85A3-A71ABC3C3CB8'
         						                                  AND attribute_name = 'wallet'
         						                                  AND value->'wallet' ? $1);`
	getUserProfileByUserIDQuery = `SELECT profile FROM "user" WHERE user_id = $1;`

	removeUserByIDQuery   = `DELETE FROM "user" WHERE user_id = $1;`
	removeUsersByIDsQuery = `DELETE FROM "user" WHERE user_id = ANY($1);`

	updateUserUserTypeIDQuery = `UPDATE "user" SET user_type_id = $2 WHERE user_id = $1;`
	updateUserOptionsQuery    = `UPDATE "user" SET options = $2 WHERE user_id = $1;`
	updateUserProfileQuery    = `UPDATE "user" SET profile = $2 WHERE user_id = $1;`

	upsertUserQuery = `INSERT INTO "user"
    						(user_id, user_type_id, profile, options, created_at, updated_at)
						VALUES
						    ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
						ON CONFLICT (user_id)
						DO UPDATE SET
							user_type_id = $2, profile = $3, options = $4, updated_at = CURRENT_TIMESTAMP;`
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

func (db *DB) UsersGetUserByID(ctx context.Context, userID uuid.UUID) (*entry.User, error) {
	var user entry.User
	if err := pgxscan.Get(ctx, db.conn, &user, getUserByIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &user, nil
}

func (db *DB) UsersGetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entry.User, error) {
	var users []*entry.User
	if err := pgxscan.Select(ctx, db.conn, &users, getUsersByIDsQuery, userIDs); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return users, nil
}

func (db *DB) UsersGetUserByWallet(ctx context.Context, wallet string) (*entry.User, error) {
	var user entry.User
	if err := pgxscan.Get(ctx, db.conn, &user, getUserByWalletQuery, wallet); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &user, nil
}

func (db *DB) UsersGetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*entry.UserProfile, error) {
	var profile *entry.UserProfile
	err := db.conn.QueryRow(ctx,
		getUserProfileByUserIDQuery, userID).Scan(&profile)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return profile, nil
}

func (db *DB) UsersUpsertUser(ctx context.Context, user *entry.User) error {
	if _, err := db.conn.Exec(
		ctx, upsertUserQuery,
		user.UserID, user.UserTypeID, user.Profile, user.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil

}

func (db *DB) UsersUpsertUsers(ctx context.Context, users []*entry.User) error {
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

func (db *DB) UsersRemoveUsersByIDs(ctx context.Context, userIDs []uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeUsersByIDsQuery, userIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UsersRemoveUserByID(ctx context.Context, userID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, removeUserByIDQuery, userID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil

}
func (db *DB) UsersUpdateUserUserTypeID(ctx context.Context, userID uuid.UUID, userTypeID uuid.UUID) error {
	if _, err := db.conn.Exec(ctx, updateUserUserTypeIDQuery, userID, userTypeID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UsersUpdateUserOptions(ctx context.Context, userID uuid.UUID, options *entry.UserOptions) error {
	if _, err := db.conn.Exec(ctx, updateUserOptionsQuery, userID, options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UsersUpdateUserProfile(ctx context.Context, userID uuid.UUID, profile *entry.UserProfile) error {
	if _, err := db.conn.Exec(ctx, updateUserProfileQuery, userID, profile); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
