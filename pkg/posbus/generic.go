package posbus

import (
	"github.com/gorilla/websocket"
)

type GenericMessage struct {
	*Message
}

type RelayToReactData struct {
	Topic string
	Data  []byte
}

func NewRelayToReactMsg(topic string, data []byte) *websocket.PreparedMessage {

	return WrapAsMessage(GenericMessageType, RelayToReactData{Topic: topic, Data: data})
}
