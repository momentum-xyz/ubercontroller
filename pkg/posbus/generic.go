package posbus

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type GenericMessage struct {
	Topic string
	Data  []byte
}

func init() {
	registerMessage(GenericMessage{})
}

func (g *GenericMessage) Type() MsgType {
	return 0xF508E4A3
}

func NewGenericMessage(topic string, data interface{}) *GenericMessage {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		err = errors.WithMessage(err, "NewGenericMessage: failed to marshal data")
		return nil
	}
	return &GenericMessage{Topic: topic, Data: dataJSON}
}
