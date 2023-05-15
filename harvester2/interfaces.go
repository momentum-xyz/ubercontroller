package harvester2

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"math/big"
)

type Address common.Address

type AdapterListener func(blockNumber uint64)

type Adapter interface {
	GetLastBlockNumber() (uint64, error)
	GetBalance(wallet string, contract string, blockNumber uint64) (*big.Int, error)
	GetLogs(fromBlock, toBlock int64, addresses []common.Address) ([]any, error)
	RegisterNewBlockListener(f AdapterListener)
	Run()
	GetInfo() (umid umid.UMID, name BCType, rpcURL string)
}

type IHarvester2 interface {
	RegisterAdapter(adapter Adapter)

	AddWallet(bcType BCType, wallet *Address) error
	RemoveWallet(bcType BCType, wallet *Address) error

	AddNFTContract(bcType BCType, contract *Address) error
	RemoveNFTContract(bcType BCType, contract *Address) error

	AddTokenContract(bcType BCType, contract *Address) error
	RemoveTokenContract(bcType BCType, contract *Address) error

	AddTokenListener(bcType BCType, contract *Address, listener TokenListener) error
	AddNFTListener(bcType BCType, contract *Address, listener NFTListener) error
	AddStakeListener(bcType BCType, contract *Address, listener StakeListener) error
}

type TokenListener func(events []TokenData)
type NFTListener func(events []*NFTData)
type StakeListener func(events []*StakeData)

type TokenData struct {
	Wallet      *Address
	Contract    *Address
	TotalAmount *big.Int
}

type NFTData struct {
	Wallet   *Address
	Contract *Address
	TokenIDs []umid.UMID
}

type StakeData struct {
	Wallet    *Address
	Contract  *Address
	OdysseyID *umid.UMID
	Stake     *Stake
}

type Stake struct {
	TotalAmount    *big.Int
	TotalDADAmount *big.Int
	TotalMOMAmount *big.Int
}

type UpdateEvent struct {
	Wallet   string
	Contract string
	Amount   *big.Int
}

type StakeEvent struct {
	TxHash    string
	Wallet    string
	OdysseyID umid.UMID
	Amount    *big.Int
}

type NftEvent struct {
	From      string
	To        string
	OdysseyID umid.UMID
	Contract  string
}
