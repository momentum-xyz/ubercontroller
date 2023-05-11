package arbitrum_nova_adapter

import (
	"log"
	"strings"

	_ "embed"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/config"
)

//go:embed abi/token.json
var erc20abi string

//go:embed abi/staking.json
var stakeABI string

//go:embed abi/nft.json
var nftABI string

type Contracts struct {
	TokenABI abi.ABI
	StakeABI abi.ABI
	NftABI   abi.ABI

	momTokenAddress common.Address
	dadTokenAddress common.Address
	stakeAddress    common.Address
	nftAddress      common.Address

	AllAddresses []common.Address
}

func NewContracts(cfg *config.Arbitrum) *Contracts {

	tokenABI, err := abi.JSON(strings.NewReader(erc20abi))
	if err != nil {
		log.Fatal(err)
	}

	stakeABI, err := abi.JSON(strings.NewReader(stakeABI))
	if err != nil {
		log.Fatal(err)
	}

	nftABI, err := abi.JSON(strings.NewReader(nftABI))
	if err != nil {
		log.Fatal(err)
	}

	contracts := &Contracts{
		TokenABI:        tokenABI,
		StakeABI:        stakeABI,
		NftABI:          nftABI,
		momTokenAddress: common.HexToAddress(cfg.MOMTokenAddress),
		dadTokenAddress: common.HexToAddress(cfg.DADTokenAddress),
		stakeAddress:    common.HexToAddress(cfg.StakeContractAddress),
		nftAddress:      common.HexToAddress(cfg.NFTContractAddress),
	}

	allAddresses := make([]common.Address, 0)
	allAddresses = append(allAddresses,
		contracts.momTokenAddress, contracts.dadTokenAddress, contracts.stakeAddress, contracts.nftAddress)
	contracts.AllAddresses = allAddresses

	return contracts
}
