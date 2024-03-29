package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/contracter"
	"github.com/momentum-xyz/ubercontroller/contracter/arbitrum_nova_adapter"
)

func main() {
	logger, _ := zap.NewProduction()
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	a := arbitrum_nova_adapter.NewArbitrumNovaAdapter(cfg, logger.Sugar())

	a.Run()

	n, err := a.GetLastBlockNumber()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Block: %+v \n", n)

	var l contracter.AdapterListener
	l = func(blockNumber uint64, diffs []*contracter.BCDiff, stakes []*contracter.BCStake) {
		fmt.Printf("Listener: %+v \n", blockNumber)
		for k, v := range diffs {
			fmt.Printf("%+v %+v %+v %+v\n", k, v.To, v.Token, v.Amount)
		}
		fmt.Printf("Diffs: %+v \n", len(diffs))
	}

	a.RegisterNewBlockListener(l)

	//token := cfg.Arbitrum.ArbitrumMOMTokenAddress // token smart contract address
	//wallet := "0x683642c22feDE752415D4793832Ab75EFdF6223c" // user address
	//wallet := "0x5ab4ef2f56001f2a21c821ef10b717d3c2dc91dd85fa823e9539e1178e5daa32" // user address
	//for i := 1; i < 100; i++ {
	//	b, err := a.GetBalance(wallet, token, n)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Balance: %+v \n", b)
	//}

	contracts := []common.Address{
		//common.HexToAddress("0x7F85fB7f42A0c0D40431cc0f7DFDf88be6495e67"),
		common.HexToAddress(cfg.Arbitrum.StakeAddress), //staking t
	}
	_ = contracts
	//diffs, err := a.GetLogs(0, 12, contracts)
	logs, err := a.GetLogs(0, 7247196, nil)

	for _, log := range logs {
		switch log.(type) {
		case *contracter.TransferERC20Log:
			l := log.(*contracter.TransferERC20Log)
			fmt.Printf("%s %s %s \n", l.From, l.To, l.Value)
			//fmt.Println(log.(*contracter.TransferERC20Log).Value)
		}
	}

	for _, log := range logs {
		switch log.(type) {
		case *contracter.StakeLog:
			l := log.(*contracter.StakeLog)
			fmt.Printf("  stake: %s %s %s %d %s %s \n", l.TxHash, l.UserWallet, l.OdysseyID, l.TokenType, l.AmountStaked, l.TotalStaked)
			//fmt.Println(log.(*contracter.TransferERC20Log).Value)
		case *contracter.UnstakeLog:
			l := log.(*contracter.UnstakeLog)
			fmt.Printf("unstake: %s %s %s %s \n", l.UserWallet, l.OdysseyID, l.AmountUnstaked, l.TotalStaked)
			//fmt.Println(log.(*contracter.TransferERC20Log).Value)
		}
	}

	time.Sleep(time.Second * 300)
}
