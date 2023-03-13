package ethereum_adapter

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/harvester"
)

type EthereumAdapter struct {
	harv      harvester.BCAdapterAPI
	uuid      uuid.UUID
	rpcURL    string
	httpURL   string
	name      string
	client    *ethclient.Client
	rpcClient *rpc.Client
}

func NewEthereumAdapter(harv harvester.BCAdapterAPI) *EthereumAdapter {
	return &EthereumAdapter{
		harv:    harv,
		uuid:    uuid.MustParse("ccccaaaa-1111-2222-3333-111111111111"),
		rpcURL:  "wss://eth.llamarpc.com",
		httpURL: "https://eth.llamarpc.com",
		name:    "ethereum",
	}
}

func (a *EthereumAdapter) GetLastBlockNumber() (uint64, error) {
	number, err := a.client.BlockNumber(context.Background())
	return number, err
}

func (a *EthereumAdapter) GetBalance(wallet string, contract string, blockNumber uint64) (*big.Int, error) {
	type request struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}

	// "0x70a08231" - crypto.Keccak256Hash([]byte("balanceOf(address)")).String()[0:10]
	data := "0x70a08231" + fmt.Sprintf("%064s", wallet[2:]) // %064s means that the string is padded with 0 to 64 bytes
	req := request{contract, data}

	var resp string
	n := hexutil.EncodeUint64(blockNumber)
	if err := a.rpcClient.Call(&resp, "eth_call", req, n); err != nil {
		log.Fatal(err)
		return nil, errors.WithMessage(err, "failed to make RPC call to ethereum:")
	}

	// remove leading zero of resp
	t := strings.TrimLeft(resp[2:], "0")
	if t == "" {
		t = "0"
	}
	s := "0x" + t
	balance, err := hexutil.DecodeBig(s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance)
	return balance, nil
}

func (a *EthereumAdapter) Run() {

	client, err := ethclient.Dial(a.rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	a.rpcClient, err = rpc.DialHTTP(a.httpURL)
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

	go func() {
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
	}()

}
