package worlds

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/world"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.Worlds = (*Worlds)(nil)

type Worlds struct {
	ctx    types.NodeContext
	log    *zap.SugaredLogger
	cfg    *config.Config
	db     database.DB
	worlds *generic.SyncMap[umid.UMID, universe.World]
}

func NewWorlds(db database.DB) *Worlds {
	return &Worlds{
		db:     db,
		worlds: generic.NewSyncMap[umid.UMID, universe.World](0),
	}
}

func (w *Worlds) Initialize(ctx types.NodeContext) error {
	w.ctx = ctx
	w.log = ctx.Logger()
	w.cfg = ctx.Config()

	return nil
}

func (w *Worlds) CreateWorld(worldID umid.UMID) (universe.World, error) {
	world := world.NewWorld(worldID, w.db)

	if err := world.Initialize(w.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize world: %s", worldID)
	}
	if err := w.AddWorld(world, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add world: %s", worldID)
	}

	return world, nil
}

func (w *Worlds) FilterWorlds(predicateFn universe.WorldsFilterPredicateFn) map[umid.UMID]universe.World {
	return w.worlds.Filter(predicateFn)
}

func (w *Worlds) GetWorld(worldID umid.UMID) (universe.World, bool) {
	world, ok := w.worlds.Load(worldID)
	return world, ok
}

func (w *Worlds) GetWorlds() map[umid.UMID]universe.World {
	w.worlds.Mu.RLock()
	defer w.worlds.Mu.RUnlock()

	worlds := make(map[umid.UMID]universe.World, len(w.worlds.Data))

	for id, world := range w.worlds.Data {
		worlds[id] = world
	}

	return worlds
}

func (w *Worlds) AddWorld(world universe.World, updateDB bool) error {
	node := universe.GetNode()

	w.worlds.Mu.Lock()
	defer w.worlds.Mu.Unlock()

	if err := world.SetParent(node, false); err != nil {
		return errors.WithMessagef(err, "failed to set parent %s to world %s", node.GetID(), world.GetID())
	}

	if updateDB {
		if err := world.Save(); err != nil {
			return errors.WithMessage(err, "failed to save world")
		}
	}

	w.worlds.Data[world.GetID()] = world

	return node.AddObjectToAllObjects(world.ToObject())
}

// TODO: optimize
func (w *Worlds) AddWorlds(worlds []universe.World, updateDB bool) error {
	for _, world := range worlds {
		if err := w.AddWorld(world, updateDB); err != nil {
			return errors.WithMessagef(err, "failed to add world: %s", world.GetID())
		}
	}
	return nil
}

// TODO: introduce "helper.RemoveWorld()" method and fix this one
func (w *Worlds) RemoveWorld(world universe.World, updateDB bool) (bool, error) {
	w.worlds.Mu.Lock()
	defer w.worlds.Mu.Unlock()

	if _, ok := w.worlds.Data[world.GetID()]; !ok {
		return false, nil
	}

	if updateDB {
		panic("not implemented")
	}

	delete(w.worlds.Data, world.GetID())

	return true, nil
}

// TODO: introduce "helper.RemoveWorld()" method and fix this one
func (w *Worlds) RemoveWorlds(worlds []universe.World, updateDB bool) (bool, error) {
	w.worlds.Mu.Lock()
	defer w.worlds.Mu.Unlock()

	for _, world := range worlds {
		if _, ok := w.worlds.Data[world.GetID()]; !ok {
			return false, nil
		}
	}

	if updateDB {
		panic("not implemented")
	}

	for _, world := range worlds {
		delete(w.worlds.Data, world.GetID())
	}

	return true, nil
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

	worldIDs, err := w.db.GetWorldsDB().GetAllWorldIDs(w.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get world ids from db")
	}

	butcher := generic.NewButcher(worldIDs)
	if err := butcher.HandleItems(
		int(w.cfg.Postgres.MAXCONNS), // modify batchSize when database consumption while loading will be changed
		func(worldID umid.UMID) error {
			world, err := w.CreateWorld(worldID)
			if err != nil {
				return errors.WithMessagef(err, "failed to create new world: %s", worldID)
			}
			if err := world.Load(); err != nil {
				return errors.WithMessagef(err, "failed to load world: %s", worldID)
			}
			return nil
		},
	); err != nil {
		return errors.WithMessage(err, "failed to load worlds")
	}

	universe.GetNode().AddAPIRegister(w)

	w.log.Infof("Worlds loaded: %d", butcher.Len())

	return nil
}

func (w *Worlds) Save() error {
	w.log.Info("Saving worlds...")

	w.worlds.Mu.RLock()
	defer w.worlds.Mu.RUnlock()

	var count int
	var errs *multierror.Error
	for _, world := range w.worlds.Data {
		if err := world.Save(); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to save world: %s", world.GetID()))
			continue
		}
		count++
	}

	w.log.Infof("Worlds saved: %d", count)

	return errs.ErrorOrNil()
}
