package assets3d

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/asset3d"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Assets3d = (*Assets3d)(nil)

type Assets3d struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	db     database.DB
	assets *generic.SyncMap[uuid.UUID, universe.Asset3d]
}

func NewAssets3d(db database.DB) *Assets3d {
	return &Assets3d{
		db:     db,
		assets: generic.NewSyncMap[uuid.UUID, universe.Asset3d](),
	}
}

func (a *Assets3d) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *Assets3d) GetAsset3d(asset3dID uuid.UUID) (universe.Asset3d, bool) {
	asset, ok := a.assets.Load(asset3dID)
	return asset, ok
}

func (a *Assets3d) GetAssets3d() map[uuid.UUID]universe.Asset3d {
	assets := make(map[uuid.UUID]universe.Asset3d)

	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	for id, asset := range a.assets.Data {
		assets[id] = asset
	}

	return assets
}

func (a *Assets3d) AddAsset3d(asset3d universe.Asset3d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if _, ok := a.assets.Data[asset3d.GetID()]; ok {
		return errors.Errorf("asset 3d already exists")
	}

	if updateDB {
		if err := a.db.Assets3dUpsetAsset(a.ctx, asset3d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.assets.Data[asset3d.GetID()] = asset3d

	return nil
}

func (a *Assets3d) AddAssets3d(assets3d []universe.Asset3d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	for i := range assets3d {
		if _, ok := a.assets.Data[assets3d[i].GetID()]; ok {
			return errors.Errorf("asset 3d already exists: %s", assets3d[i].GetID())
		}
	}

	if updateDB {
		entries := make([]*entry.Asset3d, len(assets3d))
		for i := range assets3d {
			entries[i] = assets3d[i].GetEntry()
		}
		if err := a.db.Assets3dUpsetAssets(a.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range assets3d {
		a.assets.Data[assets3d[i].GetID()] = assets3d[i]
	}

	return nil
}

func (a *Assets3d) RemoveAsset3d(asset3d universe.Asset3d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if _, ok := a.assets.Data[asset3d.GetID()]; !ok {
		return errors.Errorf("asset 3d not found")
	}

	if updateDB {
		if err := a.db.Assets3dRemoveAssetByID(a.ctx, asset3d.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	delete(a.assets.Data, asset3d.GetID())

	return nil
}

func (a *Assets3d) RemoveAssets3d(assets3d []universe.Asset3d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	for i := range assets3d {
		if _, ok := a.assets.Data[assets3d[i].GetID()]; !ok {
			return errors.Errorf("asset 3d not found: %s", assets3d[i].GetID())
		}
	}

	if updateDB {
		ids := make([]uuid.UUID, len(assets3d))
		for i := range assets3d {
			ids[i] = assets3d[i].GetID()
		}
		if err := a.db.Assets3dRemoveAssetsByIDs(a.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range assets3d {
		delete(a.assets.Data, assets3d[i].GetID())
	}

	return nil
}

func (a *Assets3d) Load() error {
	entries, err := a.db.Assets3dGetAssets(a.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get assets 3d")
	}

	for i := range entries {
		asset := asset3d.NewAsset3d(*entries[i].Asset3dID, a.db)

		if err := asset.Initialize(a.ctx); err != nil {
			return errors.WithMessagef(err, "failed to initialize asset 3d: %s", *entries[i].Asset3dID)
		}
		if err := asset.LoadFromEntry(entries[i]); err != nil {
			return errors.WithMessagef(err, "failed to load asset 3d from entry: %s", *entries[i].Asset3dID)
		}

		a.assets.Store(*entries[i].Asset3dID, asset)
	}

	universe.GetNode().AddAPIRegister(a)

	return nil
}

func (a *Assets3d) Save() error {
	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	entries := make([]*entry.Asset3d, len(a.assets.Data))
	for _, asset := range a.assets.Data {
		entries = append(entries, asset.GetEntry())
	}

	if err := a.db.Assets3dUpsetAssets(a.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert assets 3d")
	}

	return nil
}
