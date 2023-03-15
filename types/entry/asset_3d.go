package entry

import (
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"time"
)

type Asset3d struct {
	Asset3dID mid.ID          `db:"asset_3d_id" json:"asset_3d_id"`
	Meta      *Asset3dMeta    `db:"meta" json:"meta"`
	Options   *Asset3dOptions `db:"options" json:"options"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

type Asset3dMeta map[string]any

type Asset3dOptions map[string]any
