package entry

import (
	"time"
	
	"github.com/google/uuid"
)

type Asset2d struct {
	Asset2dID uuid.UUID       `db:"asset_2d_id"`
	Meta      *Asset2dMeta    `db:"meta"`
	Options   *Asset2dOptions `db:"options"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt *time.Time      `db:"updated_at"`
}

type Asset2dMeta map[string]any

type Asset2dOptions struct {
}
