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
	upsertObjectActivityQuery = `INSERT INTO object_activity
    						(object_id, activity_id,
    						created_at)
						VALUES
							($1, $2, CURRENT_TIMESTAMP)
						ON CONFLICT (object_id)
						DO UPDATE SET
							activity_id = $2;`
)

var _ database.ObjectActivitiesDB = (*DB)(nil)

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

func (db *DB) UpsertObjectActivity(ctx context.Context, objectActivity *entry.ObjectActivity) error {
	if _, err := db.conn.Exec(
		ctx, upsertObjectActivityQuery,
		objectActivity.ObjectID, objectActivity.ActivityID,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertObjectActivities(ctx context.Context, objectActivities []*entry.ObjectActivity) error {
	batch := &pgx.Batch{}
	for _, objectActivity := range objectActivities {
		batch.Queue(
			upsertObjectActivityQuery, objectActivity.ObjectID, objectActivity.ActivityID,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %s", objectActivities[i].ObjectID),
			)
		}
	}

	return errs.ErrorOrNil()
}
