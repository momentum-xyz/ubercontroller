package entry

import (
	"time"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type Stake struct {
	WalletID     []byte    `db:"wallet_id"`
	BlockchainID umid.UMID `db:"blockchain_id"`
	ObjectID     umid.UMID `db:"object_id"`
	LastComment  string    `db:"last_comment"`
	Amount       *BigInt   `db:"amount"` //TODO should be big.Int
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	Count        uint8     `json:"count"`
}
