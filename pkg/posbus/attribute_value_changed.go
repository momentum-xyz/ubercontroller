package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

// TODO: does musgo support type aliases/const like this?
type AttributeChangeType string

const (
	InvalidAttributeChangeType AttributeChangeType = ""
	ChangedAttributeChangeType AttributeChangeType = "attribute_changed"
	RemovedAttributeChangeType AttributeChangeType = "attribute_removed"
)

type AttributeValueChanged struct {
	// The plugin that owns the attribute
	PluginID umid.UMID `json:"plugin_id"`
	// Name of attribute (scoped to plugin)
	AttributeName string `json:"attribute_name"`
	// Indicate what has changed (removed or value changed)
	ChangeType string `json:"change_type"`
	// The new value, in case of change/new.
	Value *StringAnyMap `json:"value"`
	// ID of the related object/user
	TargetID umid.UMID `json:"target_id"`
}

func (r *AttributeValueChanged) GetType() MsgType {
	return 0x10DACDB7
}

func init() {
	registerMessage(AttributeValueChanged{})
}
