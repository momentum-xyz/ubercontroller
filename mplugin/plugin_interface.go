package mplugin

import (
	"github.com/google/uuid"
)

type internalPluginInterface struct {
	id       PluginID
	secretId uuid.UUID
	worldId  uuid.UUID
	pc       *PluginController
}

func (p internalPluginInterface) getPluginController() *PluginController {
	return p.pc
}
func (p internalPluginInterface) GetSecret() uuid.UUID {
	return p.secretId
}

func (p internalPluginInterface) GetWorld() uuid.UUID {
	return p.worldId
}

func (p internalPluginInterface) GetId() PluginID {
	return p.id
}

func PluginSubscribeHook[A any](p PluginInterface, name string, hook func(A) error) (HookID, error) {
	return subscribeHook(p, name, hook)
}

func (p internalPluginInterface) UnsubscribeAllHooks() error {
	return p.pc.UnsubscribeAllHooks(p)
}

func (p internalPluginInterface) UnsubscribeHooksByName(name string) error {
	return p.pc.UnsubscribeHooksByName(p, name)
}
func (p internalPluginInterface) UnsubscribeHookByHookId(id HookID) error {
	return p.pc.UnsubscribeHookByHookId(p, id)
}
