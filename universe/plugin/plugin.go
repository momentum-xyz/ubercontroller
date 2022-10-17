package plugin

import (
	"context"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/mplugin"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"plugin"
	"sync"
)

type Plugin struct {
	ctx context.Context
	log *zap.SugaredLogger
	db  database.DB
	mu  sync.RWMutex

	id                uuid.UUID
	name              string
	description       *string
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

func (p *Plugin) GetName() string {
	return p.name
}

func (p *Plugin) GetOptions() *entry.PluginOptions {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.options
}

func (p *Plugin) GetDescription() *string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.description
}

func (p *Plugin) SetName(name string, updateDB bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if updateDB {
		if err := p.db.PluginsUpdatePluginName(p.ctx, p.id, name); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	p.name = name

	return nil
}

func (p *Plugin) SetOptions(modifyFn modify.Fn[entry.PluginOptions], updateDB bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	options := modifyFn(p.options)

	if updateDB {
		if err := p.db.PluginsUpdatePluginOptions(p.ctx, p.id, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	p.options = options

	return nil
}

func (p *Plugin) SetDescription(description *string, updateDB bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if updateDB {
		if err := p.db.PluginsUpdatePluginDescription(p.ctx, p.id, description); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	p.description = description

	return nil
}

func (p *Plugin) GetEntry() *entry.Plugin {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return &entry.Plugin{
		PluginID:    p.id,
		PluginName:  p.name,
		Description: p.description,
		Options:     p.options,
	}
}

func (p *Plugin) LoadFromEntry(entry *entry.Plugin) error {
	if entry.PluginID != p.id {
		return errors.Errorf("plugin ids mismatch: %s != %s", entry.PluginID, p.id)
	}

	p.id = entry.PluginID

	var err error
	if entry.Options != nil {
		if p.object, p.definedAttributes, p.newInstance, err = p.resolveSharedLibrary(entry.Options.File); err != nil {
			return errors.WithMessage(err, "failed to resolve shared library")
		}
	}

	if err = p.SetName(entry.PluginName, false); err != nil {
		return errors.WithMessage(err, "failed to set name")
	}
	if err = p.SetDescription(entry.Description, false); err != nil {
		return errors.WithMessage(err, "failed to set description")
	}
	if err = p.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
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
