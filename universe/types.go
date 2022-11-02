package universe

import "github.com/google/uuid"

type SpaceFilterPredicateFn func(spaceID uuid.UUID, space Space) bool

const (
	// node
	NodeAttributeNodeSettingsName = "node_settings"

	// kusama
	KusamaUserAttributeWalletName      = "wallet"
	KusamaUserAttributeWalletWalletKey = "wallet"
)
