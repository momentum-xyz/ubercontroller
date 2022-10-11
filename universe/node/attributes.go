package node

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

func (n *Node) loadNodeAttributes() error {
	entries, err := n.db.NodeAttributesGetNodeAttributes(n.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get node attributes")
	}

	for _, instance := range entries {
		attr, ok := n.GetAttributes().GetAttribute(
			entry.AttributeID{
				PluginID: instance.PluginID,
				Name:     instance.Name,
			},
		)
		if ok {
			n.nodeAttributes.SetAttributeInstance(
				types.NewNodeAttributeIndex(instance.PluginID, instance.Name), attr, instance.Value, instance.Options,
			)
		}
	}

	return nil
}
