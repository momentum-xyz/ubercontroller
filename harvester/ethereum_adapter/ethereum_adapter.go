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
	listener         harvester.AdapterListener
	umid             umid.UMID
	rpcURL           string
	httpURL          string
	name             string
	client           *ethclient.Client
	rpcClient        *rpc.Client
	contractABI      abi.ABI
	stakeContractABI abi.ABI
}

func NewEthereumAdapter() *EthereumAdapter {
	return &EthereumAdapter{
		umid:    umid.MustParse("ccccaaaa-1111-2222-3333-111111111111"),
		rpcURL:  "wss://bcdev.antst.net:8546",
		httpURL: "https://bcdev.antst.net:8545",
		//rpcURL:  "ws://127.0.0.1:8545",
		//httpURL: "http://127.0.0.1:8545",
		//rpcURL:  "wss://eth.llamarpc.com",
		//httpURL: "https://eth.llamarpc.com",
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

func (a *EthereumAdapter) GetTransferLogs(fromBlock, toBlock int64, addresses []common.Address) ([]*harvester.BCDiff, []*harvester.BCStake, error) {

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: addresses,
	}

	logs, err := a.client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to filter log")
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	logStakeSigHash := crypto.Keccak256Hash([]byte(a.stakeContractABI.Events["Stake"].Sig))
	//fmt.Println(logStakeSigHash)
	//fmt.Println(a.stakeContractABI.Events["Stake"].ID.Hex())

	diffs := make([]*harvester.BCDiff, 0)
	stakes := make([]*harvester.BCStake, 0)

	//fmt.Println(a.contractABI.Events)

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			//fmt.Printf("Log Name: Transfer\n")

			var transferEvent harvester.BCDiff

			ev, err := a.contractABI.Unpack("Transfer", vLog.Data)
			if err != nil {
				return nil, nil, errors.WithMessage(err, "failed to unpack event from ABI")
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
		case logStakeSigHash.Hex():
			//fmt.Println("STAKE")
			//fmt.Println(vLog)

			ev, err := a.stakeContractABI.Unpack("Stake", vLog.Data)
			if err != nil {
				return nil, nil, errors.WithMessage(err, "failed to unpack event from ABI")
			}

			// Read and convert event params
			fromWallet := ev[0].(common.Address)

			arr := ev[1].([16]byte)
			odysseyID, err := umid.FromBytes(arr[:])
			if err != nil {
				return nil, nil, errors.WithMessage(err, "failed to parse umid from bytes")
			}

			amount := ev[2].(*big.Int)

			tokenType := ev[3].(uint8)

			stake := &harvester.BCStake{
				From:      fromWallet.Hex(),
				OdysseyID: odysseyID,
				TokenType: tokenType,
				Amount:    amount,
			}

			stakes = append(stakes, stake)

			//fmt.Printf("%+v %+v %+v %+v \n\n", fromWallet.String(), odysseyID.String(), amount, tokenType)
			//fmt.Println(ev)
		}
	}

	return diffs, stakes, nil
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
	contractABI, err := abi.JSON(strings.NewReader(erc20abi))
	if err != nil {
		log.Fatal(err)
	}
	a.contractABI = contractABI

	stakeContractABI, err := abi.JSON(strings.NewReader(stakeABI))
	if err != nil {
		log.Fatal(err)
	}
	a.stakeContractABI = stakeContractABI

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

func (a *EthereumAdapter) onNewBlock(b *harvester.BCBlock) {
	block, err := a.client.BlockByNumber(context.TODO(), big.NewInt(int64(b.Number)))
	if err != nil {
		err = errors.WithMessage(err, "failed to get block by number")
		fmt.Println(err)
	}

	diffs := make([]*harvester.BCDiff, 0)
	for _, tx := range block.Transactions() {
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
			//log.Fatal(err)
		}

		if methodName == "transfer" {
			diff := &harvester.BCDiff{}
			diff.From = strings.ToLower(a.GetTransactionMessage(tx).From.Hex())

			diff.To = strings.ToLower(methodInput["_to"].(common.Address).Hex())
			diff.Token = strings.ToLower(tx.To().Hex())
			diff.Amount = methodInput["_value"].(*big.Int)
			diffs = append(diffs, diff)
		}

		if methodName == "transferFrom" {
			diff := &harvester.BCDiff{}
			diff.From = strings.ToLower(methodInput["_from"].(common.Address).Hex())
			diff.To = strings.ToLower(methodInput["_to"].(common.Address).Hex())
			diff.Token = strings.ToLower(tx.To().Hex())
			diff.Amount = methodInput["_value"].(*big.Int)
			diffs = append(diffs, diff)
		}

		// TODO Check that tx success
		//receipt, err := a.client.TransactionReceipt(context.Background(), tx.Hash())
		//if err != nil {
		//	log.Fatal(err)
		//}

		//fmt.Println(tx.Hash())
		//fmt.Println(receipt.Status) // 1
	}

	//amount := big.NewInt(0)
	////amount.SetString("33190774000000000000000", 10)
	//amount.SetString("1", 10)
	//
	//mockDiffs := []*harvester.BCDiff{
	//	{
	//		From:   "0x2813fd17ea95b2655a7228383c5236e31090419e",
	//		To:     "0x3f363b4e038a6e43ce8321c50f3efbf460196d4b",
	//		Token:  "0xdefa4e8a7bcba345f687a2f1456f5edd9ce97202",
	//		Amount: amount,
	//	},
	//}

	a.listener(b.Number, diffs)
}

// refer https://github.com/ethereum/web3.py/blob/master/web3/contract.py#L435
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

const stakeABI = `[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "ClaimedUnstaked",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "RewardsClaimed",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "bytes16",
				"name": "",
				"type": "bytes16"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "enum MockStaking.Token",
				"name": "",
				"type": "uint8"
			}
		],
		"name": "Stake",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "Test",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "Unstake",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "MANAGER_ROLE",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "dad_token",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "mom_token",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "odysseys",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "odyssey_id",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "total_staked_into",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "total_stakers",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes16",
				"name": "odyssey_id",
				"type": "bytes16"
			},
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"internalType": "enum MockStaking.Token",
				"name": "token",
				"type": "uint8"
			}
		],
		"name": "stake",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "stakers",
		"outputs": [
			{
				"internalType": "address",
				"name": "user",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "total_rewards",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "total_staked",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "dad_amount",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "mom_amount",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "total_staked",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "unstakes",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "dad_amount",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "mom_amount",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "since",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`
