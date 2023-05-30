package node

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.ObjectUserAttributes = (*objectUserAttributes)(nil)

type objectUserAttributes struct {
	node *Node
}

func newObjectUserAttributes(node *Node) *objectUserAttributes {
	return &objectUserAttributes{
		node: node,
	}
}

func (oua *objectUserAttributes) GetPayload(objectUserAttributeID entry.ObjectUserAttributeID) (*entry.AttributePayload, bool) {
	payload, err := oua.node.db.GetObjectUserAttributesDB().GetObjectUserAttributePayloadByID(
		oua.node.ctx, objectUserAttributeID,
	)
	if err != nil {
		return nil, false
	}
	return payload, true
}

func (oua *objectUserAttributes) GetValue(objectUserAttributeID entry.ObjectUserAttributeID) (*entry.AttributeValue, bool) {
	value, err := oua.node.db.GetObjectUserAttributesDB().GetObjectUserAttributeValueByID(oua.node.ctx, objectUserAttributeID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (oua *objectUserAttributes) GetOptions(objectUserAttributeID entry.ObjectUserAttributeID) (*entry.AttributeOptions, bool) {
	options, err := oua.node.db.GetObjectUserAttributesDB().GetObjectUserAttributeOptionsByID(oua.node.ctx, objectUserAttributeID)
	if err != nil {
		return nil, false
	}
	return options, true
}

func (oua *objectUserAttributes) GetEffectiveOptions(
	objectUserAttributeID entry.ObjectUserAttributeID,
) (*entry.AttributeOptions, bool) {
	attributeType, ok := oua.node.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(objectUserAttributeID.AttributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := oua.GetOptions(objectUserAttributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		oua.node.log.Error(
			errors.WithMessagef(
				err, "Object user attributes: GetEffectiveOptions: failed to merge options: %+v", objectUserAttributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (oua *objectUserAttributes) Upsert(
	objectUserAttributeID entry.ObjectUserAttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	payload, err := oua.node.db.GetObjectUserAttributesDB().UpsertObjectUserAttribute(
		oua.node.ctx, objectUserAttributeID, modifyFn,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upsert object user attribute")
	}

	if oua.node.GetEnabled() {
		go func() {
			var value *entry.AttributeValue
			if payload != nil {
				value = payload.Value
			}
			oua.node.onObjectUserAttributeChanged(
				universe.ChangedAttributeChangeType, objectUserAttributeID, value, nil,
			)
		}()
	}

	return payload, nil
}

func (oua *objectUserAttributes) UpdateValue(
	objectUserAttributeID entry.ObjectUserAttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	value, err := oua.node.db.GetObjectUserAttributesDB().UpdateObjectUserAttributeValue(
		oua.node.ctx, objectUserAttributeID, modifyFn,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update object user attribute value")
	}

	if oua.node.GetEnabled() {
		go oua.node.onObjectUserAttributeChanged(
			universe.ChangedAttributeChangeType, objectUserAttributeID, value, nil,
		)
	}

	return value, nil
}

func (oua *objectUserAttributes) UpdateOptions(
	objectUserAttributeID entry.ObjectUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	options, err := oua.node.db.GetObjectUserAttributesDB().UpdateObjectUserAttributeOptions(
		oua.node.ctx, objectUserAttributeID, modifyFn,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update object user attribute options")
	}

	if oua.node.GetEnabled() {
		go func() {
			value, ok := oua.GetValue(objectUserAttributeID)
			if !ok {
				oua.node.log.Errorf(
					"Object user attributes: UpdateOptions: failed to get object user attribute value: %+v",
					objectUserAttributeID,
				)
				return
			}
			oua.node.onObjectUserAttributeChanged(
				universe.ChangedAttributeChangeType, objectUserAttributeID, value, nil,
			)
		}()
	}

	return options, nil
}

func (oua *objectUserAttributes) Remove(objectUserAttributeID entry.ObjectUserAttributeID, updateDB bool) (bool, error) {
	effectiveOptions, ok := oua.GetEffectiveOptions(objectUserAttributeID)
	if !ok {
		return false, nil
	}

	if err := oua.node.db.GetObjectUserAttributesDB().RemoveObjectUserAttributeByID(
		oua.node.ctx, objectUserAttributeID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.WithMessage(err, "failed to remove object user attribute")
	}

	if oua.node.GetEnabled() {
		go oua.node.onObjectUserAttributeChanged(
			universe.RemovedAttributeChangeType, objectUserAttributeID, nil, effectiveOptions,
		)
	}

	return true, nil
}

func (oua *objectUserAttributes) Len() int {
	count, err := oua.node.db.GetUserAttributesDB().GetUserAttributesCount(oua.node.ctx)
	if err != nil {
		oua.node.log.Error(
			errors.WithMessage(err, "Object user attributes: Len: failed to get object user attributes count"),
		)
		return 0
	}
	return int(count)
}

func (n *Node) onObjectUserAttributeChanged(
	changeType universe.AttributeChangeType, objectUserAttributeID entry.ObjectUserAttributeID,
	value *entry.AttributeValue, effectiveOptions *entry.AttributeOptions,
) {
	if effectiveOptions == nil {
		options, ok := n.GetObjectUserAttributes().GetEffectiveOptions(objectUserAttributeID)
		if !ok {
			n.log.Errorf(
				"Node: onObjectUserAttributeChanged: failed to get object user attribute effective options: %+v",
				objectUserAttributeID,
			)
			return
		}
		effectiveOptions = options
	}

	go func() {
		if err := n.posBusAutoOnObjectUserAttributeChanged(
			changeType, objectUserAttributeID, value, effectiveOptions,
		); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err, "Node: onObjectUserAttributeChanged: failed to handle posbus auto: %+v", objectUserAttributeID,
				),
			)
		}
	}()
}

// AttributePermissionsAuthorizer
func (oua *objectUserAttributes) GetUserRoles(
	ctx context.Context,
	attrType entry.AttributeType,
	targetID entry.ObjectUserAttributeID,
	userID umid.UMID,
) ([]entry.PermissionsRoleType, error) {
	var roles []entry.PermissionsRoleType
	object, ok := oua.node.GetObjectFromAllObjects(targetID.ObjectID)
	if !ok {
		return nil, errors.New("ObjectUserAttribute roles: object not found")
	}
	objectRoles, err := object.GetObjectAttributes().GetUserRoles(
		ctx, attrType, entry.NewAttributeID(attrType.PluginID, attrType.Name), userID)
	if err != nil {
		return nil, fmt.Errorf("ObjectUserAttribute roles: %w", err)
	}
	roles = append(roles, objectRoles...)

	if targetID.UserID == userID {
		roles = append(roles, entry.PermissionUserOwner)
	} // TODO: user members, walk the user tree for indirect ownership...
	return roles, nil
}
