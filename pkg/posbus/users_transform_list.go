package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

const UserTransformMessageSize = MsgUUIDTypeSize + cmath.Float32Bytes*6

//func init() {
//	registerMessage(&UsersTransformList{})
//}

type UserTransform struct {
	ID        umid.UMID           `json:"id"`
	Transform cmath.UserTransform `json:"transform"`
}

type UsersTransformList struct {
	Value []UserTransform `json:"value"`
}

func (s *UsersTransformList) Type() MsgType {
	return TypeSetUsersTransforms
}

//type UserTransformBuffer struct {
//	maxUsers  int
//	posBuffer []byte
//	nUsers    int
//}

//func StartUserTransformBuffer(maxUsers int) *UserTransformBuffer {
//	obj := &UserTransformBuffer{maxUsers: maxUsers}
//	obj.posBuffer = make([]byte, MsgArrTypeSize+UserTransformMessageSize*maxUsers)
//	obj.nUsers = 0
//	return obj
//}
//
//func (utb *UserTransformBuffer) AddPosition(data []byte) {
//	start := MsgArrTypeSize + utb.nUsers*UserTransformMessageSize
//	copy(utb.posBuffer[start:], data)
//	utb.nUsers++
//}
//
//func (utb *UserTransformBuffer) Finalize() {
//	binary.LittleEndian.PutUint32(utb.posBuffer, uint32(utb.nUsers))
//}
//
//func (utb *UserTransformBuffer) Buf() []byte {
//	return utb.posBuffer[:MsgArrTypeSize+utb.nUsers*UserTransformMessageSize]
//}
//
//func (utb *UserTransformBuffer) NumUsers() int {
//	return utb.nUsers
//}

//func (utb *UserTransformBuffer) Decode() UsersTransformList {
//	t := make(
//		[]UserTransform, utb.nUsers,
//	)
//	start := MsgArrTypeSize
//	for i := 0; i < utb.nUsers; i++ {
//		copy(t[i].ID[:], utb.posBuffer[start:start+MsgUUIDTypeSize])
//		start += MsgUUIDTypeSize
//		t[i].Transform.Position = &cmath.Vec3{}
//		t[i].Transform.Rotation = &cmath.Vec3{}
//		t[i].Transform.CopyFromBuffer(utb.posBuffer[start:])
//		start += cmath.Float32Bytes * 6
//	}
//	return UsersTransformList{Value: t}
//}

//func BytesToUserTransformBuffer(buf []byte) *UserTransformBuffer {
//	nUsers := int(binary.LittleEndian.Uint32(buf))
//	if len(buf) < MsgArrTypeSize+UserTransformMessageSize*nUsers {
//		return nil
//	}
//	var b UserTransformBuffer
//
//	b.posBuffer = buf
//	b.maxUsers = nUsers
//	b.nUsers = nUsers
//
//	return &b
//}

//func NewSendTransformBuffer(id umid.UMID) []byte {
//	buf := make([]byte, UserTransformMessageSize)
//	copy(buf[:MsgUUIDTypeSize], utils.BinID(id))
//	return buf
//}
