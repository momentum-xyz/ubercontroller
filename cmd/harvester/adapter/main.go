package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/ethereum_adapter"
)

func main() {
	a := ethereum_adapter.NewEthereumAdapter()
	a.Run()
	var l harvester.AdapterListener
	l = func(b *harvester.BCBlock, diffs []*harvester.BCDiff) {
		fmt.Printf("Listener: %+v %+v \n", b.Number, b.Hash)
		fmt.Printf("Diffs: %+v \n", len(diffs))
	}

	a.RegisterNewBlockListener(l)

	n, err := a.GetLastBlockNumber()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Last Block Number: %d \n", n)

	token := "0x85F17Cf997934a597031b2E18a9aB6ebD4B9f6a4"  // token smart contract address
	wallet := "0x0c3A3040075dd985F141800a1392a0Db81A09cAd" // user address
	b, err := a.GetBalance(wallet, token, n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %+v \n", b)

	contracts := []common.Address{
		common.HexToAddress("0x6de037ef9ad2725eb40118bb1702ebb27e4aeb24"),
		common.HexToAddress("0xa8c8cfb141a3bb59fea1e2ea6b79b5ecbcd7b6ca"),
	}
	diffs, err := a.GetTransferLogs(16882578, 16882878, contracts)

	for k, v := range diffs {
		fmt.Printf("%+v %+v %+v %+v  \n", k, v.Token, v.To, v.Amount.String())
	}

	time.Sleep(time.Second * 30)
}
