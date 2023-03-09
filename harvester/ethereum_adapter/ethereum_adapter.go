package ethereum_adapter

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/harvester"
)

type EthereumAdapter struct {
	harv   harvester.BCAdapterAPI
	uuid   uuid.UUID
	rpcURL string
	name   string
	client *ethclient.Client
}

func NewEthereumAdapter(harv harvester.BCAdapterAPI) *EthereumAdapter {
	return &EthereumAdapter{
		harv:   harv,
		uuid:   uuid.MustParse("ccccaaaa-1111-2222-3333-111111111111"),
		rpcURL: "wss://eth.llamarpc.com",
		name:   "ethereum",
	}
}

func (a *EthereumAdapter) GetLastBlockNumber() (uint64, error) {
	number, err := a.client.BlockNumber(context.Background())
	return number, err
}

func (a *EthereumAdapter) Run() {

	client, err := ethclient.Dial(a.rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	a.client = client

	if err := a.harv.RegisterBCAdapter(a.uuid, a.name, a.rpcURL, a); err != nil {
		log.Fatal(err)
	}

	//client, err := ethclient.Dial("wss://rinkeby.infura.io/ws")
	//url := "ws://localhost:8546"
	//url := "wss://ethg.antst.net:8546"

	fmt.Println("Connected to Ethereum Block Chain: " + a.rpcURL)

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

			a.harv.OnNewBlock(harvester.Ethereum, block)
		}
	}
}
