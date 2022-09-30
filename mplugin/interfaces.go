package mplugin

import "github.com/google/uuid"

type PluginInstance interface {
	GetId() PluginID
	GetWorld() uuid.UUID
	GetSecret() uuid.UUID
	getPluginController() *PluginController
	//SubscribeHook(name string, hook any) (HookID, error)
	UnsubscribeAllHooks() error
	UnsubscribeHooksByName(name string) error
	UnsubscribeHookByHookId(id HookID) error
}
