package entry

import (
	"time"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type UserActivityID struct {
	UserID     umid.UMID `db:"user_id" json:"user_id"`
	ActivityID umid.UMID `db:"activity_id" json:"activity_id"`
}

type UserActivity struct {
	UserActivityID
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewUserActivity(userActivityID UserActivityID) *UserActivity {
	return &UserActivity{
		UserActivityID: userActivityID,
	}
}

func NewUserActivityID(userID umid.UMID, activityID umid.UMID) UserActivityID {
	return UserActivityID{
		UserID:     userID,
		ActivityID: activityID,
	}
}
