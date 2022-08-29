package universe

import "github.com/google/uuid"

type Asset2dEntry struct {
	Asset2dID *uuid.UUID           `db:"2d_asset_id"`
	Name      *string              `db:"3d_asset_name"`
	Options   *Asset2dOptionsEntry `db:"options"`
}

type Asset2dOptionsEntry struct {
}
