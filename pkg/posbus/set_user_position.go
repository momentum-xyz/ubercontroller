package posbus

import (
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils"
)

const UserPositionsMessageSize = MsgUUIDTypeSize + cmath.Float32Bytes*6

type SetUserPositionsBuffer struct {
	maxUsers  int
	posBuffer []byte
	nUsers    uint32
}

func StartUserPositionsBuffer(maxUsers int) *SetUserPositionsBuffer {
	obj := &SetUserPositionsBuffer{maxUsers: maxUsers}
	obj.posBuffer = make([]byte, MsgArrTypeSize+UserPositionsMessageSize*maxUsers)
	obj.nUsers = 0
	return obj
}

func (m *SetUserPositionsBuffer) AddPosition(data []byte) {
	start := MsgArrTypeSize + m.nUsers*UserPositionsMessageSize
	copy(m.posBuffer[start:], data)
	m.nUsers++
}

func (m *SetUserPositionsBuffer) Finalize() {
	binary.LittleEndian.PutUint32(m.posBuffer, m.nUsers)
}

func (m *SetUserPositionsBuffer) Buf() []byte {
	return m.posBuffer[:MsgArrTypeSize+m.nUsers*UserPositionsMessageSize]
}

func NewSendPosBuffer(id uuid.UUID) []byte {
	buf := make([]byte, UserPositionsMessageSize)
	copy(buf[:MsgUUIDTypeSize], utils.BinID(id))
	return buf
}
