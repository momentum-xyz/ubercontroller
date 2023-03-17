package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type RemoveObjects struct {
	Objects []umid.UMID `json:"objects"`
}

func init() {
	registerMessage(&RemoveObjects{})
}

func (g *RemoveObjects) Type() MsgType {
	return TypeRemoveObjects
}
