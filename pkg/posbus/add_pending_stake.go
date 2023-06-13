package posbus

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type AddPendingStake struct {
	TransactionID umid.UMID    `json:"transaction_id"`
	OdysseyId     umid.UMID    `json:"odyssey_id"`
	Wallet        PBEthAddress `json:"wallet"`
	Amount        PBUint256    `json:"amount"`
	Comment       string       `json:"comment"`
	Kind          int          `json:"kind"`
}

func init() {
	registerMessage(AddPendingStake{})
}

func (g *AddPendingStake) GetType() MsgType {
	return 0xF020D682
}
