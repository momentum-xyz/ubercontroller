package node

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) UpsertNodeAttribute(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.NodeAttribute, error) {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[attributeID]
	if !ok {
		payload = nil
	}

	payload, err := modifyFn(payload)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify attribute payload")
	}

	nodeAttribute := entry.NewNodeAttribute(entry.NewNodeAttributeID(attributeID), payload)
	if updateDB {
		if err := n.db.NodeAttributesUpsertNodeAttribute(n.ctx, nodeAttribute); err != nil {
			return nil, errors.WithMessage(err, "failed to upsert node attribute")
		}
	}

	n.nodeAttributes.Data[attributeID] = payload

	return nodeAttribute, nil

}

func (n *Node) GetNodeAttributeValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	payload, ok := n.GetNodeAttributePayload(attributeID)
	if !ok {
		return nil, false
	}
	if payload == nil {
		return nil, true
	}
	return payload.Value, true
}

func (n *Node) GetNodeAttributeOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	payload, ok := n.GetNodeAttributePayload(attributeID)
	if !ok {
		return nil, false
	}
	if payload == nil {
		return nil, true
	}
	return payload.Options, true
}

func (n *Node) GetNodeAttributeEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attributeType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(attributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := n.GetNodeAttributeOptions(attributeID)
	if !ok {
		attributeOptions = nil
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		n.log.Error(
			errors.WithMessagef(
				err, "Node: GetNodeAttributeEffectiveOptions: failed to merge options: %+v", attributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (n *Node) GetNodeAttributePayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	return n.nodeAttributes.Load(attributeID)
}

func (n *Node) GetNodeAttributesValue() map[entry.NodeAttributeID]*entry.AttributeValue {
	n.nodeAttributes.Mu.RLock()
	defer n.nodeAttributes.Mu.RUnlock()

	values := make(map[entry.NodeAttributeID]*entry.AttributeValue, len(n.nodeAttributes.Data))

	for attributeID, payload := range n.nodeAttributes.Data {
		nodeAttributeID := entry.NewNodeAttributeID(attributeID)
		if payload == nil {
			values[nodeAttributeID] = nil
			continue
		}
		values[nodeAttributeID] = payload.Value
	}

	return values
}

func (n *Node) GetNodeAttributesOptions() map[entry.NodeAttributeID]*entry.AttributeOptions {
	n.nodeAttributes.Mu.RLock()
	defer n.nodeAttributes.Mu.RUnlock()

	options := make(map[entry.NodeAttributeID]*entry.AttributeOptions, len(n.nodeAttributes.Data))

	for attributeID, payload := range n.nodeAttributes.Data {
		nodeAttributeID := entry.NewNodeAttributeID(attributeID)
		if payload == nil {
			options[nodeAttributeID] = nil
			continue
		}
		options[nodeAttributeID] = payload.Options
	}

	return options
}

func (n *Node) GetNodeAttributesPayload() map[entry.NodeAttributeID]*entry.AttributePayload {
	n.nodeAttributes.Mu.RLock()
	defer n.nodeAttributes.Mu.RUnlock()

	attributes := make(map[entry.NodeAttributeID]*entry.AttributePayload, len(n.nodeAttributes.Data))

	for attributeID, payload := range n.nodeAttributes.Data {
		attributes[entry.NewNodeAttributeID(attributeID)] = payload
	}

	return attributes
}

func (n *Node) UpdateNodeAttributeValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[attributeID]
	if !ok {
		return nil, errors.Errorf("not attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify value")
	}

	if updateDB {
		if err := n.db.NodeAttributesUpdateNodeAttributeValue(
			n.ctx, attributeID, value,
		); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Value = value
	n.nodeAttributes.Data[attributeID] = payload

	return value, nil
}

func (n *Node) UpdateNodeAttributeOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[attributeID]
	if !ok {
		return nil, errors.Errorf("node attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := n.db.NodeAttributesUpdateNodeAttributeOptions(
			n.ctx, attributeID, options,
		); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Options = options
	n.nodeAttributes.Data[attributeID] = payload

	return options, nil
}

func (n *Node) RemoveNodeAttribute(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	if _, ok := n.nodeAttributes.Data[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := n.db.NodeAttributesRemoveNodeAttributeByAttributeID(n.ctx, attributeID); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(n.nodeAttributes.Data, attributeID)

	return true, nil
}

func (n *Node) RemoveNodeAttributes(attributeIDs []entry.AttributeID, updateDB bool) (bool, error) {
	res := true
	var errs *multierror.Error
	for i := range attributeIDs {
		removed, err := n.RemoveNodeAttribute(attributeIDs[i], updateDB)
		if err != nil {
			errs = multierror.Append(errs,
				errors.WithMessagef(err, "failed to remove node attribute: %+v", attributeIDs[i]),
			)
		}
		if !removed {
			res = false
		}
	}
	return res, errs.ErrorOrNil()
}

func (n *Node) loadNodeAttributes() error {
	entries, err := n.db.NodeAttributesGetNodeAttributes(n.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get node attributes")
	}

	for _, instance := range entries {
		if _, err := n.UpsertNodeAttribute(
			instance.AttributeID, modify.MergeWith(instance.AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert node attribute: %+v", instance.NodeAttributeID)
		}
	}

	return nil
}
