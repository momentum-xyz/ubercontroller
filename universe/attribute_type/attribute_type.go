package attribute_type

import (
	"context"
	"sync"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/momentum-xyz/ubercontroller/universe"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.AttributeType = (*AttributeType)(nil)

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
		db: db,
		id: id,
	}
}

func (a *AttributeType) Initialize(ctx types.LoggerContext) error {

	a.ctx = ctx
	a.log = ctx.Logger()

	return nil
}

func (a *AttributeType) GetID() entry.AttributeTypeID {
	return a.id
}

func (a *AttributeType) GetName() string {
	return a.id.Name
}

func (a *AttributeType) GetPluginID() umid.UMID {
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

func (a *AttributeType) SetOptions(
	modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	options, err := modifyFn(a.options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := a.db.GetAttributeTypesDB().UpdateAttributeTypeOptions(a.ctx, a.id, options); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	a.options = options

	return options, nil
}

func (a *AttributeType) SetDescription(description *string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.GetAttributeTypesDB().UpdateAttributeTypeDescription(a.ctx, a.id, description); err != nil {
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
	if entry.AttributeTypeID != a.GetID() {
		return errors.Errorf("attribute type ids mismatch: %s != %s", entry.AttributeTypeID, a.GetID())
	}

	if err := a.SetDescription(entry.Description, false); err != nil {
		return errors.WithMessage(err, "failed to set description")
	}
	if _, err := a.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	return nil
}
