package main

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester3"
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

	env := "anton_private_net"
	//env := "main_net"

	var mom, dad, w1 common.Address

	if env == "main_net" {
		cfg.Arbitrum2.RPCURL = "https://nova.arbitrum.io/rpc"
		mom = common.HexToAddress("0x0C270A47D5B00bb8db42ed39fa7D6152496944ca")
		dad = common.HexToAddress("0x11817050402d2bb1418753ca398fdB3A3bc7CfEA")
		w1 = common.HexToAddress("0xAdd2e75c298F34E4d66fBbD4e056DA31502Da5B0")
	}

	if env == "anton_private_net" {
		cfg.Arbitrum2.RPCURL = "https://bcdev.antst.net:8547"
		mom = common.HexToAddress("0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee")
		dad = common.HexToAddress("0xfCa1B6bD67AeF9a9E7047bf7D3949f40E8dde18d")
		w1 = common.HexToAddress("0x683642c22feDE752415D4793832Ab75EFdF6223c")
	}

	a := arbitrum_nova_adapter3.NewArbitrumNovaAdapter(&cfg.Arbitrum2, sugaredLogger)
	a.Run()

	matrix := harvester3.NewTokenMatrix(a, sugaredLogger)
	err = matrix.Run()
	if err != nil {
		log.Fatal(err)
	}
	_ = matrix

	err = matrix.AddContract(mom)
	if err != nil {
		log.Fatal(err)
	}
	err = matrix.AddWallet(w1)
	if err != nil {
		log.Fatal(err)
	}
	err = matrix.AddContract(dad)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 30)
}
