package arbitrum_nova_adapter2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"
	"math/big"
)

func addrToHash(addr *common.Address) *common.Hash {
	if addr == nil {
		return nil
	}
	res := common.HexToHash(addr.Hex())
	return &res
}

func (a *ArbitrumNovaAdapter) GetNFTBalance(block int64, wallet *common.Address, nftContract *common.Address) ([]umid.UMID, error) {
	transferString := "Transfer(address,address,uint256)"
	transferTopic := common.BytesToHash(crypto.Keccak256([]byte(transferString)))

	if nftContract == nil {
		return nil, errors.New("Failed to GetNFTBalance: NFT contract can not be nil")
	}

	contracts := []common.Address{
		//mom,
		*nftContract,
	}

	logs := make([]types.Log, 0)

	logsFrom, err := a.getLogs(transferTopic, contracts, big.NewInt(0), big.NewInt(block), wallet, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get logs for nft contract")
	}

	logsTo, err := a.getLogs(transferTopic, contracts, big.NewInt(0), big.NewInt(block), nil, wallet)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get logs for nft contract")
	}

	fmt.Println(logs)

	m := make(map[umid.UMID]int8)

	for _, l := range logsFrom {
		itemID := l.Topics[3].Big()

		var id umid.UMID
		itemID.FillBytes(id[:])
		_, ok := m[id]
		if !ok {
			m[id] = 0
		}

		m[id] -= 1
	}

	for _, l := range logsTo {
		itemID := l.Topics[3].Big()

		var id umid.UMID
		itemID.FillBytes(id[:])
		_, ok := m[id]
		if !ok {
			m[id] = 0
		}

		m[id] += 1
	}

	ids := make([]umid.UMID, 0)
	for id, i := range m {
		if i != 0 && i != 1 {
			fmt.Println("Failed to parse NFT transfers, Something wrong in blockchain history")
		}
		if i == 1 {
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func (a *ArbitrumNovaAdapter) getLogs(event common.Hash, address []common.Address, fromBlock *big.Int, toBlock *big.Int, source *common.Address, dest *common.Address) (replies []types.Log, err error) {
	args := make(map[string]interface{})
	var topix []interface{}
	topix = append(topix, event)
	topix = append(topix, addrToHash(source))
	topix = append(topix, addrToHash(dest))
	args["topics"] = topix
	args["address"] = address
	args["fromBlock"] = hexutil.EncodeBig(fromBlock)
	args["toBlock"] = hexutil.EncodeBig(toBlock)
	if err != nil {
		return
	}
	err = a.rpcClient.CallContext(context.TODO(), &replies, "eth_getLogs", args)
	return
}
