package activities

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
	getActivitiesQuery            = `SELECT * FROM activity;`
	getActivityByIDQuery          = `SELECT * FROM activity WHERE activity_id = $1;`
	getActivityIDsByParentIDQuery = `SELECT activity_id FROM activity WHERE parent_id = $1;`
	getActivitiesByUserIDQuery    = `SELECT * FROM activity WHERE user_id = $1;`
	getActivitiesByObjectIDQuery  = `SELECT * FROM activity WHERE object_id = $1;`

	upsertActivityQuery = `INSERT INTO activity
    						(activity_id, user_id, object_id, "type", "data",
    						created_at)
						VALUES
							($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
						ON CONFLICT (activity_id)
						DO UPDATE SET
							user_id = $2, object_id = $3, "type" = $4, "data" = $5;`

	updateActivityDataQuery     = `UPDATE activity SET data = $2 WHERE activity_id = $1;`
	updateActivityTypeQuery     = `UPDATE activity SET "type" = $2 WHERE activity_id = $1;`
	updateActivityUserIDQuery   = `UPDATE activity SET user_id = $2 WHERE activity_id = $1;`
	updateActivityObjectIDQuery = `UPDATE activity SET object_id = $2 WHERE activity_id = $1;`

	removeActivityByIDQuery    = `DELETE FROM activity WHERE activity_id = $1;`
	removeActivitiesByIDsQuery = `DELETE FROM activity WHERE activity_id = ANY($1);`
)

var _ database.ActivitiesDB = (*DB)(nil)

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

func (db *DB) GetActivities(ctx context.Context) ([]*entry.Activity, error) {
	var activities []*entry.Activity
	if err := pgxscan.Select(ctx, db.conn, &activities, getActivitiesQuery); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return activities, nil
}

func (db *DB) GetActivityByID(ctx context.Context, activityID umid.UMID) (*entry.Activity, error) {
	var activity entry.Activity
	if err := pgxscan.Get(ctx, db.conn, &activity, getActivityByIDQuery, activityID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return &activity, nil
}

func (db *DB) GetActivityIDsByParentID(ctx context.Context, parentID umid.UMID) ([]umid.UMID, error) {
	var ids []umid.UMID
	if err := pgxscan.Select(ctx, db.conn, &ids, getActivityIDsByParentIDQuery, parentID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return ids, nil
}

func (db *DB) GetActivitiesByUserID(ctx context.Context, userID umid.UMID) ([]*entry.Activity, error) {
	var activities []*entry.Activity
	if err := pgxscan.Select(ctx, db.conn, &activities, getActivitiesByUserIDQuery, userID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return activities, nil
}

func (db *DB) GetActivitiesByObjectID(ctx context.Context, objectID umid.UMID) ([]*entry.Activity, error) {
	var activities []*entry.Activity
	if err := pgxscan.Select(ctx, db.conn, &activities, getActivitiesByObjectIDQuery, objectID); err != nil {
		return nil, errors.WithMessage(err, "failed to query db")
	}
	return activities, nil
}

func (db *DB) RemoveActivityByID(ctx context.Context, activityID umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removeActivityByIDQuery, activityID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) RemoveActivitiesByIDs(ctx context.Context, activityIDs []umid.UMID) error {
	if _, err := db.conn.Exec(ctx, removeActivitiesByIDsQuery, activityIDs); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateActivityData(ctx context.Context, activityID umid.UMID, options *entry.ActivityData) error {
	if _, err := db.conn.Exec(ctx, updateActivityDataQuery, activityID, options); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateActivityType(ctx context.Context, activityID umid.UMID, activityType *entry.ActivityType) error {
	if _, err := db.conn.Exec(ctx, updateActivityTypeQuery, activityID, activityType); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateActivityUserID(ctx context.Context, activityID umid.UMID, userID *umid.UMID) error {
	if _, err := db.conn.Exec(ctx, updateActivityUserIDQuery, activityID, userID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpdateActivityObjectID(ctx context.Context, activityID umid.UMID, objectID *umid.UMID) error {
	if _, err := db.conn.Exec(ctx, updateActivityObjectIDQuery, activityID, objectID); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertActivity(ctx context.Context, activity *entry.Activity) error {
	if _, err := db.conn.Exec(
		ctx, upsertActivityQuery,
		activity.ActivityID, activity.UserID, activity.ObjectID, activity.Type, activity.Data,
	); err != nil {
		return errors.WithMessage(err, "failed to exec db")
	}
	return nil
}

func (db *DB) UpsertActivities(ctx context.Context, activities []*entry.Activity) error {
	batch := &pgx.Batch{}
	for _, activity := range activities {
		batch.Queue(
			upsertActivityQuery, activity.ActivityID, activity.UserID, activity.ObjectID, activity.Type, activity.Data,
		)
	}

	batchRes := db.conn.SendBatch(ctx, batch)
	defer batchRes.Close()

	var errs *multierror.Error
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchRes.Exec(); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to exec db for: %s", activities[i].ActivityID),
			)
		}
	}

	return errs.ErrorOrNil()
}
