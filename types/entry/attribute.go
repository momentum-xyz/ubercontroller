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
	AttributeID       `mapstructure:",squash"`
	*AttributePayload `mapstructure:",squash"`
}

type AttributePayload struct {
	Value   *AttributeValue   `db:"value" mapstructure:"value"`
	Options *AttributeOptions `db:"options" mapstructure:"options"`
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
