package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"

	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
)

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

	topic := option.Topic
	if topic == "" {
		topic = attributeID.PluginID.String()
	}

	data := posbus.AttributeValueChanged{
		Topic:      topic,
		ChangeType: string(changeType),
		Data: posbus.AttributeValueChangedData{
			AttributeName: attributeID.Name,
			Value:         (*posbus.StringMapAny)(value),
		},
	}

	return posbus.WSMessage(&data), nil

	//return nil, errors.Errorf("send to type is not supported yet: %d", option.SendTo)
}
