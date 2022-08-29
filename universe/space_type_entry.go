package universe

import (
	"time"

	"github.com/google/uuid"
)

type SpaceTypeEntry struct {
	SpaceTypeID   *uuid.UUID         `db:"space_type_id"`
	Asset2DID     *uuid.UUID         `db:"2d_asset_id"`
	Asset3DID     *uuid.UUID         `db:"3d_asset_id"`
	SpaceTypeName *string            `db:"space_type_name"`
	CategoryName  *string            `db:"category_name"`
	Description   *string            `db:"description"`
	Options       *SpaceOptionsEntry `db:"options"`
	CreatedAt     *time.Time         `db:"created_at"`
	UpdatedAt     *time.Time         `db:"update_at"`
}
