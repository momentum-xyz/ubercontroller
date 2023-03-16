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

	a := ethereum_adapter.NewEthereumAdapter()
	a.Run()

	t := harvester.NewTable(pool, a, func(p any) {})

	token := "0x85F17Cf997934a597031b2E18a9aB6ebD4B9f6a4"
	wallet := "0x0c3A3040075dd985F141800a1392a0Db81A09cAd"
	t.AddWalletContract(wallet, token)

	time.Sleep(time.Second * 3)
	fmt.Println(t)
}
