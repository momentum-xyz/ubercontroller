//go:generate go run gen/gen.go
package posbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"sync"

	"github.com/gobeam/stringy"
	"github.com/gorilla/websocket"
)

const (
	MsgTypeSize      = 4
	MsgArrTypeSize   = 4
	MsgUUIDTypeSize  = 16
	MsgLockStateSize = 4
)

type MsgType uint32

type Message interface {
	MarshalMUS(buf []byte) int
	UnmarshalMUS(buf []byte) (int, error)
	SizeMUS() int
	GetType() MsgType
}

type mDef struct {
	Name     string
	TypeName string
	DataType reflect.Type
}

var messageMaps = struct {
	lock       sync.Mutex
	Def        map[MsgType]mDef
	IdByName   map[string]MsgType
	IdList     []MsgType
	ExtraTypes []reflect.Type
}{Def: make(map[MsgType]mDef), IdByName: make(map[string]MsgType), IdList: nil, ExtraTypes: make([]reflect.Type, 0)}

func MessageNameById(id MsgType) string {
	def, ok := messageMaps.Def[id]
	if ok {
		return def.Name
	}
	return "none"
}

func MessageTypeNameById(id MsgType) string {
	def, ok := messageMaps.Def[id]
	if ok {
		return def.TypeName
	}
	return "none"
}

func MessageDataTypeById(id MsgType) reflect.Type {
	def, ok := messageMaps.Def[id]
	if ok {
		return def.DataType
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

func MessageType(buf []byte) MsgType {
	/*if len(m.buf) < 4 {
		// TODO: Handle
	}*/
	header := binary.LittleEndian.Uint32(buf[:MsgTypeSize])
	footer := binary.LittleEndian.Uint32(buf[len(buf)-MsgTypeSize:])
	if header == ^footer {
		return MsgType(header)
	}
	return 0
}

func NewMessageOfType(msgType MsgType) (Message, error) {
	m, ok := reflect.New(MessageDataTypeById(msgType)).Interface().(Message)
	if !ok {
		return nil, errors.New("unknown message type")
	}
	return m, nil
}

func Decode(buf []byte) (Message, error) {
	msgType := MessageType(buf)
	m, ok := NewMessageOfType(msgType)
	if ok != nil {
		return nil, errors.New("unknown message type")
	}
	err := DecodeTo(buf, m)
	return m, err
}

func DecodeTo(buf []byte, m Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic decoding Message: %v", r)
		}
	}()
	_, err = m.UnmarshalMUS(buf[MsgTypeSize:])
	return
}

func WSMessage(m Message) *websocket.PreparedMessage {
	msg, _ := websocket.NewPreparedMessage(websocket.BinaryMessage, BinMessage(m))
	return msg
}

func BinMessage(m Message) []byte {
	len := m.SizeMUS()
	buf := make([]byte, MsgTypeSize*2+len)
	msgType := m.GetType()
	binary.LittleEndian.PutUint32(buf, uint32(msgType))
	m.MarshalMUS(buf[MsgTypeSize:])
	binary.LittleEndian.PutUint32(buf[MsgTypeSize+len:], uint32(^msgType))
	return buf
}

func MsgTypeName(m Message) string {
	return stringy.New(reflect.TypeOf(m).Elem().Name()).SnakeCase().ToLower()
}

func MsgTypeId(m Message) MsgType {
	return m.GetType()
}

func registerMessage[T any](m T) {
	mType := reflect.ValueOf(&m).MethodByName("GetType").Call([]reflect.Value{})[0].Interface().(MsgType)
	mTypeName := reflect.TypeOf(m).Name()
	mName := stringy.New(mTypeName).SnakeCase().ToLower()

	if d, ok := messageMaps.Def[mType]; ok {
		fmt.Printf("Message Type ID '0x%08X' already used for '%+v'\n", mType, d.TypeName)
		os.Exit(-1)
	}
	messageMaps.Def[mType] = mDef{Name: mName, TypeName: mTypeName, DataType: reflect.TypeOf(m)}
	messageMaps.IdByName[mName] = mType

	// if not in "generate" check that all messages conforms to the Message interface
	if !isGenerate() {
		//var m1 Message
		//m1 = &m
		m1, ok := (interface{}(&m)).(Message)
		if !ok {
			fmt.Printf("Can not initialize posbus package\n")
			fmt.Printf("Message '%+v' does not conform Message interface\n", MsgTypeName(m1))
			os.Exit(-1)
		}
	}

}

// GetMessageIds : returns list of message IDs sorted according to message names
func GetMessageIds() []MsgType {

	messageMaps.lock.Lock()
	defer messageMaps.lock.Unlock()
	if messageMaps.IdList != nil {
		return messageMaps.IdList
	}
	msgNames := make([]string, 0, len(messageMaps.IdByName))
	for name, _ := range messageMaps.IdByName {
		msgNames = append(msgNames, name)
	}
	sort.Strings(msgNames)
	messageMaps.IdList = make([]MsgType, 0)
	for _, name := range msgNames {
		messageMaps.IdList = append(messageMaps.IdList, messageMaps.IdByName[name])
	}
	return messageMaps.IdList
}

func isGenerate() bool {
	_, ok1 := os.LookupEnv("GOPACKAGE")
	_, ok2 := os.LookupEnv("GOFILE")
	return ok1 && ok2
}

func addExtraType[T any](v T) {
	messageMaps.ExtraTypes = append(messageMaps.ExtraTypes, reflect.TypeOf(v))
}

func ExtraTypes() []reflect.Type {
	return messageMaps.ExtraTypes
}
