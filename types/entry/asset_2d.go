package entry

import (
	"time"

	"github.com/google/uuid"
)

type Asset2d struct {
	Asset2dID uuid.UUID       `db:"asset_2d_id" json:"asset_2d_id"`
	Meta      Asset2dMeta     `db:"meta" json:"meta"`
	Options   *Asset2dOptions `db:"options" json:"options"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

type Asset2dMeta map[string]any

type Asset2dOptions map[string]any
