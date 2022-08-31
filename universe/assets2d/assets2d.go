package assets2d

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/asset2d"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Assets2d = (*Assets2d)(nil)

type Assets2d struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	db     database.DB
	assets *generics.SyncMap[uuid.UUID, universe.Asset2d]
}

func NewAssets2d(db database.DB) *Assets2d {
	return &Assets2d{
		db:     db,
		assets: generics.NewSyncMap[uuid.UUID, universe.Asset2d](),
	}
}

func (a *Assets2d) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *Assets2d) GetAsset2d(asset2dID uuid.UUID) (universe.Asset2d, bool) {
	asset, ok := a.assets.Load(asset2dID)
	return asset, ok
}

func (a *Assets2d) GetAssets2d() map[uuid.UUID]universe.Asset2d {
	assets := make(map[uuid.UUID]universe.Asset2d)

	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	for id, asset := range a.assets.Data {
		assets[id] = asset
	}

	return assets
}

func (a *Assets2d) AddAsset2d(asset2d universe.Asset2d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if _, ok := a.assets.Data[asset2d.GetID()]; ok {
		return errors.Errorf("asset 2d already exists")
	}

	if updateDB {
		if err := a.db.Assets2dUpsetAsset(a.ctx, asset2d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.assets.Data[asset2d.GetID()] = asset2d

	return nil
}

func (a *Assets2d) AddAssets2d(assets2d []universe.Asset2d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	for i := range assets2d {
		if _, ok := a.assets.Data[assets2d[i].GetID()]; ok {
			return errors.Errorf("asset 2d already exists: %s", assets2d[i].GetID())
		}
	}

	if updateDB {
		entries := make([]*entry.Asset2d, len(assets2d))
		for i := range assets2d {
			entries[i] = assets2d[i].GetEntry()
		}
		if err := a.db.Assets2dUpsetAssets(a.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range assets2d {
		a.assets.Data[assets2d[i].GetID()] = assets2d[i]
	}

	return nil
}

func (a *Assets2d) RemoveAsset2d(asset2d universe.Asset2d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if _, ok := a.assets.Data[asset2d.GetID()]; !ok {
		return errors.Errorf("asset 2d not found")
	}

	if updateDB {
		if err := a.db.Assets2dRemoveAssetByID(a.ctx, asset2d.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	delete(a.assets.Data, asset2d.GetID())

	return nil
}

func (a *Assets2d) RemoveAssets2d(assets2d []universe.Asset2d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	for i := range assets2d {
		if _, ok := a.assets.Data[assets2d[i].GetID()]; !ok {
			return errors.Errorf("asset 2d not found: %s", assets2d[i].GetID())
		}
	}

	if updateDB {
		ids := make([]uuid.UUID, len(assets2d))
		for i := range assets2d {
			ids[i] = assets2d[i].GetID()
		}
		if err := a.db.Assets2dRemoveAssetsByIDs(a.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range assets2d {
		delete(a.assets.Data, assets2d[i].GetID())
	}

	return nil
}

func (a *Assets2d) RegisterAPI(r *gin.Engine) {

}

func (a *Assets2d) Load(ctx context.Context) error {
	assets, err := a.db.Assets2dGetAssets(ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get assets 2d")
	}

	for i := range assets {
		asset := asset2d.NewAsset2d(*assets[i].Asset2dID, a.db)

		if err := asset.Initialize(ctx); err != nil {
			return errors.WithMessagef(err, "failed to initialize asset 2d: %s", *assets[i].Asset2dID)
		}
		if err := asset.LoadFromEntry(assets[i]); err != nil {
			return errors.WithMessagef(err, "failed to load asset 2d from entry: %s", *assets[i].Asset2dID)
		}
		a.assets.Store(*assets[i].Asset2dID, asset)
	}

	return nil
}

func (a *Assets2d) Save(ctx context.Context) error {
	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	entries := make([]*entry.Asset2d, len(a.assets.Data))
	for _, asset := range a.assets.Data {
		entries = append(entries, asset.GetEntry())
	}

	if err := a.db.Assets2dUpsetAssets(ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert assets 2d")
	}

	return nil
}
