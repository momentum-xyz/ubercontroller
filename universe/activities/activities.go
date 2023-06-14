package activities

import (
	"sort"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/activity"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.Activities = (*Activities)(nil)

type Activities struct {
	ctx        types.NodeContext
	log        *zap.SugaredLogger
	db         database.DB
	activities *generic.SyncMap[umid.UMID, universe.Activity]
}

func NewActivities(db database.DB) *Activities {
	return &Activities{
		db:         db,
		activities: generic.NewSyncMap[umid.UMID, universe.Activity](0),
	}
}

func (a *Activities) Initialize(ctx types.NodeContext) error {
	a.ctx = ctx
	a.log = ctx.Logger()

	return nil
}

func (a *Activities) CreateActivity(activityID umid.UMID) (universe.Activity, error) {
	activity := activity.NewActivity(activityID, a.db)

	if err := activity.Initialize(a.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize activity: %s", activityID)
	}
	if err := a.AddActivity(activity, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add activity: %s", activityID)
	}

	return activity, nil
}

func (a *Activities) GetActivity(activityID umid.UMID) (universe.Activity, bool) {
	asset, ok := a.activities.Load(activityID)
	return asset, ok
}

func (a *Activities) GetActivities() map[umid.UMID]universe.Activity {
	return a.activities.Map(func(k umid.UMID, v universe.Activity) universe.Activity {
		return v
	})
}

func (a *Activities) GetPaginatedActivitiesByObjectID(objectID *umid.UMID, page int, pageSize int) []universe.Activity {
	a.activities.Mu.RLock()
	defer a.activities.Mu.RUnlock()

	if page < 1 {
		page = 1
	}

	const maxPageSize = 100
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	var allActivities []universe.Activity
	for _, activityD := range a.activities.Data {
		if activityD.GetObjectID() == *objectID {
			allActivities = append(allActivities, activityD)
		}
	}

	sort.Slice(allActivities, func(i, j int) bool {
		return allActivities[i].GetCreatedAt().After(allActivities[j].GetCreatedAt())
	})

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(allActivities) {
		end = len(allActivities)
	}

	if start >= end {
		return []universe.Activity{}
	}

	return allActivities[start:end]
}

func (a *Activities) GetActivitiesByUserID(userID umid.UMID) map[umid.UMID]universe.Activity {
	a.activities.Mu.RLock()
	defer a.activities.Mu.RUnlock()

	activities := make(map[umid.UMID]universe.Activity, len(a.activities.Data))

	for id, asset := range a.activities.Data {
		if asset.GetUserID() == userID {
			activities[id] = asset
		}
	}

	return activities
}

func (a *Activities) AddActivity(activity universe.Activity, updateDB bool) error {
	a.activities.Mu.Lock()
	defer a.activities.Mu.Unlock()

	if updateDB {
		if err := a.db.GetActivitiesDB().UpsertActivity(a.ctx, activity.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.activities.Data[activity.GetID()] = activity

	return nil
}

func (a *Activities) AddActivities(activities []universe.Activity, updateDB bool) error {
	a.activities.Mu.Lock()
	defer a.activities.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.Activity, len(activities))
		for i := range activities {
			entries[i] = activities[i].GetEntry()
		}
		if err := a.db.GetActivitiesDB().UpsertActivities(a.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range activities {
		a.activities.Data[activities[i].GetID()] = activities[i]
	}

	return nil
}

func (a *Activities) RemoveActivity(activity universe.Activity, updateDB bool) (bool, error) {
	a.activities.Mu.Lock()
	defer a.activities.Mu.Unlock()

	if _, ok := a.activities.Data[activity.GetID()]; !ok {
		return false, nil
	}

	if updateDB {
		if err := a.db.GetActivitiesDB().RemoveActivityByID(a.ctx, activity.GetID()); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(a.activities.Data, activity.GetID())

	return true, nil
}

func (a *Activities) RemoveActivities(activities2d []universe.Activity, updateDB bool) (bool, error) {
	a.activities.Mu.Lock()
	defer a.activities.Mu.Unlock()

	for i := range activities2d {
		if _, ok := a.activities.Data[activities2d[i].GetID()]; !ok {
			return false, nil
		}
	}

	if updateDB {
		ids := make([]umid.UMID, len(activities2d))
		for i := range activities2d {
			ids[i] = activities2d[i].GetID()
		}
		if err := a.db.GetActivitiesDB().RemoveActivitiesByIDs(a.ctx, ids); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range activities2d {
		delete(a.activities.Data, activities2d[i].GetID())
	}

	return true, nil
}

func (a *Activities) Load() error {
	a.log.Info("Loading activities...")

	entries, err := a.db.GetActivitiesDB().GetActivities(a.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get activities")
	}

	for _, assetEntry := range entries {
		activity, err := a.CreateActivity(assetEntry.ActivityID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new activity: %s", assetEntry.ActivityID)
		}
		if err := activity.LoadFromEntry(assetEntry); err != nil {
			return errors.WithMessagef(err, "failed to load activity from entry: %s", assetEntry.ActivityID)
		}
	}

	a.log.Infof("Activities loaded: %d", a.activities.Len())

	return nil
}

func (a *Activities) Save() error {
	a.log.Info("Saving activities...")

	a.activities.Mu.RLock()
	defer a.activities.Mu.RUnlock()

	entries := make([]*entry.Activity, 0, len(a.activities.Data))
	for _, activity := range a.activities.Data {
		entries = append(entries, activity.GetEntry())
	}

	if err := a.db.GetActivitiesDB().UpsertActivities(a.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert activities")
	}

	a.log.Infof("Activities saved: %d", len(a.activities.Data))

	return nil
}
