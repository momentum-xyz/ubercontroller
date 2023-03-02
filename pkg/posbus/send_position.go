package posbus

import (
	"encoding/binary"
	"math"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

type SendPosition struct {
	*Message
}

func NewSendPositionMsg(pos cmath.Vec3, rotation cmath.Vec3, scale cmath.Vec3) *SendPosition {
	obj := SendPosition{Message: NewMessage(SendPositionType)}
	obj.makeBuffer(9 * MsgOnePosSize)
	obj.SetPosition(pos, rotation, scale)
	return &obj
}

func (m *SendPosition) SetPosition(pos cmath.Vec3, rotation cmath.Vec3, scale cmath.Vec3) {
	binary.LittleEndian.PutUint32(m.Msg(), math.Float32bits(pos.X))
	binary.LittleEndian.PutUint32(m.Msg()[MsgOnePosSize:], math.Float32bits(pos.Y))
	binary.LittleEndian.PutUint32(m.Msg()[2*MsgOnePosSize:], math.Float32bits(pos.Z))
	binary.LittleEndian.PutUint32(m.Msg()[3*MsgOnePosSize:], math.Float32bits(rotation.X))
	binary.LittleEndian.PutUint32(m.Msg()[4*MsgOnePosSize:], math.Float32bits(rotation.Y))
	binary.LittleEndian.PutUint32(m.Msg()[5*MsgOnePosSize:], math.Float32bits(rotation.Z))
	binary.LittleEndian.PutUint32(m.Msg()[6*MsgOnePosSize:], math.Float32bits(scale.X))
	binary.LittleEndian.PutUint32(m.Msg()[7*MsgOnePosSize:], math.Float32bits(scale.Y))
	binary.LittleEndian.PutUint32(m.Msg()[8*MsgOnePosSize:], math.Float32bits(scale.Z))
}
