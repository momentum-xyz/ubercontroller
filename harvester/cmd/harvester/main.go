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

	var mom, dad, nft, w1, w2 common.Address
	_ = dad
	_ = w2

	if env == "main_net" {
		cfg.Arbitrum3.RPCURL = "https://nova.arbitrum.io/rpc"
		mom = common.HexToAddress("0x0C270A47D5B00bb8db42ed39fa7D6152496944ca")
		dad = common.HexToAddress("0x11817050402d2bb1418753ca398fdB3A3bc7CfEA")
		nft = common.HexToAddress("0x1F59C1db986897807d7c3eF295C3480a22FBa834")
		w1 = common.HexToAddress("0xAdd2e75c298F34E4d66fBbD4e056DA31502Da5B0")
		w2 = common.HexToAddress("0xc6220f7F21e15B8886eD38A98496E125b564c414")
	}

	if env == "anton_private_net" {
		cfg.Arbitrum3.RPCURL = "https://bcdev.antst.net:8547"
		mom = common.HexToAddress("0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee")
		dad = common.HexToAddress("0xfCa1B6bD67AeF9a9E7047bf7D3949f40E8dde18d")
		nft = common.HexToAddress("0xbc48cb82903f537614E0309CaF6fe8cEeBa3d174")
		w1 = common.HexToAddress("0x683642c22feDE752415D4793832Ab75EFdF6223c")
		w2 = common.HexToAddress("0x695c0AbC571F5F434351dAB51b92b851562a8BE1")
	}

	a := arbitrum_nova_adapter.NewArbitrumNovaAdapter(&cfg.Arbitrum3, sugaredLogger)
	a.Run()

	harvesterOutput := make(chan harvester.UpdateCell)
	go worker(harvesterOutput)

	harvester := harvester.NewHarvester(harvesterOutput, pool, a, sugaredLogger)
	err = harvester.Run()
	if err != nil {
		log.Fatal(err)
	}
	_ = harvester

	err = harvester.AddTokenContract(mom)
	if err != nil {
		log.Fatal(err)
	}

	err = harvester.AddWallet(w1)
	if err != nil {
		log.Fatal(err)
	}

	err = harvester.AddWallet(w2)
	if err != nil {
		log.Fatal(err)
	}

	err = harvester.AddNFTContract(nft)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 300)
}

func worker(output chan harvester.UpdateCell) {
	for {
		fmt.Println(<-output)
	}
}
