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
	node, err := service.LoadNode(ctx, cfg)
	if err != nil {
		return errors.WithMessage(err, "loading node")
	}
	tm2 := time.Now()
	rand.Seed(time.Now().UnixNano())

	cute.SetTitleColor(cute.BrightGreen)
	cute.SetMessageColor(cute.BrightBlue)
	cute.Println("Node loaded", "Loading time:", tm2.Sub(tm1))

	defer func() {
		if err := node.Stop(); err != nil {
			log.Error(errors.WithMessagef(err, "failed to stop node: %s", node.GetID()))
		}
	}()

	if err := node.Run(); err != nil {
		return errors.WithMessagef(err, "failed to run node: %s", node.GetID())
	}

	cute.SetTitleColor(cute.BrightPurple)
	cute.SetMessageColor(cute.BrightBlue)
	cute.Println("Node stopped", "That's all folks!")

	return nil
}
