package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// ObjectTransform is a transform to apply to a specific object.
type ObjectTransform struct {
	ID        umid.UMID       `json:"id"`
	Transform cmath.Transform `json:"object_transform"`
}

func init() {
	registerMessage(ObjectTransform{})
}

func (g *ObjectTransform) GetType() MsgType {
	return 0xEA6DA4B4
}
