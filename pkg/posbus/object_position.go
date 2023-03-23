package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type ObjectPosition struct {
	ID        umid.UMID             `json:"id"`
	Transform cmath.ObjectTransform `json:"object_transform"`
}

func init() {
	registerMessage(ObjectPosition{})
}

func (g *ObjectPosition) GetType() MsgType {
	return 0xEA6DA4B4
}
