package user_space

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	getUserSpacesQueryAll              = `SELECT * FROM user_space;`
	getUserSpacesByUserIDQuery         = `SELECT * FROM user_space WHERE user_id = $1;`
	getUserSpacesBySpaceIDQuery        = `SELECT * FROM user_space WHERE space_id = $1;`
	getUserSpaceByUserAndSpaceIDsQuery = `SELECT * FROM user_space WHERE user_id = $1 AND space_id = $2;`

	getUserSpaceValueByUserAndSpaceIDsQuery = `SELECT value FROM user_space WHERE user_id = $1 AND space_id = $2;`

	updateUserSpaceValueByUserAndSpaceIDsQuery = `UPDATE user_space SET value = $3 WHERE user_id = $1 AND space_id = $2;`
)

var _ database.UserSpaceDB = (*DB)(nil)

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

func (db *DB) UserSpaceGetUserSpaces(ctx context.Context) ([]*entry.UserSpace, error) {
	var userSpaces []*entry.UserSpace
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpacesQueryAll); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) UserSpaceGetUserSpacesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserSpace, error) {
	var userSpaces []*entry.UserSpace
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpacesByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) UserSpaceGetUserSpacesBySpaceID(ctx context.Context, spaceID uuid.UUID) ([]*entry.UserSpace, error) {
	var userSpaces []*entry.UserSpace
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpacesBySpaceIDQuery, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) UserSpaceGetUserSpaceByUserAndSpaceIDs(ctx context.Context, userID, spaceID uuid.UUID) (*entry.UserSpace, error) {
	var userSpaces *entry.UserSpace
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpaceByUserAndSpaceIDsQuery, userID, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}

func (db *DB) UserSpaceUpdateValueByUserAndSpaceIDs(ctx context.Context, userID, spaceID uuid.UUID, value *entry.UserSpaceValue) error {
	return nil
}

func (db *DB) UserSpaceGetValueByUserAndSpaceIDs(ctx context.Context, userID, spaceID uuid.UUID) (*entry.UserSpaceValue, error) {
	var userSpaces *entry.UserSpaceValue
	if err := pgxscan.Select(ctx, db.conn, &userSpaces, getUserSpaceValueByUserAndSpaceIDsQuery, userID, spaceID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return userSpaces, nil
}
