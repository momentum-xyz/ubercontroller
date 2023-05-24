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
	GetBalance(wallet *common.Address, contract *common.Address, blockNumber uint64) (*big.Int, error)
	GetNFTBalance(block int64, wallet *common.Address, nftContract *common.Address) ([]umid.UMID, error)
	GetStakeBalance(block int64, wallet *common.Address, nftContract *common.Address) (map[umid.UMID]*[3]*big.Int, error)
	GetLogs(fromBlock, toBlock int64, addresses []common.Address) ([]any, error)
	RegisterNewBlockListener(f AdapterListener)
	Run()
	GetInfo() (umid umid.UMID, name BCType, rpcURL string)
}

type IHarvester2 interface {
	RegisterAdapter(adapter Adapter)

	AddWallet(bcType BCType, wallet Address) error
	RemoveWallet(bcType BCType, wallet Address) error

	AddNFTContract(bcType BCType, contract Address) error
	RemoveNFTContract(bcType BCType, contract Address) error

	AddTokenContract(bcType BCType, contract Address) error
	RemoveTokenContract(bcType BCType, contract Address) error

	AddTokenListener(bcType BCType, contract Address, listener TokenListener) error
	AddNFTListener(bcType BCType, contract Address, listener NFTListener) error
	AddStakeListener(bcType BCType, contract Address, listener StakeListener) error
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

type TransferERC20Log struct {
	From     common.Address
	To       common.Address
	Value    *big.Int
	Contract common.Address
}

/**
 *
 * @param user User address
 * @param odyssey Odyssey ID that's being staked
 * @param amount_staked Amount being staked
 * @param token Token used (MOM or DAD)
 * @param total_staked Total being staked by the user on the Odyssey
 */
//event Stake(address user, bytes16 odyssey, uint256 amount_staked, Token token, uint256 total_staked);
type StakeLog struct {
	TxHash       string
	UserWallet   common.Address
	OdysseyID    umid.UMID
	AmountStaked *big.Int
	TokenType    uint8
	TotalStaked  *big.Int
}

/**
 *
 * @param user User address
 * @param odyssey Odyssey ID that's being unstaked
 * @param amount_unstaked Amount unstaked
 * @param token Token used (MOM or DAD)
 * @param total_staked Total remained staked by the user on that Odyssey
 */
//event Unstake(address user, bytes16 odyssey, uint256 amount_unstaked, Token token, uint256 total_staked);
type UnstakeLog struct {
	UserWallet     common.Address
	OdysseyID      umid.UMID
	AmountUnstaked *big.Int
	TokenType      uint8
	TotalStaked    *big.Int
}

/**
 *
 * @param user User address
 * @param odyssey_from Odyssey ID that the user is removing stake
 * @param odyssey_to Odyssey ID that the user is staking into
 * @param amount Amount that's being restaked
 * @param token Token used (MOM or DAD)
 * @param total_staked_from Total amount of tokens that remains staked on the `odyssey_from`
 * @param total_staked_to Total amount of tokens staked on `odyssey_to`
 */
//event Restake(address user,
//bytes16 odyssey_from,
//bytes16 odyssey_to,
//uint256 amount,
//Token token,
//uint256 total_staked_from,
//uint256 total_staked_to);
type RestakeLog struct {
	UserWallet        common.Address
	FromOdysseyID     umid.UMID
	ToOdysseyID       umid.UMID
	Amount            *big.Int
	TokenType         uint8
	TotalStakedToFrom *big.Int
	TotalStakedToTo   *big.Int
}

/**
 * @dev Emitted when `tokenId` token is transferred from `from` to `to`.
 */
//event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
type TransferNFTLog struct {
	From     common.Address
	To       common.Address
	TokenID  umid.UMID
	Contract common.Address
}
