package user

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

var _ universe.User = (*User)(nil)

type User struct {
	id        umid.UMID
	sessionID umid.UMID
	db        database.DB
	// The websocket connection.
	conn *websocket.Conn

	transform cmath.TransformNoScale
	//pos          *cmath.Vec3 // going to data part to posMsgBuffer content for simple access
	//rotation     *cmath.Vec3 // going to data part to posMsgBuffer content for simple access
	//posMsgBuffer []byte

	lastPositionUpdateTimestamp     int64
	lastSendPositionUpdateTimestamp int64
	userType                        universe.UserType
	log                             *zap.SugaredLogger
	ctx                             context.Context
	send                            chan *websocket.PreparedMessage
	mu                              deadlock.RWMutex
	object                          universe.Object
	world                           universe.World
	profile                         *entry.UserProfile
	options                         *entry.UserOptions
	bufferSends                     atomic.Bool
	numSendsQueued                  atomic.Int64
	directLock                      sync.Mutex
	offlineTimer                    *generic.TimerSet[umid.UMID]
}

func NewUser(id umid.UMID, db database.DB) *User {
	return &User{
		id:           id,
		offlineTimer: generic.NewTimerSet[umid.UMID](),
		db:           db,
	}
}

func (u *User) GetID() umid.UMID {
	return u.id
}

//func (u *User) GetTransform() cmath.UserTransform {
//	t := cmath.NewUserTransform()
//	t.CopyTo(&u.transform)
//	return t
//}

func (u *User) GetTransform() *cmath.TransformNoScale {
	return &u.transform
}

func (u *User) SetTransform(t cmath.TransformNoScale) {
	u.transform = t.Copy()
}

func (u *User) SetPosition(p cmath.Vec3) {
	fmt.Printf("Set position for user %+v, as %+v\n", u.id, p)
	u.transform.Position = p
}

func (u *User) GetPosition() cmath.Vec3 {
	return u.transform.Position
}

func (u *User) GetRotation() cmath.Vec3 {
	return u.transform.Rotation
}

func (u *User) GetUserDefinition() *posbus.UserData {
	d := new(posbus.UserData)
	d.Transform = u.transform
	d.ID = u.id
	guestUserTypeID, _ := common.GetGuestUserTypeID()
	d.IsGuest = u.userType.GetID() == guestUserTypeID
	d.Name = *u.profile.Name
	if u.profile.AvatarHash != nil {
		d.Avatar = *u.profile.AvatarHash
	}
	return d
}

//func (u *User) GetPosBuffer() []byte {
//	return u.posMsgBuffer
//}

func (u *User) GetLastPosTime() int64 {
	return u.lastPositionUpdateTimestamp
}

func (u *User) GetLastSendPosTime() int64 {
	return u.lastSendPositionUpdateTimestamp
}

func (u *User) SetLastSendPosTime(i int64) {
	u.lastSendPositionUpdateTimestamp = i
}

func (u *User) GetWorld() universe.World {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.world
}

func (u *User) SetWorld(world universe.World) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.world = world
}

func (u *User) GetObject() universe.Object {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.object
}

func (u *User) SetObject(object universe.Object) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.object = object
}

func (u *User) GetUserType() universe.UserType {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.userType
}

func (u *User) SetOfflineTimer() (bool, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	u.offlineTimer.Set(u.ctx, u.id, time.Minute*20, u.DeleteTemporaryUser)
	u.log.Infof("Timer set: %s", u.GetID())
	return true, nil
}

func (u *User) IsTemporaryUser() (bool, error) {
	guestUserTypeID, err := common.GetGuestUserTypeID()
	if err != nil {
		return false, errors.WithMessage(err, "failed to get guestUserTypeID")
	}
	if u.GetUserType().GetID() != guestUserTypeID {
		return false, nil
	}

	return true, nil
}

func (u *User) DeleteTemporaryUser(uid umid.UMID) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.log.Infof("Deleting temp user: %s", u.GetID())

	if err := u.db.GetUsersDB().RemoveUserByID(u.ctx, uid); err != nil {
		return errors.WithMessage(err, "failed to delete temporary user by id")
	}

	return nil
}

func (u *User) GetProfile() *entry.UserProfile {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.profile
}

func (u *User) Initialize(ctx types.LoggerContext) error {
	u.ctx = ctx
	u.log = ctx.Logger()
	u.bufferSends.Store(true)
	u.numSendsQueued.Store(chanIsClosed)
	//u.posMsgBuffer = posbus.NewSendTransformBuffer(u.GetID())
	//u.transform.Position = (*cmath.Vec3)(unsafe.Add(unsafe.Pointer(&u.posMsgBuffer[0]), 16))
	//u.transform.Rotation = (*cmath.Vec3)(unsafe.Add(unsafe.Pointer(&u.posMsgBuffer[0]), 16+3*4))
	return nil
}

func (u *User) SetUserType(userType universe.UserType, updateDB bool) error {
	if userType == nil {
		return errors.Errorf("user type is nil")
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if updateDB {
		if err := u.db.GetUsersDB().UpdateUserUserTypeID(u.ctx, u.GetID(), userType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	u.userType = userType

	return nil
}

func (u *User) Run() error {
	isTemporaryUser, err := u.IsTemporaryUser()
	if err != nil {
		return errors.WithMessagef(err, "failed to assess if user is temporary user: %s", u.GetID())
	}

	if isTemporaryUser && u.offlineTimer != nil {
		u.offlineTimer.StopAll()
	}

	u.StartIOPumps()
	return nil
}

func (u *User) Stop() error {
	ns := u.numSendsQueued.Add(1)
	if ns >= 0 {
		u.send <- nil
	}

	return nil
}

func (u *User) Update() error {
	//TODO implement me
	panic("implement me")
}

func (u *User) Load() error {
	u.log.Infof("Loading user: %s", u.GetID())

	entry, err := u.db.GetUsersDB().GetUserByID(u.ctx, u.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get user by umid")
	}

	if err := u.LoadFromEntry(entry); err != nil {
		return errors.WithMessage(err, "failed to load from entry")
	}

	u.log.Infof("User loaded: %s", u.GetID())

	return nil
}

func (u *User) LoadFromEntry(entry *entry.User) error {
	if entry.UserID != u.GetID() {
		return errors.Errorf("user ids mismatch: %s != %s", entry.UserID, u.GetID())
	}

	u.profile = &entry.Profile
	u.options = entry.Options

	node := universe.GetNode()
	userType, ok := node.GetUserTypes().GetUserType(entry.UserTypeID)
	if !ok {
		return errors.Errorf("failed to get user type: %s", entry.UserTypeID)
	}
	if err := u.SetUserType(userType, false); err != nil {
		return errors.WithMessagef(err, "failed to set user type: %s", entry.UserTypeID)
	}

	return nil
}

func (u *User) UpdatePosition(t *posbus.MyTransform) error {
	//u.SetTransform(t)
	u.transform = cmath.TransformNoScale(*t)
	// not locking will speed up but introduce minor data race with zero impact
	//u.world.users.positionLock.RLock()
	//copy(u.posMsgBuffer[16:40], data)
	//u.world.users.positionLock.RUnlock()
	u.lastPositionUpdateTimestamp = time.Now().Unix()

	return nil
}

//func (u *User) TeleportToWorld(umid umid.UMID) {
//	//url := universe.GetNode().ResolveNodeByWorldID(umid)
//	//u.Send(posbus.NewTeleportRequest(umid, url).WSMessage())
//	//u.close(true)
//
//}
