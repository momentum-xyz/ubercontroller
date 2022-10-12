package entry

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type Attribute struct {
	*AttributeID
	Description *string           `db:"description"`
	Options     *AttributeOptions `db:"options"`
}

type AttributeID struct {
	PluginID uuid.UUID `db:"plugin_id"`
	Name     string    `db:"attribute_name"`
}

type AttributeOptions map[string]any

type AttributeValue map[string]any

func NewAttributeOptions() *AttributeOptions {
	return utils.GetPTR(AttributeOptions(make(map[string]any)))
}

func NewAttributeValue() *AttributeValue {
	return utils.GetPTR(AttributeValue(make(map[string]any)))
}
