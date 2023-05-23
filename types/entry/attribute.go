package entry

import (
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type PosBusAutoScopeAttributeOption string

const (
	InvalidPosBusAutoScopeAttributeOption PosBusAutoScopeAttributeOption = ""
	NodePosBusAutoScopeAttributeOption    PosBusAutoScopeAttributeOption = "node"
	WorldPosBusAutoScopeAttributeOption   PosBusAutoScopeAttributeOption = "world"
	ObjectPosBusAutoScopeAttributeOption  PosBusAutoScopeAttributeOption = "object"
	UserPosBusAutoScopeAttributeOption    PosBusAutoScopeAttributeOption = "user"
)

type SlotType string

const (
	SlotTypeInvalid SlotType = ""
	SlotTypeTexture SlotType = "texture"
	SlotTypeString  SlotType = "string"
	SlotTypeNumber  SlotType = "number"
)

type SlotContentType string

const (
	SlotContentTypeInvalid SlotContentType = ""
	SlotContentTypeString  SlotContentType = "string"
	SlotContentTypeNumber  SlotContentType = "number"
	SlotContentTypeImage   SlotContentType = "image"
	SlotContentTypeText    SlotContentType = "text"
	SlotContentTypeVideo   SlotContentType = "video"
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
	Scope []PosBusAutoScopeAttributeOption `db:"scope" json:"scope"`
	Topic string                           `db:"topic" json:"topic"`
}

type RenderAutoAttributeOption struct {
	SlotType           SlotType        `db:"slot_type" json:"slot_type"`
	SlotName           string          `db:"slot_name" json:"slot_name"`
	ValueField         string          `db:"value_field" json:"value_field"`
	ContentType        SlotContentType `db:"content_type" json:"content_type"`
	TextRenderTemplate string          `db:"text_render_template" json:"text_render_template"`
}

type PermissionsRoleType string

const (
	PermissionAny        PermissionsRoleType = "any"
	PermissionUser       PermissionsRoleType = "user"
	PermissionUserOwner  PermissionsRoleType = "user_owner"
	PermissionAdmin      PermissionsRoleType = "admin"
	PermissionTargetUser PermissionsRoleType = "target_user"
)

type PermissionsAttributeOption struct {
	Read  string `json:"read" mapstructure:"read"`
	Write string `json:"write" mapstructure:"write"`
	// TODO: replace string, impl decoder for e.g 'admin+user_owner'
}

func NewAttribute(attributeID AttributeID, payload *AttributePayload) *Attribute {
	return &Attribute{
		AttributeID:      attributeID,
		AttributePayload: payload,
	}
}

func NewAttributeID(pluginID umid.UMID, name string) AttributeID {
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
