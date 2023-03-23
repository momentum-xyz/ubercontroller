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
	"github.com/momentum-xyz/ubercontroller/harvester/ethereum_adapter"
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
	ethereumAdapter := ethereum_adapter.NewEthereumAdapter()
	ethereumAdapter.Run()
	_, _, _ = ethereumAdapter.GetInfo()
	if err := harv.RegisterAdapter(ethereumAdapter); err != nil {
		log.Fatal(err)
	}

	// ** Harvester Clients
	testHandler1 := testHandler1
	ptrTestHandler1 := &testHandler1
	harv.Subscribe(harvester.Ethereum, harvester.NewBlock, ptrTestHandler1)
	harv.Subscribe(harvester.Polkadot, harvester.NewBlock, ptrTestHandler1)

	testHandler2 := testHandler2
	ptrTestHandler2 := &testHandler2
	harv.Subscribe(harvester.Ethereum, harvester.NewBlock, ptrTestHandler2)

	wallet1 := "0x9592b70a5a6c8ece2ef55547c3f07f1862372fd1"
	contract1 := "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae"
	contract2 := "0xde0b295669a9fd93d5f28d9ec85e40f4cb697ccc"

	wallet2 := "0x9Dd3f13cbacf6bd96E9757eCaceDf5236ffF787f"
	contract3 := "0x3af33bEF05C2dCb3C7288b77fe1C8d2AeBA4d789"

	err = harv.SubscribeForWalletAndContract(harvester.Ethereum, wallet1, contract1, ptrTestHandler2)
	if err != nil {
		panic(err)
	}
	err = harv.SubscribeForWalletAndContract(harvester.Ethereum, wallet1, contract2, ptrTestHandler2)
	if err != nil {
		panic(err)
	}

	err = harv.SubscribeForWalletAndContract(harvester.Ethereum, wallet2, contract3, ptrTestHandler2)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 30)
	harv.Unsubscribe(harvester.Ethereum, harvester.NewBlock, ptrTestHandler2)

	time.Sleep(time.Second * 50)
}

func testHandler1(p any) {
	fmt.Printf("testHandler1: %+v \n", p)
}

func testHandler2(p any) {
	fmt.Printf("testHandler2: %+v \n", p)
}
