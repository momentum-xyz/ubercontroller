package main

import (
	"fmt"
	"log"
	"time"

	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/ethereum_adapter"
)

func main() {
	a := ethereum_adapter.NewEthereumAdapter()
	a.Run()
	var l harvester.AdapterListener
	l = func(b *harvester.BCBlock) {
		fmt.Printf("Listener: %+v %+v \n", b.Number, b.Hash)
	}

	a.RegisterNewBlockListener(l)

	n, err := a.GetLastBlockNumber()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Last Block Number: %d \n", n)

	token := "0x85F17Cf997934a597031b2E18a9aB6ebD4B9f6a4"
	wallet := "0x0c3A3040075dd985F141800a1392a0Db81A09cAd"
	b, err := a.GetBalance(wallet, token, n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %+v \n", b)

	time.Sleep(time.Second * 30)
}
