package entry

import (
	"github.com/google/uuid"
	"time"
)

type Asset2d struct {
	Asset2dID uuid.UUID       `db:"asset_2d_id"`
	Meta      Meta            `db:"meta"`
	Options   *Asset2dOptions `db:"options"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt *time.Time      `db:"updated_at"`
}

type Asset2dOptions struct {
}
