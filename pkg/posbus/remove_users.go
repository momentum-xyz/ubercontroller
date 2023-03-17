package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type RemoveUsers struct {
	Users []umid.UMID `json:"users"`
}

func (a *RemoveUsers) Type() MsgType {
	return TypeRemoveUsers
}

func init() {
	registerMessage(&RemoveUsers{})
}
