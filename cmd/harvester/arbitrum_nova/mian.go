package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/arbitrum_nova_adapter"
)

func main() {
	a := arbitrum_nova_adapter.NewArbitrumNovaAdapter()

	a.Run()

	n, err := a.GetLastBlockNumber()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Block: %+v \n", n)

	var l harvester.AdapterListener
	l = func(blockNumber uint64, diffs []*harvester.BCDiff) {
		fmt.Printf("Listener: %+v \n", blockNumber)
		for k, v := range diffs {
			fmt.Printf("%+v %+v %+v %+v\n", k, v.To, v.Token, v.Amount)
		}
		fmt.Printf("Diffs: %+v \n", len(diffs))
	}

	a.RegisterNewBlockListener(l)

	token := "0x7F85fB7f42A0c0D40431cc0f7DFDf88be6495e67"  // token smart contract address
	wallet := "0xe2148ee53c0755215df69b2616e552154edc584f" // user address
	b, err := a.GetBalance(wallet, token, n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %+v \n", b)

	contracts := []common.Address{
		common.HexToAddress("0x7F85fB7f42A0c0D40431cc0f7DFDf88be6495e67"),
	}
	_ = contracts
	//diffs, err := a.GetTransferLogs(0, 12, contracts)
	diffs, err := a.GetTransferLogs(0, 12, []common.Address{})

	for k, v := range diffs {
		fmt.Printf("%+v %+v %+v %+v  \n", k, v.Token, v.To, v.Amount.String())
	}

	time.Sleep(time.Second * 300)
}
