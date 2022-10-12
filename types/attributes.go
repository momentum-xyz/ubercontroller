package types

import (
	"github.com/google/uuid"
)

const (
	NodeSettingsAttributeName = "node_settings"
	UserWalletAttributeName   = "wallet"
)

const (
	UserWalletAddressAttributeValueKey = "address"
)

type BaseAttributeIndex struct {
	PluginID uuid.UUID
	Name     string
}

type NodeAttributeIndex struct {
	BaseAttributeIndex
}

type SpaceAttributeIndex struct {
	BaseAttributeIndex
}

type SpaceUserAttributeIndex struct {
	BaseAttributeIndex
	UserID uuid.UUID
}

func NewNodeAttributeIndex(pluginID uuid.UUID, name string) NodeAttributeIndex {
	return NodeAttributeIndex{
		BaseAttributeIndex: BaseAttributeIndex{
			PluginID: pluginID,
			Name:     name,
		},
	}
}

func NewSpaceAttributeIndex(pluginID uuid.UUID, name string) SpaceAttributeIndex {
	return SpaceAttributeIndex{
		BaseAttributeIndex: BaseAttributeIndex{
			PluginID: pluginID,
			Name:     name,
		},
	}
}

func NewSpaceUserAttributeIndex(pluginID uuid.UUID, name string, userID uuid.UUID) SpaceUserAttributeIndex {
	return SpaceUserAttributeIndex{
		BaseAttributeIndex: BaseAttributeIndex{
			PluginID: pluginID,
			Name:     name,
		},
		UserID: userID,
	}
}
