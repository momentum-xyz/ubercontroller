package posbus

type AttributeValueChanged struct {
	Topic      string                    `json:"topic"`
	ChangeType string                    `json:"change_type"`
	Data       AttributeValueChangedData `json:"data"`
}

type AttributeValueChangedData struct {
	AttributeName string        `json:"attribute_name"`
	Value         *StringMapAny `json:"value"`
}

func (r *AttributeValueChanged) GetType() MsgType {
	return 0x10DACDB7
}

func init() {
	addExtraType(AttributeValueChangedData{})
	registerMessage(AttributeValueChanged{})
}
