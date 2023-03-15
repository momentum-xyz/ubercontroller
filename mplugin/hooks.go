package mplugin

import (
	"github.com/hashicorp/go-multierror"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"sync"
)

type HookID mid.ID

type PluginID mid.ID

type HookEntry[A any] struct {
	hook     func(A) error
	pluginID PluginID
}

type HookEntries[A any] struct {
	Mu   sync.RWMutex
	Data map[HookID]HookEntry[A]
}

func NewHookEntries[A any]() *HookEntries[A] {
	return &HookEntries[A]{
		Mu:   sync.RWMutex{},
		Data: make(map[HookID]HookEntry[A]),
	}
}

func (s *HookEntries[A]) Add(hook func(A) error, pluginID PluginID) HookID {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	hookID := HookID(mid.New())
	s.Data[hookID] = HookEntry[A]{
		hook:     hook,
		pluginID: pluginID,
	}

	return hookID
}

func (s *HookEntries[A]) Run(arg A) error {
	var wg sync.WaitGroup
	var errs *multierror.Error
	var errMu sync.Mutex
	//fmt.Println("rr")
	s.Mu.RLock()
	for _, v := range s.Data {
		//fmt.Println("rr0")
		wg.Add(1)

		go func(hook func(A) error) {
			defer wg.Done()
			//fmt.Println("rr1")
			if err := hook(arg); err != nil {
				errMu.Lock()
				defer errMu.Unlock()

				errs = multierror.Append(errs, err)
			}
		}(v.hook)
	}
	s.Mu.RUnlock()

	wg.Wait()

	return errs.ErrorOrNil()
}

func (s *HookEntries[A]) Num() int {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return len(s.Data)
}

func (s *HookEntries[A]) RemoveByHookID(hookID HookID) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Data, hookID)

	return true
}

func (s *HookEntries[A]) RemoveByPluginID(pluginID PluginID) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for k, v := range s.Data {
		if v.pluginID == pluginID {
			delete(s.Data, k)
		}
	}
	return true
}
