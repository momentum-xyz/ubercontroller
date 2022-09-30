package mplugin

import (
	"github.com/google/uuid"
)

type internalPluginInstance struct {
	id       PluginID
	secretId uuid.UUID
	worldId  uuid.UUID
	pc       *PluginController
}

func (p internalPluginInstance) getPluginController() *PluginController {
	return p.pc
}
func (p internalPluginInstance) GetSecret() uuid.UUID {
	return p.secretId
}

func (p internalPluginInstance) GetWorld() uuid.UUID {
	return p.worldId
}

func (p internalPluginInstance) GetId() PluginID {
	return p.id
}

func PluginSubscribeHook[A any](p PluginInstance, name string, hook func(A) error) (HookID, error) {
	return subscribeHook(p, name, hook)
}

func (p internalPluginInstance) UnsubscribeAllHooks() error {
	return p.pc.UnsubscribeAllHooks(p)
}

func (p internalPluginInstance) UnsubscribeHooksByName(name string) error {
	return p.pc.UnsubscribeHooksByName(p, name)
}
func (p internalPluginInstance) UnsubscribeHookByHookId(id HookID) error {
	return p.pc.UnsubscribeHookByHookId(p, id)
}
