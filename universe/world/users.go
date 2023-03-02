package world

import (
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) GetUser(userID uuid.UUID, recursive bool) (universe.User, bool) {
	return w.ToObject().GetUser(userID, false)
}

func (w *World) GetUsers(recursive bool) map[uuid.UUID]universe.User {
	return w.ToObject().GetUsers(false)
}

func (w *World) AddUser(user universe.User, updateDB bool) error {
	//w.Users.Mu.Lock()
	//defer w.Users.Mu.Unlock()
	var err error
	defer func() {
		if err != nil {
			user.Stop()
		}
	}()

	exUser, ok := w.Users.Load(user.GetID())

	if ok {
		w.log.Infof("Found existing user: %+v\n", exUser.GetID())
		if exUser != user {
			if exUser.GetSessionID() == user.GetSessionID() {
				w.log.Infof(
					"World: same session, must be teleport or network failure: world %s, user %s", w.GetID(),
					user.GetID(),
				)
			} else {
				w.log.Infof("World: double-login detected for world %s, user %s", w.GetID(), exUser.GetID())

				exUser.SendDirectly(posbus.WrapAsMessage(posbus.SignalType, posbus.SignalDualConnection))

				time.Sleep(time.Millisecond * 100)
			}
			exUser.Stop()
			//w.RemoveUser(exUser, true)
		} else {
			//TODO: handle this (if this ever can happen)
			panic("implement me")
		}
	}

	w.log.Infof("Setworld: %+v\n", user.GetID())
	user.SetWorld(w)

	w.log.Infof("AddUser: %+v\n", user.GetID())
	// effectively replace user if exists
	if err = w.ToObject().AddUser(user, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to add user %s to world: %s", user.GetID(), w.GetID())
	}

	err = w.initializeUnity(user)
	return err
}

func (w *World) RemoveUser(user universe.User, updateDB bool) (bool, error) {
	w.Users.Mu.Lock()
	defer w.Users.Mu.Unlock()

	return w.noLockRemoveUser(user, updateDB)
}

func (w *World) Send(msg *websocket.PreparedMessage, recursive bool) error {
	return w.ToObject().Send(msg, false)
}

func (w *World) GetUserSpawnPosition(userID uuid.UUID) cmath.Vec3 {
	return cmath.Vec3{X: 40, Y: 40, Z: 40}
}

func (w *World) noLockRemoveUser(user universe.User, updateDB bool) (bool, error) {
	if user.GetWorld().GetID() != w.GetID() {
		return false, errors.Errorf("worlds mismatch: %s != %s", user.GetWorld().GetID(), w.GetID())
	}

	if _, ok := w.Users.Data[user.GetID()]; !ok {
		return false, nil
	}

	user.SetWorld(nil)
	user.Stop()

	delete(w.Users.Data, user.GetID())

	// clean up all locks hold by this user,
	// temporarily here.
	w.allObjects.Mu.RLock()
	defer w.allObjects.Mu.RUnlock()

	for _, child := range w.allObjects.Data {
		if child.LockUnityObject(user, 0) {
			w.Send(
				posbus.WrapAsMessage(
					posbus.ObjectLockResultType,
					posbus.ObjectLockResultData{ID: child.GetID(), Result: 0, LockOwner: user.GetID()},
				), true,
			)
		}
	}

	return true, nil
}

func (w *World) initializeUnity(user universe.User) error {
	// TODO: rest of startup logic
	if err := user.SendDirectly(w.metaMsg.Load()); err != nil {
		return errors.WithMessage(err, "failed to send meta msg")
	}
	//user.SendDirectly(w.Object.Ge)

	// TODO: fix circular dependency
	if err := user.SendDirectly(
		posbus.NewSendPositionMsg(
			user.GetPosition(), user.GetRotation(), cmath.Vec3{X: 0, Y: 0, Z: 0},
		).WebsocketMessage(),
	); err != nil {
		return errors.WithMessage(err, "failed to send position")
	}

	//go func() {
	//	time.Sleep(30 * time.Second)
	//	user.ReleaseSendBuffer()
	//}()

	w.SendSpawnMessage(user.SendDirectly, true)
	w.log.Infof("Sent Spawn: %+v\n", user.GetID())
	time.Sleep(1 * time.Second)

	w.SendAllAutoAttributes(user.SendDirectly, true)
	w.log.Infof("Sent Textures: %+v\n", user.GetID())
	user.ReleaseSendBuffer()
	return nil
}

//func (w *World) SpawnUser(userID uuid.UUID, sessionID uuid.UUID, socketConnection *websocket.Conn) {
//
//	if exclient, ok := h.clients[x.ID]; ok && exclient.quiueID != x.quiueID {
//		if exclient.SessionID == x.SessionID {
//			h.unregister <- exclient
//		} else {
//			Logln(0, "Double-login detected for", x.ID)
//			msg := make([]byte, 24)
//			binary.LittleEndian.PutUint64(msg[0:8], msgSignal)
//			binary.LittleEndian.PutUint64(msg[8:16], SignalDualConn)
//			binary.LittleEndian.PutUint64(msg[16:24], ^msgSignal)
//			omsg, _ := websocket.NewPreparedMessage(websocket.BinaryMessage, msg)
//			exclient.conn.WritePreparedMessage(omsg)
//			// exclient.send <- omsg
//			time.Sleep(time.Millisecond * 300)
//			h.unregister <- exclient
//		}
//		go func() {
//			time.Sleep(time.Millisecond * 100)
//			h.register <- x
//		}()
//		return
//	}
//	defer func() {
//		// Logln(4, "Reg Done")
//		Logf(0, "Spawned %s on %s", x.ID, x.hub.ID)
//	}()
//
//	Logln(1, "Registering user: "+x.ID.String())
//	x.send = make(chan *websocket.PreparedMessage, 32)
//
//	binid, _ := x.ID.MarshalBinary()
//
//	copy(x.UnityBitsID[0:16], UnityUUID(binid))
//
//	x.pos = make([]byte, outPosMessageSize)
//	copy(x.pos[0:16], x.UnityBitsID[0:16])
//
//	x.send <- NewMessage(msgWorld, []byte(h.ID.String()))
//	Logln(4, x.ipos)
//	bpos := PackPos(x.ipos)
//	x.send <- NewMessage(msgSelfPos, bpos)
//	copy(x.pos[16:28], bpos[0:12])
//	x.hub = h
//	go x.IOPump()
//	h.clients[x.ID] = x
//	Logln(1, "Registration done: "+x.ID.String())
//
//}
