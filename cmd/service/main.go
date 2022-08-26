package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/momentum-xyz/controller/logger"
	"github.com/momentum-xyz/controller/types"
	"github.com/momentum-xyz/controller/universe/world"
)

func main() {
	log := logger.L()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wctx := context.WithValue(ctx, types.ContextLoggerKey, log)

	world := world.NewWorld(uuid.New())
	fmt.Printf("HERE 1: %+v\n", world.Initialize(wctx))
}
