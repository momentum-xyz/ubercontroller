package object

import (
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-multierror"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"github.com/pkg/errors"
	"github.com/zakaria-chahboun/cute"
)

func (o *Object) GetUser(userID mid.ID, recursive bool) (universe.User, bool) {
	user, ok := o.Users.Load(userID)
	if ok {
		return user, true
	}

	if !recursive {
		return nil, false
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		if user, ok := child.GetUser(userID, true); ok {
			return user, true
		}
	}

	return nil, false
}

// GetUsers return map with all nested users if recursive is true,
// otherwise the method return map with users dependent only to current object.
func (o *Object) GetUsers(recursive bool) map[mid.ID]universe.User {
	o.Users.Mu.RLock()
	users := make(map[mid.ID]universe.User, len(o.Users.Data))
	for id, user := range o.Users.Data {
		users[id] = user
	}
	o.Users.Mu.RUnlock()

	if !recursive {
		return users
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		for id, user := range child.GetUsers(true) {
			users[id] = user
		}
	}

	return users
}

// TODO: think about rollback on error
func (o *Object) AddUser(user universe.User, updateDB bool) error {
	o.Users.Mu.Lock()
	defer o.Users.Mu.Unlock()

	if user.GetWorld().GetID() != o.GetWorld().GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", user.GetWorld().GetID(), o.GetWorld().GetID())
	}

	if updateDB {
		o.log.Error("Object: AddUser: update database is not supported")
	}

	user.SetObject(o)
	o.Users.Data[user.GetID()] = user
	o.sendObjectEnterLeaveStats(user, 1)

	return nil
}

func (o *Object) RemoveUser(user universe.User, updateDB bool) (bool, error) {
	o.Users.Mu.Lock()
	defer o.Users.Mu.Unlock()

	if user.GetWorld().GetID() != o.GetWorld().GetID() {
		return false, nil
	}

	if updateDB {
		o.log.Error("Object: RemoveUser: update database is not supported")
	}

	user.SetObject(nil)
	delete(o.Users.Data, user.GetID())
	o.sendObjectEnterLeaveStats(user, 1)
	return true, nil
}

func (o *Object) SendToUser(userID mid.ID, msg *websocket.PreparedMessage, recursive bool) error {
	return errors.Errorf("implement me")
}

func (o *Object) Send(msg *websocket.PreparedMessage, recursive bool) error {
	if msg == nil {
		cute.SetTitleColor(cute.BrightRed)
		cute.SetMessageColor(cute.Red)
		cute.Printf("Object: Send", "%+v", errors.WithStack(errors.Errorf("empty message received")))
		return nil
	}

	if o.GetEnabled() {
		if o.numSendsQueued.Add(1) < 0 {
			return nil
		}
		o.broadcastPipeline <- msg
	}

	if !recursive {
		return nil
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	var errs *multierror.Error
	for _, child := range o.Children.Data {
		if err := child.Send(msg, true); err != nil {
			errs = multierror.Append(
				errs, errors.WithMessagef(err, "failed to send message to child: %s", child.GetID()),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (o *Object) performBroadcast(message *websocket.PreparedMessage) {
	o.Users.Mu.RLock()
	defer o.Users.Mu.RUnlock()

	for _, user := range o.Users.Data {
		user.Send(message)
	}
}
