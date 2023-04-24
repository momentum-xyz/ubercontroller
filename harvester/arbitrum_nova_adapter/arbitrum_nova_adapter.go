package arbitrum_nova_adapter

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type ArbitrumNovaAdapter struct {
	listener harvester.AdapterListener
	umid     umid.UMID
	wsURL    string
	httpURL  string
	name     string
	//client           *ethclient.Client
	rpcClient        *rpc.Client
	contractABI      abi.ABI
	stakeContractABI abi.ABI
	stakeContract    common.Address
	lastBlock        uint64
}

func NewArbitrumNovaAdapter(cfg *config.Config) *ArbitrumNovaAdapter {
	return &ArbitrumNovaAdapter{
		umid:          umid.MustParse("ccccaaaa-1111-2222-3333-222222222222"),
		wsURL:         cfg.Arbitrum.ArbitrumWSURL,
		httpURL:       cfg.Arbitrum.ArbitrumRPCURL,
		name:          "arbitrum_nova",
		stakeContract: common.HexToAddress(cfg.Arbitrum.ArbitrumStakeContractAddress),
	}
}

func (a *ArbitrumNovaAdapter) GetLastBlockNumber() (uint64, error) {
	//number, err := a.client.BlockNumber(context.TODO())
	//_ = number
	//fmt.Println(err)
	//fmt.Println(number)

	var resp string
	if err := a.rpcClient.Call(&resp, "eth_blockNumber"); err != nil {
		return 0, errors.WithMessage(err, "failed to make RPC call to arbitrum:")
	}

	return hex2int(resp), nil
}

func (a *ArbitrumNovaAdapter) Run() {
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

	//a.client, err = ethclient.Dial(a.rpcURL)
	//if err != nil {
	//	log.Fatal(err)
	//}

	a.rpcClient, err = rpc.DialHTTP(a.httpURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Arbitrum Block Chain: " + a.wsURL)
	///////

	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				_ = t
				//fmt.Println("Tick at", t)
				n, err := a.GetLastBlockNumber()
				if err != nil {
					fmt.Println(err)
				}
				if a.lastBlock < n {
					a.lastBlock = n
					a.listener(n, nil, nil)
				}
			}
		}
	}()

	///////

	//
	//ch := make(chan *types.Header)
	//
	//sub, err := a.client.SubscribeNewHead(context.Background(), ch)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//go func() {
	//	for {
	//		select {
	//		case err := <-sub.Err():
	//			log.Fatal(err)
	//		case vLog := <-ch:
	//
	//			fmt.Println(vLog.Number)
	//			//fmt.Println(vLog.ReceiptHash)
	//			//fmt.Println(vLog.ParentHash)
	//			//fmt.Println(vLog.Root)
	//			//fmt.Println(vLog.TxHash)
	//			//fmt.Println(vLog.Hash())
	//
	//			block := &harvester.BCBlock{
	//				Hash: vLog.Hash().String(),
	//			}
	//
	//			if vLog.Number != nil {
	//				block.Number = vLog.Number.Uint64()
	//			}
	//
	//			if a.listener != nil {
	//				a.onNewBlock(block)
	//			}
	//		}
	//	}
	//}()
}

func (a *ArbitrumNovaAdapter) RegisterNewBlockListener(f harvester.AdapterListener) {
	a.listener = f
}

//func (a *ArbitrumNovaAdapter) onNewBlock(b *harvester.BCBlock) {
//	diffs, stakes, err := a.GetTransferLogs(int64(b.Number), int64(b.Number), []common.Address{})
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	a.listener(b.Number, diffs, stakes)
//}

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
		return nil, errors.WithMessage(err, "failed to make RPC call to arbitrum:")
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

func (a *ArbitrumNovaAdapter) GetTransferLogs(fromBlock, toBlock int64, addresses []common.Address) ([]*harvester.BCDiff, []*harvester.BCStake, error) {

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: append(addresses, a.stakeContract),
	}

	logs, err := a.FilterLogs(context.TODO(), query)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to filter log")
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	logStakeSigHash := crypto.Keccak256Hash([]byte(a.stakeContractABI.Events["Stake"].Sig))
	logUnstakeSigHash := a.stakeContractABI.Events["Unstake"].ID
	logRestakeSigHash := a.stakeContractABI.Events["Restake"].ID

	diffs := make([]*harvester.BCDiff, 0)
	stakes := make([]*harvester.BCStake, 0)

	for _, vLog := range logs {
		//fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		//fmt.Printf("Log Index: %d\n", vLog.Index)

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

			totalAmount := ev[4].(*big.Int)

			stake := &harvester.BCStake{
				From:        fromWallet.Hex(),
				OdysseyID:   odysseyID,
				TokenType:   tokenType,
				Amount:      amount,
				TotalAmount: totalAmount,
			}

			if stake.OdysseyID.String() != uuid.MustParse("ccccaaaa-1111-2222-3333-222222222222").String() &&
				stake.OdysseyID.String() != uuid.MustParse("ccccaaaa-1111-2222-3333-222222222244").String() &&
				stake.OdysseyID.String() != uuid.MustParse("ccccaaaa-1111-2222-3333-222222222241").String() {

				stakes = append(stakes, stake)
			}

		//fmt.Printf("%+v %+v %+v %+v \n\n", fromWallet.String(), odysseyID.String(), amount, tokenType)
		//fmt.Println(ev)

		case logUnstakeSigHash.Hex():
			log.Println("Unstake")

			ev, err := a.stakeContractABI.Unpack("Unstake", vLog.Data)
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

			totalAmount := ev[4].(*big.Int)

			stake := &harvester.BCStake{
				From:        fromWallet.Hex(),
				OdysseyID:   odysseyID,
				TokenType:   tokenType,
				Amount:      amount,
				TotalAmount: totalAmount,
			}

			stakes = append(stakes, stake)

		case logRestakeSigHash.Hex():
			fmt.Println("Restake")
		}
	}

	return diffs, stakes, nil
}

func (a *ArbitrumNovaAdapter) GetInfo() (umid umid.UMID, name string, rpcURL string) {
	return a.umid, a.name, a.wsURL
}

func (a *ArbitrumNovaAdapter) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	var result []types.Log
	arg, err := toFilterArg(q)
	if err != nil {
		return nil, err
	}
	err = a.rpcClient.CallContext(ctx, &result, "eth_getLogs", arg)
	return result, err
}

func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = toBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = toBlockNumArg(q.ToBlock)
	}
	return arg, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	finalized := big.NewInt(int64(rpc.FinalizedBlockNumber))
	if number.Cmp(finalized) == 0 {
		return "finalized"
	}
	safe := big.NewInt(int64(rpc.SafeBlockNumber))
	if number.Cmp(safe) == 0 {
		return "safe"
	}
	return hexutil.EncodeBig(number)
}

func hex2int(hexStr string) uint64 {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(result)
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
				"name": "previousAdmin",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "newAdmin",
				"type": "address"
			}
		],
		"name": "AdminChanged",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "beacon",
				"type": "address"
			}
		],
		"name": "BeaconUpgraded",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "claim_rewards",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "claim_unstaked_tokens",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
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
		"name": "ClaimedUnstaked",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "grantRole",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_mom_token",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "_dad_token",
				"type": "address"
			}
		],
		"name": "initialize",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint8",
				"name": "version",
				"type": "uint8"
			}
		],
		"name": "Initialized",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "renounceRole",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "bytes16",
				"name": "from_odyssey_id",
				"type": "bytes16"
			},
			{
				"internalType": "bytes16",
				"name": "to_odyssey_id",
				"type": "bytes16"
			},
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"internalType": "enum Staking.Token",
				"name": "token",
				"type": "uint8"
			}
		],
		"name": "restake",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
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
				"internalType": "enum Staking.Token",
				"name": "",
				"type": "uint8"
			}
		],
		"name": "Restake",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "revokeRole",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
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
				"indexed": true,
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "previousAdminRole",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "newAdminRole",
				"type": "bytes32"
			}
		],
		"name": "RoleAdminChanged",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "RoleGranted",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "account",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "sender",
				"type": "address"
			}
		],
		"name": "RoleRevoked",
		"type": "event"
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
				"internalType": "enum Staking.Token",
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
				"internalType": "enum Staking.Token",
				"name": "",
				"type": "uint8"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "Stake",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "bytes16",
				"name": "odyssey_id",
				"type": "bytes16"
			},
			{
				"internalType": "enum Staking.Token",
				"name": "token",
				"type": "uint8"
			}
		],
		"name": "unstake",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
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
				"internalType": "enum Staking.Token",
				"name": "",
				"type": "uint8"
			}
		],
		"name": "Unstake",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_dad_token",
				"type": "address"
			}
		],
		"name": "update_dad_token_contract",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_locking_period",
				"type": "uint256"
			}
		],
		"name": "update_locking_period",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_mom_token",
				"type": "address"
			}
		],
		"name": "update_mom_token_contract",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address[]",
				"name": "addresses",
				"type": "address[]"
			},
			{
				"internalType": "uint256[]",
				"name": "amounts",
				"type": "uint256[]"
			}
		],
		"name": "update_rewards",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "implementation",
				"type": "address"
			}
		],
		"name": "Upgraded",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "newImplementation",
				"type": "address"
			}
		],
		"name": "upgradeTo",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "newImplementation",
				"type": "address"
			},
			{
				"internalType": "bytes",
				"name": "data",
				"type": "bytes"
			}
		],
		"name": "upgradeToAndCall",
		"outputs": [],
		"stateMutability": "payable",
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
		"name": "DEFAULT_ADMIN_ROLE",
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
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			}
		],
		"name": "getRoleAdmin",
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
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "role",
				"type": "bytes32"
			},
			{
				"internalType": "address",
				"name": "account",
				"type": "address"
			}
		],
		"name": "hasRole",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "locking_period",
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
				"internalType": "bytes16",
				"name": "",
				"type": "bytes16"
			}
		],
		"name": "odysseys",
		"outputs": [
			{
				"internalType": "bytes16",
				"name": "odyssey_id",
				"type": "bytes16"
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
		"inputs": [],
		"name": "proxiableUUID",
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
		"inputs": [
			{
				"internalType": "bytes4",
				"name": "interfaceId",
				"type": "bytes4"
			}
		],
		"name": "supportsInterface",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
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
