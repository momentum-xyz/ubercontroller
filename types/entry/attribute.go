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
	AttributeID
	*AttributePayload
}

type AttributePayload struct {
	Value   *AttributeValue   `db:"value"`
	Options *AttributeOptions `db:"options"`
}

type AttributeValue map[string]any

type PosBusAutoAttributeOption struct {
	SendTo PosBusDestinationType            `json:"send_to" db:"send_to" mapstructure:"send_to"`
	Scope  []PosBusAutoScopeAttributeOption `json:"scope" db:"scope" mapstructure:"scope"`
	Topic  string                           `json:"topic" db:"topic" mapstructure:"topic"`
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

type UnitySlotKind string

const (
	UnitySlotKindInvalid UnitySlotKind = ""
	UnitySlotKindTexture UnitySlotKind = "texture"
	UnitySlotKindString  UnitySlotKind = "string"
	UnitySlotKindNumber  UnitySlotKind = "number"
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
	SlotKind           UnitySlotKind    `json:"slot_kind" db:"slot_kind" mapstructure:"send_to"`
	SlotName           string           `json:"slot_name" db:"slot_name" mapstructure:"slot_name"`
	ValueField         string           `json:"value_field" db:"value_field" mapstructure:"value_field"`
	ContentType        UnityContentType `json:"content_type" db:"content_type" mapstructure:"content_type"`
	TextRenderTemplate AttributeValue   `json:"text_render_template" db:"text_render_template" mapstructure:"text_render_template"`
}
