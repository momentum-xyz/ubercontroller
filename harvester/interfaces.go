package harvester

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type BCBlock struct {
	Hash   string
	Number uint64
}

type BCDiff struct {
	From   string
	To     string
	Token  string
	Amount *big.Int
}

type UpdateEvent struct {
	Wallet   string
	Contract string
	Amount   *big.Int
}

type AdapterListener func(blockNumber uint64, diffs []*BCDiff)

type Adapter interface {
	GetLastBlockNumber() (uint64, error)
	GetBalance(wallet string, contract string, blockNumber uint64) (*big.Int, error)
	GetTransferLogs(fromBlock, toBlock int64, addresses []common.Address) ([]*BCDiff, error)
	RegisterNewBlockListener(f AdapterListener)
	Run()
	GetInfo() (umid umid.UMID, name string, rpcURL string)
}

type BCType string

const Ethereum string = "ethereum"
const Polkadot string = "polkadot"

type Event string

const NewBlock Event = "new_block"
const BalanceChange Event = "balance_change"

type IHarvester interface {
	RegisterAdapter(bcAdapter Adapter) error
	OnBalanceChange()
	Subscribe(bcType string, eventName Event, callback Callback)
	Unsubscribe(bcType string, eventName Event, callback Callback)
	SubscribeForWallet(bcType string, wallet, callback Callback)
	SubscribeForWalletAndContract(bcType string, wallet string, contract string, callback Callback) error
}
