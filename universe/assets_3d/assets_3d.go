package assets_3d

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/asset_3d"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Assets3d = (*Assets3d)(nil)

type Assets3d struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	cfg    *config.Config
	db     database.DB
	assets *generic.SyncMap[umid.UMID, universe.Asset3d]
}

func NewAssets3d(db database.DB) *Assets3d {
	return &Assets3d{
		db:     db,
		assets: generic.NewSyncMap[umid.UMID, universe.Asset3d](0),
	}
}

func (a *Assets3d) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.Errorf("failed to get config from context: %T", ctx.Value(types.ConfigContextKey))
	}

	a.ctx = ctx
	a.log = log
	a.cfg = cfg

	return nil
}

func (a *Assets3d) CreateAsset3d(asset3dID umid.UMID) (universe.Asset3d, error) {
	asset3d := asset_3d.NewAsset3d(asset3dID, a.db)

	if err := asset3d.Initialize(a.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize asset 3d: %s", asset3dID)
	}
	if err := a.AddAsset3d(asset3d, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add asset 3d: %s", asset3dID)
	}

	return asset3d, nil
}

func (a *Assets3d) FilterAssets3d(predicateFn universe.Assets3dFilterPredicateFn) map[umid.UMID]universe.Asset3d {
	return a.assets.Filter(predicateFn)
}

func (a *Assets3d) GetAsset3d(asset3dID umid.UMID) (universe.Asset3d, bool) {
	asset, ok := a.assets.Load(asset3dID)
	return asset, ok
}

func (a *Assets3d) GetAssets3d() map[umid.UMID]universe.Asset3d {
	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	assets := make(map[umid.UMID]universe.Asset3d, len(a.assets.Data))

	for id, asset := range a.assets.Data {
		assets[id] = asset
	}

	return assets
}

func (a *Assets3d) AddAsset3d(asset3d universe.Asset3d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if updateDB {
		if err := a.db.GetAssets3dDB().UpsertAsset(a.ctx, asset3d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.assets.Data[asset3d.GetID()] = asset3d

	return nil
}

func (a *Assets3d) AddAssets3d(assets3d []universe.Asset3d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.Asset3d, len(assets3d))
		for i := range assets3d {
			entries[i] = assets3d[i].GetEntry()
		}
		if err := a.db.GetAssets3dDB().UpsertAssets(a.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range assets3d {
		a.assets.Data[assets3d[i].GetID()] = assets3d[i]
	}

	return nil
}

func (a *Assets3d) RemoveAsset3d(asset3d universe.Asset3d, updateDB bool) (bool, error) {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if _, ok := a.assets.Data[asset3d.GetID()]; !ok {
		return false, nil
	}

	if updateDB {
		if err := a.db.GetAssets3dDB().RemoveAssetByID(a.ctx, asset3d.GetID()); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(a.assets.Data, asset3d.GetID())

	return true, nil
}

func (a *Assets3d) RemoveAssets3d(assets3d []universe.Asset3d, updateDB bool) (bool, error) {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	for i := range assets3d {
		if _, ok := a.assets.Data[assets3d[i].GetID()]; !ok {
			return false, nil
		}
	}

	if updateDB {
		ids := make([]umid.UMID, 0, len(assets3d))
		for i := range assets3d {
			ids[i] = assets3d[i].GetID()
		}
		if err := a.db.GetAssets3dDB().RemoveAssetsByIDs(a.ctx, ids); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range assets3d {
		delete(a.assets.Data, assets3d[i].GetID())
	}

	return true, nil
}

func (a *Assets3d) RemoveAsset3dByID(asset3dID umid.UMID, updateDB bool) (bool, error) {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if _, ok := a.assets.Data[asset3dID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := a.db.GetAssets3dDB().RemoveAssetByID(a.ctx, asset3dID); err != nil {
			return false, errors.Errorf("failed to update db")
		}
	}

	delete(a.assets.Data, asset3dID)

	return true, nil
}

func (a *Assets3d) RemoveAssets3dByIDs(assets3dIDs []umid.UMID, updateDB bool) (bool, error) {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	for i := range assets3dIDs {
		if _, ok := a.assets.Data[assets3dIDs[i]]; !ok {
			return false, nil
		}
	}

	if updateDB {
		if err := a.db.GetAssets3dDB().RemoveAssetsByIDs(a.ctx, assets3dIDs); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range assets3dIDs {
		delete(a.assets.Data, assets3dIDs[i])
	}

	return true, nil
}

func (a *Assets3d) Load() error {
	a.log.Info("Loading assets 3d...")

	entries, err := a.db.GetAssets3dDB().GetAssets(a.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get assets 3d")
	}

	for _, assetEntry := range entries {
		asset3d, err := a.CreateAsset3d(assetEntry.Asset3dID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new asset 3d: %s", assetEntry.Asset3dID)
		}
		if err := asset3d.LoadFromEntry(assetEntry); err != nil {
			return errors.WithMessagef(err, "failed to load asset 3d from entry: %s", assetEntry.Asset3dID)
		}
	}

	universe.GetNode().AddAPIRegister(a)

	a.log.Infof("Assets 3d loaded: %d", a.assets.Len())

	return nil
}

func (a *Assets3d) Save() error {
	a.log.Info("Saving assets 3d...")

	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	entries := make([]*entry.Asset3d, 0, len(a.assets.Data))
	for _, asset3d := range a.assets.Data {
		entries = append(entries, asset3d.GetEntry())
	}

	if err := a.db.GetAssets3dDB().UpsertAssets(a.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert assets 3d")
	}

	a.log.Infof("Assets 3d saved: %d", len(a.assets.Data))

	return nil
}
