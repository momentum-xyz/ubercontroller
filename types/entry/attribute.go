package entry

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type PosBusDestinationType byte

const (
	InvalidPosBusDestinationType    PosBusDestinationType = 0b00
	ReactPosBusDestinationType      PosBusDestinationType = 0b01
	UnityPosBusDestinationType      PosBusDestinationType = 0b10
	ReactUnityPosBusDestinationType PosBusDestinationType = 0b11
)

type PosBusAutoScopeAttributeOption string

const (
	InvalidPosBusAutoScopeAttributeOption PosBusAutoScopeAttributeOption = ""
	NodePosBusAutoScopeAttributeOption    PosBusAutoScopeAttributeOption = "node"
	WorldPosBusAutoScopeAttributeOption   PosBusAutoScopeAttributeOption = "world"
	SpacePosBusAutoScopeAttributeOption   PosBusAutoScopeAttributeOption = "space"
	UserPosBusAutoScopeAttributeOption    PosBusAutoScopeAttributeOption = "user"
)

type AttributeID AttributeTypeID

type Attribute struct {
	AttributeID       `json:",squash"`
	*AttributePayload `json:",squash"` // won't work, see "node.SpaceTemplateFromMap()"
}

type AttributePayload struct {
	Value   *AttributeValue   `db:"value" json:"value"`
	Options *AttributeOptions `db:"options" json:"options"`
}

type AttributeValue map[string]any

type PosBusAutoAttributeOption struct {
	SendTo PosBusDestinationType            `json:"send_to" db:"send_to"`
	Scope  []PosBusAutoScopeAttributeOption `json:"scope" db:"scope"`
	Topic  string                           `json:"topic" db:"topic"`
}

func NewAttribute(attributeID AttributeID, payload *AttributePayload) *Attribute {
	return &Attribute{
		AttributeID:      attributeID,
		AttributePayload: payload,
	}
}

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

type UnitySlotType string

const (
	UnitySlotTypeInvalid UnitySlotType = ""
	UnitySlotTypeTexture UnitySlotType = "texture"
	UnitySlotTypeString  UnitySlotType = "string"
	UnitySlotTypeNumber  UnitySlotType = "number"
)

type UnityContentType string

const (
	UnityContentTypeInvalid UnityContentType = ""
	UnityContentTypeString  UnityContentType = "string"
	UnityContentTypeNumber  UnityContentType = "number"
	UnityContentTypeImage   UnityContentType = "image"
	UnityContentTypeText    UnityContentType = "text"
	UnityContentTypeVideo   UnityContentType = "video"
)

type UnityAutoAttributeOption struct {
	SlotType           UnitySlotType    `json:"slot_type" db:"slot_type" mapstructure:"slot_type"`
	SlotName           string           `json:"slot_name" db:"slot_name" mapstructure:"slot_name"`
	ValueField         string           `json:"value_field" db:"value_field" mapstructure:"value_field"`
	ContentType        UnityContentType `json:"content_type" db:"content_type" mapstructure:"content_type"`
	TextRenderTemplate string           `json:"text_render_template" db:"text_render_template" mapstructure:"text_render_template"`
}
