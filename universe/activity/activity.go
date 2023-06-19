package activity

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.Activity = (*Activity)(nil)

type Activity struct {
	ctx        types.NodeContext
	log        *zap.SugaredLogger
	db         database.DB
	mu         sync.RWMutex
	entry      *entry.Activity
	activities universe.Activities
}

func NewActivity(id umid.UMID, db database.DB, activities universe.Activities) *Activity {
	return &Activity{
		db: db,
		entry: &entry.Activity{
			ActivityID: id,
		},
		activities: activities,
	}
}

func (a *Activity) GetID() umid.UMID {
	return a.entry.ActivityID
}

func (a *Activity) Initialize(ctx types.NodeContext) error {
	a.ctx = ctx
	a.log = ctx.Logger()

	return nil
}

func (a *Activity) GetActivities() universe.Activities {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.activities
}

func (a *Activity) GetData() *entry.ActivityData {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Data
}

func (a *Activity) SetData(modifyFn modify.Fn[entry.ActivityData], updateDB bool) (*entry.ActivityData, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	data, err := modifyFn(a.entry.Data)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify data")
	}

	if updateDB {
		if err := a.db.GetActivitiesDB().UpdateActivityData(a.ctx, a.entry.ActivityID, data); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Data = data

	return data, nil
}

func (a *Activity) GetType() *entry.ActivityType {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Type
}

func (a *Activity) SetType(activityType *entry.ActivityType, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.GetActivitiesDB().UpdateActivityType(a.ctx, a.entry.ActivityID, activityType); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Type = activityType

	return nil
}

func (a *Activity) GetObjectID() umid.UMID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.ObjectID
}

func (a *Activity) SetObjectID(objectID umid.UMID, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.GetActivitiesDB().UpdateActivityObjectID(a.ctx, a.GetID(), &objectID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.ObjectID = objectID

	return nil
}

func (a *Activity) GetUserID() umid.UMID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.UserID
}

func (a *Activity) SetUserID(userID umid.UMID, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.GetActivitiesDB().UpdateActivityUserID(a.ctx, a.GetID(), &userID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.UserID = userID

	return nil
}

func (a *Activity) GetEntry() *entry.Activity {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry
}

func (a *Activity) LoadFromEntry(entry *entry.Activity) error {
	if entry.ActivityID != a.GetID() {
		return errors.Errorf("activity ids mismatch: %s != %s", entry.ActivityID, a.GetID())
	}

	if err := a.SetObjectID(entry.ObjectID, false); err != nil {
		return errors.WithMessage(err, "failed to set object ID")
	}
	if err := a.SetUserID(entry.UserID, false); err != nil {
		return errors.WithMessage(err, "failed to set user ID")
	}
	if err := a.SetType(entry.Type, false); err != nil {
		return errors.WithMessage(err, "failed to set type")
	}
	if err := a.SetCreatedAt(entry.CreatedAt, false); err != nil {
		return errors.WithMessage(err, "failed to set created at")
	}
	if _, err := a.SetData(modify.MergeWith(entry.Data), false); err != nil {
		return errors.WithMessage(err, "failed to set data")
	}

	return nil
}

func (a *Activity) GetCreatedAt() time.Time {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.CreatedAt
}

func (a *Activity) SetCreatedAt(createdAt time.Time, updateDB bool) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if updateDB {
		if err := a.db.GetActivitiesDB().UpdateActivityCreatedAt(a.ctx, a.GetID(), createdAt); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.CreatedAt = createdAt

	return nil
}
