package posbus

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
)

type AttributeValueChangedMessage struct {
	Type universe.AttributeChangeType     `json:"type"`
	Data AttributeValueChangedMessageData `json:"data"`
}

type AttributeValueChangedMessageData struct {
	AttributeName string                `json:"attribute_name"`
	Value         *entry.AttributeValue `json:"value"`
}

func GetOptionAutoOption(options *entry.AttributeOptions) (*entry.PosBusAutoAttributeOption, error) {
	if options == nil {
		return nil, nil
	}

	autoOptionsValue, ok := (*options)["posbus_auto"]
	if !ok {
		return nil, nil
	}

	var autoOption entry.PosBusAutoAttributeOption
	if err := utils.MapDecode(autoOptionsValue, &autoOption); err != nil {
		return nil, errors.WithMessage(err, "failed to decode auto option")
	}

	return &autoOption, nil
}

func GetOptionAutoMessage(
	option *entry.PosBusAutoAttributeOption, changeType universe.AttributeChangeType,
	attributeID entry.AttributeID, value *entry.AttributeValue,
) (*websocket.PreparedMessage, error) {
	if option == nil {
		return nil, nil
	}

	data, err := json.Marshal(&AttributeValueChangedMessage{
		Type: changeType,
		Data: AttributeValueChangedMessageData{
			AttributeName: attributeID.Name,
			Value:         value,
		},
	})
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to marshal message payload")
	}

	topic := option.Topic
	if topic == "" {
		topic = attributeID.PluginID.String()
	}
	switch option.SendTo {
	case entry.ReactPosBusDestinationType:
		return posbus.NewRelayToReactMsg(topic, data).WebsocketMessage(), nil
	}

	return nil, errors.Errorf("send to type is not supported yet: %d", option.SendTo)
}
