package posbus

type AddUsers struct {
	Users []UserData `json:"users"`
}

func (a *AddUsers) Type() MsgType {
	return TypeAddUsers
}

func init() {
	registerMessage(&AddUsers{})
}
