package entry

import (
	"time"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type Activity struct {
	ActivityID umid.UMID     `db:"activity_id" json:"activity_id"`
	UserID     *umid.UMID    `db:"user_id" json:"user_id"`
	ObjectID   *umid.UMID    `db:"object_id" json:"object_id"`
	Type       *ActivityType `db:"type" json:"type"`
	Data       *ActivityData `db:"data" json:"data"`
	CreatedAt  time.Time     `db:"created_at" json:"created_at"`
}

type ActivityData struct {
	Position    *cmath.Vec3 `db:"position" json:"position"`
	Description *string     `db:"description" json:"description"`
	Hash        *string     `db:"hash" json:"hash"`
}

type ActivityType string

const (
	ActivityTypeVideo      ActivityType = "video"
	ActivityTypeScreenshot ActivityType = "screenshot"
)
