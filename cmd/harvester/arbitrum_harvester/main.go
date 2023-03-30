package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/arbitrum_nova_adapter"
)

func main() {
	fmt.Println("Harvester Debugger")

	cfg := config.GetConfig()
	logger, _ := zap.NewProduction()
	pgConfig, err := cfg.Postgres.GenConfig(logger)
	pool, err := pgxpool.ConnectConfig(context.Background(), pgConfig)
	if err != nil {
		log.Fatal("failed to create db pool")
	}
	defer pool.Close()

	// ** Harvester
	var harv harvester.IHarvester
	harv = harvester.NewHarvester(pool)

	// ** Ethereum Adapter
	adapter := arbitrum_nova_adapter.NewArbitrumNovaAdapter()
	adapter.Run()
	_, _, _ = adapter.GetInfo()
	if err := harv.RegisterAdapter(adapter); err != nil {
		log.Fatal(err)
	}

	// ** Harvester Clients
	testHandler1 := testHandler1
	ptrTestHandler1 := &testHandler1
	harv.Subscribe(harvester.ArbitrumNova, harvester.NewBlock, ptrTestHandler1)

	type pair struct {
		Wallet   string
		Contract string
	}

	pairs := []pair{
		//{
		//	Wallet:   "0x31854122F629B1B1E3b2aA85336F7b68f83924fA",
		//	Contract: "0x556353dab72b2F3223de2B2ac69700B3F280d357",
		//},
		{
			Wallet:   "0xe2148ee53c0755215df69b2616e552154edc584f",
			Contract: "0x7F85fB7f42A0c0D40431cc0f7DFDf88be6495e67",
		},
	}

	for _, pair := range pairs {
		err = harv.SubscribeForWalletAndContract(harvester.ArbitrumNova, pair.Wallet, pair.Contract, ptrTestHandler1)
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Second * 300)
	harv.Unsubscribe(harvester.ArbitrumNova, harvester.NewBlock, ptrTestHandler1)

	time.Sleep(time.Second * 500)
}

func testHandler1(p any) {
	fmt.Printf("testHandler1: %+v \n", p)
}

func testHandler2(p any) {
	fmt.Printf("testHandler2: %+v \n", p)
}
