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

	logger, _ := zap.NewProduction()
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to get config: %s", err))
	}
	pgConfig, err := cfg.Postgres.GenConfig(logger)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to create db config: %s", err))
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), pgConfig)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to create db pool: %s", err))
	}
	defer pool.Close()

	harvester.Initialise(context.Background(), logger.Sugar(), cfg, pool)

	// ** Ethereum Adapter
	adapter := arbitrum_nova_adapter.NewArbitrumNovaAdapter(cfg)
	adapter.Run()
	_, _, _ = adapter.GetInfo()
	if err := harvester.GetInstance().RegisterAdapter(adapter); err != nil {
		log.Fatal(err)
	}

	err = harvester.SubscribeAllWallets(context.Background(), harvester.GetInstance(), cfg, pool)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
	fmt.Println(err)
}
