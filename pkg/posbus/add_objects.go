package posbus

type AddObjects struct {
	Objects []ObjectDefinition `json:"objects"`
}

func init() {
	registerMessage(AddObjects{})
}

func (g *AddObjects) GetType() MsgType {
	return 0x2452A9C1
}
