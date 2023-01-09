package node

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

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

func (na *nodeAttributes) GetPayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	na.node.mu.RLock()
	defer na.node.mu.RUnlock()

	if payload, ok := na.data[attributeID]; ok {
		return payload, true
	}
	return nil, false
}

func (na *nodeAttributes) GetValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	na.node.mu.RLock()
	defer na.node.mu.RUnlock()

	if payload, ok := na.data[attributeID]; ok && payload != nil {
		return payload.Value, true
	}
	return nil, false
}

func (na *nodeAttributes) GetOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	na.node.mu.RLock()
	defer na.node.mu.RUnlock()

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
	na.node.mu.Lock()
	defer na.node.mu.Unlock()

	payload, err := modifyFn(na.data[attributeID])
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify payload")
	}

	if updateDB {
		if err := na.node.db.NodeAttributesUpsertNodeAttribute(
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
	na.node.mu.Lock()
	defer na.node.mu.Unlock()

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
		if err := na.node.db.NodeAttributesUpdateNodeAttributeValue(na.node.ctx, attributeID, value); err != nil {
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
	na.node.mu.Lock()
	defer na.node.mu.Unlock()

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
		if err := na.node.db.NodeAttributesUpdateNodeAttributeOptions(na.node.ctx, attributeID, options); err != nil {
			return nil, errors.WithMessagef(err, "failed to update node attribute options")
		}
	}

	payload.Options = options
	na.data[attributeID] = payload

	return options, nil
}

func (na *nodeAttributes) Remove(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	na.node.mu.Lock()
	defer na.node.mu.Unlock()

	if _, ok := na.data[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := na.node.db.NodeAttributesRemoveNodeAttributeByAttributeID(na.node.ctx, attributeID); err != nil {
			return false, errors.WithMessagef(err, "failed to remove space attribute")
		}
	}

	delete(na.data, attributeID)

	return true, nil
}

func (na *nodeAttributes) Len() int {
	na.node.mu.RLock()
	defer na.node.mu.RUnlock()

	return len(na.data)
}

func (n *Node) loadNodeAttributes() error {
	n.log.Infof("Loading node attributes: %s...", n.GetID())

	entries, err := n.db.NodeAttributesGetNodeAttributes(n.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get node attributes")
	}

	attributes := n.GetNodeAttributes()
	for i := range entries {
		if _, err := attributes.Upsert(
			entries[i].AttributeID, modify.MergeWith(entries[i].AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert node attribute: %+v", entries[i].NodeAttributeID)
		}
	}

	n.log.Infof("Node attributes loaded: %s: %d", n.GetID(), attributes.Len())

	return nil
}
