package space_types

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/space_type"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.ObjectTypes = (*SpaceTypes)(nil)

type SpaceTypes struct {
	ctx        context.Context
	log        *zap.SugaredLogger
	db         database.DB
	spaceTypes *generic.SyncMap[uuid.UUID, universe.ObjectType]
}

func NewSpaceTypes(db database.DB) *SpaceTypes {
	return &SpaceTypes{
		db:         db,
		spaceTypes: generic.NewSyncMap[uuid.UUID, universe.ObjectType](0),
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

func (s *SpaceTypes) CreateObjectType(spaceTypeID uuid.UUID) (universe.ObjectType, error) {
	spaceType := space_type.NewSpaceType(spaceTypeID, s.db)

	if err := spaceType.Initialize(s.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize space type: %s", spaceTypeID)
	}
	if err := s.AddObjectType(spaceType, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add space type: %s", spaceTypeID)
	}

	return spaceType, nil
}

func (s *SpaceTypes) FilterObjectTypes(predicateFn universe.ObjectTypesFilterPredicateFn) map[uuid.UUID]universe.ObjectType {
	return s.spaceTypes.Filter(predicateFn)
}

func (s *SpaceTypes) GetObjectType(spaceTypeID uuid.UUID) (universe.ObjectType, bool) {
	spaceType, ok := s.spaceTypes.Load(spaceTypeID)
	return spaceType, ok
}

func (s *SpaceTypes) GetObjectTypes() map[uuid.UUID]universe.ObjectType {
	s.spaceTypes.Mu.RLock()
	defer s.spaceTypes.Mu.RUnlock()

	spaceTypes := make(map[uuid.UUID]universe.ObjectType, len(s.spaceTypes.Data))

	for id, spaceType := range s.spaceTypes.Data {
		spaceTypes[id] = spaceType
	}

	return spaceTypes
}

func (s *SpaceTypes) AddObjectType(spaceType universe.ObjectType, updateDB bool) error {
	s.spaceTypes.Mu.Lock()
	defer s.spaceTypes.Mu.Unlock()

	if updateDB {
		if err := s.db.GetObjectTypesDB().UpsertObjectType(s.ctx, spaceType.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.spaceTypes.Data[spaceType.GetID()] = spaceType

	return nil
}

func (s *SpaceTypes) AddObjectTypes(spaceTypes []universe.ObjectType, updateDB bool) error {
	s.spaceTypes.Mu.Lock()
	defer s.spaceTypes.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.ObjectType, len(spaceTypes))
		for i := range spaceTypes {
			entries[i] = spaceTypes[i].GetEntry()
		}
		if err := s.db.GetObjectTypesDB().UpsertObjectTypes(s.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range spaceTypes {
		s.spaceTypes.Data[spaceTypes[i].GetID()] = spaceTypes[i]
	}

	return nil
}

func (s *SpaceTypes) RemoveObjectType(spaceType universe.ObjectType, updateDB bool) (bool, error) {
	s.spaceTypes.Mu.Lock()
	defer s.spaceTypes.Mu.Unlock()

	if _, ok := s.spaceTypes.Data[spaceType.GetID()]; !ok {
		return false, nil
	}

	if updateDB {
		if err := s.db.GetObjectTypesDB().RemoveObjectTypeByID(s.ctx, spaceType.GetID()); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(s.spaceTypes.Data, spaceType.GetID())

	return true, nil
}

func (s *SpaceTypes) RemoveObjectTypes(spaceTypes []universe.ObjectType, updateDB bool) (bool, error) {
	s.spaceTypes.Mu.Lock()
	defer s.spaceTypes.Mu.Unlock()

	for i := range spaceTypes {
		if _, ok := s.spaceTypes.Data[spaceTypes[i].GetID()]; !ok {
			return false, nil
		}
	}

	if updateDB {
		ids := make([]uuid.UUID, len(spaceTypes))
		for i := range spaceTypes {
			ids[i] = spaceTypes[i].GetID()
		}
		if err := s.db.GetObjectTypesDB().RemoveObjectTypesByIDs(s.ctx, ids); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range spaceTypes {
		delete(s.spaceTypes.Data, spaceTypes[i].GetID())
	}

	return true, nil
}

func (s *SpaceTypes) Load() error {
	s.log.Info("Loading space types...")

	entries, err := s.db.GetObjectTypesDB().GetObjectTypes(s.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get space types")
	}

	for i := range entries {
		spaceType, err := s.CreateObjectType(entries[i].ObjectTypeID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new space type: %s", entries[i].ObjectTypeID)
		}
		if err := spaceType.LoadFromEntry(entries[i]); err != nil {
			return errors.WithMessagef(err, "failed to load space type from entry: %s", entries[i].ObjectTypeID)
		}
	}

	universe.GetNode().AddAPIRegister(s)

	s.log.Infof("Object types loaded: %d", s.spaceTypes.Len())

	return nil
}

func (s *SpaceTypes) Save() error {
	s.log.Info("Saving spate types...")

	s.spaceTypes.Mu.RLock()
	defer s.spaceTypes.Mu.RUnlock()

	entries := make([]*entry.ObjectType, 0, len(s.spaceTypes.Data))
	for _, spaceType := range s.spaceTypes.Data {
		entries = append(entries, spaceType.GetEntry())
	}

	if err := s.db.GetObjectTypesDB().UpsertObjectTypes(s.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert space types")
	}

	s.log.Info("Object types saved")

	return nil
}
