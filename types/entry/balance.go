package entry

import "github.com/momentum-xyz/ubercontroller/utils/mid"

type Balance struct {
	WalletID                 []byte `db:"wallet_id"`
	ContractID               []byte `db:"contract_id"`
	BlockchainID             mid.ID `db:"blockchain_id"`
	LastProcessedBlockNumber uint64 `db:"last_processed_block_number"`
	Balance                  uint64 `db:"balance"` //TODO should be big.Int
}
