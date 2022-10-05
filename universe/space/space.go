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
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.Space = (*Space)(nil)

type Space struct {
	id               uuid.UUID
	world            universe.World
	ctx              context.Context
	log              *zap.SugaredLogger
	db               database.DB
	Users            *generic.SyncMap[uuid.UUID, universe.User]
	Children         *generic.SyncMap[uuid.UUID, universe.Space]
	mu               sync.RWMutex
	ownerID          uuid.UUID
	position         *cmath.Vec3
	options          *entry.SpaceOptions
	parent           universe.Space
	asset2d          universe.Asset2d
	asset3d          universe.Asset3d
	spaceType        universe.SpaceType
	entry            *entry.Space
	effectiveOptions *entry.SpaceOptions

	spaceAttributes     universe.AttributeList[entry.AttributeID]
	userSpaceAttributes universe.AttributeList[UserAttributeIndex]
}

func NewSpace(id uuid.UUID, db database.DB, world universe.World) *Space {
	return &Space{
		id:       id,
		db:       db,
		Users:    generic.NewSyncMap[uuid.UUID, universe.User](),
		Children: generic.NewSyncMap[uuid.UUID, universe.Space](),
		world:    world,
	}
}

func (s *Space) GetID() uuid.UUID {
	return s.id
}

// todo: implement this via spaceAttributes
func (s *Space) GetName() string {
	return "unknown"
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

	if parent != nil && parent.GetWorld().GetID() != s.world.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", parent.GetWorld().GetID(), s.world.GetID())
	}

	if updateDB {
		if parent == nil {
			return errors.Errorf("parent is nil")
		}
		if err := s.db.SpacesUpdateSpaceParentID(s.ctx, s.id, parent.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.parent = parent
	s.clearCache()

	return nil
}

func (s *Space) GetPosition() *cmath.Vec3 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.position
}

func (s *Space) SetPosition(position *cmath.Vec3, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpacesUpdateSpacePosition(s.ctx, s.id, position); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.position = position
	s.clearCache()

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
	s.clearCache()

	return nil
}

func (s *Space) GetAsset2D() universe.Asset2d {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.asset2d
}

func (s *Space) SetAsset2D(asset2d universe.Asset2d, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		var asset2dID *uuid.UUID
		if asset2d != nil {
			asset2dID = utils.GetPtr(asset2d.GetID())
		}
		if err := s.db.SpacesUpdateSpaceAsset2dID(s.ctx, s.id, asset2dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset2d = asset2d
	s.clearCache()

	return nil
}

func (s *Space) GetAsset3D() universe.Asset3d {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.asset3d
}

func (s *Space) SetAsset3D(asset3d universe.Asset3d, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		var asset3dID *uuid.UUID
		if asset3d != nil {
			asset3dID = utils.GetPtr(asset3d.GetID())
		}
		if err := s.db.SpacesUpdateSpaceAsset3dID(s.ctx, s.id, asset3dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset3d = asset3d
	s.clearCache()

	return nil
}

func (s *Space) GetSpaceType() universe.SpaceType {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.spaceType
}

func (s *Space) SetSpaceType(spaceType universe.SpaceType, updateDB bool) error {
	if spaceType == nil {
		return errors.Errorf("space type is nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpacesUpdateSpaceSpaceTypeID(s.ctx, s.id, spaceType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.spaceType = spaceType
	s.clearCache()

	return nil
}

func (s *Space) GetOptions() *entry.SpaceOptions {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.options
}

func (s *Space) SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	options := modifyFn(s.options)

	if updateDB {
		if err := s.db.SpacesUpdateSpaceOptions(s.ctx, s.id, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.options = options
	s.clearCache()

	return nil
}

func (s *Space) GetEffectiveOptions() *entry.SpaceOptions {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.effectiveOptions == nil {
		s.effectiveOptions = utils.MergePTRs(s.options, s.spaceType.GetOptions())
	}
	return s.effectiveOptions
}

func (s *Space) GetEntry() *entry.Space {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.entry == nil {
		s.entry = &entry.Space{
			SpaceID:  &s.id,
			OwnerID:  &s.ownerID,
			Options:  s.options,
			Position: s.position,
		}
		if s.spaceType != nil {
			s.entry.SpaceTypeID = utils.GetPtr(s.spaceType.GetID())
		}
		if s.parent != nil {
			s.entry.ParentID = utils.GetPtr(s.parent.GetID())
		}
		if s.asset2d != nil {
			s.entry.Asset2dID = utils.GetPtr(s.asset2d.GetID())
		}
		if s.asset3d != nil {
			s.entry.Asset3dID = utils.GetPtr(s.asset3d.GetID())
		}
	}

	return s.entry
}

func (s *Space) Update(recursive bool) error {
	s.mu.Lock()
	s.clearCache()
	s.mu.Unlock()

	if !recursive {
		return nil
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		if err := child.Update(recursive); err != nil {
			return errors.WithMessagef(err, "failed to update child: %s", child.GetID())
		}
	}

	return nil
}

func (s *Space) LoadFromEntry(entry *entry.Space, recursive bool) error {
	s.log.Debugf("Loading space %s...", *entry.SpaceID)

	if *entry.SpaceID != s.GetID() {
		return errors.Errorf("space ids mismatch: %s != %s", entry.SpaceID, s.GetID())
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

	entries, err := s.db.SpacesGetSpacesByParentID(s.ctx, s.GetID())
	if err != nil {
		return errors.WithMessagef(err, "failed to get spaces by parent id: %s", s.GetID())
	}

	for i := range entries {
		space, err := s.NewSpace(*entries[i].SpaceID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new space: %s", entries[i].SpaceID)
		}
		if err := space.LoadFromEntry(entries[i], recursive); err != nil {
			return errors.WithMessagef(err, "failed to load space from entry: %s", entries[i].SpaceID)
		}
		s.Children.Store(*entries[i].SpaceID, space)
	}

	return nil
}

func (s *Space) loadSelfData(entry *entry.Space) error {
	if err := s.SetOwnerID(*entry.OwnerID, false); err != nil {
		return errors.WithMessagef(err, "failed to set owner id: %s", entry.OwnerID)
	}
	if err := s.SetPosition(entry.Position, false); err != nil {
		return errors.WithMessage(err, "failed to set position")
	}
	if err := s.SetOptions(modify.ReplaceWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}
	return nil
}

func (s *Space) loadDependencies(entry *entry.Space) error {
	node := universe.GetNode()

	spaceType, ok := node.GetSpaceTypes().GetSpaceType(*entry.SpaceTypeID)
	if !ok {
		return errors.Errorf("failed to get space type: %s", entry.SpaceTypeID)
	}
	if err := s.SetSpaceType(spaceType, false); err != nil {
		return errors.WithMessagef(err, "failed to set space type: %s", entry.SpaceTypeID)
	}

	if entry.Asset2dID != nil {
		asset2d, ok := node.GetAssets2d().GetAsset2d(*entry.Asset2dID)
		if !ok {
			return errors.Errorf("failed to get asset 2d: %s", entry.Asset2dID)
		}
		if err := s.SetAsset2D(asset2d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 2d: %s", entry.Asset2dID)
		}
	}

	if entry.Asset3dID != nil {
		asset3d, ok := node.GetAssets3d().GetAsset3d(*entry.Asset3dID)
		if !ok {
			return errors.Errorf("failed to get asset 3d: %s", entry.Asset3dID)
		}
		if err := s.SetAsset3D(asset3d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 3d: %s", entry.Asset3dID)
		}
	}

	return nil
}

func (s *Space) clearCache() {
	s.entry = nil
	s.effectiveOptions = nil
}
