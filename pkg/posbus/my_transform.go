package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

type MyTransform cmath.UserTransform

func init() {
	registerMessage(&MyTransform{})
}

func (a *MyTransform) Type() MsgType {
	return TypeMyTransform
}
