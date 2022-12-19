package worlds

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/config"
	"sync"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
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
	cfg    *config.Config
	db     database.DB
	worlds *generic.SyncMap[uuid.UUID, universe.World]
}

func NewWorlds(db database.DB) *Worlds {
	return &Worlds{
		db:     db,
		worlds: generic.NewSyncMap[uuid.UUID, universe.World](0),
	}
}

func (w *Worlds) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.Errorf("failed to get config from context: %T", ctx.Value(types.ConfigContextKey))
	}

	w.ctx = ctx
	w.log = log
	w.cfg = cfg

	return nil
}

func (w *Worlds) CreateWorld(worldID uuid.UUID) (universe.World, error) {
	world := world.NewWorld(worldID, w.db)

	if err := world.Initialize(w.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize world: %s", worldID)
	}
	if err := w.AddWorld(world, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add world: %s", worldID)
	}

	return world, nil
}

func (w *Worlds) FilterWorlds(predicateFn universe.WorldsFilterPredicateFn) map[uuid.UUID]universe.World {
	return w.worlds.Filter(predicateFn)
}

func (w *Worlds) GetWorld(worldID uuid.UUID) (universe.World, bool) {
	world, ok := w.worlds.Load(worldID)
	return world, ok
}

func (w *Worlds) GetWorlds() map[uuid.UUID]universe.World {
	w.worlds.Mu.RLock()
	defer w.worlds.Mu.RUnlock()

	worlds := make(map[uuid.UUID]universe.World, len(w.worlds.Data))

	for id, world := range w.worlds.Data {
		worlds[id] = world
	}

	return worlds
}

func (w *Worlds) AddWorld(world universe.World, updateDB bool) error {
	w.worlds.Mu.Lock()
	defer w.worlds.Mu.Unlock()

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
		spaces := world.GetAllSpaces()
		ids := make([]uuid.UUID, 0, len(spaces))
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
				spaces := world.GetAllSpaces()
				ids := make([]uuid.UUID, 0, len(spaces))
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

	var errs *multierror.Error
	for _, world := range w.worlds.Data {
		if err := world.Run(); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to run world: %s", world.GetID()))
		}
		world.SetEnabled(true)
	}

	return errs.ErrorOrNil()
}

func (w *Worlds) Stop() error {
	w.worlds.Mu.RLock()
	defer w.worlds.Mu.RUnlock()

	var errs *multierror.Error
	for _, world := range w.worlds.Data {
		if err := world.Stop(); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to stop world: %s", world.GetID()))
		}
		world.SetEnabled(false)
	}

	return errs.ErrorOrNil()
}

func (w *Worlds) Load() error {
	w.log.Info("Loading worlds...")

	worldIDs, err := w.db.WorldsGetWorldIDs(w.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get world ids from db")
	}
	worldsCount := len(worldIDs)

	// modify batchSize when database consumption per world loading will be changed
	batchSize := int(w.cfg.Postgres.MAXCONNS)
	for len(worldIDs) > 0 {
		batch := worldIDs
		if len(worldIDs) > batchSize {
			batch = worldIDs[:batchSize]
			worldIDs = worldIDs[batchSize:]
		} else {
			worldIDs = nil
		}

		if err := w.loadBatch(batch); err != nil {
			return errors.WithMessage(err, "failed to load worlds batch")
		}
	}

	universe.GetNode().AddAPIRegister(w)

	w.log.Infof("Worlds loaded: %d", worldsCount)

	return nil
}

func (w *Worlds) loadBatch(worldIDs []uuid.UUID) error {
	w.log.Info("Loading worlds batch...")

	group, _ := errgroup.WithContext(w.ctx)
	for i := range worldIDs {
		worldID := worldIDs[i]

		group.Go(func() error {
			world, err := w.CreateWorld(worldID)
			if err != nil {
				return errors.WithMessagef(err, "failed to create new world: %s", worldID)
			}
			if err := world.Load(); err != nil {
				return errors.WithMessagef(err, "failed to load world: %s", worldID)
			}

			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return err
	}

	w.log.Infof("Worlds batch loaded: %d", w.worlds.Len())

	return nil
}

func (w *Worlds) Save() error {
	w.log.Info("Saving worlds...")

	var wg sync.WaitGroup
	var errs *multierror.Error
	var errsMu sync.Mutex

	w.worlds.Mu.RLock()
	defer w.worlds.Mu.RUnlock()

	for i := range w.worlds.Data {
		wg.Add(1)

		go func(world universe.World) {
			defer wg.Done()

			if err := world.Save(); err != nil {
				errsMu.Lock()
				defer errsMu.Unlock()
				errs = multierror.Append(errs, errors.WithMessagef(err, "failed to save world: %s", world.GetID()))
			}
		}(w.worlds.Data[i])
	}

	wg.Wait()

	w.log.Info("Worlds saved")

	return errs.ErrorOrNil()
}
