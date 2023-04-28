package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/arbitrum_nova_adapter"
)

func main() {
	cfg := config.GetConfig()
	a := arbitrum_nova_adapter.NewArbitrumNovaAdapter(cfg)

	a.Run()

	n, err := a.GetLastBlockNumber()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Block: %+v \n", n)

	var l harvester.AdapterListener
	l = func(blockNumber uint64, diffs []*harvester.BCDiff, stakes []*harvester.BCStake) {
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
		common.HexToAddress("0x938c38D417fD1b0a29EA1722C84Ad16fF5dD89c3"), //staking t
	}
	_ = contracts
	//diffs, err := a.GetLogs(0, 12, contracts)
	diffs, stakes, err := a.GetLogs(0, 7247196, []common.Address{})

	fmt.Println("tokens ---")
	for k, v := range diffs {
		fmt.Printf("%+v %+v %+v %+v  \n", k, v.Token, v.To, v.Amount.String())
	}

	fmt.Println("stakes ---")
	for k, v := range stakes {
		fmt.Printf("%+v %+v %+v %+v %+v  \n", k, v.From, v.OdysseyID.String(), v.Amount.String(), v.TotalAmount.String())
	}

	time.Sleep(time.Second * 300)
}
