package activities

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

const (
	upsertUserActivityQuery = `INSERT INTO user_activity
    						(user_id, activity_id,
    						created_at)
						VALUES
							($1, $2, CURRENT_TIMESTAMP)
						ON CONFLICT (user_id)
						DO UPDATE SET
							user_id = $2;`
)

var _ database.UserActivitiesDB = (*DB)(nil)

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

func (db *DB) UpsertUserActivity(ctx context.Context, userActivity *entry.UserActivity) error {
	if _, err := db.conn.Exec(
		ctx, upsertUserActivityQuery,
		userActivity.UserID, userActivity.ActivityID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertUserActivities(ctx context.Context, userActivities []*entry.UserActivity) error {
	batch := &pgx.Batch{}
	for _, userActivity := range userActivities {
		batch.Queue(
			upsertUserActivityQuery, userActivity.UserID, userActivity.ActivityID,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %s", userActivities[i].UserID),
			)
		}
	}

	return errs.ErrorOrNil()
}
