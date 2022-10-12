package attribute

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

type Attribute struct {
	ctx context.Context
	log *zap.SugaredLogger
	db  database.DB
	mu  sync.RWMutex

	id          entry.AttributeID
	description *string
	options     *entry.AttributeOptions
	entry       *entry.Attribute
}

func NewAttribute(id entry.AttributeID, db database.DB) *Attribute {
	return &Attribute{
		db:      db,
		id:      id,
		options: entry.NewAttributeOptions(),
	}
}

func NewAttributeWithNameAndPluginID(pluginID uuid.UUID, name string, db database.DB) *Attribute {
	return NewAttribute(entry.AttributeID{PluginID: pluginID, Name: name}, db)
}

func (a *Attribute) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *Attribute) GetID() entry.AttributeID {
	return a.id
}

func (a *Attribute) GetName() string {
	return a.id.Name
}

func (a *Attribute) GetPluginID() uuid.UUID {
	return a.id.PluginID
}

func (a *Attribute) GetOptions() *entry.AttributeOptions {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.options
}

func (a *Attribute) GetDescription() *string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.description
}

func (a *Attribute) SetOptions(modifyFn modify.Fn[entry.AttributeOptions], updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	options := modifyFn(a.options)

	if updateDB {
		if err := a.db.AttributesUpdateAttributeOptions(a.ctx, a.id, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.options = options

	return nil
}

func (a *Attribute) SetDescription(description *string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.AttributesUpdateAttributeDescription(a.ctx, a.id, description); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.description = description

	return nil
}

func (a *Attribute) GetEntry() *entry.Attribute {
	a.mu.Lock()
	defer a.mu.Unlock()

	return &entry.Attribute{
		AttributeID: utils.GetPTR(a.id),
		Description: a.description,
		Options:     a.options,
	}
}

func (a *Attribute) LoadFromEntry(entry *entry.Attribute) error {
	a.id = *entry.AttributeID

	if err := a.SetDescription(entry.Description, false); err != nil {
		return errors.WithMessage(err, "failed to set description")
	}
	if err := a.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	return nil
}
