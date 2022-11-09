package node

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-multierror"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func (n *Node) OnSpaceAttributeValueChanged(
	changeType universe.AttributeValueChangeType, spaceAttributeID entry.SpaceAttributeID, subAttributeKey string, newValue any,
) error {
	space, ok := n.GetSpaceFromAllSpaces(spaceAttributeID.SpaceID)
	if !ok {
		return errors.Errorf("space not found: %s", spaceAttributeID.SpaceID)
	}

	effectiveOptions, ok := space.GetSpaceAttributeEffectiveOptions(spaceAttributeID.AttributeID)
	if !ok {
		return nil
	}

	autoOption, err := n.getPosBusAutoAttributeOption(effectiveOptions)
	if err != nil {
		return errors.WithMessage(err, "failed to get option")
	}
	msg, err := n.getAttributeValueChangedMsg(autoOption, changeType, spaceAttributeID.AttributeID, subAttributeKey, newValue)
	if err != nil {
		return errors.WithMessage(err, "failed to get message")
	}
	if msg == nil {
		return nil
	}

	var errs *multierror.Error
	for i := range autoOption.Scope {
		switch autoOption.Scope[i] {
		case entry.SpacePosBusAutoScopeAttributeOption:
			if err := space.Broadcast(msg, false); err != nil {
				errs = multierror.Append(
					errs, errors.WithMessagef(
						err, "failed to broadcast message: %s", autoOption.Scope[i],
					),
				)
			}
		default:
			errs = multierror.Append(
				errs, errors.Errorf(
					"scope type in not supported yet: %s", autoOption.Scope[i],
				),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (n *Node) getPosBusAutoAttributeOption(options *entry.AttributeOptions) (*entry.PosBusAutoAttributeOption, error) {
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

func (n *Node) getAttributeValueChangedMsg(
	option *entry.PosBusAutoAttributeOption, changeType universe.AttributeValueChangeType,
	attributeID entry.AttributeID, subAttributeKey string, newValue any,
) (*websocket.PreparedMessage, error) {
	if option == nil {
		return nil, nil
	}

	data, err := json.Marshal(&AttributeValueChangedMessage{
		Type: changeType,
		Data: AttributeValueChangedMessageData{
			AttributeName: attributeID.Name,
			SubName:       subAttributeKey,
			Value:         newValue,
		},
	})
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to marshal message payload")
	}

	switch option.SendTo {
	case entry.ReactPosBusDestinationType:
		return posbus.NewRelayToReactMsg(option.Topic, data).WebsocketMessage(), nil
	}

	return nil, errors.Errorf("send to type is not supported yet: %d", option.SendTo)
}
