package node

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) InjectActivity(activity universe.Activity) error {
	err := n.activities.AddActivity(activity, true)
	if err != nil {
		return errors.WithMessage(err, "failed to inject activity")
	}

	// n.NotifyActivityProcessor()

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

	// n.NotifyActivityProcessor()

	return nil
}

func (n *Node) RemoveActivity(activity universe.Activity) error {
	ok, err := n.activities.RemoveActivity(activity, true)
	if err != nil {
		return errors.WithMessage(err, "failed to remove activity")
	}
	if ok {
		// n.NotifyActivityProcessor()
	}

	return nil
}
