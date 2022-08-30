package entry

import "github.com/google/uuid"

type Asset3d struct {
	Asset3dID *uuid.UUID      `db:"3d_asset_id"`
	Name      *string         `db:"2d_asset_name"`
	Options   *Asset3dOptions `db:"options"`
}

type Asset3dOptions struct {
}
