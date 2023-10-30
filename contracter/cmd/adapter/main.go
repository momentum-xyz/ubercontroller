package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/contracter"
	"github.com/momentum-xyz/ubercontroller/contracter/ethereum_adapter"
)

func main() {
	a := ethereum_adapter.NewEthereumAdapter()
	a.Run()
	var l contracter.AdapterListener
	l = func(blockNumber uint64, diffs []*contracter.BCDiff, stakes []*contracter.BCStake) {
		fmt.Printf("Listener: %+v \n", blockNumber)
		for k, v := range diffs {
			fmt.Printf("%+v %+v %+v %+v\n", k, v.To, v.Token, v.Amount)
		}
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

		//common.HexToAddress("0xa8F737436eCb9cef0EDA3F8855b551f41d48C6fe"), // stakes geth
		//common.HexToAddress("0x6787F4aeaEe40cDC85B49DD885B28a85F83899C8"), // stakes ganache
		common.HexToAddress("0x02FE980E3aA97A42A402a69DA38e4804898033De"), // stakes ganache
	}

	// Last Block Number: 737179
	diffs, stakes, err := a.GetLogs(1, 31, contracts)

	fmt.Println(stakes)

	for k, v := range diffs {
		fmt.Printf("%+v %+v %+v %+v  \n", k, v.Token, v.To, v.Amount.String())
	}

	time.Sleep(time.Second * 300)
}
