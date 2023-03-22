package posbus

type SignalType uint32

const (
	SignalNone SignalType = iota
	SignalDualConnection
	SignalReady
	SignalInvalidToken
	SignalSpawn
	SignalLeaveWorld
	SignalConnectionFailed
	SignalConnected
	SignalConnectionClosed
	SignalWorldDoesNotExist
)

type Signal struct {
	Value SignalType `json:"value"`
}

func init() {
	registerMessage(&Signal{})
}

func (g *Signal) Type() MsgType {
	return 0xADC1964D
}
