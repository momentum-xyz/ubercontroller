package space_types

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
)

var _ universe.SpaceTypes = (*SpaceTypes)(nil)

type SpaceTypes struct {
	db         database.DB
	spaceTypes *generics.SyncMap[uuid.UUID, universe.SpaceType]
}

func NewSpaceTypes(db database.DB) *SpaceTypes {
	return &SpaceTypes{
		db:         db,
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

func (s *SpaceTypes) Load(ctx context.Context) error {
	return errors.Errorf("implement me")
}
