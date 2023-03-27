package posbus

type EventStart map[string]string

func (r *EventStart) GetType() MsgType {
	return 0xAA854D2C
}

func init() {
	registerMessage(EventStart{})
}
