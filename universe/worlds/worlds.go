package worlds

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
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

func (w *Worlds) GetWorlds() *generics.SyncMap[uuid.UUID, universe.World] {
	return w.worlds
}

func (w *Worlds) Load(ctx context.Context) error {
	return errors.Errorf("implement me")
}
