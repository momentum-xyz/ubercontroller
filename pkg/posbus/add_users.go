package posbus

type AddUsers struct {
	Users []UserData `json:"users"`
}

func (a *AddUsers) Type() MsgType {
	return 0xF51F2AFF
}

func init() {
	registerMessage(&AddUsers{})
}
