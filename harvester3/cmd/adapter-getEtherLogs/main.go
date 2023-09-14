package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/harvester3/arbitrum_nova_adapter3"
	helper "github.com/momentum-xyz/ubercontroller/harvester3/cmd"
)

func main() {
	cfg := helper.MustGetConfig()

	logger := helper.GetZapLogger()
	sugaredLogger := logger.Sugar()

	//env := "anton_private_net"
	env := "main_net"

	var mom, dad, nft, nftOMNIA, w1, wKovi, wOMNIAHOLDER common.Address
	_ = mom
	_ = dad
	_ = nft
	_ = w1
	_ = wKovi
	_ = nftOMNIA
	_ = wOMNIAHOLDER

	if env == "main_net" {
		cfg.Arbitrum3.RPCURL = "https://nova.arbitrum.io/rpc"
		mom = common.HexToAddress("0x0C270A47D5B00bb8db42ed39fa7D6152496944ca")
		dad = common.HexToAddress("0x11817050402d2bb1418753ca398fdB3A3bc7CfEA")
		nft = common.HexToAddress("0x1F59C1db986897807d7c3eF295C3480a22FBa834")
		nftOMNIA = common.HexToAddress("0x402a928dd8342f5604a9a416d00997105c76bfa2")
		wOMNIAHOLDER = common.HexToAddress("0x9daaa0ff2be321b03b78165f5ad21a44e3c14bd6")
		//w1 = common.HexToAddress("0xAdd2e75c298F34E4d66fBbD4e056DA31502Da5B0")
		w1 = common.HexToAddress("0x42ae6199bb589cfe2df3a93cf93cf5fc1caab2e2")
		wKovi = common.HexToAddress("0xc6220f7F21e15B8886eD38A98496E125b564c414")
	}

	if env == "anton_private_net" {
		cfg.Arbitrum3.RPCURL = "https://bcdev.antst.net:8547"
		mom = common.HexToAddress("0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee")
		dad = common.HexToAddress("0xfCa1B6bD67AeF9a9E7047bf7D3949f40E8dde18d")
		nft = common.HexToAddress("0xbc48cb82903f537614E0309CaF6fe8cEeBa3d174")
		w1 = common.HexToAddress("0x683642c22feDE752415D4793832Ab75EFdF6223c")
		wKovi = common.HexToAddress("0xc6220f7F21e15B8886eD38A98496E125b564c414")
	}

	a := arbitrum_nova_adapter3.NewArbitrumNovaAdapter(&cfg.Arbitrum3, sugaredLogger)
	a.Run()

	n, err := a.GetLastBlockNumber()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Block: %+v \n", n)

	//items, err := a.GetNFTBalance(&nft, &w1, n)
	wallets := map[common.Address]bool{
		w1: true,
	}

	start := time.Now()

	//items, err := a.GetEtherLogs(n-20, n, wallets)
	items, err := a.GetEtherLogs(19100246-10, 19100246, wallets)
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range items {
		fmt.Printf("%v %v %v\n", i.Block, i.Wallet.Hex(), i.Delta.String())
	}

	fmt.Println(time.Now().Sub(start))

}
