//go:generate go run -mod vendor gen/mus.go
package posbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gobeam/stringy"
	"github.com/gorilla/websocket"
	"os"
	"reflect"
	"sort"
	"sync"
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
	Type() MsgType
}

type mDef struct {
	Name     string
	TypeName string
	DataType reflect.Type
}

var messageMaps = struct {
	lock     sync.Mutex
	Def      map[MsgType]mDef
	IdByName map[string]MsgType
	IdList   []MsgType
}{Def: make(map[MsgType]mDef), IdByName: make(map[string]MsgType), IdList: nil}

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

func Decode(buf []byte) (Message, error) {
	msgType := MessageType(buf)
	m, ok := reflect.New(MessageDataTypeById(msgType)).Interface().(Message)
	if !ok {
		return nil, errors.New("unknown message type")
	}
	err := DecodeTo(buf, m)
	return m, err
}

func DecodeTo(buf []byte, m Message) error {
	_, err := m.UnmarshalMUS(buf[MsgTypeSize:])
	return err
}

func WSMessage(m Message) *websocket.PreparedMessage {
	msg, _ := websocket.NewPreparedMessage(websocket.BinaryMessage, BinMessage(m))
	return msg
}

func BinMessage(m Message) []byte {
	len := m.SizeMUS()
	buf := make([]byte, MsgTypeSize*2+len)
	msgType := m.Type()
	binary.LittleEndian.PutUint32(buf, uint32(msgType))
	m.MarshalMUS(buf[MsgTypeSize:])
	binary.LittleEndian.PutUint32(buf[MsgTypeSize+len:], uint32(^msgType))
	return buf
}

func MsgTypeName(m Message) string {
	return stringy.New(reflect.ValueOf(m).Elem().Type().Name()).SnakeCase().ToLower()
}

func MsgTypeId(m Message) MsgType {
	return m.Type()
}

func registerMessage(m interface{}) {

	mType := reflect.ValueOf(m).MethodByName("Type").Call([]reflect.Value{})[0].Interface().(MsgType)
	mTypeName := reflect.ValueOf(m).Elem().Type().Name()
	mName := stringy.New(mTypeName).SnakeCase().ToLower()
	messageMaps.Def[mType] = mDef{Name: mName, TypeName: mTypeName, DataType: reflect.ValueOf(m).Elem().Type()}
	messageMaps.IdByName[mName] = mType

	// if not in "generate" check that all messages conforms to the Message interface
	if !isGenerate() {
		m1, ok := m.(Message)
		if !ok {
			fmt.Println("Can not initialize posbus package\n")
			fmt.Println("Message '%+v' does not conform Message interface\n", MsgTypeName(m1))
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

//TypeSetObjectData MsgType = 0xCACE197C
//TypeTriggerVisualEffects MsgType = 0xD96089C6
//TypeUserAction           MsgType = 0xEF1A2E75
