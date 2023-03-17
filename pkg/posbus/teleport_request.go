package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type TeleportRequest struct {
	Target umid.UMID `json:"target"`
}

func init() {
	registerMessage(&TeleportRequest{})
}

func (g *TeleportRequest) Type() MsgType {
	return TypeTeleportRequest
}
