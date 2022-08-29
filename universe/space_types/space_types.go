package space_types

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/types/generics"
	"github.com/momentum-xyz/controller/universe"
)

var _ universe.SpaceTypes = (*SpaceTypes)(nil)

type SpaceTypes struct {
	spaceTypes *generics.SyncMap[uuid.UUID, universe.SpaceType]
}

func NewSpaceTypes() *SpaceTypes {
	return &SpaceTypes{
		spaceTypes: generics.NewSyncMap[uuid.UUID, universe.SpaceType](),
	}
}

func (s *SpaceTypes) Initialize(ctx context.Context) error {
	return nil
}

func (s *SpaceTypes) GetSpaceType(spaceTypeID uuid.UUID) (universe.SpaceType, bool) {
	spaceType, ok := s.spaceTypes.Load(spaceTypeID)
	return spaceType, ok
}

func (s *SpaceTypes) Load() error {
	return errors.Errorf("implement me")
}
