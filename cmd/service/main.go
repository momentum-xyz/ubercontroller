package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/zakaria-chahboun/cute"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/pkg/service"
	"github.com/momentum-xyz/ubercontroller/types"
)

// Build version, overridden with flag during build.
var version = "devel"

func main() {
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(errors.WithMessage(err, "failed to run service"))
	}
}

func run(ctx context.Context) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.WithMessage(err, "failed to get config")
	}

	log := logger.L()
	log.Debugf("Version: %s", version)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	nodeCtx, err := types.NewNodeContext(ctx, log, cfg)
	if err != nil {
		return errors.WithMessage(err, "failed to create context")
	}

	tm1 := time.Now()
	pool, err := service.CreateDBConnection(nodeCtx, &cfg.Postgres)
	if err != nil {
		return errors.WithMessage(err, "failed to create db connection")
	}
	defer pool.Close()

	node, err := service.LoadNode(nodeCtx, cfg, pool)
	if err != nil {
		return errors.WithMessage(err, "loading node")
	}
	tm2 := time.Now()
	rand.Seed(time.Now().UnixNano())

	cute.SetTitleColor(cute.BrightGreen)
	cute.SetMessageColor(cute.BrightBlue)
	cute.Println("Node loaded", "Loading time:", tm2.Sub(tm1))

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
