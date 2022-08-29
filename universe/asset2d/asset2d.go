package asset2d

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/universe"
)

var _ universe.Asset2d = (*Asset2d)(nil)

type Asset2d struct {
	db    database.DB
	mu    sync.RWMutex
	id    uuid.UUID
	entry *universe.Asset2dEntry
}

func NewAsset2D(id uuid.UUID, db database.DB) *Asset2d {
	return &Asset2d{
		id: id,
		db: db,
	}
}

func (a *Asset2d) GetID() uuid.UUID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.id
}

func (a *Asset2d) Initialize(ctx context.Context) error {
	return nil
}

func (a *Asset2d) GetName() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return *a.entry.Name
}

func (a *Asset2d) SetName(name string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	*a.entry.Name = name

	return nil
}

func (a *Asset2d) LoadFromEntry(ctx context.Context, entry *universe.Asset2dEntry) error {
	return errors.Errorf("implement me")
}

func (a *Asset2d) GetOptions() *universe.Asset2dOptionsEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset2d) SetOptions(options *universe.Asset2dOptionsEntry, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.entry.Options = options

	return nil
}
