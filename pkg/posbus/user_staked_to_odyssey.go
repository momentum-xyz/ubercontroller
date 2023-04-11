package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type UserStakedToOdyssey struct {
	TransactionHash [32]byte  `json:"transaction_hash"`
	ObjectID        umid.UMID `json:"object_id"`
	Amount          uint64    `json:"amount"`
	Comment         string    `json:"comment"`
}

func (r *UserStakedToOdyssey) GetType() MsgType {
	return 0x10DACABC
}

func init() {
	registerMessage(UserStakedToOdyssey{})
}
