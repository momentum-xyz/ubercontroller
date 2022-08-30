package entry

import (
	"time"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

type Space struct {
	SpaceID     *uuid.UUID    `db:"space_id"`
	SpaceTypeID *uuid.UUID    `db:"space_type_id"`
	OwnerID     *uuid.UUID    `db:"owner_id"`
	ParentID    *uuid.UUID    `db:"parent_id"`
	Asset2dID   *uuid.UUID    `db:"2d_asset_id"`
	Asset3dID   *uuid.UUID    `db:"3d_asset_id"`
	Options     *SpaceOptions `db:"options"`
	Position    *cmath.Vec3   `db:"position"`
	CreatedAt   *time.Time    `db:"created_at"`
	UpdatedAt   *time.Time    `db:"updated_at"`
}

type SpaceOptions struct {
	Options2d            *SpaceOptions2d
	Options3d            *SpaceOptions3d
	FrameTemplate        *SpaceFrameTemplate
	ChildPlacement       *SpaceChildPlacement
	AllowedSubspaceTypes []uuid.UUID
	InfoUIID             *uuid.UUID
	Minimap              *bool
	Visible              *byte
	Secret               *bool
}

type SpaceOptions2d struct {
}

type SpaceOptions3d struct {
}

type SpaceFrameTemplate struct {
}

type SpaceChildPlacement struct {
}
