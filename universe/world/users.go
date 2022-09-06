package world

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) GetUser(userID uuid.UUID, recursive bool) (universe.User, bool) {
	user, ok := w.Users.Load(userID)
	return user, ok
}

func (w *World) GetUsers(recursive bool) map[uuid.UUID]universe.User {
	users := make(map[uuid.UUID]universe.User)

	w.Users.Mu.RLock()
	defer w.Users.Mu.RUnlock()

	for id, user := range w.Users.Data {
		users[id] = user
	}

	return users
}

func (w *World) AddUser(user universe.User, updateDB bool) error {
	w.Users.Mu.Lock()
	defer w.Users.Mu.Unlock()

	if err := user.SetWorld(w, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set world to user: %s", user.GetID())
	}
	w.Users.Data[user.GetID()] = user

	return nil
}

// RemoveUser remove user from world and space.
// TODO: think about rollback on error
func (w *World) RemoveUser(user universe.User, updateDB bool) error {
	w.Users.Mu.Lock()
	defer w.Users.Mu.Unlock()

	if err := user.SetSpace(nil, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set space to user: %s", user.GetID())
	}
	if err := user.SetWorld(nil, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set world to user: %s", user.GetID())
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
