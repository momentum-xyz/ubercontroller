package entry

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type AttributeID AttributeTypeID

type Attribute struct {
	AttributeID
	*AttributePayload
}

type AttributePayload struct {
	Value   *AttributeValue   `db:"value"`
	Options *AttributeOptions `db:"options"`
}

type AttributeValue map[string]any

func NewAttributeID(pluginID uuid.UUID, name string) AttributeID {
	return AttributeID{
		PluginID: pluginID,
		Name:     name,
	}
}

func NewAttributePayload(value *AttributeValue, options *AttributeOptions) *AttributePayload {
	return &AttributePayload{
		Value:   value,
		Options: options,
	}
}

func NewAttributeValue() *AttributeValue {
	return utils.GetPTR(AttributeValue(make(map[string]any)))
}
