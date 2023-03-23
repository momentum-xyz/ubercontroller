package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type HighFive struct {
	SenderID   umid.UMID `json:"sender_id"`
	ReceiverID umid.UMID `json:"receiver_id"`
	Message    string    `json:"message"`
}

func init() {
	registerMessage(HighFive{})
}

func (g *HighFive) GetType() MsgType {
	return 0x3D501432
}
