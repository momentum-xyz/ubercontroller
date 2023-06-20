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
	if err := a.NotifyProcessor(activity, posbus.NewActivityUpdateType); err != nil {
		return errors.WithMessage(err, "failed to notify activity processor")
	}

	return nil
}

func (a *Activities) Modify(activity universe.Activity, modifyFn modify.Fn[entry.ActivityData]) error {
	_, err := activity.SetData(modifyFn, true)
	if err != nil {
		return errors.WithMessage(err, "failed to set activity data")
	}
	if err := a.Save(); err != nil {
		return errors.WithMessage(err, "failed to save activity")
	}
	if err := a.NotifyProcessor(activity, posbus.ChangedActivityUpdateType); err != nil {
		return errors.WithMessage(err, "failed to notify activity processor")
	}

	return nil
}

func (a *Activities) Remove(activity universe.Activity) error {
	if err := a.NotifyProcessor(activity, posbus.RemovedActivityUpdateType); err != nil {
		return errors.WithMessage(err, "failed to notify activity processor")
	}
	_, err := a.RemoveActivity(activity, true)
	if err != nil {
		return errors.WithMessage(err, "failed to remove activity")
	}

	return nil
}
