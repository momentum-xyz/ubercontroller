package attributes

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
	"github.com/momentum-xyz/ubercontroller/universe/attribute"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Attributes = (*Attributes)(nil)

type Attributes struct {
	ctx        context.Context
	log        *zap.SugaredLogger
	db         database.DB
	attributes *generic.SyncMap[entry.AttributeID, universe.Attribute]
}

func NewAttributes(db database.DB) *Attributes {
	return &Attributes{
		db:         db,
		attributes: generic.NewSyncMap[entry.AttributeID, universe.Attribute](),
	}
}

func (a *Attributes) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *Attributes) NewAttributeWithNameAndPluginID(pluginId uuid.UUID, name string) (universe.Attribute, error) {
	id := entry.AttributeID{PluginID: pluginId, Name: name}
	return a.NewAttribute(id)
}

func (a *Attributes) NewAttribute(attributeId entry.AttributeID) (universe.Attribute, error) {

	attribute := attribute.NewAttribute(attributeId, a.db)

	if err := attribute.Initialize(a.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize attribute: %s", attributeId)
	}
	if err := a.AddAttribute(attribute, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add attribute: %s", attributeId)
	}

	return attribute, nil
}

func (a *Attributes) GetAttribute(attributeID entry.AttributeID) (universe.Attribute, bool) {
	attribute, ok := a.attributes.Load(attributeID)
	return attribute, ok
}

func (a *Attributes) GetAttributes() map[entry.AttributeID]universe.Attribute {
	a.attributes.Mu.RLock()
	defer a.attributes.Mu.RUnlock()

	attributes := make(map[entry.AttributeID]universe.Attribute, len(a.attributes.Data))

	for id, attribute := range a.attributes.Data {
		attributes[id] = attribute
	}

	return attributes
}

func (a *Attributes) AddAttribute(attribute universe.Attribute, updateDB bool) error {
	a.attributes.Mu.Lock()
	defer a.attributes.Mu.Unlock()

	if updateDB {
		if err := a.db.AttributesUpsertAttribute(a.ctx, attribute.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.attributes.Data[attribute.GetID()] = attribute

	return nil
}

func (a *Attributes) AddAttributes(attributes []universe.Attribute, updateDB bool) error {
	a.attributes.Mu.Lock()
	defer a.attributes.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.Attribute, len(attributes))
		for i := range attributes {
			entries[i] = attributes[i].GetEntry()
		}
		if err := a.db.AttributesUpsertAttributes(a.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range attributes {
		a.attributes.Data[attributes[i].GetID()] = attributes[i]
	}

	return nil
}

func (a *Attributes) RemoveAttribute(attribute universe.Attribute, updateDB bool) error {
	a.attributes.Mu.Lock()
	defer a.attributes.Mu.Unlock()

	if _, ok := a.attributes.Data[attribute.GetID()]; !ok {
		return errors.Errorf("attribute not found")
	}

	if updateDB {
		if err := a.db.AttributesRemoveAttributeByID(a.ctx, attribute.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	delete(a.attributes.Data, attribute.GetID())

	return nil
}

func (a *Attributes) RemoveAttributes(attributes []universe.Attribute, updateDB bool) error {
	a.attributes.Mu.Lock()
	defer a.attributes.Mu.Unlock()

	for i := range attributes {
		if _, ok := a.attributes.Data[attributes[i].GetID()]; !ok {
			return errors.Errorf("attribute not found: %s", attributes[i].GetID())
		}
	}

	if updateDB {
		ids := make([]entry.AttributeID, len(attributes))
		for i := range attributes {
			ids[i] = attributes[i].GetID()
		}
		if err := a.db.AttributesRemoveAttributesByIDs(a.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range attributes {
		delete(a.attributes.Data, attributes[i].GetID())
	}

	return nil
}

func (a *Attributes) Load() error {
	a.log.Info("Loading attributes...")

	entries, err := a.db.AttributesGetAttributes(a.ctx)

	if err != nil {
		return errors.WithMessage(err, "failed to get attributes")
	}

	for i := range entries {
		attribute, err := a.NewAttribute(*entries[i].AttributeID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new attribute: %s", entries[i].AttributeID)
		}
		if err := attribute.LoadFromEntry(entries[i]); err != nil {
			return errors.WithMessagef(err, "failed to load attribute from entry: %s", entries[i].AttributeID)
		}
		a.attributes.Store(*entries[i].AttributeID, attribute)
	}

	universe.GetNode().AddAPIRegister(a)

	a.log.Info("Attributes loaded")

	return nil
}

func (a *Attributes) Save() error {
	a.log.Info("Saving attributes...")

	a.attributes.Mu.RLock()
	defer a.attributes.Mu.RUnlock()

	entries := make([]*entry.Attribute, 0, len(a.attributes.Data))
	for _, attribute := range a.attributes.Data {
		entries = append(entries, attribute.GetEntry())
	}

	if err := a.db.AttributesUpsertAttributes(a.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert attributes")
	}

	a.log.Info("Attributes saved")

	return nil
}
