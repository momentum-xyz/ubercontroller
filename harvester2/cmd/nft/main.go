package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester2"
	"github.com/momentum-xyz/ubercontroller/harvester2/arbitrum_nova_adapter2"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger, _ := zap.NewProduction()
	pgConfig, err := cfg.Postgres.GenConfig(logger)
	if err != nil {
		log.Fatal(err)
	}
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

	mom := common.HexToAddress("0x567d4e8264dC890571D5392fDB9fbd0e3FCBEe56")
	nft := common.HexToAddress("0x97E0B10D89a494Eb5cfFCc72853FB0750BD64AcD")
	stake := common.HexToAddress("0xe9C6d7Cd04614Dde6Ca68B62E6fbf23AC2ECe2F8")
	_ = mom
	_ = nft
	_ = stake

	w04 := common.HexToAddress("0xA058Aa2fCf33993e17D074E6843202E7C94bf267")
	w78 := common.HexToAddress("0x78B00B17E7e5619113A4e922BC3c8cb290355043")
	w68 := common.HexToAddress("0x683642c22feDE752415D4793832Ab75EFdF6223c")
	_ = w04
	_ = w78
	_ = w68

	//err = harv.AddTokenContract("arbitrum_nova", (harvester2.Address)(mom))
	//if err != nil {
	//	fmt.Println(err)
	//}

	//err = harv.AddNFTContract("arbitrum_nova", (harvester2.Address)(nft))
	//if err != nil {
	//	fmt.Println(err)
	//}

	//err = harv.AddStakeContract("arbitrum_nova", (harvester2.Address)(stake))
	//if err != nil {
	//	fmt.Println(err)
	//}

	//harv.AddNFTListener("arbitrum_nova", (harvester2.Address)(stake), "a", func(events []*harvester2.NFTData) {
	//	fmt.Println("NFT LISTENER!!!")
	//})

	//err = harv.AddWallet("arbitrum_nova", (harvester2.Address)(w78))
	//err = harv.AddWallet("arbitrum_nova", (harvester2.Address)(w04))
	//err = harv.AddWallet("arbitrum_nova", (harvester2.Address)(w68))
	//if err != nil {
	//	fmt.Println(err)
	//}

	data, err := harv.GetWalletTokenData("arbitrum_nova", (harvester2.Address)(mom), (harvester2.Address)(w68))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Token Balance ", data.TotalAmount.String())

	nftData, err := harv.GetWalletNFTData("arbitrum_nova", (harvester2.Address)(nft), (harvester2.Address)(w78))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("NFT Balance ", nftData)

	time.Sleep(time.Second * 3)

	err = harv.AddTokenListener("arbitrum_nova", (harvester2.Address)(nft), "ttt", func(events []*harvester2.TokenData) {
		log.Println("TOKEN Listener")
	})

	go func() {
		for {
			time.Sleep(time.Second * 15)
			harv.Display("arbitrum_nova")
		}
	}()

	time.Sleep(time.Second * 1000)
}
