package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/arbitrum_nova_adapter"
)

func main() {
	cfg, err := config.GetConfig()
	//cfg.Arbitrum.MOMTokenAddress = "0x6ab6d61428fde76768d7b45d8bfeec19c6ef91a8" ////https://nova.arbiscan.io/token/0x6ab6d61428fde76768d7b45d8bfeec19c6ef91a8
	cfg.Arbitrum.MOMTokenAddress = "0x3c2e532a334149d6a2e43523f2427e2fa187c5f0"
	cfg.Arbitrum.RPCURL = "https://nova.arbitrum.io/rpc"
	cfg.Arbitrum.BlockchainID = "42170"
	if err != nil {
		log.Fatal(err)
	}
	a := arbitrum_nova_adapter.NewArbitrumNovaAdapter(cfg)

	a.Run()
	n, err := a.GetLastBlockNumber()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Block: %+v \n", n)

	contracts := []common.Address{
		common.HexToAddress(cfg.Arbitrum.MOMTokenAddress),
	}
	_ = contracts
	//diffs, err := a.GetLogs(0, 12, contracts)
	logs, err := a.GetLogsRecursively(0, int64(n), contracts, 0)
	if err != nil {
		log.Println(err)
	}
	log.Println("LOGS:" + strconv.Itoa(len(logs)))

	for _, log := range logs {
		switch log.(type) {
		case *harvester.TransferERC20Log:
			l := log.(*harvester.TransferERC20Log)
			fmt.Printf("%s %s %s \n", l.From, l.To, l.Value)
			//fmt.Println(log.(*harvester.TransferERC20Log).Value)
		}
	}
	log.Println("LOGS:" + strconv.Itoa(len(logs)))

	for _, log := range logs {
		switch log.(type) {
		case *harvester.StakeLog:
			l := log.(*harvester.StakeLog)
			fmt.Printf("  stake: %s %s %s %d %s %s \n", l.TxHash, l.UserWallet, l.OdysseyID, l.TokenType, l.AmountStaked, l.TotalStaked)
			//fmt.Println(log.(*harvester.TransferERC20Log).Value)
		case *harvester.UnstakeLog:
			l := log.(*harvester.UnstakeLog)
			fmt.Printf("unstake: %s %s %s %s \n", l.UserWallet, l.OdysseyID, l.AmountUnstaked, l.TotalStaked)
			//fmt.Println(log.(*harvester.TransferERC20Log).Value)
		}
	}

	time.Sleep(time.Second * 3000)
}
