package entry

import (
	"time"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type NFT struct {
	WalletID     []byte    `db:"wallet_id"`
	BlockchainID umid.UMID `db:"blockchain_id"`
	ObjectID     umid.UMID `db:"object_id"`
	ContractID   []byte    `db:"contract_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
