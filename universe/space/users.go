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
	s.Users.Mu.RLock()
	users := make(map[uuid.UUID]universe.User, len(s.Users.Data))
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

// TODO: think about rollback on error
func (s *Space) AddUser(user universe.User, updateDB bool) error {
	s.Users.Mu.Lock()
	defer s.Users.Mu.Unlock()

	if user.GetWorld().GetID() != s.GetWorld().GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", user.GetWorld().GetID(), s.GetWorld().GetID())
	}

	if updateDB {
		s.log.Error("Space: AddUser: update database is not supported")
	}

	user.SetSpace(s)
	s.Users.Data[user.GetID()] = user
	s.sendSpaceEnterLeaveStats(user, 1)

	return nil
}

func (s *Space) RemoveUser(user universe.User, updateDB bool) error {
	s.Users.Mu.Lock()
	defer s.Users.Mu.Unlock()

	if user.GetWorld().GetID() != s.GetWorld().GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", user.GetWorld().GetID(), s.GetWorld().GetID())
	}

	if updateDB {
		s.log.Error("Space: RemoveUser: update database is not supported")
	}

	user.SetSpace(nil)
	delete(s.Users.Data, user.GetID())
	s.sendSpaceEnterLeaveStats(user, 1)
	return nil
}

func (s *Space) Send(msg *websocket.PreparedMessage, recursive bool) error {
	//return errors.Errorf("implement me")
	if msg == nil {
		return nil
	}
	if s.numSendsQueued.Add(1) < 0 {
		return nil
	}
	s.broadcastPipeline <- msg

	if recursive {
		s.Children.Mu.RLock()
		defer s.Children.Mu.RUnlock()

		for _, space := range s.Children.Data {
			space.Send(msg, recursive)
		}
	}

	return nil
}

func (s *Space) performBroadcast(message *websocket.PreparedMessage) {
	s.Users.Mu.RLock()
	defer s.Users.Mu.RUnlock()

	for _, user := range s.Users.Data {
		user.Send(message)
	}
}
