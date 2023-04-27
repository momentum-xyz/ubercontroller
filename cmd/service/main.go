package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/zakaria-chahboun/cute"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/arbitrum_nova_adapter"
	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/pkg/service"
	"github.com/momentum-xyz/ubercontroller/types"
)

var log = logger.L()

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(errors.WithMessage(err, "failed to run service"))
	}
}

func run(ctx context.Context) error {
	cfg := config.GetConfig()

	ctx = context.WithValue(ctx, types.LoggerContextKey, log)
	ctx = context.WithValue(ctx, types.ConfigContextKey, cfg)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tm1 := time.Now()
	pool, err := service.CreateDBConnection(ctx, &cfg.Postgres)
	if err != nil {
		return errors.WithMessage(err, "failed to create db connection")
	}
	defer pool.Close()

	node, err := service.LoadNode(ctx, cfg, pool)
	if err != nil {
		return errors.WithMessage(err, "loading node")
	}
	tm2 := time.Now()
	rand.Seed(time.Now().UnixNano())

	cute.SetTitleColor(cute.BrightGreen)
	cute.SetMessageColor(cute.BrightBlue)
	cute.Println("Node loaded", "Loading time:", tm2.Sub(tm1))

	//harvester.Initialise(ctx, log, cfg, pool)
	//if cfg.Arbitrum.ArbitrumMOMTokenAddress != "" {
	//	arbitrumAdapter := arbitrum_nova_adapter.NewArbitrumNovaAdapter(cfg)
	//	arbitrumAdapter.Run()
	//	if err := harvester.GetInstance().RegisterAdapter(arbitrumAdapter); err != nil {
	//		return errors.WithMessage(err, "failed to register arbitrum adapter")
	//	}
	//}
	//err = harvester.SubscribeAllWallets(ctx, harvester.GetInstance(), cfg, pool)
	//if err != nil {
	//	log.Error(err)
	//}

	/**
	Simplified version of harvester
	*/
	if cfg.Arbitrum.ArbitrumMOMTokenAddress != "" {
		adapter := arbitrum_nova_adapter.NewArbitrumNovaAdapter(cfg)
		adapter.Run()

		t := harvester.NewTable2(pool, adapter, nil)
		t.Run()
	}

	if err := node.Run(); err != nil {
		return errors.WithMessagef(err, "failed to run node: %s", node.GetID())
	}
	defer func() {
		if err := node.Stop(); err != nil {
			log.Error(errors.WithMessagef(err, "failed to stop node: %s", node.GetID()))
		}
	}()

	cute.SetTitleColor(cute.BrightPurple)
	cute.SetMessageColor(cute.BrightBlue)
	cute.Println("Node stopped", "That's all folks!")

	return nil
}
