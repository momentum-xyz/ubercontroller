package space

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
)

var _ universe.Space = (*Space)(nil)

type Space struct {
	ctx       context.Context
	log       *zap.SugaredLogger
	db        database.DB
	Users     *generics.SyncMap[uuid.UUID, universe.User]
	children  *generics.SyncMap[uuid.UUID, universe.Space]
	mu        sync.RWMutex
	id        uuid.UUID
	world     universe.World
	root      universe.Space
	parent    universe.Space
	theta     float64
	position  cmath.Vec3
	ownerID   uuid.UUID
	options   *universe.SpaceOptionsEntry
	asset2d   universe.Asset2d
	asset3d   universe.Asset3d
	spaceType universe.SpaceType
}

func NewSpace(id uuid.UUID, db database.DB, world universe.World) *Space {
	return &Space{
		id:       id,
		db:       db,
		Users:    generics.NewSyncMap[uuid.UUID, universe.User](),
		children: generics.NewSyncMap[uuid.UUID, universe.Space](),
		world:    world,
	}
}

func (s *Space) GetID() uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

func (s *Space) GetOwnerID() uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ownerID
}

func (s *Space) SetOwnerID(ownerID uuid.UUID, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ownerID = ownerID

	return nil
}

func (s *Space) Update(recursive bool) error {
	return errors.Errorf("implement me")
}

func (s *Space) LoadFromEntry(entry *universe.SpaceEntry) error {
	return errors.Errorf("implement me")
}

func (s *Space) GetAsset2D() universe.Asset2d {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.asset2d
}

func (s *Space) GetAsset3D() universe.Asset3d {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.asset3d
}

func (s *Space) SetAsset2D(asset2d universe.Asset2d, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.asset2d = asset2d

	return nil
}

func (s *Space) SetAsset3D(asset3d universe.Asset3d, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.asset3d = asset3d

	return nil
}

func (s *Space) GetSpaceType() universe.SpaceType {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.spaceType
}

func (s *Space) SetSpaceType(spaceType universe.SpaceType, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.spaceType = spaceType

	return nil
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
