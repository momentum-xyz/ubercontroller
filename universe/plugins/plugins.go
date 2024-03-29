package plugins

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/plugin"
)

var _ universe.Plugins = (*Plugins)(nil)

type Plugins struct {
	ctx     types.LoggerContext
	log     *zap.SugaredLogger
	db      database.DB
	plugins *generic.SyncMap[umid.UMID, universe.Plugin]
}

func NewPlugins(db database.DB) *Plugins {
	return &Plugins{
		db:      db,
		plugins: generic.NewSyncMap[umid.UMID, universe.Plugin](0),
	}
}

func (p *Plugins) Initialize(ctx types.LoggerContext) error {
	p.ctx = ctx
	p.log = ctx.Logger()

	return nil
}

func (p *Plugins) CreatePlugin(pluginId umid.UMID) (universe.Plugin, error) {
	plugin := plugin.NewPlugin(pluginId, p.db)

	if err := plugin.Initialize(p.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize plugin: %s", pluginId)
	}
	if err := p.AddPlugin(plugin, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add plugin: %s", pluginId)
	}

	return plugin, nil
}

func (p *Plugins) FilterPlugins(predicateFn universe.PluginsFilterPredicateFn) map[umid.UMID]universe.Plugin {
	return p.plugins.Filter(predicateFn)
}

func (p *Plugins) GetPlugin(pluginID umid.UMID) (universe.Plugin, bool) {
	return p.plugins.Load(pluginID)
}

func (p *Plugins) GetPlugins() map[umid.UMID]universe.Plugin {
	p.plugins.Mu.RLock()
	defer p.plugins.Mu.RUnlock()

	plugins := make(map[umid.UMID]universe.Plugin, len(p.plugins.Data))

	for id, plugin := range p.plugins.Data {
		plugins[id] = plugin
	}

	return plugins
}

func (p *Plugins) AddPlugin(plugin universe.Plugin, updateDB bool) error {
	p.plugins.Mu.Lock()
	defer p.plugins.Mu.Unlock()

	if updateDB {
		if err := p.db.GetPluginsDB().UpsertPlugin(p.ctx, plugin.GetEntry()); err != nil {
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

func (p *Plugins) RemovePlugin(plugin universe.Plugin, updateDB bool) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugins) RemovePlugins(plugins []universe.Plugin, updateDB bool) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugins) Load() error {
	p.log.Info("Loading plugins...")

	entries, err := p.db.GetPluginsDB().GetPlugins(p.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get plugins")
	}

	group, _ := errgroup.WithContext(p.ctx)
	for i := range entries {
		pluginEntry := entries[i]

		group.Go(
			func() error {
				plugin, err := p.CreatePlugin(pluginEntry.PluginID)
				if err != nil {
					return errors.WithMessagef(err, "failed to create new plugin: %s", pluginEntry.PluginID)
				}
				if err := plugin.LoadFromEntry(pluginEntry); err != nil {
					return errors.WithMessagef(err, "failed to load plugin from entry: %s", pluginEntry.PluginID)
				}

				return nil
			},
		)
	}
	if err := group.Wait(); err != nil {
		return err
	}

	universe.GetNode().AddAPIRegister(p)

	p.log.Infof("Plugins loaded: %d", p.plugins.Len())

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

	if err := p.db.GetPluginsDB().UpsertPlugins(p.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert plugins")
	}

	p.log.Infof("Plugins saved: %d", len(p.plugins.Data))

	return nil
}
