package space_type

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/controller/types"
	"github.com/momentum-xyz/controller/universe"
)

var _ universe.SpaceType = (*SpaceType)(nil)

type SpaceType struct {
	ctx          context.Context
	log          *zap.SugaredLogger
	mu           sync.RWMutex
	id           uuid.UUID
	name         string
	categoryName string
	description  *string
	options      *universe.SpaceOptionsEntry
}

func NewSpaceType(id uuid.UUID) *SpaceType {
	return &SpaceType{
		id: id,
	}
}

func (s *SpaceType) GetID() uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.id
}

func (s *SpaceType) Initialize(ctx context.Context) error {
	log, ok := ctx.Value(types.ContextLoggerKey).(*zap.SugaredLogger)
	if !ok {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	s.log = log
	s.ctx = ctx

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

	s.description = description

	return nil
}

// LoadFromEntry loads only internal data of the space type exclude nested data like Asset2D, Asset3D, etc.
func (s *SpaceType) LoadFromEntry(entry *universe.SpaceTypeEntry) error {
	return errors.Errorf("implement me")
}

func (s *SpaceType) GetOptions() *universe.SpaceOptionsEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.options
}

func (s *SpaceType) SetOptions(options *universe.SpaceOptionsEntry, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.options = options

	return nil
}
