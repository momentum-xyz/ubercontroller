package entry

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"time"
)

type ObjectVisibleType byte

const (
	InvisibleObjectVisibleType  ObjectVisibleType = 0b00
	ReactObjectVisibleType      ObjectVisibleType = 0b01
	UnityObjectVisibleType      ObjectVisibleType = 0b10
	ReactUnityObjectVisibleType ObjectVisibleType = 0b11
)

type Object struct {
	ObjectID     mid.ID                 `db:"object_id" json:"object_id"`
	ObjectTypeID mid.ID                 `db:"object_type_id" json:"object_type_id"`
	OwnerID      mid.ID                 `db:"owner_id" json:"owner_id"`
	ParentID     mid.ID                 `db:"parent_id" json:"parent_id"`
	Asset2dID    *mid.ID                `db:"asset_2d_id" json:"asset_2d_id"`
	Asset3dID    *mid.ID                `db:"asset_3d_id" json:"asset_3d_id"`
	Options      *ObjectOptions         `db:"options" json:"options"`
	Position     *cmath.ObjectTransform `db:"position" json:"position"`
	CreatedAt    time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time              `db:"updated_at" json:"updated_at"`
}

type ObjectOptions struct {
	Asset2dOptions    any                              `db:"asset_2d_options" json:"asset_2d_options,omitempty"`
	Asset3dOptions    any                              `db:"asset_3d_options" json:"asset_3d_options,omitempty"`
	FrameTemplates    map[string]any                   `db:"frame_templates" json:"frame_templates,omitempty"`
	ChildPlacements   map[mid.ID]*ObjectChildPlacement `db:"child_placement" json:"child_placement,omitempty"`
	AllowedSubObjects []mid.ID                         `db:"allowed_subobjects" json:"allowed_subobjects,omitempty"`
	DefaultTiles      []any                            `db:"default_tiles" json:"default_tiles,omitempty"`
	InfoUIID          *mid.ID                          `db:"infoui_id" json:"infoui_id,omitempty"`
	Minimap           *bool                            `db:"minimap" json:"minimap,omitempty"`
	Visible           *ObjectVisibleType               `db:"visible" json:"visible,omitempty"`
	Editable          *bool                            `db:"editable" json:"editable,omitempty"`
	Private           *bool                            `db:"private" json:"private,omitempty"`
	DashboardPlugins  []string                         `db:"dashboard_plugins" json:"dashboard_plugins,omitempty"`
	Subs              map[string]any                   `db:"subs" json:"subs"`
}

type ObjectChildPlacement struct {
	Algo    *string        `db:"algo" json:"algo,omitempty"`
	Options map[string]any `db:"options" json:"options,omitempty"`
}

type ObjectAttributeID struct {
	AttributeID
	ObjectID mid.ID `db:"object_id" json:"object_id"`
}

type ObjectUserAttributeID struct {
	AttributeID
	ObjectID mid.ID `db:"object_id" json:"object_id"`
	UserID   mid.ID `db:"user_id" json:"user_id"`
}

type ObjectAttribute struct {
	ObjectAttributeID
	*AttributePayload
}

type ObjectUserAttribute struct {
	ObjectUserAttributeID
	*AttributePayload
}

func NewObjectAttribute(objectAttributeID ObjectAttributeID, payload *AttributePayload) *ObjectAttribute {
	return &ObjectAttribute{
		ObjectAttributeID: objectAttributeID,
		AttributePayload:  payload,
	}
}

func NewObjectUserAttribute(
	objectUserAttributeID ObjectUserAttributeID, payload *AttributePayload,
) *ObjectUserAttribute {
	return &ObjectUserAttribute{
		ObjectUserAttributeID: objectUserAttributeID,
		AttributePayload:      payload,
	}
}

func NewObjectAttributeID(attributeID AttributeID, objectID mid.ID) ObjectAttributeID {
	return ObjectAttributeID{
		AttributeID: attributeID,
		ObjectID:    objectID,
	}
}

func NewObjectUserAttributeID(attributeID AttributeID, objectID, userID mid.ID) ObjectUserAttributeID {
	return ObjectUserAttributeID{
		AttributeID: attributeID,
		ObjectID:    objectID,
		UserID:      userID,
	}
}
