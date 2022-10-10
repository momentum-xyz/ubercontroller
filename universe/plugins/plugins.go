package plugins

import (
	"context"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/plugin"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"

	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
)

var _ universe.Plugins = (*Plugins)(nil)

type Plugins struct {
	ctx     context.Context
	log     *zap.SugaredLogger
	db      database.DB
	plugins *generic.SyncMap[uuid.UUID, universe.Plugin]
}

func (a *Plugins) AddPlugin(plugin universe.Plugin, updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (a *Plugins) AddPlugins(plugins []universe.Plugin, updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (a *Plugins) RemovePlugin(plugin universe.Plugin, updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (a *Plugins) RemovePlugins(plugins []universe.Plugin, updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func NewPlugins(db database.DB) *Plugins {
	return &Plugins{
		db:      db,
		plugins: generic.NewSyncMap[uuid.UUID, universe.Plugin](),
	}
}

func (a *Plugins) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *Plugins) CreatePlugin(pluginId uuid.UUID) (universe.Plugin, error) {

	plugin := plugin.NewPlugin(pluginId, a.db)

	if err := plugin.Initialize(a.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize plugin: %s", pluginId)
	}
	if err := a.AddPlugin(plugin, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add plugin: %s", pluginId)
	}

	return plugin, nil
}

func (a *Plugins) GetPlugin(pluginID uuid.UUID) (universe.Plugin, bool) {
	plugin, ok := a.plugins.Load(pluginID)
	return plugin, ok
}

func (a *Plugins) GetPlugins() map[uuid.UUID]universe.Plugin {
	a.plugins.Mu.RLock()
	defer a.plugins.Mu.RUnlock()

	plugins := make(map[uuid.UUID]universe.Plugin, len(a.plugins.Data))

	for id, plugin := range a.plugins.Data {
		plugins[id] = plugin
	}

	return plugins
}

func (a *Plugins) Load() error {
	a.log.Info("Loading plugins...")

	entries, err := a.db.PluginsGetPlugins(a.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get plugins")
	}

	for i := range entries {
		plugin, err := a.CreatePlugin(*entries[i].PluginID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new plugin: %s", entries[i].PluginID)
		}
		if err := plugin.LoadFromEntry(entries[i]); err != nil {
			return errors.WithMessagef(err, "failed to load plugin from entry: %s", entries[i].PluginID)
		}
		a.plugins.Store(*entries[i].PluginID, plugin)
	}

	universe.GetNode().AddAPIRegister(a)

	a.log.Info("Plugins loaded")

	return nil
}

func (a *Plugins) Save() error {
	a.log.Info("Saving plugins...")

	a.plugins.Mu.RLock()
	defer a.plugins.Mu.RUnlock()

	entries := make([]*entry.Plugin, 0, len(a.plugins.Data))
	for _, plugin := range a.plugins.Data {
		entries = append(entries, plugin.GetEntry())
	}

	if err := a.db.PluginsUpsertPlugins(a.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert plugins")
	}

	a.log.Info("Plugins saved")

	return nil
}
