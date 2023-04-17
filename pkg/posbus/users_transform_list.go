package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

const UserTransformMessageSize = MsgUUIDTypeSize + cmath.Float32Bytes*6

//func init() {
//	registerMessage(&UsersTransformList{})
//}


// UserTransform is a transform to apply to a specific user.
type UserTransform struct {
	ID        umid.UMID              `json:"id"`
	Transform cmath.TransformNoScale `json:"transform"`
}

type UsersTransformList struct {
	Value []UserTransform `json:"value"`
}

func (s *UserTransform) GetType() MsgType {
	// make it Message-compatible to auto-register
	return 0x3BC97EBB
}

func (s *UsersTransformList) GetType() MsgType {
	return 0x285954B8
}

func init() {
	registerMessage(UserTransform{})
	registerMessage(UsersTransformList{})
}
