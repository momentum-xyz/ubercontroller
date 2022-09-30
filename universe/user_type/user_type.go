package user_type

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
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.UserType = (*UserType)(nil)

type UserType struct {
	id          uuid.UUID
	ctx         context.Context
	log         *zap.SugaredLogger
	db          database.DB
	mu          sync.RWMutex
	name        string
	description *string
	options     *entry.UserOptions
	entry       *entry.UserType
}

func NewUserType(id uuid.UUID, db database.DB) *UserType {
	return &UserType{
		id: id,
		db: db,
		options: &entry.UserOptions{
			IsGuest: utils.GetPtr(true),
		},
	}
}

func (u *UserType) GetID() uuid.UUID {
	return u.id
}

func (u *UserType) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	u.ctx = ctx
	u.log = log

	return nil
}

func (u *UserType) GetName() string {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.name
}

func (u *UserType) SetName(name string, updateDB bool) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if updateDB {
		if err := u.db.UserTypesUpdateUserTypeName(u.ctx, u.id, name); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	u.name = name
	u.clearCache()

	return nil
}

func (u *UserType) GetDescription() *string {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.description
}

func (u *UserType) SetDescription(description *string, updateDB bool) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if updateDB {
		if err := u.db.UserTypesUpdateUserTypeDescription(u.ctx, u.id, description); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	u.description = description
	u.clearCache()

	return nil
}

func (u *UserType) GetOptions() *entry.UserOptions {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.options
}

func (u *UserType) SetOptions(modifyFn modify.Fn[entry.UserOptions], updateDB bool) error {
	u.mu.Lock()
	options := modifyFn(u.options)

	if updateDB {
		if err := u.db.UserTypesUpdateUserTypeOptions(u.ctx, u.id, options); err != nil {
			u.mu.Unlock()
			return errors.WithMessage(err, "failed to update db")
		}
	}

	u.options = options
	u.clearCache()
	u.mu.Unlock()

	for _, world := range universe.GetNode().GetWorlds().GetWorlds() {
		for _, user := range world.GetUsers(false) {
			if user.GetUserType() == nil {
				continue
			}
			if user.GetUserType().GetID() != u.GetID() {
				continue
			}
			if err := user.Update(); err != nil {
				return errors.WithMessagef(err, "failed to update user: %v", user.GetID().String())
			}
		}
	}

	return nil
}

func (u *UserType) GetEntry() *entry.UserType {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.entry == nil {
		u.entry = &entry.UserType{
			UserTypeID:   utils.GetPtr(u.id),
			UserTypeName: &u.name,
			Description:  u.description,
			Options:      u.options,
		}
	}

	return u.entry
}

func (u *UserType) clearCache() {
	u.entry = nil
}

func (u *UserType) LoadFromEntry(entry *entry.UserType) error {

	u.id = *entry.UserTypeID
	if err := u.SetName(*entry.UserTypeName, false); err != nil {
		return errors.WithMessage(err, "failed to set name")
	}
	if err := u.SetDescription(entry.Description, false); err != nil {
		return errors.WithMessage(err, "failed to set description")
	}
	if err := u.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	return nil
}
