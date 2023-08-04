package harvester3

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type AdapterListener func(blockNumber uint64)

type BCType string

type TransferERC20Log struct {
	Block    int64
	From     common.Address
	To       common.Address
	Value    *big.Int
	Contract common.Address
}

type Adapter interface {
	GetLastBlockNumber() (uint64, error)
	GetTokenBalance(contract *common.Address, wallet *common.Address, blockNumber uint64) (*big.Int, int64, error)
	//GetNFTBalance(block int64, wallet *common.Address, nftContract *common.Address) ([]umid.UMID, error)
	//GetStakeBalance(block int64, wallet *common.Address, nftContract *common.Address) (map[umid.UMID]*[3]*big.Int, error)
	GetTokenLogs(fromBlock, toBlock int64, addresses []common.Address) ([]any, error)
	RegisterNewBlockListener(f AdapterListener)
	Run()
	GetInfo() (umid umid.UMID, name string, rpcURL string)
}
