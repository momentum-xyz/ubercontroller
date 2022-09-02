package asset2d

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
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.Asset2d = (*Asset2d)(nil)

type Asset2d struct {
	ctx   context.Context
	log   *zap.SugaredLogger
	db    database.DB
	mu    sync.RWMutex
	entry *entry.Asset2d
}

func NewAsset2d(id uuid.UUID, db database.DB) *Asset2d {
	return &Asset2d{
		db: db,
		entry: &entry.Asset2d{
			Asset2dID: &id,
		},
	}
}

func (a *Asset2d) GetID() uuid.UUID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return *a.entry.Asset2dID
}

func (a *Asset2d) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

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

	if updateDB {
		if err := a.db.Assets2dUpdateAssetName(a.ctx, *a.entry.Asset2dID, name); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	*a.entry.Name = name

	return nil
}

func (a *Asset2d) GetOptions() *entry.Asset2dOptions {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset2d) SetOptions(modifyFn modify.Fn[entry.Asset2dOptions], updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	options := modifyFn(a.entry.Options)

	if updateDB {
		if err := a.db.Assets2dUpdateAssetOptions(a.ctx, *a.entry.Asset2dID, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Options = options

	return nil
}

func (a *Asset2d) GetEntry() *entry.Asset2d {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry
}

func (a *Asset2d) LoadFromEntry(entry *entry.Asset2d) error {
	a.entry = entry

	return nil
}
