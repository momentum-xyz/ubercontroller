package activity

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.Activity = (*Activity)(nil)

type Activity struct {
	ctx   context.Context
	log   *zap.SugaredLogger
	db    database.DB
	mu    sync.RWMutex
	entry *entry.Activity
}

func NewActivity(id umid.UMID, db database.DB) *Activity {
	return &Activity{
		db: db,
		entry: &entry.Activity{
			ActivityID: id,
		},
	}
}

func (a *Activity) GetID() umid.UMID {
	return a.entry.ActivityID
}

func (a *Activity) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	a.ctx = ctx
	a.log = log

	return nil
}

func (a *Activity) GetData() *entry.ActivityData {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Data
}

func (a *Activity) GetType() *string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.Type
}

func (a *Activity) GetObjectID() *umid.UMID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.ObjectID
}

func (a *Activity) GetUserID() *umid.UMID {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry.UserID
}

func (a *Activity) GetEntry() *entry.Activity {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.entry
}

func (a *Activity) LoadFromEntry(entry *entry.Activity) error {
	if entry.ActivityID != a.GetID() {
		return errors.Errorf("activity ids mismatch: %s != %s", entry.ActivityID, a.GetID())
	}

	a.entry = entry

	return nil
}
