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
)

var _ universe.SpaceType = (*SpaceType)(nil)

type SpaceType struct {
	ctx          context.Context
	log          *zap.SugaredLogger
	db           database.DB
	mu           sync.RWMutex
	id           uuid.UUID
	name         string
	categoryName string
	description  *string
	options      *entry.SpaceOptions
}

func NewSpaceType(id uuid.UUID, db database.DB) *SpaceType {
	return &SpaceType{
		id: id,
		db: db,
	}
}

func (s *SpaceType) GetID() uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

	return nil
}

func (s *SpaceType) LoadFromEntry(ctx context.Context, entry *entry.SpaceType) error {
	return errors.Errorf("implement me")
}

func (s *SpaceType) GetOptions() *entry.SpaceOptions {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.options
}

func (s *SpaceType) SetOptions(options *entry.SpaceOptions, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpaceTypesUpdateSpaceTypeOptions(s.ctx, s.id, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.options = options

	return nil
}
