package harvester3

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type AdapterListener func(blockNumber uint64)

type BCType string

type TransferERC20Log struct {
	Block    uint64
	From     common.Address
	To       common.Address
	Value    *big.Int
	Contract common.Address
}

type TransferNFTLog struct {
	Block    uint64
	From     common.Address
	To       common.Address
	TokenID  common.Hash
	Contract common.Address
}

type Adapter interface {
	GetLastBlockNumber() (uint64, error)
	GetTokenBalance(contract *common.Address, wallet *common.Address, blockNumber uint64) (*big.Int, uint64, error)
	GetNFTBalance(nftContract *common.Address, wallet *common.Address, block uint64) ([]common.Hash, error)
	//GetStakeBalance(block int64, wallet *common.Address, nftContract *common.Address) (map[umid.UMID]*[3]*big.Int, error)
	GetTokenLogs(fromBlock, toBlock uint64, addresses []common.Address) ([]any, error)
	GetNFTLogs(fromBlock, toBlock uint64, contracts []common.Address) ([]any, error)
	RegisterNewBlockListener(f AdapterListener)
	Run()
	GetInfo() (umid umid.UMID, name string, rpcURL string)
}
