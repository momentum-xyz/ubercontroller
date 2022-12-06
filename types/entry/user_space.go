package entry

import (
	"time"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type UserSpace struct {
	SpaceID   uuid.UUID      `db:"space_id"`
	UserID    uuid.UUID      `db:"user_id"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	Value     UserSpaceValue `db:"value"`
}

type UserSpaceValue map[string]any

func NewUserSpaceValue() *UserSpaceValue {
	return utils.GetPTR(UserSpaceValue(make(map[string]any)))
}
