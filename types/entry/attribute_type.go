package entry

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type AttributeTypeID struct {
	PluginID umid.UMID `db:"plugin_id" json:"plugin_id"`
	Name     string    `db:"attribute_name" json:"attribute_name"`
}

type AttributeType struct {
	AttributeTypeID
	Description *string           `db:"description" json:"description"`
	Options     *AttributeOptions `db:"options" json:"options"`
}

func NewAttributeTypeID(pluginID umid.UMID, name string) AttributeTypeID {
	return AttributeTypeID{
		PluginID: pluginID,
		Name:     name,
	}
}
