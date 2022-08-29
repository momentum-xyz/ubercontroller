package asset3d

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/universe"
)

var _ universe.Asset3D = (*Asset3D)(nil)

type Asset3D struct {
	mu    sync.RWMutex
	id    uuid.UUID
	entry *universe.Asset3DEntry
}

func NewAsset3D(id uuid.UUID) *Asset3D {
	return &Asset3D{
		id: id,
	}
}

func (a *Asset3D) GetID() uuid.UUID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.id
}

func (a *Asset3D) Initialize(ctx context.Context) error {
	return nil
}

func (a *Asset3D) GetName() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return *a.entry.Name
}

func (a *Asset3D) SetName(name string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	*a.entry.Name = name

	return nil
}

func (a *Asset3D) GetOptions() *universe.Asset3DOptionsEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset3D) SetOptions(options *universe.Asset3DOptionsEntry, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entry.Options = options

	return nil
}

func (a *Asset3D) LoadFromEntry(entry *universe.Asset3DEntry) error {
	return errors.Errorf("implement me")
}
