package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types"
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
	nCurrentSends               atomic.Int64
	world                       *world.World
}

func (u *User) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}
	u.ctx = ctx
	u.log = log
	u.nCurrentSends.Store(0)
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
