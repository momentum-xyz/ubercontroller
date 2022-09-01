package world

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/space"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.World = (*World)(nil)

type World struct {
	*space.Space
	ctx context.Context
	log *zap.SugaredLogger
	db  database.DB
}

func NewWorld(id uuid.UUID, db database.DB) *World {
	world := &World{
		db: db,
	}
	world.Space = space.NewSpace(id, db, world)

	return world
}

func (w *World) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	w.ctx = ctx
	w.log = log

	return w.Space.Initialize(ctx)
}

// TODO: implement
func (w *World) Run() error {
	return nil
}

// TODO: implement
func (w *World) Stop() error {
	return nil
}

func (w *World) Load() error {
	entry, err := w.db.SpacesGetSpaceByID(w.ctx, w.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get space by id")
	}

	if err := w.LoadFromEntry(entry, true); err != nil {
		return errors.WithMessage(err, "failed to load from entry")
	}

	universe.GetNode().AddAPIRegister(w)

	return nil
}

func (w *World) Save() error {
	spaces := w.GetSpaces(true)

	entries := make([]*entry.Space, len(spaces))
	for _, space := range spaces {
		entries = append(entries, space.GetEntry())
	}

	if err := w.db.SpacesUpsertSpaces(w.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert spaces")
	}

	return nil
}
