package world

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
)

type DecorationMetadata struct {
	AssetID  mid.ID     `json:"asset_id" db:"asset_id"`
	Position cmath.Vec3 `json:"position" db:"position"`
	rotation cmath.Vec3
}

type Metadata struct {
	LOD              []uint32             `json:"lod" db:"lod"`
	Decorations      []DecorationMetadata `json:"decorations,omitempty" db:"decorations,omitempty"`
	AvatarController mid.ID               `json:"avatar_controller" db:"avatar_controller"`
	SkyboxController mid.ID               `json:"skybox_controller" db:"skybox_controller"`
}
