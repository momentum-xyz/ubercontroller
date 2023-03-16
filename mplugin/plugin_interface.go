package mplugin

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type internalPluginInterface struct {
	id       PluginID
	secretId umid.UMID
	worldId  umid.UMID
	pc       *PluginController
}

func (p internalPluginInterface) getPluginController() *PluginController {
	return p.pc
}
func (p internalPluginInterface) GetSecret() umid.UMID {
	return p.secretId
}

func (p internalPluginInterface) GetWorld() umid.UMID {
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
