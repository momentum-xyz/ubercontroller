package user

import (
	"fmt"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func (u *User) OnMessage(msg *posbus.Message) error {
	switch msg.Type() {
	case posbus.TypeSendTransform:
		if err := u.UpdatePosition(msg.Msg()); err != nil {
			return errors.WithMessage(err, "failed to handle: send transform")
		}

	//case posbus.T:
	//	if err := u.InteractionHandler(msg.AsTriggerInteraction()); err != nil {
	//		return errors.WithMessage(err, "failed to handle: interaction")
	//	}
	case posbus.TypeGenericMessage:
		if err := u.GenericMessageHandler(msg.Msg()); err != nil {
			return errors.WithMessage(err, "failed to handle: relay to controller")
		}
	case posbus.TypeTeleportRequest:
		var tr posbus.TeleportRequest
		err := msg.DecodeTo(&tr)
		if err != nil {
			return errors.WithMessage(err, "failed to decode: teleport")
		}
		return u.Teleport(tr.Target)
	case posbus.TypeSignal:
		var signal posbus.Signal
		err := msg.DecodeTo(&signal)
		if err != nil {
			return errors.WithMessage(err, "failed to decode: signal")
		}
		return u.SignalsHandler(signal)
	//case posbus.TypeSetObjectPosition:
	//	if err := u.UpdateObjectPosition(msg.Msg()); err != nil {
	//		return errors.WithMessage(err, "failed to update object transform")
	//	}
	case posbus.TypeSetObjectLock:
		var lock posbus.SetObjectLock
		err := msg.DecodeTo(&lock)
		if err != nil {
			return errors.WithMessage(err, "failed to decode: set object lock")
		}
		return u.LockObject(lock)
	default:
		return errors.Errorf("unknown message: %d", msg.Type())
	}

	return nil
}

func (u *User) UpdateObjectPosition(msg posbus.ObjectPosition) error {
	object, ok := universe.GetNode().GetObjectFromAllObjects(msg.ID)
	if !ok {
		return errors.Errorf("object not found: %s", msg.ID)
	}
	return object.SetTransform(utils.GetPTR(msg.ObjectTransform), true)
}

func (u *User) Teleport(target umid.UMID) error {
	world, ok := universe.GetNode().GetWorlds().GetWorld(target)
	if !ok {
		u.Send(
			posbus.NewMessageFromData(
				posbus.TypeSignal, posbus.Signal{Value: posbus.SignalWorldDoesNotExist},
			).WSMessage(),
		)
		return errors.New("Target world does not exist")
	}
	if oldWorld := u.GetWorld(); oldWorld != nil {
		oldWorld.RemoveUser(u, true)
	}
	return world.AddUser(u, true)
}

func (u *User) GenericMessageHandler(msg []byte) error {
	//if m.Topic() == "emoji" {
	//	// TODO: comes as plugin?
	//	//u.HandleEmoji(msg.AsRelayToController())
	//}
	return nil
}

func (u *User) SignalsHandler(s posbus.Signal) error {
	fmt.Printf("Got Signal %+v\n", s)
	switch s.Value {
	case posbus.SignalLeaveWorld:
		if oldWorld := u.GetWorld(); oldWorld != nil {
			oldWorld.RemoveUser(u, true)
		}
	}
	//case posbus.SignalReady:
	//	u.ReleaseSendBuffer()
	//	//u.log.Debugf("Got signalReady from %s", u.umid.String())
	//	//TODO: Do we need it?
	//	//if err := u.world.SendWorldData(u); err != nil {
	//	//	log.Error(
	//	//		errors.WithMessagef(
	//	//			err, "User: SignalsHandler: SignalReady: failed to send world data: %s", u.UMID,
	//	//		),
	//	//	)
	//	//	u.world.unregisterUser <- u
	//	//	return
	//	//}
	//	//u.connection.EnableWriting()
	//}

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

func (u *User) LockObject(lock posbus.SetObjectLock) error {
	id := lock.ID
	state := lock.State

	object, ok := u.GetWorld().GetObjectFromAllObjects(id)
	if !ok {
		return errors.Errorf("object not found: %s", id)
	}

	result := object.LockUnityObject(u, state)
	newState := state
	if !result {
		newState = 1 - state
	}

	lock.State = newState

	return u.GetWorld().Send(posbus.NewMessageFromData(posbus.TypeSetObjectLock, lock).WSMessage(), true)
}

//func (u *User) HandleHighFive(m *posbus.TriggerInteraction) error {
//	targetID := m.Target()
//	if targetID == u.GetID() {
//		return errors.New("high-five yourself not permitted")
//	}
//
//	world := u.GetWorld()
//	target, ok := world.GetUser(targetID, false)
//	if !ok {
//		u.Send(
//			posbus.NewSimpleNotificationMsg(
//				posbus.DestinationReact, posbus.NotificationTextMessage, 0, "Target user not found",
//			).WSMessage(),
//		)
//		return errors.Errorf("failed to get target: %s", targetID)
//	}
//
//	var uName string
//	var tName string
//	uProfile := u.GetProfile()
//	tProfile := target.GetProfile()
//	if uProfile != nil && uProfile.Name != nil {
//		uName = *uProfile.Name
//	}
//	if tProfile != nil && tProfile.Name != nil {
//		tName = *tProfile.Name
//	}
//
//	high5Msg := struct {
//		SenderID   string `json:"senderId"`
//		ReceiverID string `json:"receiverId"`
//		Message    string `json:"message"`
//	}{
//		SenderID:   u.GetID().String(),
//		ReceiverID: targetID.String(),
//		Message:    uName,
//	}
//	high5Data, err := json.Marshal(&high5Msg)
//	if err != nil {
//		return errors.WithMessage(err, "failed to marshal high5 message")
//	}
//
//	u.Send(
//		posbus.NewSimpleNotificationMsg(
//			posbus.DestinationReact, posbus.NotificationHighFive, 0, tName,
//		).WSMessage(),
//	)
//	target.Send(posbus.NewGenericMessage("high5", high5Data).WSMessage())
//
//	effectsEmitterID := world.GetSettings().Objects["effects_emitter"]
//	effect := posbus.NewTriggerTransitionalBridgingEffectsOnPositionMsg(1)
//	effect.SetEffect(0, effectsEmitterID, u.GetTransform(), target.GetTransform(), 1001)
//	u.GetWorld().Send(effect.WSMessage(), false)
//
//	go u.SendHighFiveStats(target)
//
//	return nil
//}
