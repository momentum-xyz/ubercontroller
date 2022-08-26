package universe

import (
	"time"

	"github.com/google/uuid"

	"github.com/momentum-xyz/controller/pkg/cmath"
)

type SpaceEntry struct {
	SpaceTypeID uuid.UUID          `db:"space_type_id"`
	OwnerID     uuid.UUID          `db:"owner_id"`
	ParentID    *uuid.UUID         `db:"parent_id"`
	Options     *SpaceOptionsEntry `db:"options"`
	Position    *cmath.Vec3        `db:"position"`
	CreatedAt   time.Time          `db:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at"`
}

type SpaceTypeEntry struct {
	Asset2DID     uuid.UUID          `db:"2d_asset_id"`
	Asset3DID     uuid.UUID          `db:"3d_asset_id"`
	Options       *SpaceOptionsEntry `db:"options"`
	SpaceTypeName string             `db:"space_type_name"`
	CategoryName  string             `db:"category_name"`
	Description   *string            `db:"description"`
	CreatedAt     time.Time          `db:"created_at"`
	UpdatedAt     time.Time          `db:"update_at"`
}

type SpaceOptionsEntry struct {
	Options2DEntry       *SpaceOptions2DEntry
	Options3DEntry       *SpaceOptions3DEntry
	FrameTemplateEntry   *SpaceFrameTemplateEntry
	ChildPlacementEntry  *SpaceChildPlacementEntry
	AllowedSubspaceTypes []uuid.UUID
	InfoUIID             *uuid.UUID
	Minimap              *bool
	Visible              *byte
	Secret               *bool
}

type SpaceOptions2DEntry struct {
}

type SpaceOptions3DEntry struct {
}

type SpaceFrameTemplateEntry struct {
}

type SpaceChildPlacementEntry struct {
}

type SpaceAsset2DEntry struct {
	Name    string                    `db:"3d_asset_name"`
	Options *SpaceAsset2DOptionsEntry `db:"options"`
}

type SpaceAsset2DOptionsEntry struct {
}

type SpaceAsset3DEntry struct {
	Name    string                    `db:"2d_asset_name"`
	Options *SpaceAsset3DOptionsEntry `db:"options"`
}

type SpaceAsset3DOptionsEntry struct {
}
