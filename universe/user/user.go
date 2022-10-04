package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/world"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

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
	name                        string
	log                         *zap.SugaredLogger
	ctx                         context.Context
	send                        chan *websocket.PreparedMessage
	mu                          sync.RWMutex
	world                       *world.World
	profile                     *entry.UserProfile
	options                     *entry.UserOptions
	sendMutex                   sync.Mutex
	readyToSend                 atomic.Bool
	quit                        atomic.Bool
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
	//TODO implement me
	panic("implement me")
}

func (u *User) SetWorld(world universe.World, updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (u *User) GetSpace() universe.Space {
	//TODO implement me
	panic("implement me")
}

func (u *User) SetSpace(space universe.Space, updateDB bool) error {
	//TODO implement me
	panic("implement me")
}

func (u *User) Update() error {
	//TODO implement me
	panic("implement me")
}

func (u *User) SetUserType(userType universe.UserType, updateDB bool) error {
	//TODO implement me
	panic("implement me")
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
	if *entry.UserID != u.GetID() {
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
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}
	u.ctx = ctx
	u.log = log
	u.quit.Store(false)
	u.readyToSend.Store(false)
	u.lastPositionUpdateTimestamp = int64(0)
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
	return u.name
}

//SetUserType(userType UserType, updateDB bool) error

//func
