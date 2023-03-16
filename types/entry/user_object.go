package entry

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"time"

	"github.com/momentum-xyz/ubercontroller/utils"
)

type UserObjectID struct {
	UserID   umid.UMID `db:"user_id" json:"user_id"`
	ObjectID umid.UMID `db:"object_id" json:"object_id"`
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

func NewUserObjectID(userID umid.UMID, objectID umid.UMID) UserObjectID {
	return UserObjectID{
		UserID:   userID,
		ObjectID: objectID,
	}
}

type UserObjectValue map[string]any

func NewUserObjectValue() *UserObjectValue {
	return utils.GetPTR(make(UserObjectValue))
}
