package activities

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (a *Activities) Inject(activity universe.Activity) error {
	if err := a.AddActivity(activity, true); err != nil {
		return errors.WithMessage(err, "failed to inject activity")
	}
	objectIDs, err := a.db.GetObjectActivitiesDB().GetObjectIDsByActivityID(a.ctx, activity.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get objectIds by activityId")
	}
	if err := a.NotifyProcessor(activity, posbus.NewActivityUpdateType, objectIDs); err != nil {
		return errors.WithMessage(err, "failed to notify activity processor")
	}

	return nil
}

func (a *Activities) Modify(activity universe.Activity, modifyFn modify.Fn[entry.ActivityData]) error {
	_, err := activity.SetData(modifyFn, true)
	if err != nil {
		return errors.WithMessage(err, "failed to set activity data")
	}
	objectIDs, err := a.db.GetObjectActivitiesDB().GetObjectIDsByActivityID(a.ctx, activity.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get objectIds by activityId")
	}
	if err := a.Save(); err != nil {
		return errors.WithMessage(err, "failed to save activity")
	}
	if err := a.NotifyProcessor(activity, posbus.ChangedActivityUpdateType, objectIDs); err != nil {
		return errors.WithMessage(err, "failed to notify activity processor")
	}

	return nil
}

func (a *Activities) Remove(activity universe.Activity) error {
	objectIDs, err := a.db.GetObjectActivitiesDB().GetObjectIDsByActivityID(a.ctx, activity.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get objectIds by activityId")
	}
	ok, err := a.RemoveActivity(activity, true)
	if err != nil {
		return errors.WithMessage(err, "failed to remove activity")
	}
	if ok {
		if err := a.NotifyProcessor(activity, posbus.RemovedActivityUpdateType, objectIDs); err != nil {
			return errors.WithMessage(err, "failed to notify activity processor")
		}
	}

	return nil
}
