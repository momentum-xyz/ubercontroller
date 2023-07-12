package arbitrum_nova_adapter

import (
	"context"
	"fmt"
	"log"
	"math"
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
	rpcClient *rpc.Client
	lastBlock uint64
	contracts *Contracts
}

func NewArbitrumNovaAdapter(cfg *config.Config) *ArbitrumNovaAdapter {
	return &ArbitrumNovaAdapter{
		umid:      umid.MustParse("ccccaaaa-1111-2222-3333-222222222222"),
		wsURL:     cfg.Arbitrum.WSURL,
		httpURL:   cfg.Arbitrum.RPCURL,
		name:      "arbitrum_nova",
		contracts: NewContracts(&cfg.Arbitrum),
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
						a.listener(n, nil, nil)
					}
				}
			}
		}
	}()
}

func (a *ArbitrumNovaAdapter) RegisterNewBlockListener(f harvester.AdapterListener) {
	a.listener = f
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

func (a *ArbitrumNovaAdapter) GetLogsRecursively(fromBlock, toBlock int64, contracts []common.Address, level int) ([]any, error) {

	fmt.Println("GET: ", strconv.Itoa(level), fromBlock, toBlock)
	if level > 3 {
		return nil, errors.New("GetLogsWrapper maximum recursion level")
	}

	logs, err := a.GetLogs(fromBlock, toBlock, contracts)
	if err != nil {
		fmt.Println(err)
		allLogs := make([]any, 0)
		parts := int64(7)
		if toBlock-fromBlock < parts {
			return nil, errors.WithMessage(err, "can not split getLogs to parts")
		}

		step, _ := math.Modf(float64((toBlock - fromBlock) / parts))
		for i := 1; i <= int(parts); i++ {

			from := int64(fromBlock) + int64(i-1)*int64(step)
			to := int64(fromBlock) + int64(i)*int64(step) - 1
			if int64(i) == parts {
				to = toBlock
			}

			l, err := a.GetLogsRecursively(from, to, contracts, level+1)
			if err != nil {
				return nil, errors.WithMessage(err, "recursive call error")
			}
			fmt.Println(i, from, to)
			allLogs = append(allLogs, l...)
		}

		return allLogs, nil
	}

	fmt.Println("RETURN: ", fromBlock, toBlock, len(logs))
	return logs, err

}

func (a *ArbitrumNovaAdapter) GetLogs(fromBlock, toBlock int64, contracts []common.Address) ([]any, error) {

	if contracts == nil {
		contracts = a.contracts.AllAddresses
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: contracts,
	}

	bcLogs, err := a.FilterLogs(context.TODO(), query)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to filter log")
	}

	logs := make([]any, 0)

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	logStakeSigHash := crypto.Keccak256Hash([]byte(a.contracts.StakeABI.Events["Stake"].Sig))
	logUnstakeSigHash := a.contracts.StakeABI.Events["Unstake"].ID
	logRestakeSigHash := a.contracts.StakeABI.Events["Restake"].ID

	logTransferNftHash := a.contracts.NftABI.Events["Transfer"].ID

	for _, vLog := range bcLogs {
		//fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		//fmt.Printf("Log Index: %d\n", vLog.Index)

		// Iterate contracts
		switch vLog.Address.Hex() {
		case a.contracts.momTokenAddress.Hex(), a.contracts.dadTokenAddress.Hex():
			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				//fmt.Printf("Log Name: Transfer\n")

				//var transferEvent harvester.BCDiff

				var e harvester.TransferERC20Log

				ev, err := a.contracts.TokenABI.Unpack("Transfer", vLog.Data)
				if err != nil {
					return nil, errors.WithMessage(err, "failed to unpack event from ABI")
				}

				e.Contract = strings.ToLower(vLog.Address.Hex())
				// Hex and Un Hex here used to remove padding zeros
				e.From = strings.ToLower(common.HexToAddress(vLog.Topics[1].Hex()).Hex())
				e.To = strings.ToLower(common.HexToAddress(vLog.Topics[2].Hex()).Hex())
				if len(ev) > 0 {
					e.Value = ev[0].(*big.Int)
				}

				logs = append(logs, &e)
			}

		case a.contracts.stakeAddress.Hex():
			switch vLog.Topics[0].Hex() {

			case logStakeSigHash.Hex():
				ev, err := a.contracts.StakeABI.Unpack("Stake", vLog.Data)
				if err != nil {
					return nil, errors.WithMessage(err, "failed to unpack event from ABI")
				}

				// Hack to remove extra zeroes
				fromWallet := common.HexToAddress(vLog.Topics[1].Hex())

				b := vLog.Topics[2].Bytes()[16:]
				odysseyID, err := umid.FromBytes(b)
				if err != nil {
					return nil, errors.WithMessage(err, "failed to parse umid from bytes")
				}
				if odysseyID == umid.MustParse("ccccaaaa-1111-2222-3333-222222222222") ||
					odysseyID == umid.MustParse("ccccaaaa-1111-2222-3333-222222222244") ||
					odysseyID == umid.MustParse("ccccaaaa-1111-2222-3333-222222222241") {
					// Skip test Odyssey IDs
					continue
				}

				transactionHash := vLog.TxHash.Hex()
				amount := ev[0].(*big.Int)
				tokenType := ev[1].(uint8)
				totalAmount := ev[2].(*big.Int)

				e := &harvester.StakeLog{
					TxHash:       transactionHash,
					LogIndex:     vLog.Index,
					UserWallet:   fromWallet.Hex(),
					OdysseyID:    odysseyID,
					AmountStaked: amount,
					TokenType:    tokenType,
					TotalStaked:  totalAmount,
				}

				logs = append(logs, e)

			case logUnstakeSigHash.Hex():
				ev, err := a.contracts.StakeABI.Unpack("Unstake", vLog.Data)
				if err != nil {
					return nil, errors.WithMessage(err, "failed to unpack event from ABI")
				}

				// Hack to remove extra zeroes
				fromWallet := common.HexToAddress(vLog.Topics[1].Hex())

				b1 := vLog.Topics[2].Bytes()
				b2 := b1[16:]
				odysseyID, err := umid.FromBytes(b2)
				if err != nil {
					return nil, errors.WithMessage(err, "failed to parse umid from bytes")
				}

				amount := ev[0].(*big.Int)
				tokenType := ev[1].(uint8)
				totalAmount := ev[2].(*big.Int)

				e := &harvester.UnstakeLog{
					UserWallet:     fromWallet.Hex(),
					OdysseyID:      odysseyID,
					AmountUnstaked: amount,
					TokenType:      tokenType,
					TotalStaked:    totalAmount,
				}

				logs = append(logs, e)

			case logRestakeSigHash.Hex():
				fmt.Println("Restake")
			}
		case a.contracts.nftAddress.Hex():
			fmt.Println("NFT")

			switch vLog.Topics[0].Hex() {
			case logTransferNftHash.Hex():
				// TODO Not sure why vLog.Data is empty
				//ev, err := a.contracts.NftABI.Unpack("Transfer", vLog.Data)
				//if err != nil {
				//	return nil, errors.WithMessage(err, "failed to unpack event from ABI")
				//}

				from := strings.ToLower(common.HexToAddress(vLog.Topics[1].Hex()).Hex())
				to := strings.ToLower(common.HexToAddress(vLog.Topics[2].Hex()).Hex())
				itemID := vLog.Topics[3].Big()

				var id umid.UMID
				itemID.FillBytes(id[:])

				if err != nil {
					return nil, errors.WithMessage(err, "failed to read umid from bytes")
				}

				e := &harvester.TransferNFTLog{
					From:     from,
					To:       to,
					TokenID:  id,
					Contract: strings.ToLower(vLog.Address.Hex()),
				}

				logs = append(logs, e)
			}
		}

	}

	return logs, nil
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
