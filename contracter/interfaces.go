package contracter

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type TransferERC20Log struct {
	From     string
	To       string
	Value    *big.Int
	Contract string
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
	LogIndex     uint
	TxHash       string
	UserWallet   string
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
	LogIndex       uint
	TxHash         string
	UserWallet     string
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
	UserWallet        string
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
	From     string
	To       string
	TokenID  umid.UMID
	Contract string
}

type TransferOdysseyLog struct {
	FromNodeID umid.UMID
	ToNodeID   umid.UMID
	OdysseyID  umid.UMID
}

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

type BCStake struct {
	From        string
	OdysseyID   umid.UMID
	TokenType   uint8 //0-MOM; 1-DAD
	Amount      *big.Int
	TotalAmount *big.Int
}

type UpdateEvent struct {
	Wallet   string
	Contract string
	Amount   *big.Int
}

type StakeEvent struct {
	TxHash       string
	LogIndex     string
	Wallet       string
	Kind         uint8
	OdysseyID    umid.UMID
	Amount       *big.Int
	ActivityType string
	//CreatedAt    time.Time
}

type NftEvent struct {
	From      string
	To        string
	OdysseyID umid.UMID
	Contract  string
}

type TransferOdysseyEvent struct {
	FromNodeID umid.UMID
	ToNodeID   umid.UMID
	OdysseyID  umid.UMID
}

type AdapterListener func(blockNumber uint64, diffs []*BCDiff, stakes []*BCStake)

type Adapter interface {
	GetLastBlockNumber() (uint64, error)
	GetBalance(wallet string, contract string, blockNumber uint64) (*big.Int, error)
	GetLogs(fromBlock, toBlock int64, addresses []common.Address) ([]any, error)
	GetLogsRecursively(fromBlock, toBlock int64, addresses []common.Address, level int) ([]any, error)
	RegisterNewBlockListener(f AdapterListener)
	Run()
	GetInfo() (umid umid.UMID, name string, rpcURL string)
}

type BCType string

type Event string
