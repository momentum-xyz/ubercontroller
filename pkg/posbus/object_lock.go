package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type LockObject struct {
	ID    umid.UMID `json:"id"`
	State uint32    `json:"state"`
}

type LockObjectResponse struct {
	ID        umid.UMID `json:"id"`
	Result    uint32    `json:"result"`
	LockOwner umid.UMID `json:"lock_owner"`
}

func init() {
	registerMessage(LockObject{})
	registerMessage(LockObjectResponse{})
}

func (l *LockObject) GetType() MsgType {
	return 0xA7DE9F59
}

func (l *LockObjectResponse) GetType() MsgType {
	return 0x0924668C
}
