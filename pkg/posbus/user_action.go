package posbus

type Trigger uint32

const (
	TriggerNone = iota
	TriggerWow
	TriggerHighFive
	TriggerEnteredObject
	TriggerLeftObject
	TriggerStake
)

type UserAction struct {
	Value Trigger `json:"value"`
}

func init() {
	registerMessage(UserAction{})
	addExtraType(Trigger(0))
}

func (g *UserAction) Type() MsgType {
	return 0xEF1A2E75
}
