package asset3d

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Asset3d = (*Asset3d)(nil)

type Asset3d struct {
	ctx   context.Context
	log   *zap.SugaredLogger
	db    database.DB
	mu    sync.RWMutex
	entry *entry.Asset3d
}

func NewAsset3d(id uuid.UUID, db database.DB) *Asset3d {
	return &Asset3d{
		db: db,
		entry: &entry.Asset3d{
			Asset3dID: &id,
		},
	}
}

func (a *Asset3d) GetID() uuid.UUID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return *a.entry.Asset3dID
}

func (a *Asset3d) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

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

	if updateDB {
		if err := a.db.Assets3dUpdateAssetName(a.ctx, *a.entry.Asset3dID, name); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	*a.entry.Name = name

	return nil
}

func (a *Asset3d) GetOptions() *entry.Asset3dOptions {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset3d) SetOptions(setFn utils.SetFn[entry.Asset3dOptions], updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	options := setFn(a.entry.Options)

	if updateDB {
		if err := a.db.Assets3dUpdateAssetOptions(a.ctx, *a.entry.Asset3dID, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Options = options

	return nil
}

func (a *Asset3d) GetEntry() *entry.Asset3d {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry
}

func (a *Asset3d) LoadFromEntry(entry *entry.Asset3d) error {
	a.entry = entry

	return nil
}
