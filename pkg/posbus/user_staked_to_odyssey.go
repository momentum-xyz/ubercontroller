package posbus

type UserStakedToOdyssey struct {
	TransactionHash string `json:"transaction_hash"`
	ObjectID        string `json:"object_id"`
	Amount          string `json:"amount"`
	Comment         string `json:"comment"`
}

func (r *UserStakedToOdyssey) GetType() MsgType {
	return 0x10DACABC
}

func init() {
	registerMessage(UserStakedToOdyssey{})
}
