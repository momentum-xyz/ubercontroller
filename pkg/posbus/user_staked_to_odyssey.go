package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type UserStakedToOdyssey struct {
	TransactionHash string    `json:"transaction_hash"`
	Wallet          string    `json:"wallet"`
	ObjectID        umid.UMID `json:"object_id"`
	Amount          string    `json:"amount"`
	Comment         string    `json:"comment"`
	Kind            int       `json:"kind"`
}

func (r *UserStakedToOdyssey) GetType() MsgType {
	return 0x10DACABC
}

func init() {
	registerMessage(UserStakedToOdyssey{})
}
