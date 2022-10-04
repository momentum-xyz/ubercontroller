package world

import "github.com/momentum-xyz/ubercontroller/mplugin"

type CorePluginInstance struct {
	PluginInterface mplugin.PluginInterface
}

func (instance CorePluginInstance) Init() error {
	return nil
}

func (instance CorePluginInstance) Run() error {
	return nil
}

func (instance CorePluginInstance) Destroy() error {
	return nil
}
