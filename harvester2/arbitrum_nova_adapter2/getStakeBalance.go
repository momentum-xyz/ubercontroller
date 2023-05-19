package arbitrum_nova_adapter2

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/momentum-xyz/ubercontroller/harvester2"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"
	"math/big"
)

func updateResultForStakeCase(result map[umid.UMID]*[3]*big.Int, odysseyID umid.UMID, amountStaked *big.Int, tokenType uint8) {
	if _, ok := result[odysseyID]; !ok {
		result[odysseyID] = &[3]*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0)}
	}
	result[odysseyID][0].Add(result[odysseyID][0], amountStaked)
	if tokenType == 0 {
		// MOM token
		result[odysseyID][1].Add(result[odysseyID][1], amountStaked)
	} else {
		result[odysseyID][2].Add(result[odysseyID][2], amountStaked)
	}
}

func updateResultForUnstakeCase(result map[umid.UMID]*[3]*big.Int, odysseyID umid.UMID, amountUnstaked *big.Int, tokenType uint8) {
	result[odysseyID][0].Sub(result[odysseyID][0], amountUnstaked)
	if tokenType == 0 {
		result[odysseyID][1].Sub(result[odysseyID][1], amountUnstaked)
	} else {
		result[odysseyID][2].Sub(result[odysseyID][2], amountUnstaked)
	}
}

func (a *ArbitrumNovaAdapter) GetStakeBalance(block int64, wallet *common.Address, nftContract *common.Address) (map[umid.UMID]*[3]*big.Int, error) {

	fmt.Println(a.contracts.NftABI.Events)
	stakeString := "Stake(address,uint256,uint256,uint8,uint256)"
	unstakeString := "Unstake(address,uint256,uint256,uint8,uint256)"
	restakeString := "Restake(address,uint256,uint256,uint256,uint8,uint256,uint256)"
	_ = stakeString
	_ = unstakeString
	_ = restakeString

	a.mu.Lock()
	defer a.mu.Unlock()

	if nftContract == nil {
		return nil, errors.New("nftContract is nil")
	}

	if wallet == nil {
		return nil, errors.New("nftContract is nil")
	}

	logs, err := a.GetLogs(0, block, []common.Address{*nftContract})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to GetLogs")
	}

	result := make(map[umid.UMID]*[3]*big.Int)

	// Filter all logs by given wallet
	for _, l := range logs {

		switch log := l.(type) {
		case *harvester2.StakeLog:
			if log.UserWallet != *wallet {
				continue
			}
			updateResultForStakeCase(result, log.OdysseyID, log.AmountStaked, log.TokenType)
		case *harvester2.UnstakeLog:
			if log.UserWallet != *wallet {
				continue
			}

			updateResultForUnstakeCase(result, log.OdysseyID, log.AmountUnstaked, log.TokenType)
		case *harvester2.RestakeLog:
			if log.UserWallet != *wallet {
				continue
			}
			updateResultForUnstakeCase(result, log.FromOdysseyID, log.Amount, log.TokenType)
			updateResultForStakeCase(result, log.ToOdysseyID, log.Amount, log.TokenType)
		}
	}

	return result, nil
}
