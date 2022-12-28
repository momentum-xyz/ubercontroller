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

type AttributeID AttributeTypeID

type Attribute struct {
	AttributeID
	*AttributePayload
}

type AttributePayload struct {
	Value   *AttributeValue   `db:"value" json:"value"`
	Options *AttributeOptions `db:"options" json:"options"`
}

type AttributeValue map[string]any

type AttributeOptions map[string]any

type PosBusAutoAttributeOption struct {
	SendTo PosBusDestinationType            `db:"send_to" json:"send_to"`
	Scope  []PosBusAutoScopeAttributeOption `db:"scope" json:"scope"`
	Topic  string                           `db:"topic" json:"topic"`
}

type UnityAutoAttributeOption struct {
	SlotType           UnitySlotType    `db:"slot_type" json:"slot_type"`
	SlotName           string           `db:"slot_name" json:"slot_name"`
	ValueField         string           `db:"value_field" json:"value_field"`
	ContentType        UnityContentType `db:"content_type" json:"content_type"`
	TextRenderTemplate string           `db:"text_render_template" json:"text_render_template"`
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
	return utils.GetPTR(make(AttributeValue))
}

func NewAttributeOptions() *AttributeOptions {
	return utils.GetPTR(make(AttributeOptions))
}
