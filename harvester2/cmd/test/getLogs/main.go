package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/momentum-xyz/ubercontroller/config"
)

func main() {

	transferString := "Transfer(address,address,uint256)"
	transferTopic := common.BytesToHash(crypto.Keccak256([]byte(transferString)))

	mom := common.HexToAddress("0x567d4e8264dC890571D5392fDB9fbd0e3FCBEe56")
	nft := common.HexToAddress("0x97E0B10D89a494Eb5cfFCc72853FB0750BD64AcD")
	_ = mom
	_ = nft

	w04 := common.HexToAddress("0xA058Aa2fCf33993e17D074E6843202E7C94bf267")
	w78 := common.HexToAddress("0x78B00B17E7e5619113A4e922BC3c8cb290355043")
	_ = w04
	_ = w78

	contracts := []common.Address{
		mom,
		//nft,
	}

	//response, err := getLogs(transferTopic, mom, big.NewInt(0), big.NewInt(1000), nil, &w04)
	//response, err := getLogs(transferTopic, mom, big.NewInt(0), big.NewInt(1000), &w04, nil)
	response, err := getLogs(transferTopic, contracts, big.NewInt(0), big.NewInt(1000), nil, &w78)
	//fmt.Println(response)

	b := big.NewInt(0)
	for k, v := range response {
		//b.FillBytes(v.Data)
		b.SetBytes(v.Data)
		fmt.Println(k, common.HexToAddress(v.Topics[1].Hex()), common.HexToAddress(v.Topics[2].Hex()), b.String())
	}
	fmt.Println(err)
}

func GetRPClient() *rpc.Client {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	rpcClient, err := rpc.DialHTTP(cfg.Arbitrum.RPCURL)
	if err != nil {
		log.Fatal(err)
	}

	return rpcClient
}

func addrToHash(addr *common.Address) *common.Hash {
	if addr == nil {
		return nil
	}
	res := common.HexToHash(addr.Hex())
	return &res
}
func getLogs(event common.Hash, address []common.Address, fromBlock *big.Int, toBlock *big.Int, source *common.Address, dest *common.Address) (replies []types.Log, err error) {
	args := make(map[string]interface{})
	var topix []interface{}
	topix = append(topix, event)
	topix = append(topix, addrToHash(source))
	topix = append(topix, addrToHash(dest))
	args["topics"] = topix
	args["address"] = address
	args["fromBlock"] = hexutil.EncodeBig(fromBlock)
	args["toBlock"] = hexutil.EncodeBig(toBlock)
	//endpoint := viper.GetString("ETH_CONNECT")
	client := GetRPClient()
	if err != nil {
		return
	}
	defer client.Close()
	err = client.CallContext(context.Background(), &replies, "eth_getLogs", args)
	return
}
