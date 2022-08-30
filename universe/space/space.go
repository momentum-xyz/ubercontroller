package space

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
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
	ownerID   uuid.UUID
	position  cmath.Vec3
	options   *entry.SpaceOptions
	world     universe.World
	parent    universe.Space
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
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	s.ctx = ctx
	s.log = log

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

	if parent != nil && parent.GetWorld().GetID() != s.GetWorld().GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", parent.GetWorld().GetID(), s.GetWorld().GetID())
	}

	if updateDB {
		var parentID uuid.UUID
		if parent != nil {
			parentID = parent.GetID()
		}
		if err := s.db.SpacesUpdateSpaceParentID(s.ctx, s.id, parentID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.parent = parent

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

	if updateDB {
		if err := s.db.SpacesUpdateSpacePosition(s.ctx, s.id, position); err != nil {
			return errors.WithMessage(err, "failed to udpate db")
		}
	}

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

	if updateDB {
		if err := s.db.SpacesUpdateSpaceOwnerID(s.ctx, s.id, ownerID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.ownerID = ownerID

	return nil
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

	if updateDB {
		var asset2dID uuid.UUID
		if asset2d != nil {
			asset2dID = asset2d.GetID()
		}
		if err := s.db.SpacesUpdateSpaceAsset2dID(s.ctx, s.id, asset2dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset2d = asset2d

	return nil
}

func (s *Space) SetAsset3D(asset3d universe.Asset3d, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		var asset3dID uuid.UUID
		if asset3d != nil {
			asset3dID = asset3d.GetID()
		}
		if err := s.db.SpacesUpdateSpaceAsset3dID(s.ctx, s.id, asset3dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

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

	if updateDB {
		var spaceTypeID uuid.UUID
		if spaceType != nil {
			spaceTypeID = spaceType.GetID()
		}
		if err := s.db.SpacesUpdateSpaceSpaceTypeID(s.ctx, s.id, spaceTypeID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.spaceType = spaceType

	return nil
}

func (s *Space) GetOptions() *entry.SpaceOptions {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.options
}

func (s *Space) GetEffectiveOptions() *entry.SpaceOptions {
	return utils.MergeStructs(s.GetOptions(), s.GetSpaceType().GetOptions())
}

func (s *Space) SetOptions(options *entry.SpaceOptions, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpacesUpdateSpaceOptions(s.ctx, s.id, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.options = options

	return nil
}

func (s *Space) Update(recursive bool) error {
	return errors.Errorf("implement me")
}

func (s *Space) LoadFromEntry(ctx context.Context, entry *entry.Space, recursive bool) error {
	if *entry.SpaceID != s.GetID() {
		return errors.Errorf("space ids mismatch: %s != %s", *entry.SpaceID, s.GetID())
	}

	if err := s.loadSelfData(entry); err != nil {
		return errors.WithMessage(err, "failed to load self data")
	}
	if err := s.loadDependencies(entry); err != nil {
		return errors.WithMessage(err, "failed to load dependencies")
	}

	if !recursive {
		return nil
	}

	spaces, err := s.db.SpacesGetSpacesByParentID(ctx, s.GetID())
	if err != nil {
		return errors.WithMessagef(err, "failed to get spaces by parent id: %s", s.GetID())
	}

	group, gctx := errgroup.WithContext(ctx)

	for i := range spaces {
		entry := spaces[i]

		group.Go(func() error {
			space := NewSpace(*entry.SpaceID, s.db, s.world)

			if err := space.LoadFromEntry(gctx, &entry, recursive); err != nil {
				return errors.WithMessagef(err, "failed to load space from entry: %s", space.GetID())
			}
			if err := space.SetParent(s, false); err != nil {
				return errors.WithMessagef(err, "failed to set parent: %s", space.GetID())
			}
			if err := space.Initialize(gctx); err != nil {
				return errors.WithMessagef(err, "failed to initialize space: %s", space.GetID())
			}

			s.children.Store(space.GetID(), space)

			return nil
		})
	}

	return group.Wait()
}

func (s *Space) loadSelfData(entry *entry.Space) error {
	if err := s.SetOwnerID(*entry.OwnerID, false); err != nil {
		return errors.WithMessagef(err, "failed to set owner id: %s", *entry.OwnerID)
	}
	if err := s.SetPosition(*entry.Position, false); err != nil {
		return errors.WithMessage(err, "failed to set position")
	}
	if err := s.SetOptions(entry.Options, false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}
	return nil
}

func (s *Space) loadDependencies(entry *entry.Space) error {
	node := universe.GetNode()

	spaceType, ok := node.GetSpaceTypes().GetSpaceType(*entry.SpaceTypeID)
	if !ok {
		return errors.Errorf("failed to get space type: %s", *entry.SpaceTypeID)
	}
	if err := s.SetSpaceType(spaceType, false); err != nil {
		return errors.WithMessagef(err, "failed to set space type: %s", *entry.SpaceTypeID)
	}

	if entry.Asset2dID != nil {
		asset2d, ok := node.GetAssets2d().GetAsset2d(*entry.Asset2dID)
		if !ok {
			return errors.Errorf("failed to get asset 2d: %s", *entry.Asset2dID)
		}
		if err := s.SetAsset2D(asset2d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 2d: %s", *entry.Asset2dID)
		}
	}

	if entry.Asset3dID != nil {
		asset3d, ok := node.GetAssets3d().GetAsset3d(*entry.Asset3dID)
		if !ok {
			return errors.Errorf("failed to get asset 3d: %s", *entry.Asset3dID)
		}
		if err := s.SetAsset3D(asset3d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 3d: %s", *entry.Asset3dID)
		}
	}

	return nil
}
