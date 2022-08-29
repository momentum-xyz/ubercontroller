package assets3d

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
)

var _ universe.Assets3d = (*Assets3d)(nil)

type Assets3d struct {
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
	return nil
}

func (a *Assets3d) GetAsset3d(asset3dID uuid.UUID) (universe.Asset3d, bool) {
	asset, ok := a.assets.Load(asset3dID)
	return asset, ok
}

func (a *Assets3d) Load(ctx context.Context) error {
	return errors.Errorf("implement me")
}
