package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type RemoveObjects struct {
	Objects []umid.UMID `json:"objects"`
}

func init() {
	registerMessage(RemoveObjects{})
}

func (g *RemoveObjects) GetType() MsgType {
	return 0x6BF88C24
}
