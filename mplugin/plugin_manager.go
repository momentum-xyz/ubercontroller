package mplugin

import "github.com/google/uuid"

type PluginManager struct {
	loadedPlugins map[uuid.UUID]any
}

func NewPluginManager() *PluginManager {
	pm := PluginManager{}
	pm.loadedPlugins = make(map[uuid.UUID]any)
	return &pm
}
