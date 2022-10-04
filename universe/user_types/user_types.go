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
	UserTypes *generic.SyncMap[uuid.UUID, universe.UserType]
}

func NewUserTypes(db database.DB) *UserTypes {
	return &UserTypes{
		db:        db,
		UserTypes: generic.NewSyncMap[uuid.UUID, universe.UserType](),
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

func (ut *UserTypes) NewUserType(UserTypeID uuid.UUID) (universe.UserType, error) {
	UserType := user_type.NewUserType(UserTypeID, ut.db)

	if err := UserType.Initialize(ut.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize user type: %ut", UserTypeID)
	}
	if err := ut.AddUserType(UserType, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add user type: %ut", UserTypeID)
	}

	return UserType, nil
}

func (ut *UserTypes) GetUserType(UserTypeID uuid.UUID) (universe.UserType, bool) {
	UserType, ok := ut.UserTypes.Load(UserTypeID)
	return UserType, ok
}

func (ut *UserTypes) GetUserTypes() map[uuid.UUID]universe.UserType {
	ut.UserTypes.Mu.RLock()
	defer ut.UserTypes.Mu.RUnlock()

	UserTypes := make(map[uuid.UUID]universe.UserType, len(ut.UserTypes.Data))

	for id, UserType := range ut.UserTypes.Data {
		UserTypes[id] = UserType
	}

	return UserTypes
}

func (ut *UserTypes) AddUserType(UserType universe.UserType, updateDB bool) error {
	ut.UserTypes.Mu.Lock()
	defer ut.UserTypes.Mu.Unlock()

	if updateDB {
		if err := ut.db.UserTypesUpsertUserType(ut.ctx, UserType.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	ut.UserTypes.Data[UserType.GetID()] = UserType

	return nil
}

func (ut *UserTypes) AddUserTypes(UserTypes []universe.UserType, updateDB bool) error {
	ut.UserTypes.Mu.Lock()
	defer ut.UserTypes.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.UserType, len(UserTypes))
		for i := range UserTypes {
			entries[i] = UserTypes[i].GetEntry()
		}
		if err := ut.db.UserTypesUpsertUserTypes(ut.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range UserTypes {
		ut.UserTypes.Data[UserTypes[i].GetID()] = UserTypes[i]
	}

	return nil
}

func (ut *UserTypes) RemoveUserType(UserType universe.UserType, updateDB bool) error {
	ut.UserTypes.Mu.Lock()
	defer ut.UserTypes.Mu.Unlock()

	if _, ok := ut.UserTypes.Data[UserType.GetID()]; !ok {
		return errors.Errorf("user type not found")
	}

	if updateDB {
		if err := ut.db.UserTypesRemoveUserTypeByID(ut.ctx, UserType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	delete(ut.UserTypes.Data, UserType.GetID())

	return nil
}

func (ut *UserTypes) RemoveUserTypes(UserTypes []universe.UserType, updateDB bool) error {
	ut.UserTypes.Mu.Lock()
	defer ut.UserTypes.Mu.Unlock()

	for i := range UserTypes {
		if _, ok := ut.UserTypes.Data[UserTypes[i].GetID()]; !ok {
			return errors.Errorf("user type not found: %ut", UserTypes[i].GetID())
		}
	}

	if updateDB {
		ids := make([]uuid.UUID, len(UserTypes))
		for i := range UserTypes {
			ids[i] = UserTypes[i].GetID()
		}
		if err := ut.db.UserTypesRemoveUserTypesByIDs(ut.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range UserTypes {
		delete(ut.UserTypes.Data, UserTypes[i].GetID())
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
		UserType, err := ut.NewUserType(*entries[i].UserTypeID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new user type: %ut", entries[i].UserTypeID)
		}
		if err := UserType.LoadFromEntry(entries[i]); err != nil {
			return errors.WithMessagef(err, "failed to load user type from entry: %ut", *entries[i].UserTypeID)
		}
		ut.UserTypes.Store(*entries[i].UserTypeID, UserType)
	}

	universe.GetNode().AddAPIRegister(ut)

	ut.log.Info("User types loaded")

	return nil
}

func (ut *UserTypes) Save() error {
	ut.log.Info("Saving spate types...")

	ut.UserTypes.Mu.RLock()
	defer ut.UserTypes.Mu.RUnlock()

	entries := make([]*entry.UserType, 0, len(ut.UserTypes.Data))
	for _, UserType := range ut.UserTypes.Data {
		entries = append(entries, UserType.GetEntry())
	}

	if err := ut.db.UserTypesUpsertUserTypes(ut.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert user types")
	}

	ut.log.Info("User types saved")

	return nil
}
