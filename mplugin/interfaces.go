package mplugin

import "github.com/google/uuid"

type PluginInterface interface {
	GetId() PluginID
	GetWorld() uuid.UUID
	GetSecret() uuid.UUID
	getPluginController() *PluginController
	//SubscribeHook(name string, hook any) (HookID, error)
	UnsubscribeAllHooks() error
	UnsubscribeHooksByName(name string) error
	UnsubscribeHookByHookId(id HookID) error
}

type PluginInstance interface {
	Init() error
	Run() error
	Destroy() error
}

type NewInstanceFunction func(pluginInterface PluginInterface) (PluginInstance, error)
