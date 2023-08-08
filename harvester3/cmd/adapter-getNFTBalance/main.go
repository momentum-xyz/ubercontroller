package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester3/arbitrum_nova_adapter3"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		// NOTE: set this false to enable stack trace
		DisableStacktrace: true,
	}

	logger, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}
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
		w1 = common.HexToAddress("0xAdd2e75c298F34E4d66fBbD4e056DA31502Da5B0")
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

	//logTransferSig := []byte("Transfer(address,address,uint256)")
	//logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	//logs, err := a.GetRawLogs(&logTransferSigHash, nil, nil, []common.Address{mom}, big.NewInt(0), big.NewInt(int64(n)))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Print(logs)

	//items, err := a.GetNFTBalance(&nft, &w1, n)
	items, err := a.GetNFTBalance(&nftOMNIA, &wOMNIAHOLDER, n)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(items)
}
