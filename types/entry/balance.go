package entry

import (
	"github.com/google/uuid"
)

type Balance struct {
	WalletID                 []byte    `db:"wallet_id"`
	ContractID               []byte    `db:"contract_id"`
	BlockchainID             uuid.UUID `db:"blockchain_id"`
	LastProcessedBlockNumber uint64    `db:"last_processed_block_number"`
	Balance                  uint64    `db:"balance"` //TODO should be big.Int
}
