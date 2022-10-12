package entry

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type AttributeID struct {
	PluginID uuid.UUID `db:"plugin_id"`
	Name     string    `db:"attribute_name"`
}

type Attribute struct {
	*AttributeID
	Description *string           `db:"description"`
	Options     *AttributeOptions `db:"options"`
}

type AttributePayload struct {
	Value   *AttributeValue   `db:"value"`
	Options *AttributeOptions `db:"options"`
}

type AttributeValue map[string]any

type AttributeOptions map[string]any

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

func NewAttributeOptions() *AttributeOptions {
	return utils.GetPTR(AttributeOptions(make(map[string]any)))
}
