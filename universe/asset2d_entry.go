package universe

import "github.com/google/uuid"

type Asset2DEntry struct {
	Asset2DID *uuid.UUID           `db:"2d_asset_id"`
	Name      *string              `db:"3d_asset_name"`
	Options   *Asset2DOptionsEntry `db:"options"`
}

type Asset2DOptionsEntry struct {
}
