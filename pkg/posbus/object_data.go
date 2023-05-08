package posbus

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type ObjectData struct {
	ID      umid.UMID                        `json:"id"`
	Entries map[entry.SlotType]*StringAnyMap `json:"entries"`
}

func init() {
	registerMessage(ObjectData{})
}

func (l *ObjectData) GetType() MsgType {
	return 0xCACE197C
}
