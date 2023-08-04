package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/zakaria-chahboun/cute"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database/migrations"
	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/pkg/service"
	"github.com/momentum-xyz/ubercontroller/types"
)

// Build version, overridden with flag during build.
var version = "devel"

func main() {
	printVersion := flag.Bool("version", false, "Print version")
	migrateOnly := flag.Bool("migrate", false, "Only migrate the database.")
	migrateSteps := flag.Int("steps", 0, `Migration steps, negative: down, positive: up. Default 0: all the way up.`)
	flag.Parse()
	if *printVersion {
		fmt.Printf("Controller version: %s\n", version)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if *migrateOnly {
		if err := runMigration(ctx, *migrateSteps); err != nil {
			log.Fatal(errors.WithMessage(err, "failed to run migration"))
		}
	} else {
		if err := run(ctx); err != nil {
			log.Fatal(errors.WithMessage(err, "failed to run service"))
		}
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

func runMigration(ctx context.Context, steps int) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.WithMessage(err, "failed to get config")
	}
	log := logger.L()
	nodeCtx, err := types.NewNodeContext(ctx, log, cfg)
	if err != nil {
		return errors.WithMessage(err, "failed to create context")
	}
	if err := migrations.MigrateDatabase(nodeCtx, &cfg.Postgres, steps); err != nil {
		return errors.WithMessage(err, "failed to migrate database")
	}
	return nil
}
