package arbitrum_nova_adapter3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

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
	"github.com/momentum-xyz/ubercontroller/harvester3"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type ArbitrumNovaAdapter struct {
	listeners []func(blockNumber uint64)
	umid      umid.UMID
	httpURL   string
	name      string
	rpcClient *rpc.Client
	lastBlock uint64

	logger *zap.SugaredLogger
}

func NewArbitrumNovaAdapter(cfg *config.Arbitrum3, logger *zap.SugaredLogger) *ArbitrumNovaAdapter {
	return &ArbitrumNovaAdapter{
		umid:      umid.MustParse("ccccaaaa-1111-2222-3333-222222222222"),
		httpURL:   cfg.RPCURL,
		name:      "arbitrum_nova",
		listeners: make([]func(blockNumber uint64), 0),

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
		a.logger.Error(err)
	}

	a.logger.Info("Connected to Arbitrum Block Chain: " + a.httpURL)

	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				_ = t
				n, err := a.GetLastBlockNumber()
				if err != nil {
					a.logger.Error(err)
				}
				if a.lastBlock < n {
					a.lastBlock = n
					for _, listener := range a.listeners {
						listener(n)
					}
				}
			}
		}
	}()
}

func (a *ArbitrumNovaAdapter) RegisterNewBlockListener(f harvester3.AdapterListener) {
	a.listeners = append(a.listeners, f)
}

func (a *ArbitrumNovaAdapter) GetTokenBalance(contract *common.Address, wallet *common.Address, blockNumber uint64) (*big.Int, uint64, error) {
	type request struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}

	w := wallet.Hex()
	c := contract.Hex()

	// "0x70a08231" - crypto.Keccak256Hash([]byte("balanceOf(address)")).String()[0:10]
	data := "0x70a08231" + fmt.Sprintf("%064s", w[2:]) // %064s means that the string is padded with 0 to 32 bytes
	req := request{c, data}

	var resp string
	n := hexutil.EncodeUint64(blockNumber)
	if err := a.rpcClient.Call(&resp, "eth_call", req, n); err != nil {
		return nil, 0, errors.WithMessage(err, "failed to make RPC call to arbitrum:")
	}

	balance := stringToBigInt(resp)

	return balance, blockNumber, nil
}

func stringToBigInt(str string) *big.Int {
	// remove leading zero of resp
	t := strings.TrimLeft(str[2:], "0")
	if t == "" {
		t = "0"
	}
	s := "0x" + t
	b, err := hexutil.DecodeBig(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
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

func (a *ArbitrumNovaAdapter) GetRawLogs(
	topic0 *common.Hash,
	topic1 *common.Hash,
	topic2 *common.Hash,
	addresses []common.Address,
	fromBlock *big.Int,
	toBlock *big.Int,
) (replies []types.Log, err error) {
	args := make(map[string]interface{})
	var topix []any
	topix = append(topix, topic0)
	topix = append(topix, topic1)
	topix = append(topix, topic2)

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

func (a *ArbitrumNovaAdapter) GetEtherBalance(wallet *common.Address, block uint64) (*big.Int, error) {
	resp := ""
	err := a.rpcClient.Call(&resp, "eth_getBalance", wallet.Hex(), hexutil.EncodeUint64(block))
	if err != nil {
		return nil, err
	}

	balance := stringToBigInt(resp)
	return balance, err
}

func (a *ArbitrumNovaAdapter) GetEtherLogs(fromBlock, toBlock uint64, wallets map[common.Address]bool) ([]harvester3.ChangeEtherLog, error) {
	logs := make([]harvester3.ChangeEtherLog, 0)
	res := make(map[string]any)

	type respTx struct {
		Hash  string  `json:"hash"`
		From  string  `json:"from"`
		To    *string `json:"to"`
		Value string  `json:"value"`
		Gas   string  `json:"gas"`
	}

	mu := sync.Mutex{}

	for b := fromBlock; b <= toBlock; b++ {
		blockNumber := hexutil.EncodeUint64(b)
		err := a.rpcClient.Call(&res, "eth_getBlockByNumber", blockNumber, true)
		if err != nil {
			return nil, err
		}
		txs, ok := res["transactions"]
		if !ok {
			continue
		}
		txsMap, ok := txs.([]any)
		for _, txMap := range txsMap {
			// TODO run in parallel
			jsonData, _ := json.Marshal(txMap)
			tx := respTx{}
			err = json.Unmarshal(jsonData, &tx)
			if err != nil {
				return nil, err
			}

			to := common.Address{}
			if tx.To != nil {
				to = common.HexToAddress(*tx.To)
			}
			from := common.HexToAddress(tx.From)

			_, hasFrom := wallets[from]
			_, hasTo := wallets[to]

			if !hasFrom && !hasTo {
				continue
			}

			if tx.Gas != "0x0" && hasFrom {
				go func() {
					receipt, err := a.eth_getTransactionReceipt(tx.Hash)
					if err != nil {
						a.logger.Error(err)
						return
					}
					if receipt.Status != "0x1" {
						return
					}

					gasUsed := big.NewInt(int64(hex2int(receipt.CumulativeGasUsed)))
					gasPrice := big.NewInt(int64(hex2int(receipt.EffectiveGasPrice)))
					delta := gasUsed.Mul(gasUsed, gasPrice)
					delta = delta.Neg(delta)

					mu.Lock()
					logs = append(logs, harvester3.ChangeEtherLog{
						Block:  b,
						Wallet: common.HexToAddress(receipt.From),
						Delta:  delta,
					})
					mu.Unlock()
				}()
			}

			if tx.Value == "0x0" {
				// Tx fee already counted in previous section
				continue
			}

			if hasTo {
				logs = append(logs, harvester3.ChangeEtherLog{
					Block:  b,
					Wallet: to,
					Delta:  big.NewInt(int64(hex2int(tx.Value))),
				})
			}
			if hasFrom {
				delta := big.NewInt(int64(hex2int(tx.Value)))
				logs = append(logs, harvester3.ChangeEtherLog{
					Block:  b,
					Wallet: common.HexToAddress(tx.From),
					Delta:  delta.Neg(delta),
				})
			}
		}
	}

	return logs, nil
}

type TransactionReceipt struct {
	From              string `json:"from"`
	Status            string `json:"status"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
}

func (a *ArbitrumNovaAdapter) eth_getTransactionReceipt(hash string) (*TransactionReceipt, error) {
	res := TransactionReceipt{}
	err := a.rpcClient.Call(&res, "eth_getTransactionReceipt", hash)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (a *ArbitrumNovaAdapter) GetNFTLogs(fromBlock, toBlock uint64, contracts []common.Address) ([]any, error) {
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	bcLogs, err := a.GetRawLogs(&logTransferSigHash, nil, nil, contracts, big.NewInt(int64(fromBlock)), big.NewInt(int64(toBlock)))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to filter log")
	}

	logs := make([]any, 0)

	for _, vLog := range bcLogs {

		if len(vLog.Topics) != 4 {
			a.logger.Error("Transfer NFT log must have 4 topic items")
			continue
		}
		//fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		//fmt.Printf("Log Index: %d\n", vLog.Index)

		var e harvester3.TransferNFTLog

		e.Block = vLog.BlockNumber
		e.Contract = vLog.Address
		// Hex and Un Hex here used to remove padding zeros
		e.From = common.HexToAddress(vLog.Topics[1].Hex())
		e.To = common.HexToAddress(vLog.Topics[2].Hex())
		e.TokenID = vLog.Topics[3]

		logs = append(logs, &e)
	}

	return logs, nil
}

func (a *ArbitrumNovaAdapter) GetTokenLogs(fromBlock, toBlock uint64, contracts []common.Address) ([]any, error) {
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	bcLogs, err := a.GetRawLogs(&logTransferSigHash, nil, nil, contracts, big.NewInt(int64(fromBlock)), big.NewInt(int64(toBlock)))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to filter log")
	}

	logs := make([]any, 0)

	for _, vLog := range bcLogs {

		if len(vLog.Topics) == 4 {
			a.logger.Error("Got Transfer NFT log from blockchain in Token contract handler")
			continue
		}
		//fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		//fmt.Printf("Log Index: %d\n", vLog.Index)

		var e harvester3.TransferERC20Log

		e.Block = vLog.BlockNumber
		e.Contract = vLog.Address
		// Hex and Un Hex here used to remove padding zeros
		e.From = common.HexToAddress(vLog.Topics[1].Hex())
		e.To = common.HexToAddress(vLog.Topics[2].Hex())

		data := common.TrimLeftZeroes(vLog.Data)
		hex := common.Bytes2Hex(data)
		hex = TrimLeftZeroes(hex)
		if hex == "" {
			a.logger.Error("Got Transfer Token log with empty data")
			continue
		}
		erc20Amount, err := hexutil.DecodeBig("0x" + hex)
		if err != nil {
			a.logger.Error(err)
		}
		e.Value = erc20Amount

		logs = append(logs, &e)
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

func (a *ArbitrumNovaAdapter) GetInfo() (umid umid.UMID, name string, rpcURL string) {
	return a.umid, a.name, a.httpURL
}

func hex2int(hexStr string) uint64 {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return result
}

func (a *ArbitrumNovaAdapter) GetNFTBalance(nftContract *common.Address, wallet *common.Address, block uint64) ([]common.Hash, error) {
	transferString := "Transfer(address,address,uint256)"
	transferTopic := common.BytesToHash(crypto.Keccak256([]byte(transferString)))

	if nftContract == nil {
		return nil, errors.New("Failed to GetNFTBalance: NFT contract can not be nil")
	}

	contracts := []common.Address{
		*nftContract,
	}

	logsFrom, err := a.GetRawLogs(&transferTopic, addrToHash(wallet), nil, contracts, big.NewInt(0), big.NewInt(int64(block)))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get logs for nft contract")
	}

	logsTo, err := a.GetRawLogs(&transferTopic, nil, addrToHash(wallet), contracts, big.NewInt(0), big.NewInt(int64(block)))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get logs for nft contract")
	}

	m := make(map[common.Hash]int8)

	for _, l := range logsFrom {
		id := l.Topics[3]

		_, ok := m[id]
		if !ok {
			m[id] = 0
		}

		m[id] -= 1
	}

	for _, l := range logsTo {
		id := l.Topics[3]

		_, ok := m[id]
		if !ok {
			m[id] = 0
		}

		m[id] += 1
	}

	ids := make([]common.Hash, 0)
	for id, i := range m {
		if i != 0 && i != 1 {
			a.logger.Error("Failed to parse NFT transfers, Something wrong in blockchain history")
		}
		if i == 1 {
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func addrToHash(addr *common.Address) *common.Hash {
	if addr == nil {
		return nil
	}
	res := common.HexToHash(addr.Hex())
	return &res
}
