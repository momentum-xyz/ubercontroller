package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/arbitrum_nova_adapter"
	"github.com/momentum-xyz/ubercontroller/logger"
)

var log = logger.L()

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

	harvester.Initialise(context.Background(), log, cfg, pool)

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
