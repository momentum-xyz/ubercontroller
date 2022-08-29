package assets2d

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
)

var _ universe.Assets2d = (*Assets2d)(nil)

type Assets2d struct {
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
	return nil
}

func (a *Assets2d) GetAsset2d(asset2dID uuid.UUID) (universe.Asset2d, bool) {
	asset, ok := a.assets.Load(asset2dID)
	return asset, ok
}

func (a *Assets2d) Load(ctx context.Context) error {
	return errors.Errorf("implement me")
}
