package asset2d

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/universe"
)

type Asset2D struct {
	id    uuid.UUID
	mu    sync.RWMutex
	entry *universe.SpaceAsset2DEntry
}

func NewAsset2D(id uuid.UUID) *Asset2D {
	return &Asset2D{
		id: id,
	}
}

func (a *Asset2D) GetID() uuid.UUID {
	return a.id
}

func (a *Asset2D) GetName() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Name
}

func (a *Asset2D) SetName(name string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entry.Name = name

	return nil
}

func (a *Asset2D) GetOptions() *universe.SpaceAsset2DOptionsEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset2D) SetOptions(options *universe.SpaceAsset2DOptionsEntry, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entry.Options = options

	return nil
}

func (a *Asset2D) Load() error {
	return errors.Errorf("implement me")
}
