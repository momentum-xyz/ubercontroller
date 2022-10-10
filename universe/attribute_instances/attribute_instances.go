package attribute_instances

import (
	"context"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"go.uber.org/zap"
)

var _ universe.AttributeInstances[uuid.UUID] = (*AttributeInstances[uuid.UUID])(nil)

type AttributeInstances[indexType comparable] struct {
	data map[indexType]AttributeInstance
	ctx  context.Context
	log  *zap.SugaredLogger
	db   database.DB
}

func (a AttributeInstances[indexType]) Initialize(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (a AttributeInstances[indexType]) GetID(id indexType) entry.AttributeID {
	//TODO implement me
	panic("implement me")
}

func (a AttributeInstances[indexType]) GetName(id indexType) string {
	//TODO implement me
	panic("implement me")
}

func (a AttributeInstances[indexType]) GetPluginID(id indexType) uuid.UUID {
	//TODO implement me
	panic("implement me")
}

func (a AttributeInstances[indexType]) GetOptions(id indexType) *entry.AttributeOptions {
	//TODO implement me
	panic("implement me")
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

func (a AttributeInstances[indexType]) GetValue(id indexType) *string {
	//TODO implement me
	panic("implement me")
}

func (a AttributeInstances[indexType]) SetValue(id indexType, modifyFn modify.Fn[string], updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (a AttributeInstances[indexType]) AddAttributeInstance(id indexType) {
	//TODO implement me
	panic("implement me")
}
