package entry

import "github.com/google/uuid"

type Attribute struct {
	AttributeID *AttributeID
	Description *string           `db:"description"`
	Options     *AttributeOptions `db:"options"`
}

type AttributeID struct {
	PluginID uuid.UUID
	Name     string
}

type AttributeOptions map[string]any

type AttributeValue map[string]any

func NewAttributeOptions() *AttributeOptions {
	o := AttributeOptions(make(map[string]any))
	return &o
}

func NewAttributeValue() *AttributeValue {
	v := AttributeValue(make(map[string]any))
	return &v
}
