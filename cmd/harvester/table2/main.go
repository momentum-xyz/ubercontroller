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
	fmt.Println("Table 2")

	cfg := config.GetConfig()
	logger, _ := zap.NewProduction()
	pgConfig, err := cfg.Postgres.GenConfig(logger)
	pool, err := pgxpool.ConnectConfig(context.Background(), pgConfig)
	if err != nil {
		log.Fatal("failed to create db pool")
	}
	defer pool.Close()

	a := arbitrum_nova_adapter.NewArbitrumNovaAdapter(cfg)
	a.Run()

	t := harvester.NewTable2(pool, a, listener)
	t.Run()

	time.Sleep(time.Hour)
}

func listener(bcName string, events []*harvester.UpdateEvent, stakeEvents []*harvester.StakeEvent, nftEvent []*harvester.NftEvent) {
	fmt.Printf("Table Listener: \n")
	for k, v := range events {
		fmt.Printf("%+v %+v %+v %+v \n", k, v.Wallet, v.Contract, v.Amount.String())
	}
}
