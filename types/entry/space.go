package entry

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"time"

	"github.com/google/uuid"
)

type SpaceVisibleType byte

const (
	ReactSpaceVisibleType      SpaceVisibleType = 0b01
	UnitySpaceVisibleType      SpaceVisibleType = 0b10
	ReactUnitySpaceVisibleType SpaceVisibleType = 0b11
)

type Space struct {
	SpaceID     uuid.UUID            `db:"space_id"`
	SpaceTypeID *uuid.UUID           `db:"space_type_id"`
	OwnerID     *uuid.UUID           `db:"owner_id"`
	ParentID    *uuid.UUID           `db:"parent_id"`
	Asset2dID   *uuid.UUID           `db:"asset_2d_id"`
	Asset3dID   *uuid.UUID           `db:"asset_3d_id"`
	Options     *SpaceOptions        `db:"options"`
	Position    *cmath.SpacePosition `db:"position"`
	CreatedAt   time.Time            `db:"created_at"`
	UpdatedAt   *time.Time           `db:"updated_at"`
}

type SpaceOptions struct {
	Asset2dOptions   any                                `db:"asset_2d_options" json:"asset_2d_options,omitempty" mapstructure:"asset_2d_options"`
	Asset3dOptions   any                                `db:"asset_3d_options" json:"asset_3d_options,omitempty" mapstructure:"asset_3d_options"`
	FrameTemplates   map[string]any                     `db:"frame_templates" json:"frame_templates,omitempty" mapstructure:"frame_templates"`
	ChildPlacements  map[uuid.UUID]*SpaceChildPlacement `db:"child_placement" json:"child_placement,omitempty" mapstructure:"child_placement"`
	AllowedSubspaces []uuid.UUID                        `db:"allowed_subspaces" json:"allowed_subspaces,omitempty" mapstructure:"allowed_subspaces"`
	DefaultTiles     []any                              `db:"default_tiles" json:"default_tiles,omitempty" mapstructure:"default_tiles"`
	InfoUIID         *uuid.UUID                         `db:"infoui_id" json:"infoui_id,omitempty" mapstructure:"infoui_id"`
	Minimap          *bool                              `db:"minimap" json:"minimap,omitempty" mapstructure:"minimap"`
	Visible          *SpaceVisibleType                  `db:"visible" json:"visible,omitempty" mapstructure:"visible"`
	Editable         *bool                              `db:"editable" json:"editable,omitempty" mapstructure:"editable"`
	Private          *bool                              `db:"private" json:"private,omitempty" mapstructure:"private"`
	DashboardPlugins []string                           `db:"dashboard_plugins" json:"dashboard_plugins,omitempty" mapstructure:"dashboard_plugins"`
	Subs             map[string]any                     `db:"subs" json:"subs" mapstructure:"subs"`
}

type SpaceChildPlacement struct {
	Algo    *string        `db:"algo" json:"algo,omitempty" mapstructure:"algo"`
	Options map[string]any `db:"options" json:"options,omitempty" mapstructure:"options"`
}

type SpaceAttributeID struct {
	AttributeID
	SpaceID uuid.UUID `db:"space_id"`
}

type SpaceUserAttributeID struct {
	AttributeID
	SpaceID uuid.UUID `db:"space_id"`
	UserID  uuid.UUID `db:"user_id"`
}

type UserSpaceID struct {
	UserID  uuid.UUID `db:"user_id"`
	SpaceID uuid.UUID `db:"space_id"`
}

type SpaceAttribute struct {
	SpaceAttributeID
	*AttributePayload
}

type SpaceUserAttribute struct {
	SpaceUserAttributeID
	*AttributePayload
}

func NewSpaceAttribute(spaceAttributeID SpaceAttributeID, payload *AttributePayload) *SpaceAttribute {
	return &SpaceAttribute{
		SpaceAttributeID: spaceAttributeID,
		AttributePayload: payload,
	}
}

func NewSpaceUserAttribute(spaceUserAttributeID SpaceUserAttributeID, payload *AttributePayload) *SpaceUserAttribute {
	return &SpaceUserAttribute{
		SpaceUserAttributeID: spaceUserAttributeID,
		AttributePayload:     payload,
	}
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

func NewUserSpaceID(userID uuid.UUID, spaceID uuid.UUID) UserSpaceID {
	return UserSpaceID{
		UserID:  userID,
		SpaceID: spaceID,
	}
}
