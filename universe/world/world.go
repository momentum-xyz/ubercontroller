package world

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/space"
)

var _ universe.World = (*World)(nil)

type World struct {
	*space.Space
	db  database.DB
	log *zap.SugaredLogger
}

func NewWorld(id uuid.UUID, db database.DB) *World {
	world := &World{
		db: db,
	}
	world.Space = space.NewSpace(id, db, world)

	return world
}

func (w *World) Initialize(ctx context.Context) error {
	log, ok := ctx.Value(types.ContextLoggerKey).(*zap.SugaredLogger)
	if !ok {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	w.log = log

	return w.Space.Initialize(ctx)
}

func (w *World) Run(ctx context.Context) error {
	return errors.Errorf("implement me")
}

func (w *World) Stop() error {
	return errors.Errorf("implement me")
}
