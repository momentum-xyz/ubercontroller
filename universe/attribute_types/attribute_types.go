package attribute_types

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/attribute_type"
)

var _ universe.AttributeTypes = (*AttributeTypes)(nil)

type AttributeTypes struct {
	ctx            types.LoggerContext
	log            *zap.SugaredLogger
	db             database.DB
	attributeTypes *generic.SyncMap[entry.AttributeTypeID, universe.AttributeType]
}

func NewAttributeTypes(db database.DB) *AttributeTypes {
	return &AttributeTypes{
		db:             db,
		attributeTypes: generic.NewSyncMap[entry.AttributeTypeID, universe.AttributeType](0),
	}
}

func (a *AttributeTypes) Initialize(ctx types.LoggerContext) error {
	a.ctx = ctx
	a.log = ctx.Logger()

	return nil
}

func (a *AttributeTypes) FilterAttributeTypes(
	predicateFn universe.AttributeTypesFilterPredicateFn,
) map[entry.AttributeTypeID]universe.AttributeType {
	return a.attributeTypes.Filter(predicateFn)
}

func (a *AttributeTypes) CreateAttributeType(attributeTypeID entry.AttributeTypeID) (universe.AttributeType, error) {
	attributeType := attribute_type.NewAttributeType(attributeTypeID, a.db)

	if err := attributeType.Initialize(a.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize attribute type: %s", attributeTypeID)
	}
	if err := a.AddAttributeType(attributeType, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add attribute type: %s", attributeTypeID)
	}

	return attributeType, nil
}

func (a *AttributeTypes) GetAttributeType(attributeTypeID entry.AttributeTypeID) (universe.AttributeType, bool) {
	attributeType, ok := a.attributeTypes.Load(attributeTypeID)
	return attributeType, ok
}

func (a *AttributeTypes) GetAttributeTypes() map[entry.AttributeTypeID]universe.AttributeType {
	a.attributeTypes.Mu.RLock()
	defer a.attributeTypes.Mu.RUnlock()

	attributeTypes := make(map[entry.AttributeTypeID]universe.AttributeType, len(a.attributeTypes.Data))

	for id, attributeType := range a.attributeTypes.Data {
		attributeTypes[id] = attributeType
	}

	return attributeTypes
}

func (a *AttributeTypes) AddAttributeType(attributeType universe.AttributeType, updateDB bool) error {
	a.attributeTypes.Mu.Lock()
	defer a.attributeTypes.Mu.Unlock()

	if updateDB {
		if err := a.db.GetAttributeTypesDB().UpsertAttributeType(a.ctx, attributeType.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.attributeTypes.Data[attributeType.GetID()] = attributeType

	return nil
}

func (a *AttributeTypes) AddAttributeTypes(attributeTypes []universe.AttributeType, updateDB bool) error {
	a.attributeTypes.Mu.Lock()
	defer a.attributeTypes.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.AttributeType, len(attributeTypes))
		for i := range attributeTypes {
			entries[i] = attributeTypes[i].GetEntry()
		}
		if err := a.db.GetAttributeTypesDB().UpsertAttributeTypes(a.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range attributeTypes {
		a.attributeTypes.Data[attributeTypes[i].GetID()] = attributeTypes[i]
	}

	return nil
}

// TODO: update node/object attributes
func (a *AttributeTypes) RemoveAttributeType(attributeType universe.AttributeType, updateDB bool) (bool, error) {
	a.attributeTypes.Mu.Lock()
	defer a.attributeTypes.Mu.Unlock()

	if _, ok := a.attributeTypes.Data[attributeType.GetID()]; !ok {
		return false, nil
	}

	if updateDB {
		if err := a.db.GetAttributeTypesDB().RemoveAttributeTypeByID(a.ctx, attributeType.GetID()); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(a.attributeTypes.Data, attributeType.GetID())

	return true, nil
}

// TODO: update node/object attributes
func (a *AttributeTypes) RemoveAttributeTypes(attributeTypes []universe.AttributeType, updateDB bool) (bool, error) {
	a.attributeTypes.Mu.Lock()
	defer a.attributeTypes.Mu.Unlock()

	for i := range attributeTypes {
		if _, ok := a.attributeTypes.Data[attributeTypes[i].GetID()]; !ok {
			return false, nil
		}
	}

	if updateDB {
		ids := make([]entry.AttributeTypeID, len(attributeTypes))
		for i := range attributeTypes {
			ids[i] = attributeTypes[i].GetID()
		}
		if err := a.db.GetAttributeTypesDB().RemoveAttributeTypesByIDs(a.ctx, ids); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range attributeTypes {
		delete(a.attributeTypes.Data, attributeTypes[i].GetID())
	}

	return true, nil
}

func (a *AttributeTypes) Load() error {
	a.log.Info("Loading attribute types...")

	entries, err := a.db.GetAttributeTypesDB().GetAttributeTypes(a.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get attribute types")
	}

	for _, atEntry := range entries {
		attributeType, err := a.CreateAttributeType(atEntry.AttributeTypeID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new attribute type: %s", atEntry.AttributeTypeID)
		}
		if err := attributeType.LoadFromEntry(atEntry); err != nil {
			return errors.WithMessagef(err, "failed to load attribute type from entry: %s", atEntry.AttributeTypeID)
		}
	}

	universe.GetNode().AddAPIRegister(a)

	a.log.Infof("Attribute types loaded: %d", a.attributeTypes.Len())

	return nil
}

func (a *AttributeTypes) Save() error {
	a.log.Info("Saving attribute types...")

	a.attributeTypes.Mu.RLock()
	defer a.attributeTypes.Mu.RUnlock()

	entries := make([]*entry.AttributeType, 0, len(a.attributeTypes.Data))
	for _, attributeType := range a.attributeTypes.Data {
		entries = append(entries, attributeType.GetEntry())
	}

	if err := a.db.GetAttributeTypesDB().UpsertAttributeTypes(a.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert attribute types")
	}

	a.log.Infof("Attribute types saved: %d", len(a.attributeTypes.Data))

	return nil
}
