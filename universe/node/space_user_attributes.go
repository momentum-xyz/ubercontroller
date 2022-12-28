package node

import (
	"github.com/jackc/pgx/v4"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/pkg/errors"
)

type spaceUserAttributes struct {
	node *Node
}

func newSpaceUserAttributes(node *Node) *spaceUserAttributes {
	return &spaceUserAttributes{
		node: node,
	}
}

func (sua *spaceUserAttributes) GetPayload(spaceUserAttributeID entry.SpaceUserAttributeID) (*entry.AttributePayload, bool) {
	payload, err := sua.node.db.SpaceUserAttributesGetSpaceUserAttributePayloadByID(sua.node.ctx, spaceUserAttributeID)
	if err != nil {
		return nil, false
	}
	return payload, true
}

func (sua *spaceUserAttributes) GetValue(spaceUserAttributeID entry.SpaceUserAttributeID) (*entry.AttributeValue, bool) {
	value, err := sua.node.db.SpaceUserAttributesGetSpaceUserAttributeValueByID(sua.node.ctx, spaceUserAttributeID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (sua *spaceUserAttributes) GetOptions(spaceUserAttributeID entry.SpaceUserAttributeID) (*entry.AttributeOptions, bool) {
	options, err := sua.node.db.SpaceUserAttributesGetSpaceUserAttributeOptionsByID(sua.node.ctx, spaceUserAttributeID)
	if err != nil {
		return nil, false
	}
	return options, true
}

func (sua *spaceUserAttributes) GetEffectiveOptions(
	spaceUserAttributeID entry.SpaceUserAttributeID,
) (*entry.AttributeOptions, bool) {
	attributeType, ok := sua.node.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(spaceUserAttributeID.AttributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := sua.GetOptions(spaceUserAttributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		sua.node.log.Error(
			errors.WithMessagef(
				err, "Space user attributes: GetEffectiveOptions: failed to merge options: %+v", spaceUserAttributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (sua *spaceUserAttributes) Upsert(
	spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	payload, err := sua.node.db.SpaceUserAttributesUpsertSpaceUserAttribute(sua.node.ctx, spaceUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upsert space user attribute")
	}

	if sua.node.GetEnabled() {
		go func() {
			var value *entry.AttributeValue
			if payload != nil {
				value = payload.Value
			}
			sua.node.onSpaceUserAttributeChanged(universe.ChangedAttributeChangeType, spaceUserAttributeID, value, nil)
		}()
	}

	return payload, nil
}

func (sua *spaceUserAttributes) UpdateValue(
	spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	value, err := sua.node.db.SpaceUserAttributesUpdateSpaceUserAttributeValue(sua.node.ctx, spaceUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update space user attribute value")
	}

	if sua.node.GetEnabled() {
		go sua.node.onSpaceUserAttributeChanged(universe.ChangedAttributeChangeType, spaceUserAttributeID, value, nil)
	}

	return value, nil
}

func (sua *spaceUserAttributes) UpdateOptions(
	spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	options, err := sua.node.db.SpaceUserAttributesUpdateSpaceUserAttributeOptions(sua.node.ctx, spaceUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update space user attribute options")
	}

	if sua.node.GetEnabled() {
		go func() {
			value, ok := sua.GetValue(spaceUserAttributeID)
			if !ok {
				sua.node.log.Errorf(
					"Space user attributes: UpdateOptions: failed to get space user attribute value: %+v",
					spaceUserAttributeID,
				)
				return
			}
			sua.node.onSpaceUserAttributeChanged(universe.ChangedAttributeChangeType, spaceUserAttributeID, value, nil)
		}()
	}

	return options, nil
}

func (sua *spaceUserAttributes) Remove(spaceUserAttributeID entry.SpaceUserAttributeID, updateDB bool) (bool, error) {
	effectiveOptions, ok := sua.GetEffectiveOptions(spaceUserAttributeID)
	if !ok {
		return false, nil
	}

	if err := sua.node.db.SpaceUserAttributesRemoveSpaceUserAttributeByID(sua.node.ctx, spaceUserAttributeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.WithMessage(err, "failed to remove space user attribute")
	}

	if sua.node.GetEnabled() {
		go sua.node.onSpaceUserAttributeChanged(universe.RemovedAttributeChangeType, spaceUserAttributeID, nil, effectiveOptions)
	}

	return true, nil
}

func (sua *spaceUserAttributes) Len() int {
	count, err := sua.node.db.UserAttributesGetUserAttributesCount(sua.node.ctx)
	if err != nil {
		sua.node.log.Error(
			errors.WithMessage(err, "Space user attributes: Len: failed to get space user attributes count"),
		)
		return 0
	}
	return int(count)
}

func (n *Node) onSpaceUserAttributeChanged(
	changeType universe.AttributeChangeType, spaceUserAttributeID entry.SpaceUserAttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) {
	if effectiveOptions == nil {
		options, ok := n.GetSpaceUserAttributes().GetEffectiveOptions(spaceUserAttributeID)
		if !ok {
			n.log.Errorf(
				"Node: onSpaceUserAttributeChanged: failed to get space user attribute effective options: %+v",
				spaceUserAttributeID,
			)
			return
		}
		effectiveOptions = options
	}

	go func() {
		if err := n.posBusAutoOnSpaceUserAttributeChanged(changeType, spaceUserAttributeID, value, effectiveOptions); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err, "Node: onSpaceUserAttributeChanged: failed to handle pos bus auto: %+v", spaceUserAttributeID,
				),
			)
		}
	}()
}
