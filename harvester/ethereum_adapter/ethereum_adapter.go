package ethereum_adapter

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/momentum-xyz/ubercontroller/harvester"
)

type EthereumAdapter struct {
	harv harvester.BCAdapterAPI
}

func NewEthereumAdapter(harvester harvester.BCAdapterAPI) *EthereumAdapter {
	return &EthereumAdapter{
		harv: harvester,
	}
}

func (ea *EthereumAdapter) Run() {
	//client, err := ethclient.Dial("wss://rinkeby.infura.io/ws")
	url := "ws://localhost:8546"
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Ethereum Block Chain: " + url)

	ch := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-ch:

			//fmt.Println(vLog.Number)
			//fmt.Println(vLog.ReceiptHash)
			//fmt.Println(vLog.ParentHash)
			//fmt.Println(vLog.Root)
			//fmt.Println(vLog.TxHash)
			//fmt.Println(vLog.Hash())

			block := &harvester.BCBlock{
				Hash: vLog.Hash().String(),
			}

			if vLog.Number != nil {
				block.Number = vLog.Number.Uint64()
			}

			ea.harv.OnNewBlock(harvester.Ethereum, block)
		}
	}
}
