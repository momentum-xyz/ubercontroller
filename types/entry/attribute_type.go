package entry

import "github.com/momentum-xyz/ubercontroller/utils/mid"

type AttributeTypeID struct {
	PluginID mid.ID `db:"plugin_id" json:"plugin_id"`
	Name     string `db:"attribute_name" json:"attribute_name"`
}

type AttributeType struct {
	AttributeTypeID
	Description *string           `db:"description" json:"description"`
	Options     *AttributeOptions `db:"options" json:"options"`
}
