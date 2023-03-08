package posbus

type GenericMessage struct {
	*Message
}

type RelayToReactData struct {
	Topic string
	Data  []byte
}

func NewRelayToReactMsg(topic string, data []byte) *Message {

	return WrapAsMessage(GenericMessageType, RelayToReactData{Topic: topic, Data: data})
}
