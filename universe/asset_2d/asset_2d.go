package asset_2d

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
			Asset2dID: id,
		},
	}
}

func (a *Asset2d) GetID() uuid.UUID {
	return a.entry.Asset2dID
}

func (a *Asset2d) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *Asset2d) GetMeta() *entry.Asset2dMeta {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Meta
}

func (a *Asset2d) SetMeta(meta *entry.Asset2dMeta, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.GetAssets2dDB().UpdateAssetMeta(a.ctx, a.entry.Asset2dID, meta); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Meta = meta

	return nil
}

func (a *Asset2d) GetOptions() *entry.Asset2dOptions {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset2d) SetOptions(modifyFn modify.Fn[entry.Asset2dOptions], updateDB bool) (*entry.Asset2dOptions, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	options, err := modifyFn(a.entry.Options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := a.db.GetAssets2dDB().UpdateAssetOptions(a.ctx, a.entry.Asset2dID, options); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Options = options

	return options, nil
}

func (a *Asset2d) GetEntry() *entry.Asset2d {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry
}

func (a *Asset2d) LoadFromEntry(entry *entry.Asset2d) error {
	if entry.Asset2dID != a.GetID() {
		return errors.Errorf("asset 2d ids mismatch: %s != %s", entry.Asset2dID, a.GetID())
	}

	a.entry = entry

	return nil
}
