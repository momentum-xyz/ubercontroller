package entry

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type AttributeTypeID struct {
	PluginID uuid.UUID `db:"plugin_id" mapstructure:"plugin_id"`
	Name     string    `db:"attribute_name" mapstructure:"attribute_name"`
}

type AttributeType struct {
	AttributeTypeID
	Description *string           `db:"description"`
	Options     *AttributeOptions `db:"options"`
}

type AttributeOptions map[string]any

func NewAttributeTypeID(pluginID uuid.UUID, name string) AttributeTypeID {
	return AttributeTypeID{
		PluginID: pluginID,
		Name:     name,
	}
}

func NewAttributeOptions() *AttributeOptions {
	return utils.GetPTR(AttributeOptions(make(map[string]any)))
}
