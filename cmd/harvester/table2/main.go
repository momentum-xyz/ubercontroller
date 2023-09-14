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

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to get config: %s", err)
	}
	logger, _ := zap.NewProduction()
	pgConfig, err := cfg.Postgres.GenConfig(logger)
	if err != nil {
		log.Fatalf("failed to create db config: %s", err)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), pgConfig)
	if err != nil {
		log.Fatalf("failed to create db pool: %s", err)
	}
	defer pool.Close()

	a := arbitrum_nova_adapter.NewArbitrumNovaAdapter(cfg, logger.Sugar())
	a.Run()

	t := harvester.NewTable(pool, a, listener)
	t.Run()

	time.Sleep(time.Hour)
}

func listener(bcName string, events []*harvester.UpdateEvent, stakeEvents []*harvester.StakeEvent, nftEvent []*harvester.NftEvent) error {
	fmt.Printf("Table Listener: \n")
	for k, v := range events {
		fmt.Printf("%+v %+v %+v %+v \n", k, v.Wallet, v.Contract, v.Amount.String())
	}

	return nil
}
