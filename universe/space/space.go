package space

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/controller/pkg/cmath"
	"github.com/momentum-xyz/controller/types"
	"github.com/momentum-xyz/controller/types/generics"
	"github.com/momentum-xyz/controller/universe"
)

var _ universe.Space = (*Space)(nil)

type Space struct {
	id       uuid.UUID
	ctx      context.Context
	log      *zap.SugaredLogger
	Users    *generics.SyncMap[uuid.UUID, universe.User]
	children *generics.SyncMap[uuid.UUID, universe.Space]
	mu       sync.RWMutex
	owner    universe.User
	world    universe.World
	root     universe.Space
	parent   universe.Space
	theta    float64
	position cmath.Vec3
	options  *universe.SpaceOptionsEntry
}

func NewSpace(id uuid.UUID, world universe.World) *Space {
	return &Space{
		id:       id,
		Users:    generics.NewSyncMap[uuid.UUID, universe.User](),
		children: generics.NewSyncMap[uuid.UUID, universe.Space](),
		world:    world,
	}
}

func (s *Space) GetID() uuid.UUID {
	return s.id
}

func (s *Space) Initialize(ctx context.Context) error {
	log, ok := ctx.Value(types.ContextLoggerKey).(*zap.SugaredLogger)
	if !ok {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	s.log = log
	s.ctx = ctx

	return nil
}

func (s *Space) GetWorld() universe.World {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.world
}

func (s *Space) GetRoot() universe.Space {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.root
}

func (s *Space) GetParent() universe.Space {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.parent
}

func (s *Space) SetParent(parent universe.Space, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if parent == nil {
		s.root = nil
		s.parent = nil
		return nil
	}

	if parent.GetWorld().GetID() != s.GetWorld().GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", parent.GetWorld().GetID(), s.GetWorld().GetID())
	}

	s.root = parent.GetRoot()
	s.parent = parent

	return nil
}

func (s *Space) GetPosition() cmath.Vec3 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.position
}

func (s *Space) GetTheta() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.theta
}

func (s *Space) SetPosition(position cmath.Vec3, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.position = position

	return nil
}

func (s *Space) SetTheta(theta float64, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.theta = theta

	return nil
}

func (s *Space) Load(recursive bool) error {
	return errors.Errorf("implement me")
}

func (s *Space) GetOptions() *universe.SpaceOptionsEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.options
}

func (s *Space) SetOptions(options *universe.SpaceOptionsEntry, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.options = options

	return nil
}

func (s *Space) GetSpace(spaceID uuid.UUID, recursive bool) (universe.Space, bool) {
	space, ok := s.children.Load(spaceID)
	if ok {
		return space, true
	}

	if !recursive {
		return nil, false
	}

	s.children.Mu.RLock()
	defer s.children.Mu.RUnlock()

	for _, child := range s.children.Data {
		space, ok := child.GetSpace(spaceID, recursive)
		if ok {
			return space, true
		}
	}

	return nil, false
}

// GetSpaces returns new sync map with all nested children if recursive is true,
// otherwise the method returns existing sync map with children dependent only to current space.
func (s *Space) GetSpaces(recursive bool) *generics.SyncMap[uuid.UUID, universe.Space] {
	if !recursive {
		return s.children
	}

	spaces := generics.NewSyncMap[uuid.UUID, universe.Space]()

	s.children.Mu.RLock()
	defer s.children.Mu.RUnlock()

	// maybe we will need lock here in future
	for id, child := range s.children.Data {
		spaces.Data[id] = child

		for id, child := range child.GetSpaces(recursive).Data {
			spaces.Data[id] = child
		}
	}

	return spaces
}

func (s *Space) AttachSpace(space universe.Space, updateDB bool) error {
	s.children.Mu.Lock()
	defer s.children.Mu.Unlock()

	if err := space.SetParent(s, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set parent to space: %s", space.GetID())
	}
	s.children.Data[space.GetID()] = space

	return nil
}

func (s *Space) AttachSpaces(spaces []universe.Space, updateDB bool) error {
	var errs *multierror.Error
	for i := range spaces {
		if err := s.AttachSpace(spaces[i], updateDB); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to attach space: %s", spaces[i].GetID()))
		}
	}
	return errs.ErrorOrNil()
}

func (s *Space) DetachSpace(spaceID uuid.UUID, recursive, updateDB bool) (bool, error) {
	s.children.Mu.Lock()
	space, ok := s.children.Data[spaceID]
	if ok {
		defer s.children.Mu.Unlock()

		if err := space.SetParent(nil, updateDB); err != nil {
			return false, errors.WithMessagef(err, "failed to set parent to space: %s", spaceID)
		}
		delete(s.children.Data, spaceID)

		return true, nil
	}
	s.children.Mu.Unlock()

	if !recursive {
		return true, nil
	}

	s.children.Mu.RLock()
	defer s.children.Mu.RUnlock()

	for _, child := range s.children.Data {
		detached, err := child.DetachSpace(spaceID, recursive, updateDB)
		if err != nil {
			return false, errors.WithMessagef(err, "failed to detach space: %s", spaceID)
		}
		if detached {
			return true, nil
		}
	}

	return false, nil
}

// DetachSpaces returns true in first value if any of spaces with space ids was detached.
func (s *Space) DetachSpaces(spaceIDs []uuid.UUID, recursive, updateDB bool) (bool, error) {
	var res bool
	for i := range spaceIDs {
		detached, err := s.DetachSpace(spaceIDs[i], recursive, updateDB)
		if err != nil {
			return false, errors.WithMessagef(err, "failed to detach space: %s", spaceIDs[i])
		}
		if detached {
			res = true
		}
	}
	return res, nil
}

func (s *Space) GetOwner() universe.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.owner
}

func (s *Space) SetOwner(owner universe.User, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.owner = owner

	return nil
}

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

// GetUsers returns new sync map with all nested Users if recursive is true,
// otherwise the method returns existing sync map with Users dependent only to current space.
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

func (s *Space) AttachUser(user universe.User, updateDB bool) error {
	s.Users.Mu.Lock()
	defer s.Users.Mu.Unlock()

	if err := user.SetSpace(s, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set space to user: %s", user.GetID())
	}
	s.Users.Data[user.GetID()] = user

	return nil
}

func (s *Space) DetachUser(user universe.User, updateDB bool) error {
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

func (s *Space) SendToUsers(msg *websocket.PreparedMessage, recursive bool) error {
	return errors.Errorf("implement me")
}
