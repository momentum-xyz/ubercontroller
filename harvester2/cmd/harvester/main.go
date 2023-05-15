package main

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester2"
	"github.com/momentum-xyz/ubercontroller/harvester2/arbitrum_nova_adapter2"
	"go.uber.org/zap"
	"log"
)

func main() {
	cfg := config.GetConfig()
	logger, _ := zap.NewProduction()
	pgConfig, err := cfg.Postgres.GenConfig(logger)
	pool, err := pgxpool.ConnectConfig(context.TODO(), pgConfig)
	if err != nil {
		log.Fatal("failed to create db pool")
	}
	defer pool.Close()

	harv := harvester2.NewHarvester(pool)

	a := arbitrum_nova_adapter2.NewArbitrumNovaAdapter(cfg)
	a.Run()

	err = harv.RegisterAdapter(a)
	if err != nil {
	}

	walletAddress := common.HexToAddress("0x683642c22feDE752415D4793832Ab75EFdF6223c")
	wallet := (*harvester2.Address)(&walletAddress)
	harv.AddWallet("arbitrum_nova", wallet)
}
