package world

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

type DecorationMetadata struct {
	AssetID  uuid.UUID  `json:"asset_id" db:"asset_id"`
	Position cmath.Vec3 `json:"position" db:"position"`
	rotation cmath.Vec3
}

type Metadata struct {
	LOD              []uint32             `json:"lod" db:"lod"`
	Decorations      []DecorationMetadata `json:"decorations,omitempty" db:"decorations,omitempty"`
	AvatarController uuid.UUID            `json:"avatar_controller" db:"avatar_controller"`
	SkyboxController uuid.UUID            `json:"skybox_controller" db:"skybox_controller"`
}
