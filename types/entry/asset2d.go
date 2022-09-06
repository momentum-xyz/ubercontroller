package entry

import "github.com/google/uuid"

type Asset2d struct {
	Asset2dID *uuid.UUID      `db:"asset_2d_id"`
	Name      *string         `db:"asset_2d_name"`
	Options   *Asset2dOptions `db:"options"`
}

type Asset2dOptions struct {
}
