package entry

import (
	"time"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type ObjectActivityID struct {
	ObjectID   umid.UMID `db:"object_id" json:"object_id"`
	ActivityID umid.UMID `db:"activity_id" json:"activity_id"`
}

type ObjectActivity struct {
	ObjectActivityID
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewObjectActivity(objectActivityID ObjectActivityID) *ObjectActivity {
	return &ObjectActivity{
		ObjectActivityID: objectActivityID,
	}
}

func NewObjectActivityID(objectID umid.UMID, activityID umid.UMID) ObjectActivityID {
	return ObjectActivityID{
		ObjectID:   objectID,
		ActivityID: activityID,
	}
}
