package space_type

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.SpaceType = (*SpaceType)(nil)

type SpaceType struct {
	id           uuid.UUID
	ctx          context.Context
	log          *zap.SugaredLogger
	db           database.DB
	mu           sync.RWMutex
	name         string
	categoryName string
	description  *string
	options      *entry.SpaceOptions
	asset2d      universe.Asset2d
	asset3d      universe.Asset3d
	entry        *entry.SpaceType
}

func NewSpaceType(id uuid.UUID, db database.DB) *SpaceType {
	return &SpaceType{
		id: id,
		db: db,
		options: &entry.SpaceOptions{
			AllowedSubspaces: []uuid.UUID{},
			Minimap:          utils.GetPtr(true),
			Visible:          utils.GetPtr(entry.ReactUnitySpaceVisibleType),
			Private:          utils.GetPtr(false),
		},
	}
}

func (s *SpaceType) GetID() uuid.UUID {
	return s.id
}

func (s *SpaceType) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	s.ctx = ctx
	s.log = log

	return nil
}

func (s *SpaceType) GetName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.name
}

func (s *SpaceType) SetName(name string, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpaceTypesUpdateSpaceTypeName(s.ctx, s.id, name); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.name = name
	s.clearCache()

	return nil
}

func (s *SpaceType) GetCategoryName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.categoryName
}

func (s *SpaceType) SetCategoryName(categoryName string, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpaceTypesUpdateSpaceTypeCategoryName(s.ctx, s.id, categoryName); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.categoryName = categoryName
	s.clearCache()

	return nil
}

func (s *SpaceType) GetDescription() *string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.description
}

func (s *SpaceType) SetDescription(description *string, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpaceTypesUpdateSpaceTypeDescription(s.ctx, s.id, description); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.description = description
	s.clearCache()

	return nil
}

func (s *SpaceType) GetAsset2d() universe.Asset2d {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.asset2d
}

func (s *SpaceType) SetAsset2d(asset2d universe.Asset2d, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.Assets2dUpsertAsset(s.ctx, asset2d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to upsert asset 2d")
		}
	}

	s.asset2d = asset2d
	s.clearCache()

	return nil
}

func (s *SpaceType) GetAsset3d() universe.Asset3d {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.asset3d
}

func (s *SpaceType) SetAsset3d(asset3d universe.Asset3d, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.Assets3dUpsertAsset(s.ctx, asset3d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to upsert asset 3d")
		}
	}

	s.asset3d = asset3d
	s.clearCache()

	return nil
}

func (s *SpaceType) GetOptions() *entry.SpaceOptions {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.options
}

func (s *SpaceType) SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) error {
	s.mu.Lock()
	options := modifyFn(s.options)

	if updateDB {
		if err := s.db.SpaceTypesUpdateSpaceTypeOptions(s.ctx, s.id, options); err != nil {
			s.mu.Unlock()
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.options = options
	s.clearCache()
	s.mu.Unlock()

	for _, world := range universe.GetNode().GetWorlds().GetWorlds() {
		for _, space := range world.GetSpaces(true) {
			if space.GetSpaceType() == nil {
				continue
			}
			if space.GetSpaceType().GetID() != s.GetID() {
				continue
			}
			if err := space.Update(true); err != nil {
				return errors.WithMessagef(err, "failed to update space: %s", space.GetID())
			}
		}
	}

	return nil
}

func (s *SpaceType) GetEntry() *entry.SpaceType {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.entry == nil {
		s.entry = &entry.SpaceType{
			SpaceTypeID:   utils.GetPtr(s.id),
			SpaceTypeName: &s.name,
			CategoryName:  &s.categoryName,
			Description:   s.description,
			Options:       s.options,
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

func (s *SpaceType) clearCache() {
	s.entry = nil
}

func (s *SpaceType) LoadFromEntry(entry *entry.SpaceType) error {
	node := universe.GetNode()

	s.id = *entry.SpaceTypeID
	if err := s.SetName(*entry.SpaceTypeName, false); err != nil {
		return errors.WithMessage(err, "failed to set name")
	}
	if err := s.SetCategoryName(*entry.CategoryName, false); err != nil {
		return errors.WithMessage(err, "failed to set category name")
	}
	if err := s.SetDescription(entry.Description, false); err != nil {
		return errors.WithMessage(err, "failed to set description")
	}

	if err := s.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	if entry.Asset2dID != nil {
		asset2d, ok := node.GetAssets2d().GetAsset2d(*entry.Asset2dID)
		if !ok {
			return errors.Errorf("asset 2d not found: %s", *entry.Asset2dID)
		}
		if err := s.SetAsset2d(asset2d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 2d: %s", *entry.Asset2dID)
		}
	}

	if entry.Asset3dID != nil {
		asset3d, ok := node.GetAssets3d().GetAsset3d(*entry.Asset3dID)
		if !ok {
			return errors.Errorf("asset 3d not found: %s", *entry.Asset3dID)
		}
		if err := s.SetAsset3d(asset3d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 3d: %s", *entry.Asset3dID)
		}
	}

	return nil
}
