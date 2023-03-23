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

func (g *GenericMessage) GetType() MsgType {
	return 0xF508E4A3
}

func NewGenericMessage(topic string, data interface{}) *GenericMessage {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		err = errors.WithMessage(err, "NewGenericMessage: failed to marshal data")
		return nil
	}
	return &GenericMessage{Topic: topic, Data: dataJSON}
	//d := make(map[string]interface{})
	//utils.MapDecode(data, d)
	//fmt.Printf("%+v\n", d)
	//return &GenericMessage{Topic: topic, Data: d}
}
