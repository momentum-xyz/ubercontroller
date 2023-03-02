package posbus

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/utils"
)

const UserPositionsMessageSize = MsgUUIDTypeSize + MsgOnePosSize*6

type SetUserPositionsMessage struct {
	*Message
	maxUsers  int
	posBuffer []byte
	nUsers    int
}

func NewSetUserPositionsMsg(maxUsers int) *SetUserPositionsMessage {
	obj := &SetUserPositionsMessage{Message: NewMessage(SetUsersPositionsType), maxUsers: maxUsers}
	obj.posBuffer = make([]byte, UserPositionsMessageSize*maxUsers)
	obj.nUsers = 0
	return obj
}

func (m *SetUserPositionsMessage) AddPosition(data []byte) {
	start := MsgArrTypeSize + m.nUsers*UserPositionsMessageSize
	copy(m.posBuffer[start:], data)
	m.nUsers++
}

func (m *SetUserPositionsMessage) Finalize() *websocket.PreparedMessage {
	m.makeBuffer(len(m.posBuffer))
	copy(m.Msg(), m.posBuffer)
	omsg, _ := websocket.NewPreparedMessage(websocket.BinaryMessage, m.Buf())
	return omsg
}

func NewSendPosBuffer(id uuid.UUID) []byte {
	buf := make([]byte, UserPositionsMessageSize)
	copy(buf[:MsgUUIDTypeSize], utils.BinID(id))
	return buf
}
