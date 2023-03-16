//go:generate go run gen/mus.go
package posbus

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"log"
	"reflect"
)

type PosbusDataType interface {
	MarshalMUS(buf []byte) int
	UnmarshalMUS(buf []byte) (int, error)
	SizeMUS() int
}

var mapMessageNameById map[MsgType]string
var mapMessageDataTypeById map[MsgType]reflect.Type
var mapMessageIdByName map[string]MsgType

func addToMaps[T any](id MsgType, name string, v T) {
	mapMessageNameById[id] = name
	mapMessageIdByName[name] = id
	mapMessageDataTypeById[id] = reflect.TypeOf(v)
}

func init() {
	mapMessageNameById = make(map[MsgType]string)
	mapMessageIdByName = make(map[string]MsgType)
	mapMessageDataTypeById = make(map[MsgType]reflect.Type)
	addToMaps(NONEType, "none", -1)
	addToMaps(TypeSetUsersTransforms, "set_users_transforms", -1)
	addToMaps(TypeSendTransform, "send_transform", -1)
	addToMaps(TypeGenericMessage, "generic_message", []byte{})
	addToMaps(TypeHandShake, "handshake", HandShake{})
	addToMaps(TypeSetWorld, "set_world", SetWorldData{})
	addToMaps(TypeAddObjects, "add_objects", AddObjects{})
	addToMaps(TypeRemoveObjects, "remove_objects", RemoveObjects{})
	addToMaps(TypeSetObjectPosition, "set_object_position", cmath.ObjectTransform{})
	addToMaps(TypeSetObjectData, "set_object_data", ObjectData{})

	addToMaps(TypeAddUsers, "add_users", AddUsers{})
	addToMaps(TypeRemoveUsers, "remove_users", RemoveUsers{})
	addToMaps(TypeSetUserData, "set_user_data", UserDefinition{})

	addToMaps(TypeSetObjectLock, "set_object_lock", SetObjectLock{})
	addToMaps(TypeObjectLockResult, "object_lock_result", ObjectLockResultData{})
	addToMaps(TypeTriggerVisualEffects, "trigger_visual_effects", -1)
	addToMaps(TypeUserAction, "user_action", -1)
	addToMaps(TypeSignal, "signal", Signal{})
	addToMaps(TypeNotification, "notification", -1)
	addToMaps(TypeTeleportRequest, "teleport_request", TeleportRequest{})
}

func MessageNameById(id MsgType) string {
	name, ok := mapMessageNameById[id]
	if ok {
		return name
	}
	return "none"
}

func MessageDataTypeById(id MsgType) reflect.Type {
	t, ok := mapMessageDataTypeById[id]
	if ok {
		return t
	}
	return nil
}

func MessageIdByName(name string) MsgType {
	id, ok := mapMessageIdByName[name]
	if ok {
		return id
	}
	return 0
}

func BytesToMessage(b []byte) *Message {
	return &Message{
		buf: b,
	}
}

func NewMessageFromData(msgid MsgType, data interface{}) *Message {
	obj := &Message{msgType: msgid}
	var mData bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&mData)
	err := enc.Encode(data)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	obj.makeBuffer(mData.Len())
	mData.Read(obj.Msg())
	return obj
}

func NewMessageFromBuffer(msgid MsgType, buf []byte) *Message {
	obj := &Message{msgType: msgid}
	obj.makeBuffer(len(buf))
	copy(obj.Msg(), buf)
	return obj
}

func NewPreallocatedMessage(msgid MsgType, n int) *Message {
	obj := &Message{msgType: msgid}
	obj.makeBuffer(n)
	return obj
}

func (m *Message) Buf() []byte {
	return m.buf
}

func (m *Message) Msg() []byte {
	return m.buf[MsgTypeSize : len(m.buf)-MsgTypeSize]
}

func (m *Message) Type() MsgType {
	/*if len(m.buf) < 4 {
		// TODO: Handle
	}*/
	if m.msgType != 0 {
		return m.msgType
	} else {
		header := binary.LittleEndian.Uint32(m.buf[:MsgTypeSize])
		footer := binary.LittleEndian.Uint32(m.buf[len(m.buf)-MsgTypeSize:])
		if header == ^footer {
			m.msgType = MsgType(header)
			return m.msgType
		}
	}
	return NONEType
}

func (m *Message) makeBuffer(len int) {
	m.buf = make([]byte, MsgTypeSize*2+len)
	binary.LittleEndian.PutUint32(m.buf, uint32(m.msgType))
	binary.LittleEndian.PutUint32(m.buf[MsgTypeSize+len:], uint32(^m.msgType))
}

func (m *Message) WSMessage() *websocket.PreparedMessage {
	omsg, _ := websocket.NewPreparedMessage(websocket.BinaryMessage, m.Buf())
	return omsg
}

func (m *Message) Decode() (interface{}, error) {
	v := reflect.New(MessageDataTypeById(m.Type())).Interface()
	err := m.DecodeTo(v)
	return v, err
}

func (m *Message) DecodeTo(result interface{}) error {
	var mData bytes.Buffer
	mData.Write(m.Msg())
	dec := gob.NewDecoder(&mData)
	return dec.Decode(result)
}
