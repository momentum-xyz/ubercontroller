package space_types

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.SpaceTypes = (*SpaceTypes)(nil)

type SpaceTypes struct {
	ctx        context.Context
	log        *zap.SugaredLogger
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
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	s.ctx = ctx
	s.log = log

	return nil
}

func (s *SpaceTypes) GetSpaceType(spaceTypeID uuid.UUID) (universe.SpaceType, bool) {
	spaceType, ok := s.spaceTypes.Load(spaceTypeID)
	return spaceType, ok
}

func (s *SpaceTypes) GetSpaceTypes(spaceTypeIDs []uuid.UUID) (*generics.SyncMap[uuid.UUID, universe.SpaceType], error) {
	s.spaceTypes.Mu.RLock()
	defer s.spaceTypes.Mu.RUnlock()

	spaceTypes := generics.NewSyncMap[uuid.UUID, universe.SpaceType]()

	// maybe we will need lock here in future
	for i := range spaceTypeIDs {
		spaceType, ok := s.spaceTypes.Data[spaceTypeIDs[i]]
		if !ok {
			return nil, errors.Errorf("space type not found: %s", spaceTypeIDs[i])
		}
		spaceTypes.Data[spaceTypeIDs[i]] = spaceType
	}

	return spaceTypes, nil
}

func (s *SpaceTypes) AddSpaceType(spaceType universe.SpaceType, updateDB bool) error {
	s.spaceTypes.Mu.Lock()
	defer s.spaceTypes.Mu.Unlock()

	if _, ok := s.spaceTypes.Data[spaceType.GetID()]; ok {
		return errors.Errorf("space type already exists")
	}

	if err := spaceType.Update(updateDB); err != nil {
		return errors.WithMessage(err, "failed to update space type")
	}

	s.spaceTypes.Data[spaceType.GetID()] = spaceType

	return nil
}

func (s *SpaceTypes) AddSpaceTypes(spaceTypes []universe.SpaceType, updateDB bool) error {
	var errs *multierror.Error
	for i := range spaceTypes {
		if err := s.AddSpaceType(spaceTypes[i], updateDB); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to add space type: %s", spaceTypes[i].GetID()))
		}
	}
	return errs.ErrorOrNil()
}

func (s *SpaceTypes) RemoveSpaceType(spaceType universe.SpaceType, updateDB bool) error {
	s.spaceTypes.Mu.Lock()
	defer s.spaceTypes.Mu.Unlock()

	if _, ok := s.spaceTypes.Data[spaceType.GetID()]; !ok {
		return errors.Errorf("space type not found")
	}

	if updateDB {
		if err := s.db.SpaceTypesRemoveSpaceTypeByID(s.ctx, spaceType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	delete(s.spaceTypes.Data, spaceType.GetID())

	return nil
}

func (s *SpaceTypes) RemoveSpaceTypes(spaceTypes []universe.SpaceType, updateDB bool) error {
	s.spaceTypes.Mu.Lock()
	defer s.spaceTypes.Mu.Unlock()

	for i := range spaceTypes {
		if _, ok := s.spaceTypes.Data[spaceTypes[i].GetID()]; !ok {
			return errors.Errorf("space type not found: %s", spaceTypes[i].GetID())
		}
	}

	if updateDB {
		ids := make([]uuid.UUID, len(spaceTypes))
		for i := range spaceTypes {
			ids[i] = spaceTypes[i].GetID()
		}
		if err := s.db.SpaceTypesRemoveSpaceTypeByIDs(s.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range spaceTypes {
		delete(s.spaceTypes.Data, spaceTypes[i].GetID())
	}

	return nil
}

func (s *SpaceTypes) Load(ctx context.Context) error {
	return errors.Errorf("implement me")
}

func (s *SpaceTypes) Update(updateDB bool) error {
	s.spaceTypes.Mu.RLock()
	defer s.spaceTypes.Mu.RUnlock()

	for _, spaceType := range s.spaceTypes.Data {
		if err := spaceType.Update(updateDB); err != nil {
			return errors.WithMessagef(err, "failed to update space type: %s", spaceType.GetID())
		}
	}

	return nil
}
