package node

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) InjectActivity(activity universe.Activity) error {
	if err := n.activities.AddActivity(activity, true); err != nil {
		return errors.WithMessage(err, "failed to inject activity")
	}
	if err := n.NotifyActivityProcessor(activity, posbus.NewActivityUpdateType); err != nil {
		return errors.WithMessage(err, "failed to notify activity processor")
	}

	return nil
}

func (n *Node) ModifyActivity(activity universe.Activity, modifyFn modify.Fn[entry.ActivityData]) error {
	_, err := activity.SetData(modifyFn, true)
	if err != nil {
		return errors.WithMessage(err, "failed to set activity data")
	}
	if err := n.activities.Save(); err != nil {
		return errors.WithMessage(err, "failed to save activity")
	}
	if err := n.NotifyActivityProcessor(activity, posbus.ChangedActivityUpdateType); err != nil {
		return errors.WithMessage(err, "failed to notify activity processor")
	}

	return nil
}

func (n *Node) RemoveActivity(activity universe.Activity) error {
	ok, err := n.activities.RemoveActivity(activity, true)
	if err != nil {
		return errors.WithMessage(err, "failed to remove activity")
	}
	if ok {
		if err := n.NotifyActivityProcessor(activity, posbus.RemovedActivityUpdateType); err != nil {
			return errors.WithMessage(err, "failed to notify activity processor")
		}
	}

	return nil
}
