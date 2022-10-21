package user

import (
	"fmt"
	"github.com/momentum-xyz/posbus-protocol/posbus"
)

func (u *User) OnMessage(msg *posbus.Message) {
	fmt.Printf("%+v\n", msg.Type())
	switch msg.Type() {
	case posbus.MsgTypeSendPosition:
		u.UpdatePosition(msg.AsSendPos())
	case posbus.MsgTypeFlatBufferMessage:
		switch msg.AsFlatBufferMessage().MsgType() {
		default:
			u.log.Warn("Got unknown Flatbuffer message for user:", u.id, "msg:", msg.AsFlatBufferMessage().MsgType())
			return
		}
	case posbus.MsgTriggerInteraction:
		u.InteractionHandler(msg.AsTriggerInteraction())
	case posbus.MsgTypeRelayToController:
		u.RelayToControllerHandler(msg.AsRelayToController())
	case posbus.MsgTypeSwitchWorld:
		u.Teleport(msg.AsSwitchWorld())
	case posbus.MsgTypeSignal:
		u.SignalsHandler(msg.AsSignal().Signal())
	default:
		u.log.Warn("Got unknown message for user:", u.id, "msg:", msg)
	}
}

func (u *User) Teleport(msg *posbus.SwitchWorld) {
	// TODO: teleport function
	//if err := u.SwitchWorld(msg.AsSwitchWorld().World()); err != nil {
	//	u.log.Error(errors.WithMessage(err, "User: OnMessage: failed to switch world"))
	//}
}

func (u *User) RelayToControllerHandler(m *posbus.RelayToController) {
	if m.Topic() == "emoji" {
		// TODO: comes as plugin?
		//u.HandleEmoji(msg.AsRelayToController())
	}
}

func (u *User) SignalsHandler(s posbus.Signal) {
	fmt.Printf("Got Signal %+v\n", s)
	switch s {
	case posbus.SignalReady:

		u.ReleaseSendBuffer()
		//u.log.Debugf("Got signalReady from %s", u.id.String())
		//TODO: Do we need it?
		//if err := u.world.SendWorldData(u); err != nil {
		//	log.Error(
		//		errors.WithMessagef(
		//			err, "User: SignalsHandler: SignalReady: failed to send world data: %s", u.ID,
		//		),
		//	)
		//	u.world.unregisterUser <- u
		//	return
		//}
		//u.connection.EnableWriting()
	}
}

func (u *User) InteractionHandler(m *posbus.TriggerInteraction) {
	kind := m.Kind()
	targetUUID := m.Target()
	flag := m.Flag()
	label := m.Label()
	u.log.Info(
		"Incoming interaction for user", u.id, "kind:", kind, "target:", targetUUID, "flag:", flag, "label:", label,
	)
	switch kind {
	//case posbus.TriggerEnteredSpace:
	//	targetUUID := m.Target()
	//	u.currentSpace.Store(targetUUID)
	//	if err := u.world.hub.DB.InsertOnline(u.ID, targetUUID); err != nil {
	//		u.log.Warn(errors.WithMessage(err, "InteractionHandler: trigger entered space: failed to insert one"))
	//	}
	//	if _, err := u.world.UpdateOnlineBySpaceId(targetUUID); err != nil {
	//		u.log.Warn(
	//			errors.WithMessage(
	//				err, "InteractionHandler: trigger entered space: failed to update online by space id",
	//			),
	//		)
	//	}
	//	u.world.hub.CancelCleanupSpace(targetUUID)
	//	go func() {
	//		if err := u.sendSpaceEnterLeaveStats(targetUUID, 1); err != nil {
	//			u.log.Warn(
	//				errors.WithMessagef(
	//					err, "InteractionHandler: trigger entered space: failed to update users on space: %s",
	//					targetUUID,
	//				),
	//			)
	//		}
	//	}()
	//case posbus.TriggerLeftSpace:
	//	targetUUID := m.Target()
	//	u.currentSpace.Store(uuid.Nil)
	//	if err := u.world.hub.DB.RemoveOnline(u.id, targetUUID); err != nil {
	//		u.log.Warn(
	//			errors.WithMessage(
	//				err, "InteractionHandler: trigger left space: failed to remove online from db",
	//			),
	//		)
	//	}
	//	if _, err := u.world.UpdateOnlineBySpaceId(targetUUID); err != nil {
	//		u.log.Warn(
	//			errors.WithMessagef(
	//				err, "InteractionHandler: trigger left space: failed to update online by space id",
	//			),
	//		)
	//	}
	//	if ok, err := u.world.hub.SpaceStorage.CheckOnlineSpaceByID(targetUUID); err != nil {
	//		u.log.Warn(
	//			errors.WithMessagef(
	//				err, "InteractionHandler: trigger left space: failed to check online space by id",
	//			),
	//		)
	//	} else if !ok {
	//		if err := u.world.hub.CleanupSpace(targetUUID); err != nil {
	//			u.log.Warn(errors.WithMessage(err, "InteractionHandler: trigger left space: failed to cleanup space"))
	//		}
	//	}
	//	go func() {
	//		if err := u.sendSpaceEnterLeaveStats(targetUUID, 0); err != nil {
	//			u.log.Warn(
	//				errors.WithMessagef(
	//					err, "InteractionHandler: trigger left space: failed to update users on space: %s", targetUUID,
	//				),
	//			)
	//		}
	//	}()
	//case posbus.TriggerHighFive:
	//	if err := u.HandleHighFive(m); err != nil {
	//		u.log.Warn(errors.WithMessage(err, "InteractionHandler: trigger high fives: failed to handle high five"))
	//	}
	//case posbus.TriggerStake:
	//	u.HandleStake(m)
	default:
		u.log.Warn("InteractionHandler: got unknown interaction for user:", u.id, "kind:", kind)
	}
}
