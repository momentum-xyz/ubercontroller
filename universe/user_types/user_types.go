package user_types

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/user_type"
)

var _ universe.UserTypes = (*UserTypes)(nil)

type UserTypes struct {
	ctx       types.LoggerContext
	log       *zap.SugaredLogger
	db        database.DB
	userTypes *generic.SyncMap[umid.UMID, universe.UserType]
}

func NewUserTypes(db database.DB) *UserTypes {
	return &UserTypes{
		db:        db,
		userTypes: generic.NewSyncMap[umid.UMID, universe.UserType](0),
	}
}

func (u *UserTypes) Initialize(ctx types.LoggerContext) error {
	u.ctx = ctx
	u.log = ctx.Logger()

	return nil
}

func (u *UserTypes) CreateUserType(userTypeID umid.UMID) (universe.UserType, error) {
	userType := user_type.NewUserType(userTypeID, u.db)

	if err := userType.Initialize(u.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize user type: %s", userTypeID)
	}
	if err := u.AddUserType(userType, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add user type: %s", userTypeID)
	}

	return userType, nil
}

func (u *UserTypes) FilterUserTypes(predicateFn universe.UserTypesFilterPredicateFn) map[umid.UMID]universe.UserType {
	return u.userTypes.Filter(predicateFn)
}

func (u *UserTypes) GetUserType(userTypeID umid.UMID) (universe.UserType, bool) {
	userType, ok := u.userTypes.Load(userTypeID)
	return userType, ok
}

func (u *UserTypes) GetUserTypes() map[umid.UMID]universe.UserType {
	u.userTypes.Mu.RLock()
	defer u.userTypes.Mu.RUnlock()

	userTypes := make(map[umid.UMID]universe.UserType, len(u.userTypes.Data))

	for id, userType := range u.userTypes.Data {
		userTypes[id] = userType
	}

	return userTypes
}

func (u *UserTypes) AddUserType(userType universe.UserType, updateDB bool) error {
	u.userTypes.Mu.Lock()
	defer u.userTypes.Mu.Unlock()

	if updateDB {
		if err := u.db.GetUserTypesDB().UpsertUserType(u.ctx, userType.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	u.userTypes.Data[userType.GetID()] = userType

	return nil
}

func (u *UserTypes) AddUserTypes(userTypes []universe.UserType, updateDB bool) error {
	u.userTypes.Mu.Lock()
	defer u.userTypes.Mu.Unlock()

	if updateDB {
		entries := make([]*entry.UserType, len(userTypes))
		for i := range userTypes {
			entries[i] = userTypes[i].GetEntry()
		}
		if err := u.db.GetUserTypesDB().UpsertUserTypes(u.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range userTypes {
		u.userTypes.Data[userTypes[i].GetID()] = userTypes[i]
	}

	return nil
}

func (u *UserTypes) RemoveUserType(userType universe.UserType, updateDB bool) (bool, error) {
	u.userTypes.Mu.Lock()
	defer u.userTypes.Mu.Unlock()

	if _, ok := u.userTypes.Data[userType.GetID()]; !ok {
		return false, nil
	}

	if updateDB {
		if err := u.db.GetUserTypesDB().RemoveUserTypeByID(u.ctx, userType.GetID()); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(u.userTypes.Data, userType.GetID())

	return true, nil
}

func (u *UserTypes) RemoveUserTypes(userTypes []universe.UserType, updateDB bool) (bool, error) {
	u.userTypes.Mu.Lock()
	defer u.userTypes.Mu.Unlock()

	for i := range userTypes {
		if _, ok := u.userTypes.Data[userTypes[i].GetID()]; !ok {
			return false, nil
		}
	}

	if updateDB {
		ids := make([]umid.UMID, len(userTypes))
		for i := range userTypes {
			ids[i] = userTypes[i].GetID()
		}
		if err := u.db.GetUserTypesDB().RemoveUserTypesByIDs(u.ctx, ids); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range userTypes {
		delete(u.userTypes.Data, userTypes[i].GetID())
	}

	return true, nil
}

func (u *UserTypes) Load() error {
	u.log.Info("Loading user types...")

	entries, err := u.db.GetUserTypesDB().GetUserTypes(u.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get user types")
	}

	for _, utEntry := range entries {
		userType, err := u.CreateUserType(utEntry.UserTypeID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new user type: %s", utEntry.UserTypeID)
		}
		if err := userType.LoadFromEntry(utEntry); err != nil {
			return errors.WithMessagef(err, "failed to load user type from entry: %s", utEntry.UserTypeID)
		}
	}

	universe.GetNode().AddAPIRegister(u)

	u.log.Infof("User types loaded: %d", u.userTypes.Len())

	return nil
}

func (u *UserTypes) Save() error {
	u.log.Info("Saving user types...")

	u.userTypes.Mu.RLock()
	defer u.userTypes.Mu.RUnlock()

	entries := make([]*entry.UserType, 0, len(u.userTypes.Data))
	for _, userType := range u.userTypes.Data {
		entries = append(entries, userType.GetEntry())
	}

	if err := u.db.GetUserTypesDB().UpsertUserTypes(u.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert user types")
	}

	u.log.Infof("User types saved: %d", len(u.userTypes.Data))

	return nil
}
