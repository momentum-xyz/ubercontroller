package node

import (
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type userAttributes struct {
	node *Node
}

func newUserAttributes(node *Node) *userAttributes {
	return &userAttributes{
		node: node,
	}
}

func (ua *userAttributes) GetPayload(userAttributeID entry.UserAttributeID) (*entry.AttributePayload, bool) {
	payload, err := ua.node.db.UserAttributesGetUserAttributePayloadByID(ua.node.ctx, userAttributeID)
	if err != nil {
		return nil, false
	}
	return payload, true
}

func (ua *userAttributes) GetValue(userAttributeID entry.UserAttributeID) (*entry.AttributeValue, bool) {
	value, err := ua.node.db.UserAttributesGetUserAttributeValueByID(ua.node.ctx, userAttributeID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (ua *userAttributes) GetOptions(userAttributeID entry.UserAttributeID) (*entry.AttributeOptions, bool) {
	options, err := ua.node.db.UserAttributesGetUserAttributeOptionsByID(ua.node.ctx, userAttributeID)
	if err != nil {
		return nil, false
	}
	return options, true
}

func (ua *userAttributes) GetEffectiveOptions(userAttributeID entry.UserAttributeID) (*entry.AttributeOptions, bool) {
	attributeType, ok := ua.node.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(userAttributeID.AttributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := ua.GetOptions(userAttributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		ua.node.log.Error(
			errors.WithMessagef(
				err, "User attributes: GetEffectiveOptions: failed to merge options: %+v", userAttributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (ua *userAttributes) Upsert(
	userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	payload, err := ua.node.db.UserAttributesUpsertUserAttribute(ua.node.ctx, userAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upsert user attribute")
	}

	if ua.node.GetEnabled() {
		go func() {
			var value *entry.AttributeValue
			if payload != nil {
				value = payload.Value
			}
			ua.node.onUserAttributeChanged(universe.ChangedAttributeChangeType, userAttributeID, value, nil)
		}()
	}

	return payload, nil
}

func (ua *userAttributes) UpdateValue(
	userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	value, err := ua.node.db.UserAttributesUpdateUserAttributeValue(ua.node.ctx, userAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user attribute value")
	}

	if ua.node.GetEnabled() {
		go ua.node.onUserAttributeChanged(universe.ChangedAttributeChangeType, userAttributeID, value, nil)
	}

	return value, nil
}

func (ua *userAttributes) UpdateOptions(
	userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	options, err := ua.node.db.UserAttributesUpdateUserAttributeOptions(ua.node.ctx, userAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user attribute options")
	}

	if ua.node.GetEnabled() {
		go func() {
			value, ok := ua.GetValue(userAttributeID)
			if !ok {
				ua.node.log.Errorf(
					"User attributes: UpdateOptions: failed to get user attribute value: %+v", userAttributeID,
				)
				return
			}
			ua.node.onUserAttributeChanged(universe.ChangedAttributeChangeType, userAttributeID, value, nil)
		}()
	}

	return options, nil
}

func (ua *userAttributes) Remove(userAttributeID entry.UserAttributeID, updateDB bool) (bool, error) {
	effectiveOptions, ok := ua.GetEffectiveOptions(userAttributeID)
	if !ok {
		return false, nil
	}

	if err := ua.node.db.UserAttributesRemoveUserAttributeByID(ua.node.ctx, userAttributeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.WithMessage(err, "failed to remove user attribute by id")
	}

	if ua.node.GetEnabled() {
		go ua.node.onUserAttributeChanged(universe.RemovedAttributeChangeType, userAttributeID, nil, effectiveOptions)
	}

	return true, nil
}

func (ua *userAttributes) Len() int {
	count, err := ua.node.db.UserAttributesGetUserAttributesCount(ua.node.ctx)
	if err != nil {
		ua.node.log.Error(errors.WithMessage(err, "User attributes: Len: failed to get user attributes count"))
		return 0
	}
	return int(count)
}

func (n *Node) onUserAttributeChanged(
	changeType universe.AttributeChangeType, userAttributeID entry.UserAttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) {
	if effectiveOptions == nil {
		options, ok := n.GetUserAttributes().GetEffectiveOptions(userAttributeID)
		if !ok {
			n.log.Errorf(
				"Node: onUserAttributeChanged: failed to get user attribute effective options: %+v",
				userAttributeID,
			)
			return
		}
		effectiveOptions = options
	}

	go func() {
		if err := n.posBusAutoOnUserAttributeChanged(changeType, userAttributeID, value, effectiveOptions); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err, "Node: onUserAttributeChanged: failed to handle posbus auto: %+v", userAttributeID,
				),
			)
		}
	}()
}
