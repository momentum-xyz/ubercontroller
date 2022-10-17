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

func NewPlugins(db database.DB) *Plugins {
	return &Plugins{
		db:      db,
		plugins: generic.NewSyncMap[uuid.UUID, universe.Plugin](),
	}
}

func (p *Plugins) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	p.ctx = ctx
	p.log = log

	return nil
}

func (p *Plugins) CreatePlugin(pluginId uuid.UUID) (universe.Plugin, error) {
	plugin := plugin.NewPlugin(pluginId, p.db)

	if err := plugin.Initialize(p.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize plugin: %s", pluginId)
	}
	if err := p.AddPlugin(plugin, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add plugin: %s", pluginId)
	}

	return plugin, nil
}

func (p *Plugins) GetPlugin(pluginID uuid.UUID) (universe.Plugin, bool) {
	plugin, ok := p.plugins.Load(pluginID)
	return plugin, ok
}

func (p *Plugins) GetPlugins() map[uuid.UUID]universe.Plugin {
	p.plugins.Mu.RLock()
	defer p.plugins.Mu.RUnlock()

	plugins := make(map[uuid.UUID]universe.Plugin, len(p.plugins.Data))

	for id, plugin := range p.plugins.Data {
		plugins[id] = plugin
	}

	return plugins
}

func (p *Plugins) AddPlugin(plugin universe.Plugin, updateDB bool) error {
	p.plugins.Mu.Lock()
	defer p.plugins.Mu.Unlock()

	if updateDB {
		if err := p.db.PluginsUpsertPlugin(p.ctx, plugin.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	p.plugins.Data[plugin.GetID()] = plugin

	return nil
}

func (p *Plugins) AddPlugins(plugins []universe.Plugin, updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (p *Plugins) RemovePlugin(plugin universe.Plugin, updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (p *Plugins) RemovePlugins(plugins []universe.Plugin, updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (p *Plugins) Load() error {
	p.log.Info("Loading plugins...")

	entries, err := p.db.PluginsGetPlugins(p.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get plugins")
	}

	for i := range entries {
		plugin, err := p.CreatePlugin(entries[i].PluginID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new plugin: %s", entries[i].PluginID)
		}
		if err := plugin.LoadFromEntry(entries[i]); err != nil {
			return errors.WithMessagef(err, "failed to load plugin from entry: %s", entries[i].PluginID)
		}
		p.plugins.Store(entries[i].PluginID, plugin)
	}

	universe.GetNode().AddAPIRegister(p)

	p.log.Info("Plugins loaded")

	return nil
}

func (p *Plugins) Save() error {
	p.log.Info("Saving plugins...")

	p.plugins.Mu.RLock()
	defer p.plugins.Mu.RUnlock()

	entries := make([]*entry.Plugin, 0, len(p.plugins.Data))
	for _, plugin := range p.plugins.Data {
		entries = append(entries, plugin.GetEntry())
	}

	if err := p.db.PluginsUpsertPlugins(p.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert plugins")
	}

	p.log.Info("Plugins saved")

	return nil
}
