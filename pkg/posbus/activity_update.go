package posbus

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// TODO: does musgo support type aliases/const like this?
type ActivityUpdateType string

const (
	ActivityUpdateChangeType  ActivityUpdateType = ""
	NewActivityUpdateType     ActivityUpdateType = "new"
	ChangedActivityUpdateType ActivityUpdateType = "changed"
	RemovedActivityUpdateType ActivityUpdateType = "removed"
)

type ActivityUpdate struct {
	ActivityId umid.UMID           `json:"activity_id"`
	ChangeType string              `json:"change_type"`
	Type       *entry.ActivityType `json:"type"`
	Data       *entry.ActivityData `json:"data"`
	UserId     umid.UMID           `json:"user_id"`
	ObjectId   umid.UMID           `json:"object_id"`
}

func (r *ActivityUpdate) GetType() MsgType {
	return 0xCA57695D
}

func init() {
	registerMessage(ActivityUpdate{})
}
