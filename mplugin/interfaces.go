package mplugin

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type PluginInterface interface {
	GetId() PluginID
	GetWorld() umid.UMID
	GetSecret() umid.UMID
	// TODO: make sure it is either non private or removed
	getPluginController() *PluginController
	//SubscribeHook(name string, hook any) (HookID, error)
	UnsubscribeAllHooks() error
	UnsubscribeHooksByName(name string) error
	UnsubscribeHookByHookId(id HookID) error
}

type PluginInstance any

type NewInstanceFunction func(pluginInterface PluginInterface) (PluginInstance, error)

type PluginSubscriberInterface[A any] interface {
	PluginSubscribeHook(PluginInterface, string, func(A) error) (HookID, error)
}
