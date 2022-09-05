package worlds

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/world"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Worlds = (*Worlds)(nil)

type Worlds struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	db     database.DB
	worlds *generic.SyncMap[uuid.UUID, universe.World]
}

func NewWorlds(db database.DB) *Worlds {
	return &Worlds{
		db:     db,
		worlds: generic.NewSyncMap[uuid.UUID, universe.World](),
	}
}

func (w *Worlds) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	w.ctx = ctx
	w.log = log

	return nil
}

func (w *Worlds) GetWorld(worldID uuid.UUID) (universe.World, bool) {
	world, ok := w.worlds.Load(worldID)
	return world, ok
}

func (w *Worlds) GetWorlds() map[uuid.UUID]universe.World {
	worlds := make(map[uuid.UUID]universe.World)

	w.worlds.Mu.RLock()
	defer w.worlds.Mu.RUnlock()

	for id, world := range w.worlds.Data {
		worlds[id] = world
	}

	return worlds
}

func (w *Worlds) AddWorld(world universe.World, updateDB bool) error {
	w.worlds.Mu.Lock()
	defer w.worlds.Mu.Unlock()

	if _, ok := w.worlds.Data[world.GetID()]; ok {
		return errors.Errorf("world already exists")
	}

	if updateDB {
		if err := world.Save(); err != nil {
			return errors.WithMessage(err, "failed to save world")
		}
	}

	w.worlds.Data[world.GetID()] = world

	return nil
}

func (w *Worlds) AddWorlds(worlds []universe.World, updateDB bool) error {
	w.worlds.Mu.Lock()
	defer w.worlds.Mu.Unlock()

	for i := range worlds {
		if _, ok := w.worlds.Data[worlds[i].GetID()]; ok {
			return errors.Errorf("world already exists: %s", worlds[i].GetID())
		}
	}

	if updateDB {
		group, _ := errgroup.WithContext(w.ctx)

		for i := range worlds {
			world := worlds[i]

			group.Go(func() error {
				if err := world.Save(); err != nil {
					return errors.WithMessagef(err, "failed to save world: %s", world.GetID())
				}
				return nil
			})
		}

		if err := group.Wait(); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range worlds {
		w.worlds.Data[worlds[i].GetID()] = worlds[i]
	}

	return nil
}

func (w *Worlds) RemoveWorld(world universe.World, updateDB bool) error {
	w.worlds.Mu.Lock()
	defer w.worlds.Mu.Unlock()

	if _, ok := w.worlds.Data[world.GetID()]; !ok {
		return errors.Errorf("world not found")
	}

	if updateDB {
		spaces := world.GetSpaces(true)
		ids := make([]uuid.UUID, len(spaces))
		for _, space := range spaces {
			ids = append(ids, space.GetID())
		}
		if err := w.db.SpacesRemoveSpacesByIDs(w.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to remove spaces by ids")
		}
	}

	delete(w.worlds.Data, world.GetID())

	return nil
}

func (w *Worlds) RemoveWorlds(worlds []universe.World, updateDB bool) error {
	w.worlds.Mu.Lock()
	defer w.worlds.Mu.Unlock()

	for i := range worlds {
		if _, ok := w.worlds.Data[worlds[i].GetID()]; !ok {
			return errors.Errorf("world not found: %s", worlds[i].GetID())
		}
	}

	if updateDB {
		group, _ := errgroup.WithContext(w.ctx)

		for i := range worlds {
			world := worlds[i]

			group.Go(func() error {
				spaces := world.GetSpaces(true)
				ids := make([]uuid.UUID, len(spaces))
				for i := range spaces {
					ids = append(ids, spaces[i].GetID())
				}

				if err := w.db.SpacesRemoveSpacesByIDs(w.ctx, ids); err != nil {
					return errors.WithMessagef(err, "failed to remove spaces by ids: %s", world.GetID())
				}

				return nil
			})
		}

		if err := group.Wait(); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range worlds {
		delete(w.worlds.Data, worlds[i].GetID())
	}

	return nil
}

func (w *Worlds) Run() error {
	w.worlds.Mu.RLock()
	defer w.worlds.Mu.RUnlock()

	for _, world := range w.worlds.Data {
		world := world

		go func() {
			if err := world.Run(); err != nil {
				w.log.Error(errors.WithMessagef(err, "failed to run world: %s", world.GetID()))
			}
		}()
	}

	return nil
}

func (w *Worlds) Stop() error {
	w.worlds.Mu.RLock()
	defer w.worlds.Mu.RUnlock()

	group, _ := errgroup.WithContext(w.ctx)

	for _, world := range w.worlds.Data {
		world := world

		group.Go(func() error {
			if err := world.Stop(); err != nil {
				return errors.WithMessagef(err, "failed to stop world: %s", world.GetID())
			}
			return nil
		})
	}

	return group.Wait()
}

func (w *Worlds) Load() error {
	w.log.Info("Loading worlds...")

	worldIDs, err := w.db.WorldsGetWorldIDs(w.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get world ids from db")
	}

	group, _ := errgroup.WithContext(w.ctx)

	for i := range worldIDs {
		worldID := worldIDs[i]

		group.Go(func() error {
			world := world.NewWorld(worldID, w.db)

			if err := world.Initialize(w.ctx); err != nil {
				return errors.WithMessagef(err, "failed to initialize world: %s", world.GetID())
			}
			if err := world.Load(); err != nil {
				return errors.WithMessagef(err, "failed to load world: %s", world.GetID())
			}

			w.worlds.Store(world.GetID(), world)

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	universe.GetNode().AddAPIRegister(w)

	return nil
}

func (w *Worlds) Save() error {
	w.worlds.Mu.RLock()
	defer w.worlds.Mu.RUnlock()

	group, _ := errgroup.WithContext(w.ctx)

	for _, world := range w.worlds.Data {
		world := world

		group.Go(func() error {
			if err := world.Save(); err != nil {
				return errors.WithMessagef(err, "failed to save world: %s", world.GetID())
			}
			return nil
		})
	}

	return group.Wait()
}
