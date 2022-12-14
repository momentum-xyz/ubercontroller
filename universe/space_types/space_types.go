package space_types

import (
	"context"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/space_type"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var _ universe.SpaceTypes = (*SpaceTypes)(nil)

type SpaceTypes struct {
	ctx        context.Context
	log        *zap.SugaredLogger
	db         database.DB
	spaceTypes *generic.SyncMap[uuid.UUID, universe.SpaceType]
}

func NewSpaceTypes(db database.DB) *SpaceTypes {
	return &SpaceTypes{
		db:         db,
		spaceTypes: generic.NewSyncMap[uuid.UUID, universe.SpaceType](0),
	}
}

func (s *SpaceTypes) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	s.ctx = ctx
	s.log = log

	return nil
}

func (s *SpaceTypes) CreateSpaceType(spaceTypeID uuid.UUID) (universe.SpaceType, error) {
	spaceType := space_type.NewSpaceType(spaceTypeID, s.db)

	if err := spaceType.Initialize(s.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize space type: %s", spaceTypeID)
	}
	if err := s.AddSpaceType(spaceType, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add space type: %s", spaceTypeID)
	}

	return spaceType, nil
}

func (s *SpaceTypes) FilterSpaceTypes(predicateFn universe.SpaceTypesFilterPredicateFn) map[uuid.UUID]universe.SpaceType {
	return s.spaceTypes.Filter(predicateFn)
}

func (s *SpaceTypes) GetSpaceType(spaceTypeID uuid.UUID) (universe.SpaceType, bool) {
	spaceType, ok := s.spaceTypes.Load(spaceTypeID)
	return spaceType, ok
}

func (s *SpaceTypes) GetSpaceTypes() map[uuid.UUID]universe.SpaceType {
	s.spaceTypes.Mu.RLock()
	defer s.spaceTypes.Mu.RUnlock()

	spaceTypes := make(map[uuid.UUID]universe.SpaceType, len(s.spaceTypes.Data))

	for id, spaceType := range s.spaceTypes.Data {
		spaceTypes[id] = spaceType
	}

	return spaceTypes
}

func (s *SpaceTypes) AddSpaceType(spaceType universe.SpaceType, updateDB bool) error {
	s.spaceTypes.Mu.Lock()
	defer s.spaceTypes.Mu.Unlock()

	if updateDB {
		if err := s.db.SpaceTypesUpsertSpaceType(s.ctx, spaceType.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.spaceTypes.Data[spaceType.GetID()] = spaceType

	return nil
}

func (s *SpaceTypes) AddSpaceTypes(spaceTypes []universe.SpaceType, updateDB bool) error {
	s.spaceTypes.Mu.Lock()
	defer s.spaceTypes.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.SpaceType, len(spaceTypes))
		for i := range spaceTypes {
			entries[i] = spaceTypes[i].GetEntry()
		}
		if err := s.db.SpaceTypesUpsertSpaceTypes(s.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range spaceTypes {
		s.spaceTypes.Data[spaceTypes[i].GetID()] = spaceTypes[i]
	}

	return nil
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
		if err := s.db.SpaceTypesRemoveSpaceTypesByIDs(s.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range spaceTypes {
		delete(s.spaceTypes.Data, spaceTypes[i].GetID())
	}

	return nil
}

func (s *SpaceTypes) Load() error {
	s.log.Info("Loading space types...")

	entries, err := s.db.SpaceTypesGetSpaceTypes(s.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get space types")
	}

	for i := range entries {
		spaceType, err := s.CreateSpaceType(entries[i].SpaceTypeID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new space type: %s", entries[i].SpaceTypeID)
		}
		if err := spaceType.LoadFromEntry(entries[i]); err != nil {
			return errors.WithMessagef(err, "failed to load space type from entry: %s", entries[i].SpaceTypeID)
		}
	}

	universe.GetNode().AddAPIRegister(s)

	s.log.Infof("Space types loaded: %d", s.spaceTypes.Len())

	return nil
}

func (s *SpaceTypes) Save() error {
	s.log.Info("Saving spate types...")

	s.spaceTypes.Mu.RLock()
	defer s.spaceTypes.Mu.RUnlock()

	entries := make([]*entry.SpaceType, 0, len(s.spaceTypes.Data))
	for _, spaceType := range s.spaceTypes.Data {
		entries = append(entries, spaceType.GetEntry())
	}

	if err := s.db.SpaceTypesUpsertSpaceTypes(s.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert space types")
	}

	s.log.Info("Space types saved")

	return nil
}
