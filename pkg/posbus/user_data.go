package posbus

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type UserData struct {
	ID        umid.UMID              `json:"id"`
	Name      string                 `json:"name"`
	Avatar    string                 `json:"avatar"`
	Transform cmath.TransformNoScale `json:"transform"`
	IsGuest   bool                   `json:"is_guest"`
}

func init() {
	registerMessage(UserData{})
}

func (a *UserData) GetType() MsgType {
	return 0xF702EF5F
}
