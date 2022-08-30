package asset3d

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Asset3d = (*Asset3d)(nil)

type Asset3d struct {
	ctx   context.Context
	log   *zap.SugaredLogger
	db    database.DB
	mu    sync.RWMutex
	id    uuid.UUID
	entry *entry.Asset3d
}

func NewAsset3D(id uuid.UUID, db database.DB) *Asset3d {
	return &Asset3d{
		id: id,
		db: db,
	}
}

func (a *Asset3d) GetID() uuid.UUID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.id
}

func (a *Asset3d) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *Asset3d) GetName() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return *a.entry.Name
}

func (a *Asset3d) SetName(name string, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.Assets3dUpdateAssetName(a.ctx, a.id, name); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	*a.entry.Name = name

	return nil
}

func (a *Asset3d) GetOptions() *entry.Asset3dOptions {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Options
}

func (a *Asset3d) SetOptions(options *entry.Asset3dOptions, updateDB bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if updateDB {
		if err := a.db.Assets3dUpdateAssetOptions(a.ctx, a.id, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	a.entry.Options = options

	return nil
}

func (a *Asset3d) LoadFromEntry(entry *entry.Asset3d) error {
	return errors.Errorf("implement me")
}

func (a *Asset3d) Update(updateDB bool) error {
	return errors.Errorf("implement me")
}
