package entry

import "github.com/google/uuid"

type Wallet struct {
	WalletID     []byte    `db:"wallet_id"`
	BlockchainID uuid.UUID `db:"blockchain_id"`
}
