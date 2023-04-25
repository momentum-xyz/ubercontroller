package world

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) GetUser(userID umid.UMID, recursive bool) (universe.User, bool) {
	return w.ToObject().GetUser(userID, false)
}

func (w *World) GetUsers(recursive bool) map[umid.UMID]universe.User {
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

				exUser.SendDirectly(
					posbus.WSMessage(&posbus.Signal{Value: posbus.SignalDualConnection}),
				)

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

	initPos := cmath.TransformNoScale{Position: cmath.Vec3{X: 0, Y: 0, Z: 0}}

	val, ok := universe.GetNode().GetObjectUserAttributes().GetValue(
		entry.ObjectUserAttributeID{
			AttributeID: entry.AttributeID{PluginID: universe.GetSystemPluginID(), Name: "last_known_position"},
			UserID:      user.GetID(), ObjectID: w.GetID()},
	)

	if ok {
		var pos cmath.TransformNoScale
		err := utils.MapDecode(val, &pos)
		if err != nil {
			initPos = pos
		}
	}

	user.SetTransform(initPos)

	w.log.Infof("AddUser: %+v\n", user.GetID())
	// effectively replace user if exists
	user.LockSendBuffer()
	if err = w.ToObject().AddUser(user, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to add user %s to world: %s", user.GetID(), w.GetID())
	}
	w.Send(
		posbus.WSMessage(&posbus.AddUsers{Users: []posbus.UserData{*user.GetUserDefinition()}}),
		true,
	)

	err = w.initializeUI(user)
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

func (w *World) GetUserSpawnPosition(userID umid.UMID) cmath.Vec3 {
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

	var mp entry.AttributeValue
	utils.MapEncode(user.GetTransform(), &mp)
	// fmt.Printf("PMAP: %+v, %+v\n", user.GetTransform(), mp)

	_, err := universe.GetNode().GetObjectUserAttributes().Upsert(
		entry.ObjectUserAttributeID{
			AttributeID: entry.AttributeID{PluginID: universe.GetSystemPluginID(), Name: "last_known_position"},
			UserID:      user.GetID(), ObjectID: w.GetID()},
		func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
			v := entry.AttributePayload{}
			v.Value = &mp
			return &v, nil
		},
		true,
	)
	if err != nil {
		return false, fmt.Errorf("storing last known position: %w", err)
	}
	user.Stop()

	delete(w.Users.Data, user.GetID())

	// clean up all locks hold by this user,
	// temporarily here.
	w.allObjects.Mu.RLock()
	defer w.allObjects.Mu.RUnlock()

	for _, child := range w.allObjects.Data {
		if child.LockUIObject(user, 0) {
			w.Send(
				posbus.WSMessage(&posbus.LockObjectResponse{ID: child.GetID(), State: 0, LockOwner: user.GetID()}),
				true,
			)
		}
	}

	w.Send(
		posbus.WSMessage(&posbus.RemoveUsers{Users: []umid.UMID{user.GetID()}}),
		true,
	)

	return true, nil
}

func (w *World) initializeUI(user universe.User) error {
	// TODO: rest of startup logic
	if err := user.SendDirectly(w.metaMsg.Load()); err != nil {
		return errors.WithMessage(err, "failed to send meta msg")
	}
	//user.SendDirectly(w.Object.Ge)

	// TODO: fix circular dependency
	if err := user.SendDirectly(
		posbus.WSMessage((*posbus.MyTransform)(user.GetTransform())),
	); err != nil {
		return errors.WithMessage(err, "failed to send position")
	}

	//go func() {
	//	time.Sleep(30 * time.Second)
	//	user.ReleaseSendBuffer()
	//}()

	w.SendSpawnMessage(user.SendDirectly, true)
	w.log.Infof("Sent Spawn: %+v\n", user.GetID())
	//time.Sleep(1 * time.Second)

	w.SendAllAutoAttributes(user.SendDirectly, true)
	w.log.Infof("Sent Textures: %+v\n", user.GetID())
	user.ReleaseSendBuffer()

	w.SendUsersSpawnMessage(user)
	w.log.Infof("Sent all users: %+v\n", user.GetID())
	return nil
}

//func (w *World) SpawnUser(userID umid.UMID, sessionID umid.UMID, socketConnection *websocket.Conn) {
//
//	if exclient, ok := h.clients[x.UMID]; ok && exclient.quiueID != x.quiueID {
//		if exclient.SessionID == x.SessionID {
//			h.unregister <- exclient
//		} else {
//			Logln(0, "Double-login detected for", x.UMID)
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
//		Logf(0, "Spawned %s on %s", x.UMID, x.hub.UMID)
//	}()
//
//	Logln(1, "Registering user: "+x.UMID.String())
//	x.send = make(chan *websocket.PreparedMessage, 32)
//
//	binid, _ := x.UMID.MarshalBinary()
//
//	copy(x.UnityBitsID[0:16], UnityUUID(binid))
//
//	x.pos = make([]byte, outPosMessageSize)
//	copy(x.pos[0:16], x.UnityBitsID[0:16])
//
//	x.send <- NewMessage(msgWorld, []byte(h.UMID.String()))
//	Logln(4, x.ipos)
//	bpos := PackPos(x.ipos)
//	x.send <- NewMessage(msgSelfPos, bpos)
//	copy(x.pos[16:28], bpos[0:12])
//	x.hub = h
//	go x.IOPump()
//	h.clients[x.UMID] = x
//	Logln(1, "Registration done: "+x.UMID.String())
//
//}
