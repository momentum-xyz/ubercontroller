package generic

import (
	"github.com/sasha-s/go-deadlock"
)

type SyncMap[K comparable, V any] struct {
	Mu deadlock.RWMutex
	//Mu   sync.RWMutex
	Data map[K]V
}

func NewSyncMap[K comparable, V any](size int) *SyncMap[K, V] {
	return &SyncMap[K, V]{
		//Mu:   sync.RWMutex{},
		Mu:   deadlock.RWMutex{},
		Data: make(map[K]V, size),
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

func (m *SyncMap[K, V]) Map(mapFn func(k K, v V) V) map[K]V {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	data := make(map[K]V, len(m.Data))

	for k, v := range m.Data {
		data[k] = mapFn(k, v)
	}

	return data
}

func (m *SyncMap[K, V]) Filter(predicateFn func(k K, v V) bool) map[K]V {
	data := make(map[K]V)

	m.Mu.RLock()
	defer m.Mu.RUnlock()

	for k, v := range m.Data {
		if predicateFn(k, v) {
			data[k] = v
		}
	}

	return data
}

func (m *SyncMap[K, V]) Purge() {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	m.Data = make(map[K]V)
}

func (m *SyncMap[K, V]) Len() int {
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	return len(m.Data)
}
