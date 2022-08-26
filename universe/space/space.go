package space

import (
	"context"
	"sync"

	"github.com/google/uuid"
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

func (s *Space) GetTheta() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.theta
}

func (s *Space) SetTheta(theta float64, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.theta = theta

	return nil
}

func (s *Space) GetPosition() cmath.Vec3 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.position
}

func (s *Space) SetPosition(position cmath.Vec3, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.position = position

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
