package entry

import (
	"time"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type Vec3 struct {
	X float32 `json:"x" db:"x"`
	Y float32 `json:"y" db:"y"`
	Z float32 `json:"z" db:"z"`
}

type Activity struct {
	ActivityID umid.UMID     `db:"activity_id" json:"activity_id"`
	UserID     *umid.UMID    `db:"user_id" json:"user_id"`
	ObjectID   *umid.UMID    `db:"object_id" json:"object_id"`
	Type       *string       `db:"type" json:"type"`
	Data       *ActivityData `db:"data" json:"data"`
	CreatedAt  time.Time     `db:"created_at" json:"created_at"`
}

type ActivityData struct {
	Position       Vec3 `db:"position" json:"position"`
	Asset3dOptions any  `db:"asset_3d_options" json:"asset_3d_options,omitempty"`
}
