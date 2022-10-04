package mplugin

import (
	"github.com/google/uuid"
	"github.com/momentum-xyz/controller/utils"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/pkg/errors"
	"reflect"
)

//type SubscriberInterface[A any] interface {
//	SubscribeHook(mplugin PluginInterface, name string, hook func(A) error) (HookID, error)
//}
//
//type Subscriber[A any] struct {
//
//}
//
//func (s Subscriber[A]) SubscribeHook(mplugin PluginInterface, name string, hook func(A) error) (HookID, error) {
//	subscribeHook()
//}

type PluginController struct {
	secretList      map[PluginID]uuid.UUID
	pluginInstances map[PluginID]internalPluginInterface
	hooksMap        *generic.SyncMap[string, any]
	parent          uuid.UUID
	loadedPlugins   *utils.SyncMap[uuid.UUID, PluginInstance]
}

func NewPluginController(parent uuid.UUID) *PluginController {
	pc := PluginController{parent: parent,
		pluginInstances: make(map[PluginID]internalPluginInterface),
		hooksMap:        generic.NewSyncMap[string, any](),
		secretList:      make(map[PluginID]uuid.UUID),
		loadedPlugins:   utils.NewSyncMap[uuid.UUID, PluginInstance](),
	}
	return &pc
}

func (p *PluginController) AddPlugin(pluginID uuid.UUID, f NewInstanceFunction) (PluginInstance, error) {
	pi := p.NewPluginInterface(PluginID(pluginID))
	instance, err := f(pi)
	if err != nil {
		return nil, errors.WithMessage(err, "To initialize plugin "+pluginID.String())
	}
	p.loadedPlugins.Store(pluginID, instance)
	return instance, nil
}

func registerHook[A any](p *PluginController, name string, F A) (*HookEntries[A], error) {
	//ft := reflect.TypeOf(F)
	p.hooksMap.Mu.Lock()
	defer p.hooksMap.Mu.Unlock()

	if _, ok := p.hooksMap.Data[name]; ok {
		return nil, errors.Errorf("hook %q already registered", name)
	}

	he := NewHookEntries[A]()
	p.hooksMap.Data[name] = he

	return he, nil
}

func (p *PluginController) NewPluginInterface(id PluginID) PluginInterface {
	secret := uuid.New()
	p.secretList[id] = secret
	pi := internalPluginInterface{id: id, secretId: secret, worldId: p.parent, pc: p}
	return &pi
}

func (p *PluginController) ValidPlugin(id PluginID, secret uuid.UUID) bool {
	if s, ok := p.secretList[id]; ok && s == secret {
		return true
	}
	return false
}

//func (p *PluginController) checkIfFunction(F any, argType reflect.Type) (reflect.Type, error) {
//	ft := reflect.TypeOf(F)
//	if ft.Kind() != reflect.Func {
//		return ft, errors.New("Hook is not a function")
//	}
//
//	if (ft.NumOut() != 1) || (ft.Out(0).String() != "error") {
//		return ft, errors.New("Hook must be a function which returns 'error'")
//	}
//	if (ft.NumIn() != 1) || (ft.In(0) != argType) {
//		return ft, errors.New("Hook must be a function with argument " + argType.String())
//	}
//
//	fmt.Println(ft)
//	return ft, nil
//}

func subscribeHook[A any](plugin PluginInterface, name string, hook func(A) error) (
	HookID, error,
) {
	if plugin.getPluginController().ValidPlugin(plugin.GetId(), plugin.GetSecret()) {
		hooks, ok := plugin.getPluginController().hooksMap.Load(name)
		if !ok {
			return HookID(uuid.Nil), errors.Errorf("hook %q is not registered", name)
		}

		list, ok := hooks.(*HookEntries[A])
		if !ok {
			return HookID(uuid.Nil), errors.Errorf("invalid hook type for %q", name)
		}

		//ft, err := p.checkIfFunction(hook, list.ArgType)

		//if err != nil {
		//	return HookID(uuid.Nil), errors.WithMessage(
		//		err,
		//		"invalid hook type for: '"+name+"', type:"+ft.String()+" expected: func("+list.ArgType.String()+") error",
		//	)
		//}
		//fmt.Println("hook")
		return list.Add(hook, plugin.GetId()), nil
	}
	return HookID(uuid.Nil), errors.New("Plugin is not authorized")
}

func (p *PluginController) UnsubscribeAllHooks(plugin PluginInterface) error {
	if p.ValidPlugin(plugin.GetId(), plugin.GetSecret()) {
		p.hooksMap.Mu.RLock()
		defer p.hooksMap.Mu.RUnlock()
		for _, hooks := range p.hooksMap.Data {
			reflect.ValueOf(hooks).MethodByName("RemoveByPluginID").Call([]reflect.Value{reflect.ValueOf(plugin.GetId())})
		}
	}
	return nil
}
func (p *PluginController) UnsubscribeHooksByName(plugin PluginInterface, name string) error {
	if p.ValidPlugin(plugin.GetId(), plugin.GetSecret()) {
		p.hooksMap.Mu.RLock()
		defer p.hooksMap.Mu.RUnlock()
		for hookName, hooks := range p.hooksMap.Data {
			if name == hookName {
				reflect.ValueOf(hooks).MethodByName("RemoveByPluginID").Call([]reflect.Value{reflect.ValueOf(plugin.GetId())})
			}
		}
	}
	return nil
}
func (p *PluginController) UnsubscribeHookByHookId(plugin PluginInterface, id HookID) error {
	if p.ValidPlugin(plugin.GetId(), plugin.GetSecret()) {
		p.hooksMap.Mu.RLock()
		defer p.hooksMap.Mu.RUnlock()
		for _, hooks := range p.hooksMap.Data {
			reflect.ValueOf(hooks).MethodByName("RemoveByHookID").Call([]reflect.Value{reflect.ValueOf(id)})
		}
	}
	return nil
}

// TODO: fix the bloody thing
func TriggerHook[A any](p *PluginController, name string, arg A) error {
	hooks, ok := p.hooksMap.Load(name)
	if !ok {
		return errors.Errorf("hook %q is not registered", name)
	}

	list, ok := hooks.(*HookEntries[A])
	if !ok {
		return errors.Errorf("invalid arg type for %q", name)
	}

	return list.Run(arg)
}
