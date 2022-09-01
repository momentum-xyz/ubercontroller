package generic

import (
	"sync"
)

type SyncMap[K comparable, V any] struct {
	Mu   sync.RWMutex
	Data map[K]V
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		Mu:   sync.RWMutex{},
		Data: make(map[K]V),
	}
}

func (m *SyncMap[K, V]) Store(k K, v V) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	m.Data[k] = v
}

func (m *SyncMap[K, V]) Load(k K) (V, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	v, ok := m.Data[k]
	return v, ok
}

func (m *SyncMap[K, V]) Remove(k K) bool {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	if _, ok := m.Data[k]; !ok {
		return false
	}

	delete(m.Data, k)

	return true
}

func (m *SyncMap[K, V]) Purge() {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	m.Data = make(map[K]V)
}
