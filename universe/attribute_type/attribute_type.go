package attribute_type

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type AttributeType struct {
	ctx context.Context
	log *zap.SugaredLogger
	db  database.DB
	mu  sync.RWMutex

	id          entry.AttributeTypeID
	description *string
	options     *entry.AttributeOptions
	entry       *entry.AttributeType
}

func NewAttributeType(id entry.AttributeTypeID, db database.DB) *AttributeType {
	return &AttributeType{
		db:      db,
		id:      id,
		options: entry.NewAttributeOptions(),
	}
}

func NewAttributeWithNameAndPluginID(pluginID uuid.UUID, name string, db database.DB) *AttributeType {
	return NewAttributeType(
		entry.AttributeTypeID{
			PluginID: pluginID,
			Name:     name,
		},
		db,
	)
}

func (a *AttributeType) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *AttributeType) GetID() entry.AttributeTypeID {
	return a.id
}

func (a *AttributeType) GetName() string {
	return a.id.Name
}

func (a *AttributeType) GetPluginID() uuid.UUID {
	return a.id.PluginID
}

func (a *AttributeType) GetOptions() *entry.AttributeOptions {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.options
}

func (a *AttributeType) GetDescription() *string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.description
}

func (a *AttributeType) SetOptions(modifyFn modify.Fn[entry.AttributeOptions], updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	options := modifyFn(a.options)

	if updateDB {
		if err := a.db.AttributeTypesUpdateAttributeTypeOptions(a.ctx, a.id, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.options = options

	return nil
}

func (a *AttributeType) SetDescription(description *string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.AttributeTypesUpdateAttributeTypeDescription(a.ctx, a.id, description); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.description = description

	return nil
}

func (a *AttributeType) GetEntry() *entry.AttributeType {
	a.mu.RLock()
	defer a.mu.RLock()

	return &entry.AttributeType{
		AttributeTypeID: a.id,
		Description:     a.description,
		Options:         a.options,
	}
}

func (a *AttributeType) LoadFromEntry(entry *entry.AttributeType) error {
	a.id = entry.AttributeTypeID

	if err := a.SetDescription(entry.Description, false); err != nil {
		return errors.WithMessage(err, "failed to set description")
	}
	if err := a.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	return nil
}
