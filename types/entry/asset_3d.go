package entry

import (
	"time"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type Asset3d struct {
	Asset3dID umid.UMID       `db:"asset_3d_id" json:"asset_3d_id"`
	Meta      *Asset3dMeta    `db:"meta" json:"meta"`
	Options   *Asset3dOptions `db:"options" json:"options"`
	Private   bool            `db:"is_private" json:"is_private"`
	UserID    umid.UMID       `db:"user_id" json:"user_id"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

type Asset3dMeta map[string]any

type Asset3dOptions map[string]any
