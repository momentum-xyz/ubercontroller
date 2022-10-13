package node

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) GetNodeAttributeValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	payload, ok := n.nodeAttributes.Load(entry.NewNodeAttributeID(attributeID))
	if !ok {
		return nil, false
	}
	return payload.Value, true
}

func (n *Node) GetNodeAttributeOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	payload, ok := n.nodeAttributes.Load(entry.NewNodeAttributeID(attributeID))
	if !ok {
		return nil, false
	}
	return payload.Options, true
}

func (n *Node) GetNodeAttributeEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attr, ok := n.GetAttributes().GetAttribute(attributeID)
	if !ok {
		return nil, false
	}
	payload, ok := n.nodeAttributes.Load(entry.NewNodeAttributeID(attributeID))
	if !ok {
		return nil, false
	}
	return utils.MergePTRs(payload.Options, attr.GetOptions()), true
}

func (n *Node) SetNodeAttributeValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) error {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[entry.NewNodeAttributeID(attributeID)]
	if !ok {
		return errors.Errorf("node attribute not found")
	}

	payload.Value = modifyFn(payload.Value)

	if updateDB {
		if err := n.db.NodeAttributesUpdateNodeAttributeValue(
			n.ctx, attributeID.PluginID, attributeID.Name, payload.Value,
		); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	return nil
}

func (n *Node) SetNodeAttributeOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) error {
	n.nodeAttributes.Mu.Lock()
	defer n.nodeAttributes.Mu.Unlock()

	payload, ok := n.nodeAttributes.Data[entry.NewNodeAttributeID(attributeID)]
	if !ok {
		return errors.Errorf("node attribute not found")
	}

	payload.Options = modifyFn(payload.Options)

	if updateDB {
		if err := n.db.NodeAttributesUpdateNodeAttributeOptions(
			n.ctx, attributeID.PluginID, attributeID.Name, payload.Options,
		); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	return nil
}

func (n *Node) loadNodeAttributes() error {
	entries, err := n.db.NodeAttributesGetNodeAttributes(n.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get node attributes")
	}

	for _, instance := range entries {
		if _, ok := n.GetAttributes().GetAttribute(instance.AttributeID); ok {
			n.nodeAttributes.Store(
				instance.NodeAttributeID,
				entry.NewAttributePayload(instance.Value, instance.Options),
			)
		}
	}

	return nil
}
