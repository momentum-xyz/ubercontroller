package arbitrum_nova_adapter

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type ArbitrumNovaAdapter struct {
	listener    harvester.AdapterListener
	umid        umid.UMID
	rpcURL      string
	httpURL     string
	name        string
	client      *ethclient.Client
	rpcClient   *rpc.Client
	contractABI abi.ABI
}

func NewArbitrumNovaAdapter() *ArbitrumNovaAdapter {
	return &ArbitrumNovaAdapter{
		umid:    umid.MustParse("ccccaaaa-1111-2222-3333-222222222222"),
		rpcURL:  "wss://bcdev.antst.net:8548",
		httpURL: "https://bcdev.antst.net:8547",
		name:    "arbitrum_nova",
	}
}

func (a *ArbitrumNovaAdapter) GetLastBlockNumber() (uint64, error) {
	number, err := a.client.BlockNumber(context.TODO())
	return number, err
}

func (a *ArbitrumNovaAdapter) Run() {
	contractABI, err := abi.JSON(strings.NewReader(erc20abi))
	if err != nil {
		log.Fatal(err)
	}
	a.contractABI = contractABI

	a.client, err = ethclient.Dial(a.rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	a.rpcClient, err = rpc.DialHTTP(a.httpURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Ethereum Block Chain: " + a.rpcURL)

	ch := make(chan *types.Header)

	sub, err := a.client.SubscribeNewHead(context.Background(), ch)
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

				if a.listener != nil {
					a.onNewBlock(block)
				}
			}
		}
	}()
}

func (a *ArbitrumNovaAdapter) RegisterNewBlockListener(f harvester.AdapterListener) {
	a.listener = f
}

func (a *ArbitrumNovaAdapter) onNewBlock(b *harvester.BCBlock) {
	diffs, err := a.GetTransferLogs(int64(b.Number), int64(b.Number), []common.Address{})
	if err != nil {
		fmt.Println(err)
	}

	a.listener(b.Number, diffs)
}

func (a *ArbitrumNovaAdapter) GetBalance(wallet string, contract string, blockNumber uint64) (*big.Int, error) {
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

	return balance, nil
}

func (a *ArbitrumNovaAdapter) GetTransactionMessage(tx *types.Transaction) *core.Message {
	msg, err := core.TransactionToMessage(tx, types.LatestSignerForChainID(tx.ChainId()), nil)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

// refer https://github.com/ethereum/web3.py/blob/master/web3/contract.py#L435
func (a *ArbitrumNovaAdapter) DecodeTransactionInputData(contractABI *abi.ABI, data []byte) (string, map[string]any, error) {
	// The first 4 bytes of the txn represent the ID of the method in the ABI
	//fmt.Println(len(data))
	methodSigData := data[:4]
	method, err := contractABI.MethodById(methodSigData)
	if err != nil {
		err = errors.WithMessage(err, "failed to get ABI contract method by id")
		return "", nil, err
	}

	// parse the inputs to this method
	inputsSigData := data[4:]
	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		err = errors.WithMessage(err, "failed to unpack ABI contract method into map")
		return "", nil, err
	}
	//fmt.Printf("Method Name: %s\n", method.Name)
	//fmt.Printf("Method inputs: %v\n", MapToJson(inputsMap))

	return method.Name, inputsMap, nil
}

func (a *ArbitrumNovaAdapter) GetTransferLogs(fromBlock, toBlock int64, addresses []common.Address) ([]*harvester.BCDiff, error) {

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: addresses,
	}

	logs, err := a.client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	diffs := make([]*harvester.BCDiff, 0)

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			//fmt.Printf("Log Name: Transfer\n")

			var transferEvent harvester.BCDiff

			ev, err := a.contractABI.Unpack("Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(ev)

			transferEvent.Token = strings.ToLower(vLog.Address.Hex())
			// Hex and Un Hex here used to remove padding zeros
			transferEvent.From = strings.ToLower(common.HexToAddress(vLog.Topics[1].Hex()).Hex())
			transferEvent.To = strings.ToLower(common.HexToAddress(vLog.Topics[2].Hex()).Hex())
			if len(ev) > 0 {
				transferEvent.Amount = ev[0].(*big.Int)
			}

			//fmt.Printf("Contract: %s\n", transferEvent.Token)
			//fmt.Printf("From: %s\n", transferEvent.From)
			//fmt.Printf("To: %s\n", transferEvent.To)
			//fmt.Printf("Tokens: %s\n", transferEvent.Amount.String())
			diffs = append(diffs, &transferEvent)
		}
	}

	return diffs, nil
}

func (a *ArbitrumNovaAdapter) GetInfo() (umid umid.UMID, name string, rpcURL string) {
	return a.umid, a.name, a.rpcURL
}

const erc20abi = `[
    {
        "constant": true,
        "inputs": [],
        "name": "name",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_spender",
                "type": "address"
            },
            {
                "name": "_value",
                "type": "uint256"
            }
        ],
        "name": "approve",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "totalSupply",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_from",
                "type": "address"
            },
            {
                "name": "_to",
                "type": "address"
            },
            {
                "name": "_value",
                "type": "uint256"
            }
        ],
        "name": "transferFrom",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [
            {
                "name": "",
                "type": "uint8"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_owner",
                "type": "address"
            }
        ],
        "name": "balanceOf",
        "outputs": [
            {
                "name": "balance",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_to",
                "type": "address"
            },
            {
                "name": "_value",
                "type": "uint256"
            }
        ],
        "name": "transfer",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_owner",
                "type": "address"
            },
            {
                "name": "_spender",
                "type": "address"
            }
        ],
        "name": "allowance",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "payable": true,
        "stateMutability": "payable",
        "type": "fallback"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "owner",
                "type": "address"
            },
            {
                "indexed": true,
                "name": "spender",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "Approval",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "from",
                "type": "address"
            },
            {
                "indexed": true,
                "name": "to",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "Transfer",
        "type": "event"
    }
]`
