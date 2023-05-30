package node

import (
	"context"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.NodeAttributes = (*nodeAttributes)(nil)

type nodeAttributes struct {
	node *Node
	data map[entry.AttributeID]*entry.AttributePayload
}

func newNodeAttributes(node *Node) *nodeAttributes {
	return &nodeAttributes{
		node: node,
		data: make(map[entry.AttributeID]*entry.AttributePayload),
	}
}

func (na *nodeAttributes) GetAll() map[entry.AttributeID]*entry.AttributePayload {
	na.node.Mu.RLock()
	defer na.node.Mu.RUnlock()

	attributes := make(map[entry.AttributeID]*entry.AttributePayload, len(na.data))
	for id, payload := range na.data {
		attributes[id] = payload
	}

	return attributes
}

func (na *nodeAttributes) Load() error {
	na.node.log.Infof("Loading node attributes: %s...", na.node.GetID())

	entries, err := na.node.db.GetNodeAttributesDB().GetNodeAttributes(na.node.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get node attributes")
	}

	for i := range entries {
		if _, err := na.Upsert(
			entries[i].AttributeID, modify.MergeWith(entries[i].AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert node attribute: %+v", entries[i].NodeAttributeID)
		}
	}

	na.node.log.Infof("Node attributes loaded: %s: %d", na.node.GetID(), na.Len())

	return nil
}

func (na *nodeAttributes) Save() error {
	na.node.log.Infof("Saving node attributes: %s...", na.node.GetID())

	na.node.Mu.RLock()
	defer na.node.Mu.RUnlock()

	attributes := make([]*entry.NodeAttribute, 0, len(na.data))
	for id, payload := range na.data {
		attributes = append(attributes, entry.NewNodeAttribute(entry.NewNodeAttributeID(id), payload))
	}

	if err := na.node.db.GetNodeAttributesDB().UpsertNodeAttributes(na.node.ctx, attributes); err != nil {
		return errors.WithMessage(err, "failed to upsert node attributes")
	}

	na.node.log.Infof("Node attributes saved: %s: %d", na.node.GetID(), len(na.data))

	return nil
}

func (na *nodeAttributes) GetPayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	na.node.Mu.RLock()
	defer na.node.Mu.RUnlock()

	if payload, ok := na.data[attributeID]; ok {
		return payload, true
	}
	return nil, false
}

func (na *nodeAttributes) GetValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	na.node.Mu.RLock()
	defer na.node.Mu.RUnlock()

	if payload, ok := na.data[attributeID]; ok && payload != nil {
		return payload.Value, true
	}
	return nil, false
}

func (na *nodeAttributes) GetOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	na.node.Mu.RLock()
	defer na.node.Mu.RUnlock()

	if payload, ok := na.data[attributeID]; ok && payload != nil {
		return payload.Options, true
	}
	return nil, false
}

func (na *nodeAttributes) GetEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attributeType, ok := na.node.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(attributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := na.GetOptions(attributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		na.node.log.Error(
			errors.WithMessagef(
				err,
				"Node attributes: GetEffectiveOptions: failed to merge options: %s: %+v",
				na.node.GetID(), attributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (na *nodeAttributes) Upsert(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	na.node.Mu.Lock()
	defer na.node.Mu.Unlock()

	payload, err := modifyFn(na.data[attributeID])
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify payload")
	}

	if updateDB {
		if err := na.node.db.GetNodeAttributesDB().UpsertNodeAttribute(
			na.node.ctx, entry.NewNodeAttribute(entry.NewNodeAttributeID(attributeID), payload),
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to upsert node attribute")
		}
	}

	na.data[attributeID] = payload

	return payload, nil
}

func (na *nodeAttributes) UpdateValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	na.node.Mu.Lock()
	defer na.node.Mu.Unlock()

	payload, ok := na.data[attributeID]
	if !ok {
		return nil, errors.Errorf("node attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify value")
	}

	if updateDB {
		if err := na.node.db.GetNodeAttributesDB().UpdateNodeAttributeValue(na.node.ctx, attributeID, value); err != nil {
			return nil, errors.WithMessagef(err, "failed to update node attribute value")
		}
	}

	payload.Value = value
	na.data[attributeID] = payload

	return value, nil
}

func (na *nodeAttributes) UpdateOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	na.node.Mu.Lock()
	defer na.node.Mu.Unlock()

	payload, ok := na.data[attributeID]
	if !ok {
		return nil, errors.Errorf("node attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify options")
	}

	if updateDB {
		if err := na.node.db.GetNodeAttributesDB().UpdateNodeAttributeOptions(na.node.ctx, attributeID, options); err != nil {
			return nil, errors.WithMessagef(err, "failed to update node attribute options")
		}
	}

	payload.Options = options
	na.data[attributeID] = payload

	return options, nil
}

func (na *nodeAttributes) Remove(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	na.node.Mu.Lock()
	defer na.node.Mu.Unlock()

	if _, ok := na.data[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := na.node.db.GetNodeAttributesDB().RemoveNodeAttributeByAttributeID(na.node.ctx, attributeID); err != nil {
			return false, errors.WithMessagef(err, "failed to remove node attribute")
		}
	}

	delete(na.data, attributeID)

	return true, nil
}

func (na *nodeAttributes) Len() int {
	na.node.Mu.RLock()
	defer na.node.Mu.RUnlock()

	return len(na.data)
}

// AttributePermissionsAuthorizer
func (na *nodeAttributes) GetUserRoles(
	ctx context.Context,
	attrType entry.AttributeType,
	targetID entry.AttributeID,
	userID umid.UMID,
) ([]entry.PermissionsRoleType, error) {
	var roles []entry.PermissionsRoleType
	// owner is always considered an admin, TODO: add this to check function
	if na.node.GetOwnerID() == userID {
		roles = append(roles, entry.PermissionAdmin)
	} else { // we have to lookup through the db user tree
		userObjectID := entry.NewUserObjectID(userID, na.node.GetID())
		isAdmin, err := na.node.db.GetUserObjectsDB().CheckIsIndirectAdminByID(ctx, userObjectID)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to check admin status")
		}
		if isAdmin {
			roles = append(roles, entry.PermissionAdmin)
		}
	}
	return roles, nil
	na.node.GetOwnerID()
	return nil, nil
}
