package worlds

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/world"
)

var _ universe.Worlds = (*Worlds)(nil)

type Worlds struct {
	db     database.DB
	worlds *generics.SyncMap[uuid.UUID, universe.World]
}

func NewWorlds(db database.DB) *Worlds {
	return &Worlds{
		db:     db,
		worlds: generics.NewSyncMap[uuid.UUID, universe.World](),
	}
}

func (w *Worlds) Initialize(ctx context.Context) error {
	return nil
}

func (w *Worlds) GetWorld(worldID uuid.UUID) (universe.World, bool) {
	world, ok := w.worlds.Load(worldID)
	return world, ok
}

// GetWorlds returns existing sync map with all stored worlds.
func (w *Worlds) GetWorlds() *generics.SyncMap[uuid.UUID, universe.World] {
	return w.worlds
}

func (w *Worlds) Load(ctx context.Context) error {
	worlds, err := w.db.WorldsGetWorlds(ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get world entries")
	}

	group, gctx := errgroup.WithContext(ctx)

	for i := range worlds {
		entry := worlds[i]

		group.Go(func() error {
			world := world.NewWorld(*entry.SpaceID, w.db)

			if err := world.LoadFromEntry(gctx, &entry, true); err != nil {
				return errors.WithMessagef(err, "failed to load world from entry: %s", world.GetID())
			}

			w.worlds.Store(world.GetID(), world)

			return nil
		})
	}

	return group.Wait()
}
