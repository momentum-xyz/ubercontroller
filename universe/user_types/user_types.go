package user_types

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/user_type"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.UserTypes = (*UserTypes)(nil)

type UserTypes struct {
	ctx       context.Context
	log       *zap.SugaredLogger
	db        database.DB
	userTypes *generic.SyncMap[uuid.UUID, universe.UserType]
}

func NewUserTypes(db database.DB) *UserTypes {
	return &UserTypes{
		db:        db,
		userTypes: generic.NewSyncMap[uuid.UUID, universe.UserType](),
	}
}

func (ut *UserTypes) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	ut.ctx = ctx
	ut.log = log

	return nil
}

func (ut *UserTypes) CreateUserType(userTypeID uuid.UUID) (universe.UserType, error) {
	userType := user_type.NewUserType(userTypeID, ut.db)

	if err := userType.Initialize(ut.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize user type: %s", userTypeID)
	}
	if err := ut.AddUserType(userType, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add user type: %s", userTypeID)
	}

	return userType, nil
}

func (ut *UserTypes) GetUserType(userTypeID uuid.UUID) (universe.UserType, bool) {
	userType, ok := ut.userTypes.Load(userTypeID)
	return userType, ok
}

func (ut *UserTypes) GetUserTypes() map[uuid.UUID]universe.UserType {
	ut.userTypes.Mu.RLock()
	defer ut.userTypes.Mu.RUnlock()

	userTypes := make(map[uuid.UUID]universe.UserType, len(ut.userTypes.Data))

	for id, userType := range ut.userTypes.Data {
		userTypes[id] = userType
	}

	return userTypes
}

func (ut *UserTypes) AddUserType(userType universe.UserType, updateDB bool) error {
	ut.userTypes.Mu.Lock()
	defer ut.userTypes.Mu.Unlock()

	if updateDB {
		if err := ut.db.UserTypesUpsertUserType(ut.ctx, userType.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	ut.userTypes.Data[userType.GetID()] = userType

	return nil
}

func (ut *UserTypes) AddUserTypes(userTypes []universe.UserType, updateDB bool) error {
	ut.userTypes.Mu.Lock()
	defer ut.userTypes.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.UserType, len(userTypes))
		for i := range userTypes {
			entries[i] = userTypes[i].GetEntry()
		}
		if err := ut.db.UserTypesUpsertUserTypes(ut.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range userTypes {
		ut.userTypes.Data[userTypes[i].GetID()] = userTypes[i]
	}

	return nil
}

func (ut *UserTypes) RemoveUserType(userType universe.UserType, updateDB bool) error {
	ut.userTypes.Mu.Lock()
	defer ut.userTypes.Mu.Unlock()

	if _, ok := ut.userTypes.Data[userType.GetID()]; !ok {
		return errors.Errorf("user type not found: %s", userType.GetID())
	}

	if updateDB {
		if err := ut.db.UserTypesRemoveUserTypeByID(ut.ctx, userType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	delete(ut.userTypes.Data, userType.GetID())

	return nil
}

func (ut *UserTypes) RemoveUserTypes(userTypes []universe.UserType, updateDB bool) error {
	ut.userTypes.Mu.Lock()
	defer ut.userTypes.Mu.Unlock()

	for i := range userTypes {
		if _, ok := ut.userTypes.Data[userTypes[i].GetID()]; !ok {
			return errors.Errorf("user type not found: %s", userTypes[i].GetID())
		}
	}

	if updateDB {
		ids := make([]uuid.UUID, len(userTypes))
		for i := range userTypes {
			ids[i] = userTypes[i].GetID()
		}
		if err := ut.db.UserTypesRemoveUserTypesByIDs(ut.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range userTypes {
		delete(ut.userTypes.Data, userTypes[i].GetID())
	}

	return nil
}

func (ut *UserTypes) Load() error {
	ut.log.Info("Loading user types...")

	entries, err := ut.db.UserTypesGetUserTypes(ut.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get user types")
	}

	for i := range entries {
		userType, err := ut.CreateUserType(*entries[i].UserTypeID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new user type: %s", entries[i].UserTypeID)
		}
		if err := userType.LoadFromEntry(entries[i]); err != nil {
			return errors.WithMessagef(err, "failed to load user type from entry: %s", *entries[i].UserTypeID)
		}
		ut.userTypes.Store(*entries[i].UserTypeID, userType)
	}

	universe.GetNode().AddAPIRegister(ut)

	ut.log.Info("User types loaded")

	return nil
}

func (ut *UserTypes) Save() error {
	ut.log.Info("Saving spate types...")

	ut.userTypes.Mu.RLock()
	defer ut.userTypes.Mu.RUnlock()

	entries := make([]*entry.UserType, 0, len(ut.userTypes.Data))
	for _, UserType := range ut.userTypes.Data {
		entries = append(entries, UserType.GetEntry())
	}

	if err := ut.db.UserTypesUpsertUserTypes(ut.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert user types")
	}

	ut.log.Info("User types saved")

	return nil
}
