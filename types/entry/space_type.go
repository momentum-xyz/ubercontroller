package entry

import (
	"time"

	"github.com/google/uuid"
)

type ObjectType struct {
	SpaceTypeID   uuid.UUID      `db:"space_type_id"`
	Asset2dID     *uuid.UUID     `db:"asset_2d_id"`
	Asset3dID     *uuid.UUID     `db:"asset_3d_id"`
	SpaceTypeName string         `db:"space_type_name"`
	CategoryName  string         `db:"category_name"`
	Description   *string        `db:"description"`
	Options       *ObjectOptions `db:"options"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     *time.Time     `db:"updated_at"`
}
