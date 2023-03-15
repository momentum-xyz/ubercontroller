package entry

import (
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"time"
)

type Blockchain struct {
	BlockchainID             mid.ID    `db:"blockchain_id"`
	LastProcessedBlockNumber uint64    `db:"last_processed_block_number"`
	BlockchainName           string    `db:"blockchain_name"`
	RPCURL                   string    `db:"rpc_url"`
	UpdatedAt                time.Time `db:"updated_at"`
}
