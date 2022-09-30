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

func (s *UserTypes) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	s.ctx = ctx
	s.log = log

	return nil
}

func (s *UserTypes) NewUserType(UserTypeID uuid.UUID) (universe.UserType, error) {
	UserType := user_type.NewUserType(UserTypeID, s.db)

	if err := UserType.Initialize(s.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize user type: %s", UserTypeID)
	}
	if err := s.AddUserType(UserType, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add user type: %s", UserTypeID)
	}

	return UserType, nil
}

func (s *UserTypes) GetUserType(UserTypeID uuid.UUID) (universe.UserType, bool) {
	UserType, ok := s.UserTypes.Load(UserTypeID)
	return UserType, ok
}

func (s *UserTypes) GetUserTypes() map[uuid.UUID]universe.UserType {
	s.UserTypes.Mu.RLock()
	defer s.UserTypes.Mu.RUnlock()

	UserTypes := make(map[uuid.UUID]universe.UserType, len(s.UserTypes.Data))

	for id, UserType := range s.UserTypes.Data {
		UserTypes[id] = UserType
	}

	return UserTypes
}

func (s *UserTypes) AddUserType(UserType universe.UserType, updateDB bool) error {
	s.UserTypes.Mu.Lock()
	defer s.UserTypes.Mu.Unlock()

	if updateDB {
		if err := s.db.UserTypesUpsertUserType(s.ctx, UserType.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.UserTypes.Data[UserType.GetID()] = UserType

	return nil
}

func (s *UserTypes) AddUserTypes(UserTypes []universe.UserType, updateDB bool) error {
	s.UserTypes.Mu.Lock()
	defer s.UserTypes.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.UserType, len(UserTypes))
		for i := range UserTypes {
			entries[i] = UserTypes[i].GetEntry()
		}
		if err := s.db.UserTypesUpsertUserTypes(s.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range UserTypes {
		s.UserTypes.Data[UserTypes[i].GetID()] = UserTypes[i]
	}

	return nil
}

func (s *UserTypes) RemoveUserType(UserType universe.UserType, updateDB bool) error {
	s.UserTypes.Mu.Lock()
	defer s.UserTypes.Mu.Unlock()

	if _, ok := s.UserTypes.Data[UserType.GetID()]; !ok {
		return errors.Errorf("user type not found")
	}

	if updateDB {
		if err := s.db.UserTypesRemoveUserTypeByID(s.ctx, UserType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	delete(s.UserTypes.Data, UserType.GetID())

	return nil
}

func (s *UserTypes) RemoveUserTypes(UserTypes []universe.UserType, updateDB bool) error {
	s.UserTypes.Mu.Lock()
	defer s.UserTypes.Mu.Unlock()

	for i := range UserTypes {
		if _, ok := s.UserTypes.Data[UserTypes[i].GetID()]; !ok {
			return errors.Errorf("user type not found: %s", UserTypes[i].GetID())
		}
	}

	if updateDB {
		ids := make([]uuid.UUID, len(UserTypes))
		for i := range UserTypes {
			ids[i] = UserTypes[i].GetID()
		}
		if err := s.db.UserTypesRemoveUserTypesByIDs(s.ctx, ids); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range UserTypes {
		delete(s.UserTypes.Data, UserTypes[i].GetID())
	}

	return nil
}

func (s *UserTypes) Load() error {
	s.log.Info("Loading user types...")

	entries, err := s.db.UserTypesGetUserTypes(s.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get user types")
	}

	for i := range entries {
		UserType, err := s.NewUserType(*entries[i].UserTypeID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new user type: %s", entries[i].UserTypeID)
		}
		if err := UserType.LoadFromEntry(entries[i]); err != nil {
			return errors.WithMessagef(err, "failed to load user type from entry: %s", *entries[i].UserTypeID)
		}
		s.UserTypes.Store(*entries[i].UserTypeID, UserType)
	}

	universe.GetNode().AddAPIRegister(s)

	s.log.Info("User types loaded")

	return nil
}

func (s *UserTypes) Save() error {
	s.log.Info("Saving spate types...")

	s.UserTypes.Mu.RLock()
	defer s.UserTypes.Mu.RUnlock()

	entries := make([]*entry.UserType, 0, len(s.UserTypes.Data))
	for _, UserType := range s.UserTypes.Data {
		entries = append(entries, UserType.GetEntry())
	}

	if err := s.db.UserTypesUpsertUserTypes(s.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert user types")
	}

	s.log.Info("User types saved")

	return nil
}
