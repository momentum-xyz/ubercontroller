package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

type MyTransform cmath.TransformNoScale

func init() {
	registerMessage(MyTransform{})
}

func (a *MyTransform) GetType() MsgType {
	return 0xF878C4BF
}
