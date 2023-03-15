package posbus

import (
	"encoding/binary"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
)

const UserTransformMessageSize = MsgUUIDTypeSize + cmath.Float32Bytes*6

type UserTransform struct {
	ID        mid.ID              `json:"mid"`
	Transform cmath.UserTransform `json:"transform"`
}

type SetUsersTransforms struct {
	Value []UserTransform `json:"value"`
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

func (utb *UserTransformBuffer) NumUsers() int {
	return utb.nUsers
}

func (utb *UserTransformBuffer) Decode() SetUsersTransforms {
	t := make(
		[]UserTransform, utb.nUsers,
	)
	start := MsgArrTypeSize
	for i := 0; i < utb.nUsers; i++ {
		copy(t[i].ID[:], utb.posBuffer[start:start+16])
		start += 16
		t[i].Transform.Position = &cmath.Vec3{}
		t[i].Transform.Rotation = &cmath.Vec3{}
		t[i].Transform.CopyFromBuffer(utb.posBuffer[start:])
		start += UserTransformMessageSize
	}
	return SetUsersTransforms{Value: t}
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

func NewSendTransformBuffer(id mid.ID) []byte {
	buf := make([]byte, UserTransformMessageSize)
	copy(buf[:MsgUUIDTypeSize], utils.BinID(id))
	return buf
}
