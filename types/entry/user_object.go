package entry

import (
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"time"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type UserObjectID struct {
	UserID   mid.ID `db:"user_id" json:"user_id"`
	ObjectID mid.ID `db:"object_id" json:"object_id"`
}

type UserObject struct {
	UserObjectID
	Value     *UserObjectValue `db:"value" json:"value"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt time.Time        `db:"updated_at" json:"updated_at"`
}

func NewUserObject(userObjectID UserObjectID, value *UserObjectValue) *UserObject {
	return &UserObject{
		UserObjectID: userObjectID,
		Value:        value,
	}
}

func NewUserObjectID(userID mid.ID, objectID mid.ID) UserObjectID {
	return UserObjectID{
		UserID:   userID,
		ObjectID: objectID,
	}
}

type UserObjectValue map[string]any

func NewUserObjectValue() *UserObjectValue {
	return utils.GetPTR(make(UserObjectValue))
}
