package entry

import (
	"time"

	"github.com/google/uuid"
)

type Blockchain struct {
	BlockchainID             uuid.UUID `db:"blockchain_id"`
	LastProcessedBlockNumber uint64    `db:"last_processed_block_number"`
	BlockchainName           string    `db:"blockchain_name"`
	RPCURL                   string    `db:"rpc_url"`
	UpdatedAt                time.Time `db:"updated_at"`
}
