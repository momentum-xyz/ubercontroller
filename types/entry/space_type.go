package entry

import (
	"time"

	"github.com/google/uuid"
)

type SpaceType struct {
	SpaceTypeID   *uuid.UUID    `db:"space_type_id"`
	Asset2dID     *uuid.UUID    `db:"2d_asset_id"`
	Asset3dID     *uuid.UUID    `db:"3d_asset_id"`
	SpaceTypeName *string       `db:"space_type_name"`
	CategoryName  *string       `db:"category_name"`
	Description   *string       `db:"description"`
	Options       *SpaceOptions `db:"options"`
	CreatedAt     *time.Time    `db:"created_at"`
	UpdatedAt     *time.Time    `db:"update_at"`
}
