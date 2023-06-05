package user_asset_3d

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.UserAsset3d = (*UserAsset3d)(nil)

type UserAsset3d struct {
	ctx     context.Context
	log     *zap.SugaredLogger
	db      database.DB
	mu      sync.RWMutex
	entry   *entry.UserAsset3d
	asset3d *universe.Asset3d
}

func NewAsset3d(assetUserID universe.AssetUserIDPair, db database.DB, asset3d *universe.Asset3d, isPrivate bool) *UserAsset3d {
	return &UserAsset3d{
		db: db,
		entry: &entry.UserAsset3d{
			Asset3dID: assetUserID.AssetID,
			UserID:    assetUserID.UserID,
			Private:   isPrivate,
		},
		asset3d: asset3d,
	}
}

func (a *UserAsset3d) GetAssetUserIDPair() universe.AssetUserIDPair {
	return universe.AssetUserIDPair{
		AssetID: a.entry.Asset3dID,
		UserID:  a.entry.UserID,
	}
}

func (a *UserAsset3d) GetAssetID() umid.UMID {
	return a.entry.Asset3dID
}

func (a *UserAsset3d) GetUserID() umid.UMID {
	return a.entry.UserID
}

func (a *UserAsset3d) GetAsset3d() *universe.Asset3d {
	return a.asset3d
}

func (a *UserAsset3d) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *UserAsset3d) GetMeta() *entry.Asset3dMeta {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Meta
}

func (a *UserAsset3d) SetMeta(meta *entry.Asset3dMeta, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		assetUserId := universe.AssetUserIDPair{
			AssetID: a.entry.Asset3dID,
			UserID:  a.entry.UserID,
		}
		if err := a.db.GetAssets3dDB().UpdateUserAssetMeta(a.ctx, assetUserId, meta); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Meta = meta

	return nil
}

func (a *UserAsset3d) IsPrivate() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.entry.Private
}

func (a *UserAsset3d) SetIsPrivate(isPrivate bool, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		assetUserId := universe.AssetUserIDPair{
			AssetID: a.entry.Asset3dID,
			UserID:  a.entry.UserID,
		}
		if err := a.db.GetAssets3dDB().UpdateUserAssetIsPrivate(a.ctx, assetUserId, isPrivate); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Private = isPrivate

	return nil
}

func (a *UserAsset3d) GetEntry() *entry.UserAsset3d {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry
}

func (a *UserAsset3d) LoadFromEntry(entry *entry.UserAsset3d) error {
	if entry.Asset3dID != a.GetAssetUserIDPair().AssetID || entry.UserID != a.GetAssetUserIDPair().UserID {
		return errors.Errorf("asset 3d ids mismatch: (%s, %s) != (%s, %s)",
			entry.Asset3dID, entry.UserID,
			a.GetAssetUserIDPair().AssetID, a.GetAssetUserIDPair().UserID)
	}

	a.entry = entry

	return nil
}
