package posbus

import (
	"encoding/binary"
	"math"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

func NewSendPositionMsg(pos cmath.Vec3, rotation cmath.Vec3) *Message {
	obj := NewPreallocatedMessage(TypeSendPosition, 6*MsgOnePosSize)
	copyUserPosition(obj.Msg(), pos, rotation)
	return obj
}

func copyUserPosition(b []byte, pos cmath.Vec3, rotation cmath.Vec3) {
	binary.LittleEndian.PutUint32(b, math.Float32bits(pos.X))
	binary.LittleEndian.PutUint32(b[MsgOnePosSize:], math.Float32bits(pos.Y))
	binary.LittleEndian.PutUint32(b[2*MsgOnePosSize:], math.Float32bits(pos.Z))
	binary.LittleEndian.PutUint32(b[3*MsgOnePosSize:], math.Float32bits(rotation.X))
	binary.LittleEndian.PutUint32(b[4*MsgOnePosSize:], math.Float32bits(rotation.Y))
	binary.LittleEndian.PutUint32(b[5*MsgOnePosSize:], math.Float32bits(rotation.Z))
}
