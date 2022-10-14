package entry

import "github.com/google/uuid"

type Asset3d struct {
	Asset3dID uuid.UUID       `db:"asset_3d_id"`
	Name      string          `db:"asset_3d_name"`
	Options   *Asset3dOptions `db:"options"`
}

type Asset3dOptions struct {
}
