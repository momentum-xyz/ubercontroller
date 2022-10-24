package entry

import (
	"time"

	"github.com/google/uuid"
)

type Asset3d struct {
	Asset3dID uuid.UUID       `db:"asset_3d_id"`
	Meta      *Asset3dMeta    `db:"meta"`
	Options   *Asset3dOptions `db:"options"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt *time.Time      `db:"updated_at"`
}

type Asset3dMeta map[string]any

type Asset3dOptions struct {
}
