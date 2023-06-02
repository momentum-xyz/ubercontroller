package asset_3d

import (
	"context"
	"sync"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.Asset3d = (*Asset3d)(nil)

type Asset3d struct {
	ctx   context.Context
	log   *zap.SugaredLogger
	db    database.DB
	mu    sync.RWMutex
	entry *entry.Asset3d
}

func NewAsset3d(assetID umid.UMID, db database.DB) *Asset3d {
	return &Asset3d{
		db: db,
		entry: &entry.Asset3d{
			Asset3dID: assetID,
		},
	}
}

func (a *Asset3d) GetID() umid.UMID {
	return a.entry.Asset3dID
}

func (a *Asset3d) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *Asset3d) GetMeta() *entry.Asset3dMeta {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Meta
}

func (a *Asset3d) SetMeta(meta *entry.Asset3dMeta, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.GetAssets3dDB().UpdateAssetMeta(a.ctx, a.entry.Asset3dID, meta); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Meta = meta

	return nil
}

func (a *Asset3d) GetOptions() *entry.Asset3dOptions {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset3d) SetOptions(modifyFn modify.Fn[entry.Asset3dOptions], updateDB bool) (*entry.Asset3dOptions, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	options, err := modifyFn(a.entry.Options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := a.db.GetAssets3dDB().UpdateAssetOptions(a.ctx, a.entry.Asset3dID, options); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Options = options

	return options, nil
}

func (a *Asset3d) GetEntry() *entry.Asset3d {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry
}

func (a *Asset3d) LoadFromEntry(entry *entry.Asset3d) error {
	if entry.Asset3dID != a.GetID() {
		return errors.Errorf("asset 3d ids mismatch: %s != %s", entry.Asset3dID, a.GetID())
	}

	a.entry = entry

	return nil
}
