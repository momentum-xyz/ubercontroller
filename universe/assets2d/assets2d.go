package assets2d

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

var _ universe.Assets2d = (*Assets2d)(nil)

type Assets2d struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	db     database.DB
	assets *generics.SyncMap[uuid.UUID, universe.Asset2d]
}

func NewAssets2D(db database.DB) *Assets2d {
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

func (a *Assets2d) GetAssets2d(asset2dIDs []uuid.UUID) (*generics.SyncMap[uuid.UUID, universe.Asset2d], error) {
	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	assets := generics.NewSyncMap[uuid.UUID, universe.Asset2d]()

	// maybe we will need lock here in future
	for i := range asset2dIDs {
		asset2d, ok := a.assets.Data[asset2dIDs[i]]
		if !ok {
			return nil, errors.Errorf("asset 2d not found: %s", asset2dIDs[i])
		}
		assets.Data[asset2dIDs[i]] = asset2d
	}

	return assets, nil
}

func (a *Assets2d) AddAsset2d(asset2d universe.Asset2d, updateDB bool) error {
	a.assets.Mu.Lock()
	defer a.assets.Mu.Unlock()

	if _, ok := a.assets.Data[asset2d.GetID()]; ok {
		return errors.Errorf("asset 2d already exists")
	}

	if err := asset2d.Update(updateDB); err != nil {
		return errors.WithMessage(err, "failed to update asset 2d")
	}

	a.assets.Data[asset2d.GetID()] = asset2d

	return nil
}

func (a *Assets2d) AddAssets2d(assets2d []universe.Asset2d, updateDB bool) error {
	var errs *multierror.Error
	for i := range assets2d {
		if err := a.AddAsset2d(assets2d[i], updateDB); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to add asset 2d: %s", assets2d[i].GetID()))
		}
	}
	return errs.ErrorOrNil()
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
		if err := a.db.Assets2dRemoveAssetByIDs(a.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range assets2d {
		delete(a.assets.Data, assets2d[i].GetID())
	}

	return nil
}

func (a *Assets2d) Load(ctx context.Context) error {
	return errors.Errorf("implement me")
}

func (a *Assets2d) Update(updateDB bool) error {
	a.assets.Mu.RLock()
	defer a.assets.Mu.RUnlock()

	for _, asset2d := range a.assets.Data {
		if err := asset2d.Update(updateDB); err != nil {
			return errors.WithMessagef(err, "failed to update asset 2d: %s", asset2d.GetID())
		}
	}

	return nil
}
