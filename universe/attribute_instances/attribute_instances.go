package attribute_instances

import (
	"context"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
)

var _ universe.AttributeInstances[uuid.UUID] = (*AttributeInstances[uuid.UUID])(nil)

type AttributeInstances[indexType comparable] struct {
	data map[indexType]*AttributeInstance
	ctx  context.Context
	log  *zap.SugaredLogger
	db   database.DB
	mu   sync.RWMutex
}

func NewAttributeInstances[indexType comparable](db database.DB) *AttributeInstances[indexType] {
	return &AttributeInstances[indexType]{
		db:   db,
		data: make(map[indexType]*AttributeInstance),
	}
}

func (a AttributeInstances[indexType]) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

//func (a AttributeInstances[indexType]) GetID(id indexType) entry.AttributeID {
//	a.mu.RLock()
//	defer a.mu.RUnlock()
//	return id.(universe.AttributeIndexType).GetAttributeID()
//
//}
//
//func (a AttributeInstances[indexType]) GetName(id indexType) string {
//	a.mu.RLock()
//	defer a.mu.RUnlock()
//	return id.GetAttributeID().Name
//}
//
//func (a AttributeInstances[indexType]) GetPluginID(id indexType) uuid.UUID {
//	a.mu.RLock()
//	defer a.mu.RUnlock()
//	return id.GetAttributeID().PluginID
//}

func (a AttributeInstances[indexType]) GetOptions(id indexType) *entry.AttributeOptions {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.data[id].GetOptions()
}

func (a AttributeInstances[indexType]) GetEffectiveOptions(id indexType) *entry.AttributeOptions {
	//TODO implement me
	panic("implement me")
}

func (a AttributeInstances[indexType]) SetOptions(
	id indexType, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) error {
	//TODO implement me
	panic("implement me")
}

func (a AttributeInstances[indexType]) GetValue(id indexType) *entry.AttributeValue {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if instance, ok := a.data[id]; ok {
		return instance.GetValue()
	}

	return nil
}

func (a AttributeInstances[indexType]) SetValue(id indexType, modifyFn modify.Fn[string], updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (a AttributeInstances[indexType]) SetAttributeInstance(
	id indexType, attr universe.Attribute, value *entry.AttributeValue, options *entry.AttributeOptions,
) universe.AttributeInstance {
	a.mu.Lock()
	defer a.mu.Unlock()
	na := &AttributeInstance{attribute: attr, value: value, options: options}
	a.data[id] = na
	return *na
}
