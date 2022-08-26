package world

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/types/generics"
	"github.com/momentum-xyz/controller/universe"
)

func (w *World) GetUser(userID uuid.UUID, recursive bool) (universe.User, bool) {
	user, ok := w.Users.Load(userID)
	return user, ok
}

// GetUsers always returns existing sync map with all users in all dependent spaces.
func (w *World) GetUsers(recursive bool) *generics.SyncMap[uuid.UUID, universe.User] {
	return w.Users
}

func (w *World) AttachUser(user universe.User, updateDB bool) error {
	w.Users.Mu.Lock()
	defer w.Users.Mu.Unlock()

	if err := user.SetWorld(w, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set world to user: %s", user.GetID())
	}
	w.Users.Data[user.GetID()] = user

	return nil
}

// DetachUser detaches user from world and space too.
// TODO: think about rollback on error
func (w *World) DetachUser(user universe.User, updateDB bool) error {
	w.Users.Mu.Lock()
	defer w.Users.Mu.Unlock()

	if err := user.SetWorld(nil, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set world to user: %s", user.GetID())
	}
	if err := user.SetSpace(nil, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set space to user: %s", user.GetID())
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
