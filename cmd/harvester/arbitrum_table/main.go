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
	fmt.Println("Arbitrum Table Debugger")

	cfg := config.GetConfig()
	logger, _ := zap.NewProduction()
	pgConfig, err := cfg.Postgres.GenConfig(logger)
	pool, err := pgxpool.ConnectConfig(context.Background(), pgConfig)
	if err != nil {
		log.Fatal("failed to create db pool")
	}
	defer pool.Close()

	a := arbitrum_nova_adapter.NewArbitrumNovaAdapter()
	a.Run()

	t := harvester.NewTable(pool, a, listener)
	t.Run()

	//token := "0x85F17Cf997934a597031b2E18a9aB6ebD4B9f6a4"
	//wallet := "0x0c3A3040075dd985F141800a1392a0Db81A09cAd"
	//t.AddWalletContract(wallet, token)

	wallet2 := "0x2813fd17ea95b2655a7228383c5236e31090419e"
	token2 := "0xdefa4e8a7bcba345f687a2f1456f5edd9ce97202"
	//wallet2Receiver := "0x3f363b4e038a6e43ce8321c50f3efbf460196d4b"
	//amount := "33190774000000000000000"
	t.AddWalletContract(wallet2, token2)
	t.AddWalletContract("0x3f363b4e038a6e43ce8321c50f3efbf460196d4b", token2)

	time.Sleep(time.Second * 30)
	t.Display()

}

func listener(bcName string, events []*harvester.UpdateEvent, stakes []*harvester.StakeEvent) {
	fmt.Printf("Table Listener: \n")
	//for k, v := range events {
	//	fmt.Printf("%+v %+v %+v %+v \n", k, v.Wallet, v.Contract, v.Amount.String())
	//}
	for k, v := range stakes {
		fmt.Printf("%+v %+v %+v %+v \n", k, v.Wallet, v.OdysseyID, v.Amount.String())
	}
}
