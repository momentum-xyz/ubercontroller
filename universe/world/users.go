package world

import (
	"fmt"
	cmath2 "github.com/momentum-xyz/controller/pkg/cmath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) GetUser(userID uuid.UUID, recursive bool) (universe.User, bool) {
	return w.Space.GetUser(userID, false)
}

func (w *World) GetUsers(recursive bool) map[uuid.UUID]universe.User {
	return w.Space.GetUsers(false)
}

func (w *World) AddUser(user universe.User, updateDB bool) error {
	//w.Users.Mu.Lock()
	//defer w.Users.Mu.Unlock()
	var err error
	defer func() {
		if err != nil {
			user.Shutdown()
		}
	}()
	fmt.Println("ee")

	exUser, ok := w.Users.Load(user.GetID())

	if ok {
		if exUser != user {
			if exUser.GetSessionID() == user.GetSessionID() {
				w.log.Info(
					"World: same session, must be teleport or network failure: world %s, user %s", w.GetID(),
					user.GetID(),
				)
			} else {
				w.log.Info("World: double-login detected for world %s, user %s", w.GetID(), exUser.GetID())

				exUser.SendDirectly(posbus.NewSignalMsg(posbus.SignalDualConnection).WebsocketMessage())

				time.Sleep(time.Millisecond * 100)
			}
			exUser.Shutdown()
			//w.RemoveUser(exUser, true)
		} else {
			//TODO: handle this (if this ever can happen)
			panic("implement me")
		}
	}

	if err = user.SetWorld(w, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set world %s to user: %s", w.GetID(), user.GetID())
	}

	// effectively replace user if exists
	if err = w.Space.AddUser(user, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to add user %s to world: %s", user.GetID(), w.GetID())
	}

	err = w.initializeUnity(user)
	return err
}

func (w *World) initializeUnity(user universe.User) error {
	// TODO: rest of startup logic
	fmt.Println(w.metaMsg.Load())
	if err := user.SendDirectly(w.metaMsg.Load()); err != nil {
		return errors.WithMessage(err, "failed to send meta msg")
	}

	// TODO: fix circular dependency
	if err := user.SendDirectly(posbus.NewSendPositionMsg(cmath2.Vec3(user.GetPosition())).WebsocketMessage()); err != nil {
		return errors.WithMessage(err, "failed to send position")
	}

	//go func() {
	//	time.Sleep(30 * time.Second)
	//	user.ReleaseSendBuffer()
	//}()

	w.Space.SendSpawnMessage(user.SendDirectly, true)
	time.Sleep(1 * time.Second)
	user.SendDirectly(
		posbus.NewSignalMsg(
			posbus.SignalSpawn,
		).WebsocketMessage(),
	)

	fmt.Println("ee10")
	return nil
}

func (w *World) RemoveUser(user universe.User, updateDB bool) error {
	w.Users.Mu.Lock()
	defer w.Users.Mu.Unlock()

	return w.noLockRemoveUser(user, updateDB)
}

func (w *World) noLockRemoveUser(user universe.User, updateDB bool) error {

	if user.GetWorld().GetID() != w.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", user.GetWorld().GetID(), w.GetID())
	}
	if err := user.SetWorld(nil, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set world nil to user: %s", user.GetID())
	}
	delete(w.Users.Data, user.GetID())
	user.Shutdown()

	return nil
}

func (w *World) SendToUser(userID uuid.UUID, msg *websocket.PreparedMessage, recursive bool) error {
	return errors.Errorf("implement me")
}

func (w *World) SendToUsers(msg *websocket.PreparedMessage, recursive bool) error {
	return errors.Errorf("implement me")
}

func (w *World) GetUserSpawnPosition(userID uuid.UUID) cmath.Vec3 {
	return cmath.Vec3{X: 40, Y: 40, Z: 40}
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

func (w *World) DisconnectUser(userID uuid.UUID) {

}
