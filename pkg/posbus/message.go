//go:generate go run gen/mus.go
package posbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gobeam/stringy"
	"github.com/gorilla/websocket"
	"os"
	"reflect"
)

type PosbusDataType interface {
	MarshalMUS(buf []byte) int
	UnmarshalMUS(buf []byte) (int, error)
	SizeMUS() int
	Type() MsgType
}

var messageMaps = struct {
	NameById     map[MsgType]string
	DataTypeById map[MsgType]reflect.Type
	IdByName     map[string]MsgType
	//initialized  bool
	//lock         sync.Mutex
}{NameById: make(map[MsgType]string), IdByName: make(map[string]MsgType), DataTypeById: make(map[MsgType]reflect.Type)}

func addToMaps(id MsgType, name string, v PosbusDataType) {
	messageMaps.NameById[id] = name
	messageMaps.IdByName[name] = id
	messageMaps.DataTypeById[id] = reflect.TypeOf(v)
}

func MessageNameById(id MsgType) string {
	name, ok := messageMaps.NameById[id]
	if ok {
		return name
	}
	return "none"
}

func MessageDataTypeById(id MsgType) reflect.Type {
	t, ok := messageMaps.DataTypeById[id]
	if ok {
		return t
	}
	return nil
}

func MessageIdByName(name string) MsgType {
	id, ok := messageMaps.IdByName[name]
	if ok {
		return id
	}
	return 0
}

//func NewMessageFromBuffer(msgid MsgType, buf []byte) *Message {
//	obj := &Message{msgType: msgid}
//	obj.makeBuffer(len(buf))
//	copy(obj.Msg(), buf)
//	return obj
//}

//func (m *Message) Msg() []byte {
//	return m.buf[MsgTypeSize : len(m.buf)-MsgTypeSize]
//}

func MessageType(buf []byte) MsgType {
	/*if len(m.buf) < 4 {
		// TODO: Handle
	}*/
	header := binary.LittleEndian.Uint32(buf[:MsgTypeSize])
	footer := binary.LittleEndian.Uint32(buf[len(buf)-MsgTypeSize:])
	if header == ^footer {
		return MsgType(header)
	}
	return TypeNONE
}

//func (m *Message) makeBuffer(len int) {
//	m.buf = make([]byte, MsgTypeSize*2+len)
//	binary.LittleEndian.PutUint32(m.buf, uint32(m.msgType))
//	binary.LittleEndian.PutUint32(m.buf[MsgTypeSize+len:], uint32(^m.msgType))
//}

func Decode(buf []byte) (PosbusDataType, error) {
	msgType := MessageType(buf)
	m, ok := reflect.New(MessageDataTypeById(msgType)).Interface().(PosbusDataType)
	if !ok {
		return nil, errors.New("unknown message type")
	}
	err := DecodeTo(buf, m)
	return m, err
}

func DecodeTo(buf []byte, m PosbusDataType) error {
	_, err := m.UnmarshalMUS(buf[MsgTypeSize:])
	return err
}

func WSMessage(m PosbusDataType) *websocket.PreparedMessage {
	msg, _ := websocket.NewPreparedMessage(websocket.BinaryMessage, BinMessage(m))
	return msg
}

func BinMessage(m PosbusDataType) []byte {
	len := m.SizeMUS()
	buf := make([]byte, MsgTypeSize*2+len)
	msgType := m.Type()
	binary.LittleEndian.PutUint32(buf, uint32(msgType))
	m.MarshalMUS(buf[MsgTypeSize:])
	binary.LittleEndian.PutUint32(buf[MsgTypeSize+len:], uint32(^msgType))
	return buf
}

func MsgTypeName(m PosbusDataType) string {
	return stringy.New(reflect.ValueOf(m).Elem().Type().Name()).SnakeCase().ToLower()
}

func MsgTypeId(m PosbusDataType) MsgType {
	return m.Type()
}

func registerMessage(m interface{}) {
	if !isGenerate() {
		m1, ok := m.(PosbusDataType)
		if !ok {
			fmt.Println("Can not initialize posbus package\n")
			fmt.Println("Message '%+v' does not conform PosbusDataType interface\n", MsgTypeName(m1))
			os.Exit(-1)
		}
		addToMaps(m1.Type(), MsgTypeName(m1), m1)
	}
}

func isGenerate() bool {
	_, ok1 := os.LookupEnv("GOPACKAGE")
	_, ok2 := os.LookupEnv("GOFILE")
	return ok1 && ok2
}
