package node

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) UpsertNodeAttribute(
	nodeAttribute *entry.NodeAttribute, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) error {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[nodeAttribute.AttributeID]
	if !ok {
		payload = (*entry.AttributePayload)(nil)
	}

	payload, err := modifyFn(payload)
	if err != nil {
		return errors.WithMessage(err, "failed to modify attribute payload")
	}

	if updateDB {
		if err := n.db.NodeAttributesUpsertNodeAttribute(
			n.ctx, entry.NewNodeAttribute(nodeAttribute.NodeAttributeID, payload),
		); err != nil {
			return errors.WithMessage(err, "failed to upsert node attribute")
		}
	}

	nodeAttribute.AttributePayload = payload
	n.nodeAttributes.Data[nodeAttribute.AttributeID] = nodeAttribute.AttributePayload

	return nil

}

func (n *Node) GetNodeAttributeValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	payload, ok := n.GetNodeAttributePayload(attributeID)
	if !ok {
		return nil, false
	}
	if payload == nil {
		return nil, false
	}
	return payload.Value, true
}

func (n *Node) GetNodeAttributeOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	payload, ok := n.GetNodeAttributePayload(attributeID)
	if !ok {
		return nil, false
	}
	if payload == nil {
		return nil, false
	}
	return payload.Options, true
}

func (n *Node) GetNodeAttributePayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	return n.nodeAttributes.Load(attributeID)
}

func (n *Node) UpdateNodeAttributeValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) error {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[attributeID]
	if !ok {
		return errors.Errorf("not attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return errors.WithMessage(err, "failed to modify value")
	}

	if updateDB {
		if err := n.db.NodeAttributesUpdateNodeAttributeValue(
			n.ctx, attributeID, value,
		); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Value = value
	n.nodeAttributes.Data[attributeID] = payload

	return nil
}

func (n *Node) UpdateNodeAttributeOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) error {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[attributeID]
	if !ok {
		return errors.Errorf("node attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := n.db.NodeAttributesUpdateNodeAttributeOptions(
			n.ctx, attributeID, options,
		); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Options = options
	n.nodeAttributes.Data[attributeID] = payload

	return nil
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
		if err := n.UpsertNodeAttribute(
			instance, modify.MergeWith(instance.AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert node attribute: %+v", instance.NodeAttributeID)
		}
	}

	return nil
}
