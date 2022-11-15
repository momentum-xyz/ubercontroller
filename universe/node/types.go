package node

import "github.com/momentum-xyz/ubercontroller/universe"

type AttributeValueChangedMessage struct {
	Type universe.AttributeValueChangeType `json:"type"`
	Data AttributeValueChangedMessageData  `json:"data"`
}

type AttributeValueChangedMessageData struct {
	AttributeName string `json:"attribute_name"`
	SubName       string `json:"sub_name"`
	Value         any    `json:"value"`
}
