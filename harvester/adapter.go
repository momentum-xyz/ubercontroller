package harvester

import "math/big"

type BCAdapter interface {
	GetLastBlockNumber() (uint64, error)
	GetBalance(wallet string, contract string, blockNumber uint64) (*big.Int, error)
}
