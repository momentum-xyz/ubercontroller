package main

import (
	"context"
	"encoding/hex"
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
	harv := harvester.NewHarvester(pool)
	var harvForClient harvester.HarvesterAPI
	harvForClient = harv
	var harvForAdapter harvester.BCAdapterAPI
	harvForAdapter = harv

	// ** Ethereum Adapter
	ethereumAdapter := ethereum_adapter.NewEthereumAdapter(harvForAdapter)
	ethereumAdapter.Run()

	// ** Harvester Clients
	testHandler1 := testHandler1
	ptrTestHandler1 := &testHandler1
	harvForClient.Subscribe(harvester.Ethereum, harvester.NewBlock, ptrTestHandler1)
	harvForClient.Subscribe(harvester.Polkadot, harvester.NewBlock, ptrTestHandler1)

	testHandler2 := testHandler2
	ptrTestHandler2 := &testHandler2
	harvForClient.Subscribe(harvester.Ethereum, harvester.NewBlock, ptrTestHandler2)

	wallet := "9592b70a5a6c8ece2ef55547c3f07f1862372fd1"
	contract := "de0b295669a9fd93d5f28d9ec85e40f4cb697bae"
	contract2 := "de0b295669a9fd93d5f28d9ec85e40f4cb697ccc"

	err = harv.SubscribeForWalletAndContract(harvester.Ethereum, HexToAddress(wallet), HexToAddress(contract), ptrTestHandler2)
	if err != nil {
		panic(err)
	}
	err = harv.SubscribeForWalletAndContract(harvester.Ethereum, HexToAddress(wallet), HexToAddress(contract2), ptrTestHandler2)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 30)
	harvForClient.Unsubscribe(harvester.Ethereum, harvester.NewBlock, ptrTestHandler2)

	time.Sleep(time.Second * 50)
}

func HexToAddress(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func testHandler1(p any) {
	fmt.Printf("testHandler1: %+v \n", p)
}

func testHandler2(p any) {
	fmt.Printf("testHandler2: %+v \n", p)
}
