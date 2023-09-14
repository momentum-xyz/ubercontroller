package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/arbitrum_nova_adapter"
	helper "github.com/momentum-xyz/ubercontroller/harvester/cmd"
)

func main() {
	cfg := helper.MustGetConfig()

	logger := helper.GetZapLogger()
	sugaredLogger := logger.Sugar()

	pgConfig, err := cfg.Postgres.GenConfig(logger)
	if err != nil {
		log.Fatalf("failed to get db config: %s", err)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), pgConfig)
	if err != nil {
		log.Fatalf("failed to create db pool: %s", err)
	}
	defer pool.Close()

	//env := "anton_private_net"
	env := "main_net"

	var mom, dad, w1 common.Address

	if env == "main_net" {
		cfg.Arbitrum3.RPCURL = "https://nova.arbitrum.io/rpc"
		mom = common.HexToAddress("0x0C270A47D5B00bb8db42ed39fa7D6152496944ca")
		dad = common.HexToAddress("0x11817050402d2bb1418753ca398fdB3A3bc7CfEA")
		w1 = common.HexToAddress("0xAdd2e75c298F34E4d66fBbD4e056DA31502Da5B0")
	}

	if env == "anton_private_net" {
		cfg.Arbitrum3.RPCURL = "https://bcdev.antst.net:8547"
		mom = common.HexToAddress("0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee")
		dad = common.HexToAddress("0xfCa1B6bD67AeF9a9E7047bf7D3949f40E8dde18d")
		w1 = common.HexToAddress("0x683642c22feDE752415D4793832Ab75EFdF6223c")
	}

	a := arbitrum_nova_adapter.NewArbitrumNovaAdapter(&cfg.Arbitrum3, sugaredLogger)
	a.Run()

	output := make(chan harvester.UpdateCell)
	go worker(output)

	matrix := harvester.NewTokens(pool, a, sugaredLogger, output)
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

func worker(c <-chan harvester.UpdateCell) {
	for {
		u := <-c
		fmt.Println("worker", u.Contract, u.Wallet, u.Block, u.Value)
	}
}
