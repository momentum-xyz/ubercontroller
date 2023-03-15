package object_types

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/utils/mid"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/object_type"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.ObjectTypes = (*ObjectTypes)(nil)

type ObjectTypes struct {
	ctx         context.Context
	log         *zap.SugaredLogger
	db          database.DB
	objectTypes *generic.SyncMap[mid.ID, universe.ObjectType]
}

func NewObjectTypes(db database.DB) *ObjectTypes {
	return &ObjectTypes{
		db:          db,
		objectTypes: generic.NewSyncMap[mid.ID, universe.ObjectType](0),
	}
}

func (ot *ObjectTypes) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	ot.ctx = ctx
	ot.log = log

	return nil
}

func (ot *ObjectTypes) CreateObjectType(objectTypeID mid.ID) (universe.ObjectType, error) {
	objectType := object_type.NewObjectType(objectTypeID, ot.db)

	if err := objectType.Initialize(ot.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize object type: %s", objectTypeID)
	}
	if err := ot.AddObjectType(objectType, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add object type: %s", objectTypeID)
	}

	return objectType, nil
}

func (ot *ObjectTypes) FilterObjectTypes(predicateFn universe.ObjectTypesFilterPredicateFn) map[mid.ID]universe.ObjectType {
	return ot.objectTypes.Filter(predicateFn)
}

func (ot *ObjectTypes) GetObjectType(objectTypeID mid.ID) (universe.ObjectType, bool) {
	return ot.objectTypes.Load(objectTypeID)
}

func (ot *ObjectTypes) GetObjectTypes() map[mid.ID]universe.ObjectType {
	ot.objectTypes.Mu.RLock()
	defer ot.objectTypes.Mu.RUnlock()

	objectTypes := make(map[mid.ID]universe.ObjectType, len(ot.objectTypes.Data))

	for id, objectType := range ot.objectTypes.Data {
		objectTypes[id] = objectType
	}

	return objectTypes
}

func (ot *ObjectTypes) AddObjectType(objectType universe.ObjectType, updateDB bool) error {
	ot.objectTypes.Mu.Lock()
	defer ot.objectTypes.Mu.Unlock()

	if updateDB {
		if err := ot.db.GetObjectTypesDB().UpsertObjectType(ot.ctx, objectType.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	ot.objectTypes.Data[objectType.GetID()] = objectType

	return nil
}

func (ot *ObjectTypes) AddObjectTypes(objectTypes []universe.ObjectType, updateDB bool) error {
	ot.objectTypes.Mu.Lock()
	defer ot.objectTypes.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.ObjectType, len(objectTypes))
		for i := range objectTypes {
			entries[i] = objectTypes[i].GetEntry()
		}
		if err := ot.db.GetObjectTypesDB().UpsertObjectTypes(ot.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range objectTypes {
		ot.objectTypes.Data[objectTypes[i].GetID()] = objectTypes[i]
	}

	return nil
}

func (ot *ObjectTypes) RemoveObjectType(objectType universe.ObjectType, updateDB bool) (bool, error) {
	ot.objectTypes.Mu.Lock()
	defer ot.objectTypes.Mu.Unlock()

	if _, ok := ot.objectTypes.Data[objectType.GetID()]; !ok {
		return false, nil
	}

	if updateDB {
		if err := ot.db.GetObjectTypesDB().RemoveObjectTypeByID(ot.ctx, objectType.GetID()); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(ot.objectTypes.Data, objectType.GetID())

	return true, nil
}

func (ot *ObjectTypes) RemoveObjectTypes(objectTypes []universe.ObjectType, updateDB bool) (bool, error) {
	ot.objectTypes.Mu.Lock()
	defer ot.objectTypes.Mu.Unlock()

	for i := range objectTypes {
		if _, ok := ot.objectTypes.Data[objectTypes[i].GetID()]; !ok {
			return false, nil
		}
	}

	if updateDB {
		ids := make([]mid.ID, len(objectTypes))
		for i := range objectTypes {
			ids[i] = objectTypes[i].GetID()
		}
		if err := ot.db.GetObjectTypesDB().RemoveObjectTypesByIDs(ot.ctx, ids); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range objectTypes {
		delete(ot.objectTypes.Data, objectTypes[i].GetID())
	}

	return true, nil
}

func (ot *ObjectTypes) Load() error {
	ot.log.Info("Loading object types...")

	entries, err := ot.db.GetObjectTypesDB().GetObjectTypes(ot.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get object types")
	}

	for _, otEntry := range entries {
		objectType, err := ot.CreateObjectType(otEntry.ObjectTypeID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new object type: %s", otEntry.ObjectTypeID)
		}
		if err := objectType.LoadFromEntry(otEntry); err != nil {
			return errors.WithMessagef(err, "failed to load object type from entry: %s", otEntry.ObjectTypeID)
		}
	}

	universe.GetNode().AddAPIRegister(ot)

	ot.log.Infof("Object types loaded: %d", ot.objectTypes.Len())

	return nil
}

func (ot *ObjectTypes) Save() error {
	ot.log.Info("Saving object types...")

	ot.objectTypes.Mu.RLock()
	defer ot.objectTypes.Mu.RUnlock()

	entries := make([]*entry.ObjectType, 0, len(ot.objectTypes.Data))
	for _, objectType := range ot.objectTypes.Data {
		entries = append(entries, objectType.GetEntry())
	}

	if err := ot.db.GetObjectTypesDB().UpsertObjectTypes(ot.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert object types")
	}

	ot.log.Infof("Object types saved: %d", len(ot.objectTypes.Data))

	return nil
}
