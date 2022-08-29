package universe

import (
	"time"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

type SpaceEntry struct {
	SpaceID     *uuid.UUID         `db:"space_id"`
	SpaceTypeID *uuid.UUID         `db:"space_type_id"`
	OwnerID     *uuid.UUID         `db:"owner_id"`
	ParentID    *uuid.UUID         `db:"parent_id"`
	Asset2dID   *uuid.UUID         `db:"2d_asset_id"`
	Asset3dID   *uuid.UUID         `db:"3d_asset_id"`
	Options     *SpaceOptionsEntry `db:"options"`
	Position    *cmath.Vec3        `db:"position"`
	CreatedAt   *time.Time         `db:"created_at"`
	UpdatedAt   *time.Time         `db:"updated_at"`
}

type SpaceOptionsEntry struct {
	Options2dEntry       *SpaceOptions2dEntry
	Options3dEntry       *SpaceOptions3dEntry
	FrameTemplateEntry   *SpaceFrameTemplateEntry
	ChildPlacementEntry  *SpaceChildPlacementEntry
	AllowedSubspaceTypes []uuid.UUID
	InfoUIID             *uuid.UUID
	Minimap              *bool
	Visible              *byte
	Secret               *bool
}

type SpaceOptions2dEntry struct {
}

type SpaceOptions3dEntry struct {
}

type SpaceFrameTemplateEntry struct {
}

type SpaceChildPlacementEntry struct {
}
