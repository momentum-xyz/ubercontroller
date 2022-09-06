package entry

import (
	"time"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

type SpaceVisibleType byte

const (
	ReactSpaceVisibleType      SpaceVisibleType = 0b01
	UnitySpaceVisibleType      SpaceVisibleType = 0b10
	ReactUnitySpaceVisibleType SpaceVisibleType = 0b11
)

type Space struct {
	SpaceID     *uuid.UUID    `db:"space_id"`
	SpaceTypeID *uuid.UUID    `db:"space_type_id"`
	OwnerID     *uuid.UUID    `db:"owner_id"`
	ParentID    *uuid.UUID    `db:"parent_id"`
	Asset2dID   *uuid.UUID    `db:"asset_2d_id"`
	Asset3dID   *uuid.UUID    `db:"asset_3d_id"`
	Options     *SpaceOptions `db:"options"`
	Position    *cmath.Vec3   `db:"position"`
	CreatedAt   *time.Time    `db:"created_at"`
	UpdatedAt   *time.Time    `db:"updated_at"`
}

type SpaceOptions struct {
	Asset2dOptions   *SpaceAsset2dOptions               `db:"asset_2d_options"`
	Asset3dOptions   *SpaceAsset3dOptions               `db:"asset_3d_options"`
	FrameTemplates   map[string]*SpaceFrameTemplate     `db:"frame_templates"`
	ChildPlacements  map[uuid.UUID]*SpaceChildPlacement `db:"child_placement"`
	AllowedSubspaces []uuid.UUID                        `db:"allowed_subspaces"`
	InfoUIID         *uuid.UUID                         `db:"infoui_id"`
	Minimap          *bool                              `db:"minimap"`
	Visible          *SpaceVisibleType                  `db:"visible"`
	Private          *bool                              `db:"private"`
	// should gone
	DefaultTiles []DefaultTile `db:"default_tiles"`
}

type SpaceAsset2dOptions struct {
}

type SpaceAsset3dOptions struct {
}

type SpaceFrameTemplate struct {
}

type SpaceChildPlacement struct {
	Algo    *string        `db:"algo"`
	Options map[string]any `db:"options"`
}

type DefaultTile struct {
}
