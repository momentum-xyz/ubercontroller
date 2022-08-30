package assets3d

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Assets3d = (*Assets3d)(nil)

type Assets3d struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	db     database.DB
	assets *generics.SyncMap[uuid.UUID, universe.Asset3d]
}

func NewAssets3D(db database.DB) *Assets3d {
	return &Assets3d{
		db:     db,
		assets: generics.NewSyncMap[uuid.UUID, universe.Asset3d](),
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

func (a *Assets3d) GetAssets3d(asset3dIDs []uuid.UUID) (*generics.SyncMap[uuid.UUID, universe.Asset3d], error) {
	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	assets := generics.NewSyncMap[uuid.UUID, universe.Asset3d]()

	// maybe we will need lock here in future
	for i := range asset3dIDs {
		asset3d, ok := a.assets.Data[asset3dIDs[i]]
		if !ok {
			return nil, errors.Errorf("asset 2d not found: %s", asset3dIDs[i])
		}
		assets.Data[asset3dIDs[i]] = asset3d
	}

	return assets, nil
}

func (a *Assets3d) AddAsset3d(asset3d universe.Asset3d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if _, ok := a.assets.Data[asset3d.GetID()]; ok {
		return errors.Errorf("asset 3d already exists")
	}

	if err := asset3d.Update(updateDB); err != nil {
		return errors.WithMessage(err, "failed to update asset 3d")
	}

	a.assets.Data[asset3d.GetID()] = asset3d

	return nil
}

func (a *Assets3d) AddAssets3d(assets3d []universe.Asset3d, updateDB bool) error {
	var errs *multierror.Error
	for i := range assets3d {
		if err := a.AddAsset3d(assets3d[i], updateDB); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to add asset 3d: %s", assets3d[i].GetID()))
		}
	}
	return errs.ErrorOrNil()
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
		if err := a.db.Assets3dRemoveAssetByIDs(a.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range assets3d {
		delete(a.assets.Data, assets3d[i].GetID())
	}

	return nil
}

func (a *Assets3d) Load(ctx context.Context) error {
	return errors.Errorf("implement me")
}

func (a *Assets3d) Update(updateDB bool) error {
	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	for _, asset3d := range a.assets.Data {
		if err := asset3d.Update(updateDB); err != nil {
			return errors.WithMessagef(err, "failed to update asset 3d: %s", asset3d.GetID())
		}
	}

	return nil
}
