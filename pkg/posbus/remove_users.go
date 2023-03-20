package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type RemoveUsers struct {
	Users []umid.UMID `json:"users"`
}

func (a *RemoveUsers) Type() MsgType {
	return 0xF5A14BB0
}

func init() {
	registerMessage(&RemoveUsers{})
}
