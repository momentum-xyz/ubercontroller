package user

import (
	"fmt"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func (u *User) OnMessage(buf []byte) error {
	msg, err := posbus.Decode(buf)
	if err != nil {
		return err
	}
	switch msg.GetType() {
	case posbus.TypeMyTransform:
		return u.UpdatePosition(msg.(*posbus.MyTransform))
		//FIXME
		//if err := u.UpdatePosition(msg.Msg()); err != nil {
		//	return errors.WithMessage(err, "failed to handle: send transform")
		//}

	//case posbus.T:
	//	if err := u.InteractionHandler(msg.AsTriggerInteraction()); err != nil {
	//		return errors.WithMessage(err, "failed to handle: interaction")
	//	}
	case posbus.TypeTeleportRequest:
		return u.Teleport(msg.(*posbus.TeleportRequest).Target)
	case posbus.TypeSignal:
		return u.SignalsHandler(msg.(*posbus.Signal))
	case posbus.TypeObjectTransform:
		if err := u.UpdateObjectTransform(msg.(*posbus.ObjectTransform)); err != nil {
			return errors.WithMessage(err, "failed to update object transform")
		}
	case posbus.TypeLockObject:
		return u.LockObject(msg.(*posbus.LockObject))
	case posbus.TypeUnlockObject:
		return u.UnlockObject(msg.(*posbus.UnlockObject))
	case posbus.TypeHighFive:
		return u.HandleHighFive(msg.(*posbus.HighFive))
	default:
		return errors.Errorf("unknown message: %d", msg.GetType())
	}

	return nil
}

func (u *User) UpdateObjectTransform(msg *posbus.ObjectTransform) error {
	object, ok := universe.GetNode().GetObjectFromAllObjects(msg.ID)
	if !ok {
		return errors.Errorf("object not found: %s", msg.ID)
	}
	return object.SetTransform(utils.GetPTR(msg.Transform), true)
}

func (u *User) Teleport(target umid.UMID) error {
	world, ok := universe.GetNode().GetWorlds().GetWorld(target)
	if !ok {
		// send buffer is locked at this point, so direct:
		u.SendDirectly(posbus.WSMessage(&posbus.Signal{Value: posbus.SignalWorldDoesNotExist}))
		return errors.New("Target world does not exist")
	}
	u.leaveCurrentWorld()
	return world.AddUser(u, true)
}

func (u *User) leaveCurrentWorld() {
	if oldWorld := u.GetWorld(); oldWorld != nil {
		oldWorld.RemoveUser(u, true)
	}
}

func (u *User) SignalsHandler(s *posbus.Signal) error {
	fmt.Printf("Got Signal %+v\n", s)
	switch s.Value {
	case posbus.SignalLeaveWorld:
		u.leaveCurrentWorld()
	}

	return nil
}

//func (u *User) InteractionHandler(m *posbus.TriggerInteraction) error {
//	kind := m.Kind()
//	targetUUID := m.Target()
//	flag := m.Flag()
//	label := m.Label()
//	u.log.Infof(
//		"Incoming interaction for user: %s, kind: %d, target: %s, flag: %d, label: %s",
//		u.GetID(), kind, targetUUID, flag, label,
//	)
//
//	switch kind {
//	case posbus.TriggerEnteredObject:
//		object, ok := universe.GetNode().GetObjectFromAllObjects(targetUUID)
//		if !ok {
//			return errors.WithMessage(
//				errors.Errorf("object not found: %s", targetUUID), "failed to handle: enter object",
//			)
//		}
//		if err := object.AddUser(u, true); err != nil {
//			return errors.WithMessage(
//				errors.Errorf("failed to add user to object: %s", targetUUID), "failed to handle: enter object",
//			)
//		}
//		return nil
//	case posbus.TriggerLeftObject:
//		object, ok := universe.GetNode().GetObjectFromAllObjects(targetUUID)
//		if !ok {
//			return errors.WithMessage(
//				errors.Errorf("object not found: %s", targetUUID), "failed to handle: left object",
//			)
//		}
//		if _, err := object.RemoveUser(u, true); err != nil {
//			return errors.WithMessage(
//				errors.Errorf("failed to remove user from object: %s", targetUUID), "failed to handle: left object",
//			)
//		}
//		return nil
//	case posbus.TriggerHighFive:
//		if err := u.HandleHighFive(m); err != nil {
//			u.log.Warn(errors.WithMessage(err, "InteractionHandler: trigger high fives: failed to handle high five"))
//		}
//		return nil
//	}
//	// case posbus.TriggerStake:
//	// 	u.HandleStake(m)
//
//	return errors.Errorf("unknown message: %d", kind)
//}

func (u *User) LockObject(lock *posbus.LockObject) error {
	objectId := lock.ID
	object, ok := u.GetWorld().GetObjectFromAllObjects(objectId)
	if !ok {
		return errors.Errorf("object not found: %s", objectId)
	}

	result := object.LockUIObject(u, 1)
	if result {
		return u.GetWorld().Send(
			posbus.WSMessage(&posbus.LockObjectResponse{ID: objectId, State: 1, LockOwner: u.GetID()}),
			true,
		)
	}
	return u.Send(posbus.WSMessage(&posbus.LockObjectResponse{ID: objectId, State: 0, LockOwner: u.GetID()}))
}

func (u *User) UnlockObject(lock *posbus.UnlockObject) error {
	objectId := lock.ID
	object, ok := u.GetWorld().GetObjectFromAllObjects(objectId)
	if !ok {
		return errors.Errorf("object not found: %s", objectId)
	}

	result := object.LockUIObject(u, 0)

	if result {
		return u.GetWorld().Send(
			posbus.WSMessage(&posbus.LockObjectResponse{ID: objectId, State: 1, LockOwner: u.GetID()}),
			true,
		)
	}
	return u.Send(posbus.WSMessage(&posbus.LockObjectResponse{ID: objectId, State: 1, LockOwner: u.GetID()}))
}

func (u *User) HandleHighFive(m *posbus.HighFive) error {
	targetID := m.ReceiverID
	if targetID == u.GetID() {
		return errors.New("high-five yourself not permitted")
	}

	world := u.GetWorld()
	_, ok := world.GetUser(targetID, false)
	if !ok {
		u.Send(posbus.WSMessage(&posbus.Notification{NotifyType: posbus.NotificationTextMessage, Value: "Target user not found"}))
		return errors.Errorf("failed to get target: %s", targetID)
	}

	/* TODO: implement as generic notification message
	u.Send(posbus.WSMessage(&posbus.Notification{NotifyType: posbus.NotificationHighFive, Value: tName}))
	target.Send(posbus.WSMessage(&msg))
	*/
	// For now, just broadcast the HighFive to the world
	world.Send(posbus.WSMessage(m), false)

	/* TODO: implement as generic (3D) effect
	effectsEmitterID := world.GetSettings().Objects["effects_emitter"]
	effect := posbus.NewTriggerTransitionalBridgingEffectsOnPositionMsg(1)
	effect.SetEffect(0, effectsEmitterID, u.GetTransform(), target.GetTransform(), 1001)
	u.GetWorld().Send(effect.WSMessage(), false)
	*/

	// TODO: stats tracking.
	//go u.SendHighFiveStats(target)

	return nil
}
