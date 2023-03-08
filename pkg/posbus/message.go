package posbus

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"log"
	"reflect"
)

type Signal uint32

const (
	SignalNone Signal = iota
	SignalDualConnection
	SignalReady
	SignalInvalidToken
	SignalSpawn
	SignalLeaveWorld
)

type Trigger uint32

const (
	TriggerNone = iota
	TriggerWow
	TriggerHighFive
	TriggerEnteredObject
	TriggerLeftObject
	TriggerStake
)

type Notification uint32

const (
	NotificationNone     Notification = 0
	NotificationWow      Notification = 1
	NotificationHighFive Notification = 2

	NotificationStageModeAccept        Notification = 10
	NotificationStageModeInvitation    Notification = 11
	NotificationStageModeSet           Notification = 12
	NotificationStageModeStageJoin     Notification = 13
	NotificationStageModeStageRequest  Notification = 14
	NotificationStageModeStageDeclined Notification = 15

	NotificationTextMessage Notification = 500
	NotificationRelay       Notification = 501

	NotificationGeneric Notification = 999
	NotificationLegacy  Notification = 1000
)

const (
	MsgTypeSize      = 4
	MsgArrTypeSize   = 4
	MsgUUIDTypeSize  = 16
	MsgOnePosSize    = 4
	MsgLockStateSize = 4
)

type Destination byte

const (
	DestinationUnity Destination = 0b01
	DestinationReact Destination = 0b10
	DestinationBoth  Destination = 0b11
)

type MsgType uint32

/* can use fmt.Sprintf("%x", int) to display hex */
const (
	NONEType              MsgType = 0x00000000
	SetUsersPositionsType MsgType = 0x285954B8
	SendPositionType      MsgType = 0xF878C4BF
	GenericMessageType    MsgType = 0xF508E4A3
	HandShakeType         MsgType = 0x7C41941A
	SetWorldType          MsgType = 0xCCDF2E49

	AddObjectsType        MsgType = 0x2452A9C1
	RemoveObjectsType     MsgType = 0x6BF88C24
	SetObjectPositionType MsgType = 0xEA6DA4B4

	SetObjectDataType MsgType = 0xCACE197C

	AddUsersType    MsgType = 0xF51F2AFF
	RemoveUsersType MsgType = 0xF5A14BB0
	SetUserDataType MsgType = 0xF702EF5F

	SetObjectLockType    MsgType = 0xA7DE9F59
	ObjectLockResultType MsgType = 0x0924668C

	TriggerVisualEffectsType MsgType = 0xD96089C6
	UserActionType           MsgType = 0xEF1A2E75

	SignalType       MsgType = 0xADC1964D
	NotificationType MsgType = 0xC1FB41D7

	TeleportRequestType MsgType = 0x78DA55D9
)

type ObjectDefinition struct {
	ID               uuid.UUID
	ParentID         uuid.UUID
	AssetType        uuid.UUID
	AssetFormat      dto.Asset3dType // TODO: Rename AssetType to AssetID, so Type can be used for this.
	Name             string
	ObjectTransform  cmath.ObjectPosition
	IsEditable       bool
	TetheredToParent bool
	ShowOnMiniMap    bool
	//InfoUI           uuid.UUID
}

type UserDefinition struct {
	ID              uuid.UUID
	Name            string
	Avatar          uuid.UUID
	ObjectTransform cmath.ObjectPosition
	IsGuest         bool
}

type SetWorld struct {
	ID              uuid.UUID
	Name            string
	Avatar          uuid.UUID
	Owner           uuid.UUID
	Avatar3DAssetID uuid.UUID
}

type ObjectDataIndex struct {
	Kind     entry.UnitySlotType
	SlotName string
}

type ObjectData struct {
	ID      uuid.UUID
	Entries map[ObjectDataIndex]interface{}
}

type SetObjectLock struct {
	ID    uuid.UUID
	State uint32
}

type ObjectLockResultData struct {
	ID        uuid.UUID
	Result    uint32
	LockOwner uuid.UUID
}

type ObjectPosition struct {
	ID              uuid.UUID
	ObjectTransform cmath.ObjectPosition
}

type Message struct {
	buf     []byte
	msgType MsgType
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
	addToMaps(SetUsersPositionsType, "set_user_position", -1)
	addToMaps(SendPositionType, "send_position", -1)
	addToMaps(GenericMessageType, "generic_message", []byte{})
	addToMaps(HandShakeType, "handshake", HandShake{})
	addToMaps(SetWorldType, "set_world", SetWorld{})
	addToMaps(AddObjectsType, "add_objects", []ObjectDefinition{})
	addToMaps(RemoveObjectsType, "remove_objects", []uuid.UUID{})
	addToMaps(SetObjectPositionType, "set_object_position", cmath.ObjectPosition{})
	addToMaps(SetObjectDataType, "set_object_data", ObjectData{})
	addToMaps(AddUsersType, "add_users", UserDefinition{})
	addToMaps(RemoveUsersType, "remove_users", []uuid.UUID{})
	addToMaps(SetUserDataType, "set_user_data", UserDefinition{})
	addToMaps(SetObjectLockType, "set_object_lock", SetObjectLock{})
	addToMaps(ObjectLockResultType, "object_lock_result", ObjectLockResultData{})
	addToMaps(TriggerVisualEffectsType, "trigger_visual_effects", -1)
	addToMaps(UserActionType, "user_action", -1)
	addToMaps(SignalType, "signal", int(0))
	addToMaps(NotificationType, "notification", -1)
	addToMaps(TeleportRequestType, "teleport_request", uuid.UUID{})
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

func NewMessage(msgid MsgType) *Message {
	obj := &Message{msgType: msgid}
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

func (m *Message) AsSendPos() []byte {
	// how 16 have been calculated?
	return m.buf[MsgTypeSize : MsgTypeSize+(6*MsgOnePosSize)]
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

func MsgFromBytes(b []byte) *Message {
	return &Message{
		buf: b,
	}
}

func WrapAsMessage(msgid MsgType, data interface{}) *Message {
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

func (m *Message) AsSignal() Signal {
	if m.Type() != SignalType {
		fmt.Printf("Unexpected (not a Signal) message type %+v\n", m.Type())
		return 0
	}
	var mData bytes.Buffer // Stand-in for a network connection
	enc := gob.NewDecoder(&mData)
	var sig Signal
	err := enc.Decode(&sig)
	if err != nil {
		fmt.Printf("Unexpected (not a Signal) message type %+v\n", m.Type())
		return 0
	}
	return sig
}
