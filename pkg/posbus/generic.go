package posbus

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type GenericMessage struct {
	*Message
}

type GenericMessageData struct {
	Topic string
	Data  []byte
}

func NewGenericMessage(topic string, data interface{}) *Message {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		err = errors.WithMessage(err, "NewGenericMessage: failed to marshal data")
		return nil
	}
	return NewMessageFromData(TypeGenericMessage, GenericMessageData{Topic: topic, Data: dataJSON})
}
