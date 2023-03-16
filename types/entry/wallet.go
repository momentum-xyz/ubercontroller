package entry

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type Wallet struct {
	WalletID     []byte    `db:"wallet_id"`
	BlockchainID umid.UMID `db:"blockchain_id"`
}
