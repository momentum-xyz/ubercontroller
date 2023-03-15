package entry

import "github.com/momentum-xyz/ubercontroller/utils/mid"

type Wallet struct {
	WalletID     []byte `db:"wallet_id"`
	BlockchainID mid.ID `db:"blockchain_id"`
}
