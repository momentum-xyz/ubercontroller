package entry

import (
	"time"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type UserObject struct {
	UserID    uuid.UUID       `db:"user_id" json:"user_id"`
	ObjectID  uuid.UUID       `db:"space_id" json:"object_id"`
	Value     UserObjectValue `db:"value" json:"value"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

type UserObjectValue map[string]any

func NewUserSpaceValue() *UserObjectValue {
	return utils.GetPTR(make(UserObjectValue))
}
