package entry

import (
	"github.com/google/uuid"
	"time"
)

type Asset3d struct {
	Asset3dID uuid.UUID       `db:"asset_3d_id"`
	Meta      Meta            `db:"meta"`
	Options   *Asset3dOptions `db:"options"`
	CreatedAt time.Time       `db:"created_at"`
	UpdateAt  *time.Time      `db:"updated_at"`
}

type Asset3dOptions struct {
}
