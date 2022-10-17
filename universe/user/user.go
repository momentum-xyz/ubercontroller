package user

import (
	"context"
	"github.com/sasha-s/go-deadlock"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.User = (*User)(nil)

type User struct {
	id        uuid.UUID
	sessionID uuid.UUID
	db        database.DB
	// The websocket connection.
	conn *websocket.Conn

	pos          *cmath.Vec3 // going to data part to posMsgBuffer content for simple access
	posMsgBuffer []byte

	lastPositionUpdateTimestamp int64
	userType                    universe.UserType
	log                         *zap.SugaredLogger
	ctx                         context.Context
	send                        chan *websocket.PreparedMessage
	mu                          deadlock.RWMutex
	world                       universe.World
	profile                     *entry.UserProfile
	options                     *entry.UserOptions
	bufferSends                 atomic.Bool
	numSendsQueued              atomic.Int64
	directLock                  sync.Mutex
}

func (u *User) GetPosBuffer() []byte {
	return u.posMsgBuffer
}

func (u *User) Run() error {
	//TODO implement me
	panic("implement me")
}

func (u *User) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (u *User) GetWorld() universe.World {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.world
}

func (u *User) SetWorld(world universe.World, updateDB bool) error {
	//TODO implement me
	u.mu.Lock()
	defer u.mu.Unlock()
	u.world = world

	return nil
}

func (u *User) GetSpace() universe.Space {
	//TODO implement me
	panic("implement me")
}

func (u *User) SetSpace(space universe.Space, updateDB bool) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (u *User) Update() error {
	//TODO implement me
	panic("implement me")
}

func (u *User) SetUserType(userType universe.UserType, updateDB bool) error {
	if userType == nil {
		return errors.Errorf("user type is nil")
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if updateDB {
		if err := u.db.UsersUpdateUserUserTypeID(u.ctx, u.id, userType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	u.userType = userType

	return nil
}

func NewUser(id uuid.UUID, db database.DB) *User {
	return &User{
		id: id,
		db: db,
	}
}

func (u *User) Load() error {
	u.log.Infof("Loading user: %s", u.GetID())

	entry, err := u.db.UsersGetUserByID(u.ctx, u.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get user by id")
	}

	if err := u.LoadFromEntry(entry); err != nil {
		return errors.WithMessage(err, "failed to load from entry")
	}

	//universe.GetNode().AddAPIRegister(u)

	u.log.Infof("User loaded: %s", u.GetID())

	return nil
}

func (u *User) LoadFromEntry(entry *entry.User) error {
	if entry.UserID != u.GetID() {
		return errors.Errorf("user ids mismatch: %s != %s", entry.UserID, u.GetID())
	}

	u.profile = entry.Profile
	u.options = entry.Options

	node := universe.GetNode()
	userType, ok := node.GetUserTypes().GetUserType(*entry.UserTypeID)
	if !ok {
		return errors.Errorf("failed to get user type: %s", entry.UserTypeID)
	}
	if err := u.SetUserType(userType, false); err != nil {
		return errors.WithMessagef(err, "failed to set user type: %s", entry.UserTypeID)
	}

	return nil
}

func (u *User) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}
	u.ctx = ctx
	u.log = log
	u.bufferSends.Store(true)
	u.numSendsQueued.Store(chanIsClosed)
	u.posMsgBuffer = message.NewSendPosBuffer(u.id)
	u.pos = (*cmath.Vec3)(unsafe.Add(unsafe.Pointer(&u.posMsgBuffer[0]), 16))
	return nil
}

func (u *User) UpdatePosition(data []byte) {

	// not locking will speed up but introduce minor data race with zero impact
	//u.world.users.positionLock.RLock()
	copy(u.posMsgBuffer[16:28], data)
	//u.world.users.positionLock.RUnlock()

	currentTime := time.Now().Unix()
	u.lastPositionUpdateTimestamp = currentTime
}

func (u *User) GetUserType() universe.UserType {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.userType
}

func (u *User) GetID() uuid.UUID {
	return u.id
}

func (u *User) GetName() string {
	return *u.profile.Name
}

func (u *User) SetPosition(p cmath.Vec3) {
	(*u.pos).X = p.X
	(*u.pos).Y = p.Y
	(*u.pos).Z = p.Z

}

func (u *User) GetPosition() cmath.Vec3 {
	return *u.pos
}

//SetUserType(userType UserType, updateDB bool) error

//func
