package assets3d

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/types/generics"
	"github.com/momentum-xyz/controller/universe"
)

var _ universe.Assets3D = (*Assets3D)(nil)

type Assets3D struct {
	assets *generics.SyncMap[uuid.UUID, universe.Asset3D]
}

func NewAssets3D() *Assets3D {
	return &Assets3D{
		assets: generics.NewSyncMap[uuid.UUID, universe.Asset3D](),
	}
}

func (a *Assets3D) Initialize(ctx context.Context) error {
	return nil
}

func (a *Assets3D) GetAsset3D(asset3DID uuid.UUID) (universe.Asset3D, bool) {
	asset, ok := a.assets.Load(asset3DID)
	return asset, ok
}

func (a *Assets3D) Load() error {
	return errors.Errorf("implement me")
}
