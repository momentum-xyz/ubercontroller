package node

import (
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) GetNodeAttributeValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	payload, ok := n.nodeAttributes.Load(attributeID)
	if !ok {
		return nil, false
	}
	return payload.Value, true
}

func (n *Node) GetNodeAttributeOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	payload, ok := n.nodeAttributes.Load(attributeID)
	if !ok {
		return nil, false
	}
	return payload.Options, true
}

func (n *Node) GetNodeAttributeEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attr, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(attributeID))
	if !ok {
		return nil, false
	}
	payload, ok := n.nodeAttributes.Load(attributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(payload.Options, attr.GetOptions())
	if err != nil {
		n.log.Error(
			err, "Node: GetNodeAttributeEffectiveOptions: failed to merge effective options: %+v", attributeID,
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (n *Node) SetNodeAttributeValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) error {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[attributeID]
	if !ok {
		return errors.Errorf("node attribute not found")
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return errors.WithMessage(err, "failed to modify value")
	}

	if updateDB {
		if err := n.db.NodeAttributesUpdateNodeAttributeValue(n.ctx, attributeID, value); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Value = value

	return nil
}

func (n *Node) SetNodeAttributeOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) error {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[attributeID]
	if !ok {
		return errors.Errorf("node attribute not found")
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := n.db.NodeAttributesUpdateNodeAttributeOptions(n.ctx, attributeID, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Options = options

	return nil
}

func (n *Node) loadNodeAttributes() error {
	entries, err := n.db.NodeAttributesGetNodeAttributes(n.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get node attributes")
	}

	for _, instance := range entries {
		if _, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(instance.AttributeID)); ok {
			n.nodeAttributes.Store(
				instance.AttributeID,
				entry.NewAttributePayload(instance.Value, instance.Options),
			)
		}
	}

	return nil
}
