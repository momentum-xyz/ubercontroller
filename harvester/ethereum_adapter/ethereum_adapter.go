package ethereum_adapter

import (
	"context"
	"encoding/json"
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

type EthereumAdapter struct {
	listener  harvester.AdapterListener
	umid      umid.UMID
	rpcURL    string
	httpURL   string
	name      string
	client    *ethclient.Client
	rpcClient *rpc.Client
}

func NewEthereumAdapter() *EthereumAdapter {
	return &EthereumAdapter{
		umid:   umid.MustParse("ccccaaaa-1111-2222-3333-111111111111"),
		rpcURL: "wss://eth.llamarpc.com",
		//rpcURL: "wss://ethereum-mainnet-rpc.allthatnode.com",
		httpURL: "https://eth.llamarpc.com",
		//httpURL: "https://ethereum-mainnet-rpc.allthatnode.com",
		name: "ethereum",
	}
}

func (a *EthereumAdapter) GetInfo() (umid umid.UMID, name string, rpcURL string) {
	return a.umid, a.name, a.rpcURL
}

func (a *EthereumAdapter) RegisterNewBlockListener(f harvester.AdapterListener) {
	a.listener = f
}

func (a *EthereumAdapter) GetLastBlockNumber() (uint64, error) {
	number, err := a.client.BlockNumber(context.Background())
	return number, err
}

func (a *EthereumAdapter) GetTransferLogs(fromBlock, toBlock int64, addresses []common.Address) ([]*harvester.BCDiff, error) {

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: addresses,
	}

	logs, err := a.client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractABI, err := abi.JSON(strings.NewReader(erc20abi))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	diffs := make([]*harvester.BCDiff, 0)

	for _, vLog := range logs {
		//fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		//fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			//fmt.Printf("Log Name: Transfer\n")

			var transferEvent harvester.BCDiff

			ev, err := contractABI.Unpack("Transfer", vLog.Data)
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

				fmt.Println(vLog.Number)
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

func (a *EthereumAdapter) onNewBlock(b *harvester.BCBlock) {
	block, err := a.client.BlockByNumber(context.TODO(), big.NewInt(int64(b.Number)))
	if err != nil {
		err = errors.WithMessage(err, "failed to get block by number")
		fmt.Println(err)
	}

	diffs := make([]*harvester.BCDiff, 0)
	for _, tx := range block.Transactions() {
		//fmt.Println(tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		//fmt.Println(tx.Value().String())    // 10000000000000000
		//fmt.Println(tx.Gas())               // 105000
		//fmt.Println(tx.GasPrice().Uint64()) // 102000000000
		//fmt.Println(tx.Nonce())             // 110644
		//fmt.Println(tx.Data())              // []
		//fmt.Println(tx.To().Hex())          // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e

		//fmt.Println(tx.Hash().Hex())

		contractABI, err := abi.JSON(strings.NewReader(erc20abi))
		if err != nil {
			log.Fatal(err)
		}

		if len(tx.Data()) < 4 {
			continue
		}

		methodName, methodInput, err := a.DecodeTransactionInputData(&contractABI, tx.Data())
		if err != nil {
			//fmt.Println(err)
			//log.Fatal(err)
		}

		if methodName == "transfer" {
			//fmt.Println(methodName)
			//fmt.Println(MapToJson(methodInput))
			//fmt.Printf("From: %s\n", a.GetTransactionMessage(tx).From().Hex()) // from field is not inside of transation
			diff := &harvester.BCDiff{}
			diff.From = strings.ToLower(a.GetTransactionMessage(tx).From.Hex())

			diff.To = strings.ToLower(methodInput["_to"].(common.Address).Hex())
			diff.Token = strings.ToLower(tx.To().Hex())
			diff.Amount = methodInput["_value"].(*big.Int)
			diffs = append(diffs, diff)
		}

		if methodName == "transferFrom" {
			//fmt.Println(methodName)
			//fmt.Println(MapToJson(methodInput))
			diff := &harvester.BCDiff{}
			diff.From = strings.ToLower(methodInput["_from"].(common.Address).Hex())
			diff.To = strings.ToLower(methodInput["_to"].(common.Address).Hex())
			diff.Token = strings.ToLower(tx.To().Hex())
			diff.Amount = methodInput["_value"].(*big.Int)
			diffs = append(diffs, diff)
		}

		//if tx.Hash().Hex() == "0x3dc59fec84347cf1929b81c0fce68a3511b6660f22d898edafdf7f8247175dbb" {
		//	//fmt.Println(tx)
		//	fmt.Println(tx.To().Hex()) // token contract
		//	fmt.Println(tx.Value().String())
		//	fmt.Println(tx.Data())
		//}

		// TODO Check that tx success
		//receipt, err := a.client.TransactionReceipt(context.Background(), tx.Hash())
		//if err != nil {
		//	log.Fatal(err)
		//}

		//fmt.Println(tx.Hash())
		//fmt.Println(receipt.Status) // 1
	}

	amount := big.NewInt(0)
	//amount.SetString("33190774000000000000000", 10)
	amount.SetString("1", 10)

	mockDiffs := []*harvester.BCDiff{
		&harvester.BCDiff{
			From:   "0x2813fd17ea95b2655a7228383c5236e31090419e",
			To:     "0x3f363b4e038a6e43ce8321c50f3efbf460196d4b",
			Token:  "0xdefa4e8a7bcba345f687a2f1456f5edd9ce97202",
			Amount: amount,
		},
	}

	a.listener(b.Number, mockDiffs)
	//a.listener(b, diffs)
}

// refer
// https://github.com/ethereum/web3.py/blob/master/web3/contract.py#L435
func (a *EthereumAdapter) DecodeTransactionInputData(contractABI *abi.ABI, data []byte) (string, map[string]any, error) {
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

func MapToJson(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func (a *EthereumAdapter) GetTransactionMessage(tx *types.Transaction) *core.Message {
	msg, err := core.TransactionToMessage(tx, types.LatestSignerForChainID(tx.ChainId()), nil)
	if err != nil {
		log.Fatal(err)
	}
	return msg
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
