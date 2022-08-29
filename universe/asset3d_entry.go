package universe

import "github.com/google/uuid"

type Asset3DEntry struct {
	Asset3DID *uuid.UUID           `db:"3d_asset_id"`
	Name      *string              `db:"2d_asset_name"`
	Options   *Asset3DOptionsEntry `db:"options"`
}

type Asset3DOptionsEntry struct {
}
