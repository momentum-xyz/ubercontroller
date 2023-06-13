package node

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func (n *Node) NotifyActivityProcessor(activity universe.Activity, updateType posbus.ActivityUpdateType) error {
	switch updateType {
	case posbus.NewActivityUpdateType:
		objectActivityID := entry.NewObjectActivityID(activity.GetObjectID(), activity.GetID())
		objectActivity := entry.NewObjectActivity(objectActivityID)
		if err := n.db.GetObjectActivitiesDB().UpsertObjectActivity(n.ctx, objectActivity); err != nil {
			return errors.WithMessage(err, "failed to upsert object activity")
		}

		userActivityID := entry.NewUserActivityID(activity.GetUserID(), activity.GetID())
		userActivity := entry.NewUserActivity(userActivityID)
		if err := n.db.GetUserActivitiesDB().UpsertUserActivity(n.ctx, userActivity); err != nil {
			return errors.WithMessage(err, "failed to upsert user activity")
		}

		object, ok := n.GetObjectFromAllObjects(activity.GetObjectID())
		if !ok {
			return errors.New("failed to get object from all objects")
		}
		var activityData posbus.StringAnyMap
		if err := utils.MapDecode(activity.GetData(), activityData); err != nil {
			return errors.WithMessage(err, "failed to marshal activity data")
		}

		var strActivityType string
		if activity.GetType() != nil {
			strActivityType = string(*activity.GetType())
		} else {
			strActivityType = ""
		}

		msg := posbus.WSMessage(&posbus.ActivityUpdate{
			ActivityId: activity.GetID(),
			UserId:     activity.GetUserID(),
			ObjectId:   activity.GetObjectID(),
			ChangeType: string(updateType),
			Type:       strActivityType,
			Data:       &activityData,
		})
		if err := object.Send(msg, true); err != nil {
			return errors.WithMessage(err, "failed to send ws message")
		}

		return nil
	case posbus.ChangedActivityUpdateType:
	case posbus.RemovedActivityUpdateType:
	default:
		return errors.New("no valid updateType supplied")
	}

	return nil
}
