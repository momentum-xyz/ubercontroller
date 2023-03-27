package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type FlyToMe struct {
	Pilot     umid.UMID `json:"pilot"`
	PilotName string    `json:"pilot_name"`
	ObjectID  umid.UMID `json:"object_id"`
}

func init() {
	registerMessage(FlyToMe{})
}

func (g *FlyToMe) GetType() MsgType {
	return 0xA6EB70C6
}
