package arbitrum_nova_adapter2

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester2"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type ArbitrumNovaAdapter struct {
	listener  harvester2.AdapterListener
	umid      umid.UMID
	httpURL   string
	name      harvester2.BCType
	rpcClient *rpc.Client
	lastBlock uint64

	// tmp solution, remove it after indexed params added to contract events
	cache  []types.Log
	block  int64
	mu     sync.Mutex
	logger *zap.SugaredLogger
}

func NewArbitrumNovaAdapter(cfg *config.Arbitrum3, logger *zap.SugaredLogger) *ArbitrumNovaAdapter {
	return &ArbitrumNovaAdapter{
		umid:    umid.MustParse("ccccaaaa-1111-2222-3333-222222222222"),
		httpURL: cfg.RPCURL,
		name:    "arbitrum_nova",
		cache:   make([]types.Log, 0),
		block:   0,

		logger: logger,
	}
}

func (a *ArbitrumNovaAdapter) GetLastBlockNumber() (uint64, error) {
	var resp string
	if err := a.rpcClient.Call(&resp, "eth_blockNumber"); err != nil {
		return 0, errors.WithMessage(err, "failed to make RPC call to arbitrum:")
	}

	return hex2int(resp), nil
}

func (a *ArbitrumNovaAdapter) Run() {
	var err error
	a.rpcClient, err = rpc.DialHTTP(a.httpURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Arbitrum Block Chain: " + a.httpURL)
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
					if a.listener != nil {
						a.listener(n)
					}
				}
			}
		}
	}()
}

func (a *ArbitrumNovaAdapter) RegisterNewBlockListener(f harvester2.AdapterListener) {
	a.listener = f
}

func (a *ArbitrumNovaAdapter) GetBalance(wallet *common.Address, contract *common.Address, blockNumber uint64) (*big.Int, error) {
	type request struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}

	w := wallet.Hex()
	c := contract.Hex()

	// "0x70a08231" - crypto.Keccak256Hash([]byte("balanceOf(address)")).String()[0:10]
	data := "0x70a08231" + fmt.Sprintf("%064s", w[2:]) // %064s means that the string is padded with 0 to 64 bytes
	req := request{c, data}

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

func (a *ArbitrumNovaAdapter) getRawLogs(
	topic0 *common.Hash,
	topic1 *common.Address,
	topic2 *common.Address,
	addresses []common.Address,
	fromBlock *big.Int,
	toBlock *big.Int,
) (replies []types.Log, err error) {
	args := make(map[string]interface{})
	var topix []any
	topix = append(topix, topic0)
	topix = append(topix, topic1)
	topix = append(topix, topic2)
	//topix = append(topix, addrToHash(source))
	//topix = append(topix, addrToHash(dest))
	args["topics"] = topix
	args["address"] = addresses
	args["fromBlock"] = hexutil.EncodeBig(fromBlock)
	args["toBlock"] = hexutil.EncodeBig(toBlock)
	if err != nil {
		return
	}
	err = a.rpcClient.CallContext(context.TODO(), &replies, "eth_getLogs", args)
	return
}

func (a *ArbitrumNovaAdapter) GetLogs(fromBlock, toBlock int64, contracts []common.Address) ([]any, error) {
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	bcLogs, err := a.getRawLogs(&logTransferSigHash, nil, nil, contracts, big.NewInt(fromBlock), big.NewInt(toBlock))
	fmt.Println(bcLogs)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to filter log")
	}

	logs := make([]any, 0)

	c := 0
	for _, vLog := range bcLogs {
		//fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		//fmt.Printf("Log Index: %d\n", vLog.Index)

		c++
		//fmt.Println(vLog.Data)
		//fmt.Println(len(vLog.Data))
		//fmt.Println(len(vLog.Topics))
		fmt.Println(c, vLog.Topics[0].Hex())

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			//fmt.Printf("Log Name: Transfer\n")

			//var transferEvent harvester.BCDiff

			var e harvester2.TransferERC20Log

			e.Contract = vLog.Address
			// Hex and Un Hex here used to remove padding zeros
			e.From = common.HexToAddress(vLog.Topics[1].Hex())
			e.To = common.HexToAddress(vLog.Topics[2].Hex())

			data := common.TrimLeftZeroes(vLog.Data)
			hex := common.Bytes2Hex(data)
			hex = TrimLeftZeroes(hex)
			if hex != "" {
				erc20Amount, err := hexutil.DecodeBig("0x" + hex)
				if err != nil {
					fmt.Println(err)
				}
				e.Value = erc20Amount
				//fmt.Println(erc20Amount)
			}

			logs = append(logs, &e)
		}

	}

	return logs, nil
}

func TrimLeftZeroes(hex string) string {
	idx := 0
	for ; idx < len(hex); idx++ {
		if hex[idx] != '0' {
			break
		}
	}
	return hex[idx:]
}

func (a *ArbitrumNovaAdapter) GetInfo() (umid umid.UMID, name harvester2.BCType, rpcURL string) {
	return a.umid, a.name, a.httpURL
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
