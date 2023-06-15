package activities

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

func (a *Activities) NotifyProcessor(activity universe.Activity, updateType posbus.ActivityUpdateType) error {
	switch updateType {
	case posbus.NewActivityUpdateType:
		if err := a.handleNewActivity(activity); err != nil {
			return err
		}
	case posbus.ChangedActivityUpdateType:
	case posbus.RemovedActivityUpdateType:
		if err := a.handleChangedRemovedActivity(activity); err != nil {
			return err
		}
	default:
		return errors.New("no valid updateType supplied")
	}

	return nil
}

func (a *Activities) handleNewActivity(activity universe.Activity) error {
	objectActivityID := entry.NewObjectActivityID(activity.GetObjectID(), activity.GetID())
	objectActivity := entry.NewObjectActivity(objectActivityID)
	if err := a.db.GetObjectActivitiesDB().UpsertObjectActivity(a.ctx, objectActivity); err != nil {
		return errors.WithMessage(err, "failed to upsert object activity")
	}

	userActivityID := entry.NewUserActivityID(activity.GetUserID(), activity.GetID())
	userActivity := entry.NewUserActivity(userActivityID)
	if err := a.db.GetUserActivitiesDB().UpsertUserActivity(a.ctx, userActivity); err != nil {
		return errors.WithMessage(err, "failed to upsert user activity")
	}

	return a.sendMessageToPosBus(activity, activity.GetObjectID(), posbus.NewActivityUpdateType)
}

func (a *Activities) handleChangedRemovedActivity(activity universe.Activity) error {
	objectIDs, err := a.db.GetObjectActivitiesDB().GetObjectIDsByActivityID(a.ctx, activity.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get objectIds by activityId")
	}

	errCh := make(chan error)
	for _, objectID := range objectIDs {
		go func(objectID umid.UMID) {
			if err := a.sendMessageToPosBus(activity, objectID, posbus.RemovedActivityUpdateType); err != nil {
				errCh <- err
			}
		}(objectID)
	}

	for range objectIDs {
		if err := <-errCh; err != nil {
			return err
		}
	}

	return nil
}

func (a *Activities) sendMessageToPosBus(activity universe.Activity, objectID umid.UMID, updateType posbus.ActivityUpdateType) error {
	n := universe.GetNode()
	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		return errors.Errorf("failed to get object from all objects: %v", objectID)
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
}
