package space

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

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

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		if user, ok := child.GetUser(userID, recursive); ok {
			return user, true
		}
	}

	return nil, false
}

// GetUsers return map with all nested users if recursive is true,
// otherwise the method return map with users dependent only to current space.
func (s *Space) GetUsers(recursive bool) map[uuid.UUID]universe.User {
	users := make(map[uuid.UUID]universe.User)

	s.Users.Mu.RLock()
	for id, user := range s.Users.Data {
		users[id] = user
	}
	s.Users.Mu.RUnlock()

	if !recursive {
		return users
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, space := range s.Children.Data {
		for id, user := range space.GetUsers(recursive) {
			users[id] = user
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
