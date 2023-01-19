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

var _ universe.ObjectType = (*SpaceType)(nil)

type SpaceType struct {
	id           uuid.UUID
	ctx          context.Context
	log          *zap.SugaredLogger
	db           database.DB
	mu           sync.RWMutex
	name         string
	categoryName string
	description  *string
	options      *entry.ObjectOptions
	asset2d      universe.Asset2d
	asset3d      universe.Asset3d
}

func NewSpaceType(id uuid.UUID, db database.DB) *SpaceType {
	return &SpaceType{
		id: id,
		db: db,
		options: &entry.ObjectOptions{
			AllowedSubObjects: []uuid.UUID{},
			Minimap:           utils.GetPTR(true),
			Visible:           utils.GetPTR(entry.ReactUnityObjectVisibleType),
			Editable:          utils.GetPTR(true),
			Private:           utils.GetPTR(false),
		},
	}
}

func (s *SpaceType) GetID() uuid.UUID {
	return s.id
}

func (s *SpaceType) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
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
		if err := s.db.GetObjectTypesDB().UpdateObjectTypeName(s.ctx, s.GetID(), name); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.name = name

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
		if err := s.db.GetObjectTypesDB().UpdateObjectTypeCategoryName(s.ctx, s.GetID(), categoryName); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.categoryName = categoryName

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
		if err := s.db.GetObjectTypesDB().UpdateObjectTypeDescription(s.ctx, s.GetID(), description); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.description = description

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
		if err := s.db.GetAssets2dDB().UpsertAsset(s.ctx, asset2d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset2d = asset2d

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
		if err := s.db.GetAssets3dDB().UpsertAsset(s.ctx, asset3d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset3d = asset3d

	return nil
}

func (s *SpaceType) GetOptions() *entry.ObjectOptions {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.options
}

func (s *SpaceType) SetOptions(modifyFn modify.Fn[entry.ObjectOptions], updateDB bool) (*entry.ObjectOptions, error) {
	options, err := func() (*entry.ObjectOptions, error) {
		s.mu.Lock()
		defer s.mu.Unlock()

		options, err := modifyFn(s.options)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to modify options")
		}

		if updateDB {
			if err := s.db.GetObjectTypesDB().UpdateObjectTypeOptions(s.ctx, s.GetID(), options); err != nil {
				return nil, errors.WithMessage(err, "failed to update db")
			}
		}

		s.options = options

		return options, nil
	}()
	if err != nil {
		return nil, err
	}

	for _, space := range universe.GetNode().GetAllObjects() {
		spaceType := space.GetObjectType()
		if spaceType == nil {
			continue
		}
		if spaceType.GetID() != s.GetID() {
			continue
		}
		space.DropCache()
	}

	return options, nil
}

func (s *SpaceType) GetEntry() *entry.ObjectType {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := &entry.ObjectType{
		ObjectTypeID:   s.id,
		ObjectTypeName: s.name,
		CategoryName:   s.categoryName,
		Description:    s.description,
		Options:        s.options,
	}
	if s.asset2d != nil {
		entry.Asset2dID = utils.GetPTR(s.asset2d.GetID())
	}
	if s.asset3d != nil {
		entry.Asset3dID = utils.GetPTR(s.asset3d.GetID())
	}

	return entry
}

func (s *SpaceType) LoadFromEntry(entry *entry.ObjectType) error {
	if entry.ObjectTypeID != s.GetID() {
		return errors.Errorf("space type ids mismatch: %s != %s", entry.ObjectTypeID, s.GetID())
	}

	if err := s.SetName(entry.ObjectTypeName, false); err != nil {
		return errors.WithMessage(err, "failed to set name")
	}
	if err := s.SetCategoryName(entry.CategoryName, false); err != nil {
		return errors.WithMessage(err, "failed to set category name")
	}
	if err := s.SetDescription(entry.Description, false); err != nil {
		return errors.WithMessage(err, "failed to set description")
	}
	if _, err := s.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	node := universe.GetNode()
	if entry.Asset2dID != nil {
		asset2d, ok := node.GetAssets2d().GetAsset2d(*entry.Asset2dID)
		if !ok {
			return errors.Errorf("asset 2d not found: %s", entry.Asset2dID)
		}
		if err := s.SetAsset2d(asset2d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 2d: %s", entry.Asset2dID)
		}
	}

	if entry.Asset3dID != nil {
		asset3d, ok := node.GetAssets3d().GetAsset3d(*entry.Asset3dID)
		if !ok {
			return errors.Errorf("asset 3d not found: %s", entry.Asset3dID)
		}
		if err := s.SetAsset3d(asset3d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 3d: %s", entry.Asset3dID)
		}
	}

	return nil
}
