package user_type

import (
	"context"
	"sync"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

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
	id          umid.UMID
	ctx         context.Context
	log         *zap.SugaredLogger
	db          database.DB
	mu          sync.RWMutex
	name        string
	description string
	options     *entry.UserOptions
}

func NewUserType(id umid.UMID, db database.DB) *UserType {
	return &UserType{
		id: id,
		db: db,
		options: &entry.UserOptions{
			IsGuest: utils.GetPTR(true),
		},
	}
}

func (u *UserType) GetID() umid.UMID {
	return u.id
}

func (u *UserType) Initialize(ctx types.LoggerContext) error {
	u.ctx = ctx
	u.log = ctx.Logger()

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
		if err := u.db.GetUserTypesDB().UpdateUserTypeName(u.ctx, u.GetID(), name); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	u.name = name

	return nil
}

func (u *UserType) GetDescription() string {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.description
}

func (u *UserType) SetDescription(description string, updateDB bool) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if updateDB {
		if err := u.db.GetUserTypesDB().UpdateUserTypeDescription(u.ctx, u.GetID(), description); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	u.description = description

	return nil
}

func (u *UserType) GetOptions() *entry.UserOptions {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.options
}

func (u *UserType) SetOptions(modifyFn modify.Fn[entry.UserOptions], updateDB bool) (*entry.UserOptions, error) {
	u.mu.Lock()
	options, err := modifyFn(u.options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := u.db.GetUserTypesDB().UpdateUserTypeOptions(u.ctx, u.GetID(), options); err != nil {
			u.mu.Unlock()
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	u.options = options
	u.mu.Unlock()

	for _, world := range universe.GetNode().GetWorlds().GetWorlds() {
		for _, user := range world.GetUsers(true) {
			userType := user.GetUserType()
			if userType == nil {
				continue
			}
			if userType.GetID() != u.GetID() {
				continue
			}
			if err := user.Update(); err != nil {
				return nil, errors.WithMessagef(err, "failed to update user: %s", user.GetID())
			}
		}
	}

	return options, nil
}

func (u *UserType) GetEntry() *entry.UserType {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return &entry.UserType{
		UserTypeID:   u.id,
		UserTypeName: u.name,
		Description:  u.description,
		Options:      u.options,
	}
}

func (u *UserType) LoadFromEntry(entry *entry.UserType) error {
	if entry.UserTypeID != u.GetID() {
		return errors.Errorf("user type ids mismatch: %s != %s", entry.UserTypeID, u.GetID())
	}

	if err := u.SetName(entry.UserTypeName, false); err != nil {
		return errors.WithMessage(err, "failed to set name")
	}
	if err := u.SetDescription(entry.Description, false); err != nil {
		return errors.WithMessage(err, "failed to set description")
	}
	if _, err := u.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	return nil
}
