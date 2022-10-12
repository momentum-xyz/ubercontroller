package node

import (
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
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
		panic("implement")
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
		panic("implement")
	}

	return nil
}

func (n *Node) loadNodeAttributes() error {
	entries, err := n.db.NodeAttributesGetNodeAttributes(n.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get node attributes")
	}

	for _, instance := range entries {
		attributeID := entry.NewAttributeID(instance.PluginID, instance.Name)
		if _, ok := n.GetAttributes().GetAttribute(attributeID); ok {
			n.nodeAttributes.Store(
				entry.NewNodeAttributeID(attributeID),
				entry.NewAttributePayload(instance.Value, instance.Options),
			)
		}
	}

	return nil
}
