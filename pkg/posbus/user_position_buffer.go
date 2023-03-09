package posbus

import (
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils"
)

const UserTransformMessageSize = MsgUUIDTypeSize + cmath.Float32Bytes*6

type UserPosition struct {
	ID        uuid.UUID
	Transform cmath.UserTransform
}

type UserTransformBuffer struct {
	maxUsers  int
	posBuffer []byte
	nUsers    int
}

func StartUserTransformBuffer(maxUsers int) *UserTransformBuffer {
	obj := &UserTransformBuffer{maxUsers: maxUsers}
	obj.posBuffer = make([]byte, MsgArrTypeSize+UserTransformMessageSize*maxUsers)
	obj.nUsers = 0
	return obj
}

func (utb *UserTransformBuffer) AddPosition(data []byte) {
	start := MsgArrTypeSize + utb.nUsers*UserTransformMessageSize
	copy(utb.posBuffer[start:], data)
	utb.nUsers++
}

func (utb *UserTransformBuffer) Finalize() {
	binary.LittleEndian.PutUint32(utb.posBuffer, uint32(utb.nUsers))
}

func (utb *UserTransformBuffer) Buf() []byte {
	return utb.posBuffer[:MsgArrTypeSize+utb.nUsers*UserTransformMessageSize]
}

func (utb *UserTransformBuffer) Decode() []struct {
	ID        uuid.UUID
	Transofrm cmath.UserTransform
} {
	t := make(
		[]struct {
			ID        uuid.UUID
			Transofrm cmath.UserTransform
		}, utb.nUsers,
	)
	start := MsgArrTypeSize
	for i := 0; i < utb.nUsers; i++ {
		//t[i].ID=uuid.UUID(utb.posBuffer[start:start+16])
		copy(t[i].ID[:], utb.posBuffer[start:start+16])
		start += 16
		t[i].Transofrm.Position = &cmath.Vec3{}
		t[i].Transofrm.Rotation = &cmath.Vec3{}
		t[i].Transofrm.CopyFromBuffer(utb.posBuffer[start:])
		start += UserTransformMessageSize
	}
	return t
}

func BytesToUserTransformBuffer(buf []byte) *UserTransformBuffer {
	nUsers := int(binary.LittleEndian.Uint32(buf))
	if len(buf) < MsgArrTypeSize+UserTransformMessageSize*nUsers {
		return nil
	}
	var b UserTransformBuffer

	b.posBuffer = buf
	b.maxUsers = nUsers
	b.nUsers = nUsers

	return &b
}

func NewSendTransformBuffer(id uuid.UUID) []byte {
	buf := make([]byte, UserTransformMessageSize)
	copy(buf[:MsgUUIDTypeSize], utils.BinID(id))
	return buf
}
