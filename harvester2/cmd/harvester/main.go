package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester2"
	"github.com/momentum-xyz/ubercontroller/harvester2/arbitrum_nova_adapter2"
	"go.uber.org/zap"
	"log"
	"time"
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
		fmt.Println(err)
	}

	momAddress := common.HexToAddress("0x567d4e8264dC890571D5392fDB9fbd0e3FCBEe56")
	mom := (*harvester2.Address)(&momAddress)
	err = harv.AddTokenContract("arbitrum_nova", mom)
	if err != nil {
		fmt.Println(err)
	}

	nftAddress := common.HexToAddress("0x97E0B10D89a494Eb5cfFCc72853FB0750BD64AcD")
	nft := (*harvester2.Address)(&nftAddress)
	err = harv.AddNFTContract("arbitrum_nova", nft)
	if err != nil {
		fmt.Println(err)
	}

	stakeAddress := common.HexToAddress("0x047C0A154271498ee718162b718b3D4F464855e0")
	stake := (*harvester2.Address)(&stakeAddress)
	err = harv.AddStakeContract("arbitrum_nova", stake)
	if err != nil {
		fmt.Println(err)
	}

	walletAddress := common.HexToAddress("0x683642c22feDE752415D4793832Ab75EFdF6223c")
	wallet := (*harvester2.Address)(&walletAddress)
	err = harv.AddWallet("arbitrum_nova", wallet)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second)
}
