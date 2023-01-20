package node

import (
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.UserUserAttributes = (*userUserAttributes)(nil)

type userUserAttributes struct {
	node *Node
}

func newUserUserAttributes(node *Node) *userUserAttributes {
	return &userUserAttributes{
		node: node,
	}
}

func (uua *userUserAttributes) GetPayload(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributePayload, bool) {
	payload, err := uua.node.db.GetUserUserAttributesDB().GetUserUserAttributePayloadByID(uua.node.ctx, userUserAttributeID)
	if err != nil {
		return nil, false
	}
	return payload, true
}

func (uua *userUserAttributes) GetValue(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributeValue, bool) {
	value, err := uua.node.db.GetUserUserAttributesDB().GetUserUserAttributeValueByID(uua.node.ctx, userUserAttributeID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (uua *userUserAttributes) GetOptions(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributeOptions, bool) {
	options, err := uua.node.db.GetUserUserAttributesDB().GetUserUserAttributeOptionsByID(uua.node.ctx, userUserAttributeID)
	if err != nil {
		return nil, false
	}
	return options, true
}

func (uua *userUserAttributes) GetEffectiveOptions(
	userUserAttributeID entry.UserUserAttributeID,
) (*entry.AttributeOptions, bool) {
	attributeType, ok := uua.node.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(userUserAttributeID.AttributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := uua.GetOptions(userUserAttributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		uua.node.log.Error(
			errors.WithMessagef(
				err, "User user attributes: GetEffectiveOptions: failed to merge options: %+v", userUserAttributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (uua *userUserAttributes) Upsert(
	userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	payload, err := uua.node.db.GetUserUserAttributesDB().UpsertUserUserAttribute(uua.node.ctx, userUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upsert user user attribute")
	}

	if uua.node.GetEnabled() {
		go func() {
			var value *entry.AttributeValue
			if payload != nil {
				value = payload.Value
			}
			uua.node.onUserUserAttributeChanged(universe.ChangedAttributeChangeType, userUserAttributeID, value, nil)
		}()
	}

	return payload, nil
}

func (uua *userUserAttributes) UpdateValue(
	userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	value, err := uua.node.db.GetUserUserAttributesDB().UpdateUserUserAttributeValue(uua.node.ctx, userUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user user attribute value")
	}

	if uua.node.GetEnabled() {
		go uua.node.onUserUserAttributeChanged(universe.ChangedAttributeChangeType, userUserAttributeID, value, nil)
	}

	return value, nil
}

func (uua *userUserAttributes) UpdateOptions(
	userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	options, err := uua.node.db.GetUserUserAttributesDB().UpdateUserUserAttributeOptions(uua.node.ctx, userUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user user attribute options")
	}

	if uua.node.GetEnabled() {
		go func() {
			value, ok := uua.GetValue(userUserAttributeID)
			if !ok {
				uua.node.log.Errorf(
					"User user attributes: UpdateOptions: failed to get user use attribute value: %+v",
					userUserAttributeID,
				)
				return
			}
			uua.node.onUserUserAttributeChanged(universe.ChangedAttributeChangeType, userUserAttributeID, value, nil)
		}()
	}

	return options, nil
}

func (uua *userUserAttributes) Remove(userUserAttributeID entry.UserUserAttributeID, updateDB bool) (bool, error) {
	effectiveOptions, ok := uua.GetEffectiveOptions(userUserAttributeID)
	if !ok {
		return false, nil
	}

	if err := uua.node.db.GetUserUserAttributesDB().RemoveUserUserAttributeByID(uua.node.ctx, userUserAttributeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.WithMessage(err, "failed to remove user user attribute")
	}

	if uua.node.GetEnabled() {
		go uua.node.onUserUserAttributeChanged(universe.RemovedAttributeChangeType, userUserAttributeID, nil, effectiveOptions)
	}

	return true, nil
}

func (uua *userUserAttributes) Len() int {
	count, err := uua.node.db.GetUserUserAttributesDB().GetUserUserAttributesCount(uua.node.ctx)
	if err != nil {
		uua.node.log.Error(
			errors.WithMessage(err, "User user attributes: Len: failed to get user user attributes count"),
		)
		return 0
	}
	return int(count)
}

func (n *Node) onUserUserAttributeChanged(
	changeType universe.AttributeChangeType, userUserAttributeID entry.UserUserAttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) {
	if effectiveOptions == nil {
		options, ok := n.GetUserUserAttributes().GetEffectiveOptions(userUserAttributeID)
		if !ok {
			n.log.Errorf(
				"Node: onUserUserAttributeChanged: failed to get user user attribute effective options: %+v",
				userUserAttributeID,
			)
			return
		}
		effectiveOptions = options
	}

	go func() {
		if err := n.posBusAutoOnUserUserAttributeChanged(changeType, userUserAttributeID, value, effectiveOptions); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err, "Node: onUserUserAttributeChanged: failed to handle posbus auto: %+v", userUserAttributeID,
				),
			)
		}
	}()
}
