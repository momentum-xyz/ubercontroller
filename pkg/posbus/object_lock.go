package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type LockObject struct {
	ID    umid.UMID `json:"id"`
	State uint32    `json:"state"`
}

type ObjectLockResult struct {
	ID        umid.UMID `json:"id"`
	Result    uint32    `json:"result"`
	LockOwner umid.UMID `json:"lock_owner"`
}

func init() {
	registerMessage(&LockObject{})
	registerMessage(&ObjectLockResult{})
}

func (l *LockObject) Type() MsgType {
	return 0xA7DE9F59
}

func (l *ObjectLockResult) Type() MsgType {
	return 0x0924668C
}
