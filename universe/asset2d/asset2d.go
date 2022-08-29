package asset2d

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/universe"
)

var _ universe.Asset2D = (*Asset2D)(nil)

type Asset2D struct {
	mu    sync.RWMutex
	id    uuid.UUID
	entry *universe.Asset2DEntry
}

func NewAsset2D(id uuid.UUID) *Asset2D {
	return &Asset2D{
		id: id,
	}
}

func (a *Asset2D) GetID() uuid.UUID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.id
}

func (a *Asset2D) Initialize(ctx context.Context) error {
	return nil
}

func (a *Asset2D) GetName() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return *a.entry.Name
}

func (a *Asset2D) SetName(name string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	*a.entry.Name = name

	return nil
}

func (a *Asset2D) LoadFromEntry(entry *universe.Asset2DEntry) error {
	return errors.Errorf("implement me")
}

func (a *Asset2D) GetOptions() *universe.Asset2DOptionsEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset2D) SetOptions(options *universe.Asset2DOptionsEntry, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entry.Options = options

	return nil
}
