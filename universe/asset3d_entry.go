package universe

import "github.com/google/uuid"

type Asset3dEntry struct {
	Asset3dID *uuid.UUID           `db:"3d_asset_id"`
	Name      *string              `db:"2d_asset_name"`
	Options   *Asset3dOptionsEntry `db:"options"`
}

type Asset3dOptionsEntry struct {
}
