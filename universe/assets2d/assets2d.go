package assets2d

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/types/generics"
	"github.com/momentum-xyz/controller/universe"
)

var _ universe.Assets2D = (*Assets2D)(nil)

type Assets2D struct {
	assets *generics.SyncMap[uuid.UUID, universe.Asset2D]
}

func NewAssets2D() *Assets2D {
	return &Assets2D{
		assets: generics.NewSyncMap[uuid.UUID, universe.Asset2D](),
	}
}

func (a *Assets2D) Initialize(ctx context.Context) error {
	return nil
}

func (a *Assets2D) GetAsset2D(asset2DID uuid.UUID) (universe.Asset2D, bool) {
	asset, ok := a.assets.Load(asset2DID)
	return asset, ok
}

func (a *Assets2D) Load() error {
	return errors.Errorf("implement me")
}