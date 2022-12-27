package entry

import (
	"github.com/google/uuid"
)

type AttributeTypeID struct {
	PluginID uuid.UUID `db:"plugin_id" json:"plugin_id"`
	Name     string    `db:"attribute_name" json:"attribute_name"`
}

type AttributeType struct {
	AttributeTypeID
	Description *string           `db:"description"`
	Options     *AttributeOptions `db:"options"`
}
