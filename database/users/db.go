package users

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
)

const (
	getUserByIDQuery = `SELECT * FROM "user" WHERE "user_id" = $1;`
)

var _ database.UsersDB = (*DB)(nil)

type DB struct {
	conn   *pgxpool.Pool
	common database.CommonDB
}

func (db *DB) UsersUpsertUser(ctx context.Context, user *entry.User) error {
	//TODO implement me
	panic("implement me")
}

func (db *DB) UsersGetUserByID(ctx context.Context, userID uuid.UUID) (*entry.User, error) {
	var user entry.User
	if err := pgxscan.Get(ctx, db.conn, &user, getUserByIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &user, nil
}

func NewDB(conn *pgxpool.Pool, commonDB database.CommonDB) *DB {
	return &DB{
		conn:   conn,
		common: commonDB,
	}
}
