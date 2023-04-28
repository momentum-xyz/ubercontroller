package entry

import (
	"time"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type Nft struct {
	WalletID     []byte    `db:"wallet_id"`
	BlockchainID umid.UMID `db:"blockchain_id"`
	ObjectID     umid.UMID `db:"object_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
