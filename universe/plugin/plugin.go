package plugin

import (
	"context"
	"plugin"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/mplugin"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.Plugin = (*Plugin)(nil)

type Plugin struct {
	ctx context.Context
	log *zap.SugaredLogger
	db  database.DB
	mu  sync.RWMutex

	id                uuid.UUID
	meta              *entry.PluginMeta
	options           *entry.PluginOptions
	object            *plugin.Plugin
	newInstance       mplugin.NewInstanceFunction
	definedAttributes *[]string
}

func NewPlugin(id uuid.UUID, db database.DB) *Plugin {
	return &Plugin{
		db:      db,
		id:      id,
		options: new(entry.PluginOptions),
	}
}

func (p *Plugin) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	p.ctx = ctx
	p.log = log

	return nil
}

func (p *Plugin) GetID() uuid.UUID {
	return p.id
}

func (p *Plugin) GetMeta() *entry.PluginMeta {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.meta
}

func (p *Plugin) SetMeta(meta *entry.PluginMeta, updateDB bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if updateDB {
		if err := p.db.GetPluginsDB().UpdatePluginMeta(p.ctx, p.id, meta); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	p.meta = meta

	return nil
}

func (p *Plugin) GetOptions() *entry.PluginOptions {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.options
}

func (p *Plugin) SetOptions(modifyFn modify.Fn[entry.PluginOptions], updateDB bool) (*entry.PluginOptions, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	options, err := modifyFn(p.options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := p.db.GetPluginsDB().UpdatePluginOptions(p.ctx, p.id, options); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	p.options = options

	return options, nil
}

func (p *Plugin) GetEntry() *entry.Plugin {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return &entry.Plugin{
		PluginID: p.id,
		Meta:     p.meta,
		Options:  p.options,
	}
}

func (p *Plugin) LoadFromEntry(entry *entry.Plugin) error {
	if entry.PluginID != p.id {
		return errors.Errorf("plugin ids mismatch: %s != %s", entry.PluginID, p.GetID())
	}

	var err error
	if entry.Options != nil {
		filepath := utils.GetFromAnyMap(*entry.Options, "file", "")
		if filepath != "" {
			if p.object, p.definedAttributes, p.newInstance, err = p.resolveSharedLibrary(filepath); err != nil {
				return errors.WithMessage(err, "failed to resolve shared library")
			}
		}
	}
	if _, err = p.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}
	if err := p.SetMeta(entry.Meta, false); err != nil {
		return errors.WithMessage(err, "failed to set meta")
	}

	if err := p.RegisterAttributes(); err != nil {
		return errors.WithMessage(err, "failed to register attributes")
	}

	return nil
}

func (p *Plugin) RegisterAttributes() error {
	//TODO: register list of attributes plugins uses
	return nil
}

func (p *Plugin) resolveSharedLibrary(filename string) (*plugin.Plugin, *[]string, mplugin.NewInstanceFunction, error) {
	obj, err := plugin.Open(filename)
	if err != nil {
		return nil, nil, nil, errors.WithMessage(err, "failed to load plugin binary")
	}
	v, err := obj.Lookup("AttributesList")
	if err != nil {
		return nil, nil, nil, errors.WithMessage(err, "plugin object does not have AttributeList")
	}
	attrs, ok := v.(*[]string)
	if !ok {
		return nil, nil, nil, errors.WithMessage(err, "plugin's AttributeList has wrong type")
	}

	v, err = obj.Lookup("NewInstance")
	if err != nil {
		return nil, nil, nil, errors.WithMessage(err, "plugin object does not have NewInstance function")
	}
	newFunc := v.(mplugin.NewInstanceFunction)

	return obj, attrs, newFunc, nil
}
