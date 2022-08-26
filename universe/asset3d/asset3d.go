package asset3d

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/universe"
)

type Asset3D struct {
	id    uuid.UUID
	mu    sync.RWMutex
	entry *universe.SpaceAsset3DEntry
}

func NewAsset3D(id uuid.UUID) *Asset3D {
	return &Asset3D{
		id: id,
	}
}

func (a *Asset3D) GetID() uuid.UUID {
	return a.id
}

func (a *Asset3D) GetName() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Name
}

func (a *Asset3D) SetName(name string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entry.Name = name

	return nil
}

func (a *Asset3D) GetOptions() *universe.SpaceAsset3DOptionsEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset3D) SetOptions(options *universe.SpaceAsset3DOptionsEntry, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entry.Options = options

	return nil
}

func (a *Asset3D) Load() error {
	return errors.Errorf("implement me")
}
