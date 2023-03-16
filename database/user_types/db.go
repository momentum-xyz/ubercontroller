package user_types

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getUserTypesQuery = `SELECT * FROM user_type;`

	upsertUserTypeQuery = `INSERT INTO user_type
											(user_type_id, user_type_name, description, options)
										VALUES
											($1, $2, $3, $4)
										ON CONFLICT (user_type_id)
										DO UPDATE SET
											user_type_name = $2, description = $3, options = $4;`

	updateUserTypeNameQuery        = `UPDATE user_type SET user_type_name = $2 WHERE user_type_id = $1;`
	updateUserTypeDescriptionQuery = `UPDATE user_type SET description = $2 WHERE user_type_id = $1;`
	updateUserTypeOptionsQuery     = `UPDATE user_type SET options = $2 WHERE user_type_id = $1;`

	removeUserTypeByIDQuery   = `DELETE FROM user_type WHERE user_type_id = $1;`
	removeUserTypesByIDsQuery = `DELETE FROM user_type WHERE user_type_id = ANY($1);`
)

var _ database.UserTypesDB = (*DB)(nil)

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

func (db *DB) GetUserTypes(ctx context.Context) ([]*entry.UserType, error) {
	var userTypes []*entry.UserType
	if err := pgxscan.Select(ctx, db.conn, &userTypes, getUserTypesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userTypes, nil
}

func (db *DB) UpsertUserType(ctx context.Context, userType *entry.UserType) error {
	if _, err := db.conn.Exec(
		ctx, upsertUserTypeQuery, userType.UserTypeID, userType.UserTypeName,
		userType.Description, userType.Options,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertUserTypes(ctx context.Context, userTypes []*entry.UserType) error {
	batch := &pgx.Batch{}
	for _, userType := range userTypes {
		batch.Queue(
			upsertUserTypeQuery, userType.UserTypeID, userType.UserTypeName,
			userType.Description, userType.Options,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %s", userTypes[i].UserTypeID),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (db *DB) RemoveUserTypeByID(ctx context.Context, userTypeID umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removeUserTypeByIDQuery, userTypeID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveUserTypesByIDs(ctx context.Context, userTypeIDs []umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removeUserTypesByIDsQuery, userTypeIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateUserTypeName(ctx context.Context, userTypeID umid.UMID, name string) error {
	if _, err := db.conn.Exec(ctx, updateUserTypeNameQuery, userTypeID, name); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateUserTypeDescription(
	ctx context.Context, userTypeID umid.UMID, description string,
) error {
	if _, err := db.conn.Exec(ctx, updateUserTypeDescriptionQuery, userTypeID, description); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateUserTypeOptions(
	ctx context.Context, userTypeID umid.UMID, options *entry.UserOptions,
) error {
	if _, err := db.conn.Exec(ctx, updateUserTypeOptionsQuery, userTypeID, options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}
