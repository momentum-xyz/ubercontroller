package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

func NewSendPositionMsg(t cmath.UserTransform) *Message {
	obj := NewPreallocatedMessage(TypeSendPosition, 6*cmath.Float32Bytes)
	t.CopyToBuffer(obj.Msg())
	return obj
}
