package asset3d

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/universe"
)

var _ universe.Asset3d = (*Asset3d)(nil)

type Asset3d struct {
	db    database.DB
	mu    sync.RWMutex
	id    uuid.UUID
	entry *universe.Asset3dEntry
}

func NewAsset3D(id uuid.UUID, db database.DB) *Asset3d {
	return &Asset3d{
		id: id,
		db: db,
	}
}

func (a *Asset3d) GetID() uuid.UUID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.id
}

func (a *Asset3d) Initialize(ctx context.Context) error {
	return nil
}

func (a *Asset3d) GetName() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return *a.entry.Name
}

func (a *Asset3d) SetName(name string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	*a.entry.Name = name

	return nil
}

func (a *Asset3d) GetOptions() *universe.Asset3dOptionsEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset3d) SetOptions(options *universe.Asset3dOptionsEntry, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entry.Options = options

	return nil
}

func (a *Asset3d) LoadFromEntry(ctx context.Context, entry *universe.Asset3dEntry) error {
	return errors.Errorf("implement me")
}
