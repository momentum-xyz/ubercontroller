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
	Asset2dOptions   any                                `db:"asset_2d_options" json:"asset_2d_options"`
	Asset3dOptions   any                                `db:"asset_3d_options" json:"asset_3d_options"`
	FrameTemplates   map[string]any                     `db:"frame_templates" json:"frame_templates"`
	ChildPlacements  map[uuid.UUID]*SpaceChildPlacement `db:"child_placement" json:"child_placement"`
	AllowedSubspaces []uuid.UUID                        `db:"allowed_subspaces" json:"allowed_subspaces"`
	DefaultTiles     []any                              `db:"default_tiles" json:"default_tiles"`
	InfoUIID         *uuid.UUID                         `db:"infoui_id" json:"infoui_id"`
	Minimap          *bool                              `db:"minimap" json:"minimap"`
	Visible          *SpaceVisibleType                  `db:"visible" json:"visible"`
	Private          *bool                              `db:"private" json:"private"`
	DashboardPlugins []string                           `db:"dashboard_plugins" json:"dashboard_plugins"`
}

type SpaceChildPlacement struct {
	Algo    *string        `db:"algo" json:"algo"`
	Options map[string]any `db:"options" json:"options"`
}

type SpaceAttributeID struct {
	AttributeID
	SpaceID uuid.UUID `db:"space_id"`
}

type SpaceAttribute struct {
	SpaceAttributeID
	AttributePayload
}

type SpaceUserAttributeID struct {
	AttributeID
	SpaceID uuid.UUID `db:"space_id"`
	UserID  uuid.UUID `db:"user_id"`
}

type SpaceUserAttribute struct {
	SpaceUserAttributeID
	AttributePayload
}

func NewSpaceAttributeID(attributeID AttributeID, spaceID uuid.UUID) SpaceAttributeID {
	return SpaceAttributeID{
		AttributeID: attributeID,
		SpaceID:     spaceID,
	}
}

func NewSpaceUserAttributeID(attributeID AttributeID, spaceID, userID uuid.UUID) SpaceUserAttributeID {
	return SpaceUserAttributeID{
		AttributeID: attributeID,
		SpaceID:     spaceID,
		UserID:      userID,
	}
}
