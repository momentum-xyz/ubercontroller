package activities

import (
	"sync"

	"github.com/hashicorp/go-multierror"
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
	case posbus.ChangedActivityUpdateType, posbus.RemovedActivityUpdateType:
		if err := a.handleChangedRemovedActivity(activity, updateType); err != nil {
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

func (a *Activities) handleChangedRemovedActivity(activity universe.Activity, updateType posbus.ActivityUpdateType) error {
	objectIDs, err := a.db.GetObjectActivitiesDB().GetObjectIDsByActivityID(a.ctx, activity.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get objectIds by activityId")
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(objectIDs))

	for _, objectID := range objectIDs {
		wg.Add(1)
		go func(objectID umid.UMID) {
			defer wg.Done()
			if err := a.sendMessageToPosBus(activity, objectID, updateType); err != nil {
				errCh <- err
			}
		}(objectID)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var result *multierror.Error
	for err := range errCh {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
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
