package space

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func (s *Space) GetUser(userID uuid.UUID, recursive bool) (universe.User, bool) {
	user, ok := s.Users.Load(userID)
	if ok {
		return user, true
	}

	if !recursive {
		return nil, false
	}

	s.children.Mu.RLock()
	defer s.children.Mu.RUnlock()

	for _, child := range s.children.Data {
		if user, ok := child.GetUser(userID, recursive); ok {
			return user, true
		}
	}

	return nil, false
}

// GetUsers return new sync map with all nested users if recursive is true,
// otherwise the method return existing sync map with users dependent only to current space.
func (s *Space) GetUsers(recursive bool) *generics.SyncMap[uuid.UUID, universe.User] {
	if !recursive {
		return s.Users
	}

	users := generics.NewSyncMap[uuid.UUID, universe.User]()

	s.Users.Mu.RLock()
	for id, user := range s.Users.Data {
		users.Data[id] = user
	}
	defer s.Users.Mu.RUnlock()

	s.children.Mu.RLock()
	defer s.children.Mu.RUnlock()

	// maybe we will need lock here in future
	for _, space := range s.children.Data {
		for id, user := range space.GetUsers(recursive).Data {
			users.Data[id] = user
		}
	}

	return users
}

func (s *Space) AddUser(user universe.User, updateDB bool) error {
	s.Users.Mu.Lock()
	defer s.Users.Mu.Unlock()

	if err := user.SetSpace(s, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set space to user: %s", user.GetID())
	}
	s.Users.Data[user.GetID()] = user

	return nil
}

func (s *Space) RemoveUser(user universe.User, updateDB bool) error {
	s.Users.Mu.Lock()
	defer s.Users.Mu.Unlock()

	if err := user.SetSpace(nil, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set space to user: %s", user.GetID())
	}
	delete(s.Users.Data, user.GetID())

	return nil
}

func (s *Space) SendToUser(userID uuid.UUID, msg *websocket.PreparedMessage, recursive bool) error {
	return errors.Errorf("implement me")
}

func (s *Space) Broadcast(msg *websocket.PreparedMessage, recursive bool) error {
	return errors.Errorf("implement me")
}
