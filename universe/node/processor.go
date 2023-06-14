package node

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
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

		msg := posbus.WSMessage(&posbus.ActivityUpdate{
			ActivityId: activity.GetID(),
			UserId:     activity.GetUserID(),
			ObjectId:   activity.GetObjectID(),
			ChangeType: string(updateType),
			Type:       activity.GetType(),
			Data:       activity.GetData(),
		})
		if err := object.Send(msg, true); err != nil {
			return errors.WithMessage(err, "failed to send ws message")
		}

		return nil
	case posbus.ChangedActivityUpdateType:
	case posbus.RemovedActivityUpdateType:
		objectIDs, err := n.db.GetObjectActivitiesDB().GetObjectIDsByActivityID(n.ctx, activity.GetID())
		if err != nil {
			return errors.WithMessage(err, "failed to get objectIds by activityId")
		}

		for _, objectID := range objectIDs {
			object, ok := n.GetObjectFromAllObjects(objectID)
			if !ok {
				return errors.Errorf("failed to get object from all objects: %s", objectID)
			}

			msg := posbus.WSMessage(&posbus.ActivityUpdate{
				ActivityId: activity.GetID(),
				UserId:     activity.GetUserID(),
				ObjectId:   activity.GetObjectID(),
				ChangeType: string(updateType),
				Type:       activity.GetType(),
				Data:       activity.GetData(),
			})
			if err := object.Send(msg, true); err != nil {
				return errors.WithMessage(err, "failed to send ws message")
			}
		}
		return nil
	default:
		return errors.New("no valid updateType supplied")
	}

	return nil
}
