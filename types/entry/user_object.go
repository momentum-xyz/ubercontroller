package entry

import (
	"time"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type UserObjectID struct {
	UserID   uuid.UUID `db:"user_id" json:"user_id"`
	ObjectID uuid.UUID `db:"object_id" json:"object_id"`
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

func NewUserObjectID(userID uuid.UUID, objectID uuid.UUID) UserObjectID {
	return UserObjectID{
		UserID:   userID,
		ObjectID: objectID,
	}
}

type UserObjectValue map[string]any

func NewUserObjectValue() *UserObjectValue {
	return utils.GetPTR(make(UserObjectValue))
}
