package assets_3d

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/asset_3d"
	"github.com/momentum-xyz/ubercontroller/universe/user_asset_3d"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.Assets3d = (*Assets3d)(nil)

type Assets3d struct {
	ctx        context.Context
	log        *zap.SugaredLogger
	cfg        *config.Config
	db         database.DB
	assets     *generic.SyncMap[umid.UMID, universe.Asset3d]
	userAssets *generic.SyncMap[universe.AssetUserIDPair, universe.UserAsset3d]
}

func NewAssets3d(db database.DB) *Assets3d {
	return &Assets3d{
		db:         db,
		assets:     generic.NewSyncMap[umid.UMID, universe.Asset3d](0),
		userAssets: generic.NewSyncMap[universe.AssetUserIDPair, universe.UserAsset3d](0),
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

func (a *Assets3d) CreateAsset3dIfMissing(asset3dID umid.UMID) (universe.Asset3d, error, bool) {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	asset3d := a.assets.Data[asset3dID]
	if asset3d != nil {
		return asset3d, nil, false
	}

	asset3d = asset_3d.NewAsset3d(asset3dID, a.db)

	if err := asset3d.Initialize(a.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize asset 3d: %s", asset3dID), false
	}

	a.assets.Data[asset3d.GetID()] = asset3d

	return asset3d, nil, true
}

func (a *Assets3d) CreateUserAsset3d(assetID umid.UMID, userID umid.UMID, isPrivate bool) (universe.UserAsset3d, error) {
	assetUserID := universe.AssetUserIDPair{
		AssetID: assetID,
		UserID:  userID,
	}

	asset3d, ok := a.GetAsset3d(assetUserID.AssetID)
	if !ok {
		return nil, errors.Errorf("failed to get asset 3d: %s", assetUserID.AssetID)
	}

	userAsset3d := user_asset_3d.NewAsset3d(assetUserID, a.db, &asset3d, isPrivate)

	if err := userAsset3d.Initialize(a.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize asset 3d: %s", assetUserID)
	}
	if err := a.AddUserAsset3d(userAsset3d, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add asset 3d: %s", assetUserID)
	}

	return userAsset3d, nil
}

func (a *Assets3d) FilterUserAssets3d(predicateFn universe.Assets3dFilterPredicateFn) map[universe.AssetUserIDPair]universe.UserAsset3d {
	return a.userAssets.Filter(predicateFn)
}

func (a *Assets3d) GetAsset3d(id umid.UMID) (universe.Asset3d, bool) {
	asset, ok := a.assets.Load(id)
	return asset, ok
}

func (a *Assets3d) GetUserAsset3d(assetID umid.UMID, userID umid.UMID) (universe.UserAsset3d, bool) {
	assetUserID := universe.AssetUserIDPair{
		AssetID: assetID,
		UserID:  userID,
	}
	userAsset3d, ok := a.userAssets.Load(assetUserID)
	return userAsset3d, ok
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

func (a *Assets3d) GetUserAssets3d() map[universe.AssetUserIDPair]universe.UserAsset3d {
	a.userAssets.Mu.RLock()
	defer a.userAssets.Mu.RUnlock()

	assetsUserInstances := make(map[universe.AssetUserIDPair]universe.UserAsset3d, len(a.userAssets.Data))

	for id, item := range a.userAssets.Data {
		assetsUserInstances[id] = item
	}

	return assetsUserInstances
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

func (a *Assets3d) AddUserAsset3d(userAsset3d universe.UserAsset3d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if updateDB {
		if err := a.db.GetAssets3dDB().UpsertUserAsset(a.ctx, userAsset3d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.userAssets.Data[userAsset3d.GetAssetUserIDPair()] = userAsset3d

	return nil
}

func (a *Assets3d) RemoveUserAsset3dByID(assetUserID universe.AssetUserIDPair, updateDB bool) (bool, error) {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if _, ok := a.userAssets.Data[assetUserID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := a.db.GetAssets3dDB().RemoveUserAssetByID(a.ctx, assetUserID); err != nil {
			return false, errors.Errorf("failed to update db")
		}
	}

	delete(a.userAssets.Data, assetUserID)

	return true, nil
}

func (a *Assets3d) Load() error {
	a.log.Info("Loading assets 3d...")

	entriesAssets, err := a.db.GetAssets3dDB().GetAssets(a.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get assets 3d")
	}

	entriesUserAssets, err := a.db.GetAssets3dDB().GetUserAssets(a.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get assets 3d")
	}

	for _, assetEntry := range entriesAssets {
		asset3d, err := a.CreateAsset3d(assetEntry.Asset3dID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new asset 3d: %s", assetEntry.Asset3dID)
		}
		if err := asset3d.LoadFromEntry(assetEntry); err != nil {
			return errors.WithMessagef(err, "failed to load asset 3d from entry: %s", assetEntry.Asset3dID)
		}
	}

	for _, userAssetEntry := range entriesUserAssets {
		asset3d, err := a.CreateUserAsset3d(userAssetEntry.Asset3dID, userAssetEntry.UserID, userAssetEntry.Private)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new user asset 3d: %s", userAssetEntry.Asset3dID)
		}
		if err := asset3d.LoadFromEntry(userAssetEntry); err != nil {
			return errors.WithMessagef(err, "failed to load user asset 3d from entry: %s", userAssetEntry.Asset3dID)
		}
	}

	universe.GetNode().AddAPIRegister(a)

	a.log.Infof("Assets 3d loaded: %d", a.assets.Len())
	a.log.Infof("User Assets 3d loaded: %d", a.userAssets.Len())

	return nil
}

func (a *Assets3d) Save() error {
	a.log.Info("Saving assets 3d...")

	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	entriesAssets := make([]*entry.Asset3d, 0, len(a.assets.Data))
	for _, asset3d := range a.assets.Data {
		entriesAssets = append(entriesAssets, asset3d.GetEntry())
	}

	entriesUserAssets := make([]*entry.UserAsset3d, 0, len(a.assets.Data))
	for _, userAsset3d := range a.userAssets.Data {
		entriesUserAssets = append(entriesUserAssets, userAsset3d.GetEntry())
	}

	if err := a.db.GetAssets3dDB().UpsertAssets(a.ctx, entriesAssets, entriesUserAssets); err != nil {
		return errors.WithMessage(err, "failed to upsert assets 3d")
	}

	a.log.Infof("Assets 3d saved: %d", len(a.assets.Data))
	a.log.Infof("User Assets 3d saved: %d", len(a.userAssets.Data))

	return nil
}
