package posbus

type AddObjects struct {
	Objects []ObjectDefinition `json:"objects"`
}

func init() {
	registerMessage(AddObjects{})
}

func (g *AddObjects) Type() MsgType {
	return 0x2452A9C1
}
