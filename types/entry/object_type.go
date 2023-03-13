package entry

import (
	"time"

	"github.com/google/uuid"
)

type ObjectType struct {
	ObjectTypeID   uuid.UUID      `db:"object_type_id" json:"object_type_id"`
	Asset2dID      *uuid.UUID     `db:"asset_2d_id" json:"asset_2d_id"`
	Asset3dID      *uuid.UUID     `db:"asset_3d_id" json:"asset_3d_id"`
	ObjectTypeName string         `db:"object_type_name" json:"object_type_name"`
	CategoryName   string         `db:"category_name" json:"category_name"`
	Description    *string        `db:"description" json:"description"`
	Options        *ObjectOptions `db:"options" json:"options"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}
