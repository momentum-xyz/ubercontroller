package entry

import "github.com/google/uuid"

type Asset2d struct {
	Asset2dID *uuid.UUID      `db:"2d_asset_id"`
	Name      *string         `db:"3d_asset_name"`
	Options   *Asset2dOptions `db:"options"`
}

type Asset2dOptions struct {
}
