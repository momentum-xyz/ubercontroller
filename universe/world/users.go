package world

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) GetUser(userID uuid.UUID, recursive bool) (universe.User, bool) {
	return w.Space.GetUser(userID, false)
}

func (w *World) GetUsers(recursive bool) map[uuid.UUID]universe.User {
	return w.Space.GetUsers(false)
}

func (w *World) AddUser(user universe.User, updateDB bool) error {
	w.Users.Mu.Lock()
	defer w.Users.Mu.Unlock()

	if user.GetSpace() != nil && user.GetSpace().GetWorld().GetID() != w.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", user.GetSpace().GetWorld().GetID(), w.GetID())
	}
	if err := user.SetWorld(w, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set world %s to user: %s", w.GetID(), user.GetID())
	}
	w.Users.Data[user.GetID()] = user

	return nil
}

func (w *World) RemoveUser(user universe.User, updateDB bool) error {
	w.Users.Mu.Lock()
	defer w.Users.Mu.Unlock()

	if user.GetWorld().GetID() != w.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", user.GetWorld().GetID(), w.GetID())
	}
	if err := user.SetWorld(nil, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set world nil to user: %s", user.GetID())
	}
	delete(w.Users.Data, user.GetID())

	return nil
}

func (w *World) SendToUser(userID uuid.UUID, msg *websocket.PreparedMessage, recursive bool) error {
	return errors.Errorf("implement me")
}

func (w *World) SendToUsers(msg *websocket.PreparedMessage, recursive bool) error {
	return errors.Errorf("implement me")
}
