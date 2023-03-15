package harvester

import (
	"math/big"

	"github.com/google/uuid"
)

type BCBlock struct {
	Hash   string
	Number uint64
}

type AdapterListener func(block *BCBlock)

type Adapter interface {
	GetLastBlockNumber() (uint64, error)
	GetBalance(wallet string, contract string, blockNumber uint64) (*big.Int, error)
	RegisterNewBlockListener(f AdapterListener)
	Run()
	GetInfo() (uuid uuid.UUID, name string, rpcURL string)
}

type BCType string

const Ethereum string = "ethereum"
const Polkadot string = "polkadot"

type Event string

const NewBlock Event = "new_block"
const BalanceChange Event = "balance_change"

type IHarvester interface {
	RegisterAdapter(uuid uuid.UUID, bcType string, rpcURL string, bcAdapter Adapter) error
	OnBalanceChange()
	Subscribe(bcType string, eventName Event, callback Callback)
	Unsubscribe(bcType string, eventName Event, callback Callback)
	SubscribeForWallet(bcType string, wallet, callback Callback)
	SubscribeForWalletAndContract(bcType string, wallet []byte, contract []byte, callback Callback) error
}
