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

	// **  Adapter
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

	testHandler2 := testHandler2
	ptrTestHandler2 := &testHandler2

	//wallet1 := "0x9592b70a5a6c8ece2ef55547c3f07f1862372fd1"
	//contract1 := "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae"
	//contract2 := "0xde0b295669a9fd93d5f28d9ec85e40f4cb697ccc"
	//
	//wallet2 := "0x31854122F629B1B1E3b2aA85336F7b68f83924fA"
	//contract3 := "0x556353dab72b2F3223de2B2ac69700B3F280d357"

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
			Wallet:   "0x15c7152B3b02324d17e753E4cfF65C0f1759812B",
			Contract: "0x556353dab72b2F3223de2B2ac69700B3F280d357",
		},
	}

	for _, pair := range pairs {
		err = harv.SubscribeForWalletAndContract(harvester.ArbitrumNova, pair.Wallet, pair.Contract, ptrTestHandler2)
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Second * 30)
	harv.Unsubscribe(harvester.ArbitrumNova, harvester.NewBlock, ptrTestHandler2)

	time.Sleep(time.Second * 500)
}

func testHandler1(p any) {
	fmt.Printf("testHandler1: %+v \n", p)
}

func testHandler2(p any) {
	fmt.Printf("testHandler2: %+v \n", p)
}
