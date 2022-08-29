package worlds

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/types/generics"
	"github.com/momentum-xyz/controller/universe"
)

var _ universe.Worlds = (*Worlds)(nil)

type Worlds struct {
	worlds *generics.SyncMap[uuid.UUID, universe.World]
}

func NewWorlds() *Worlds {
	return &Worlds{
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

func (w *Worlds) Load() error {
	return errors.Errorf("implement me")
}
