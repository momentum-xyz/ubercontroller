package world

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/mplugin"

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
	ctx              context.Context
	log              *zap.SugaredLogger
	db               database.DB
	pluginController *mplugin.PluginController
	//corePluginInstance  mplugin.PluginInstance
	corePluginInterface mplugin.PluginInterface
}

func NewWorld(id uuid.UUID, db database.DB) *World {
	world := &World{
		db: db,
	}
	world.Space = space.NewSpace(id, db, world)
	world.pluginController = mplugin.NewPluginController(id)
	//world.corePluginInstance, _ = world.pluginController.AddPlugin(world.GetID(), world.corePluginInitFunc)
	world.pluginController.AddPlugin(world.GetID(), world.corePluginInitFunc)
	return world
}

func (w *World) corePluginInitFunc(pi mplugin.PluginInterface) (mplugin.PluginInstance, error) {
	instance := CorePluginInstance{PluginInterface: pi}
	w.corePluginInterface = pi
	return instance, nil
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
	w.log.Infof("Loading world: %s", w.GetID())

	entry, err := w.db.SpacesGetSpaceByID(w.ctx, w.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get space by id")
	}

	if err := w.LoadFromEntry(entry, true); err != nil {
		return errors.WithMessage(err, "failed to load from entry")
	}

	universe.GetNode().AddAPIRegister(w)

	w.log.Infof("World loaded: %s", w.GetID())

	return nil
}

func (w *World) Save() error {
	w.log.Infof("Saving world: %s", w.GetID())

	spaces := w.GetSpaces(true)

	entries := make([]*entry.Space, 0, len(spaces))
	for _, space := range spaces {
		entries = append(entries, space.GetEntry())
	}

	if err := w.db.SpacesUpsertSpaces(w.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert spaces")
	}

	w.log.Infof("World saved: %s", w.GetID())

	return nil
}
