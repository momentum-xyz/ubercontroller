package user

import (
	"fmt"

	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
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
	case posbus.TriggerEnteredSpace:
		spaceID := m.Target()
		space, ok := universe.GetNode().GetSpaceFromAllSpaces(spaceID)
		if !ok {
			u.log.Errorf("User: InteractionHandler: TriggerEnteredSpace: space not found: %s", spaceID)
			return
		}
		if err := space.AddUser(u, true); err != nil {
			u.log.Error(
				errors.WithMessagef(
					err, "User: InteractionHandler: TriggerEnteredSpace: failed to add user to space: %s", spaceID,
				),
			)
			return
		}
	case posbus.TriggerLeftSpace:
		spaceID := m.Target()
		space, ok := universe.GetNode().GetSpaceFromAllSpaces(spaceID)
		if !ok {
			u.log.Errorf("User: InteractionHandler: TriggerLeftSpace: space not found: %s", spaceID)
			return
		}
		if err := space.RemoveUser(u, true); err != nil {
			u.log.Error(
				errors.WithMessagef(
					err, "User: InteractionHandler: TriggerEnteredSpace: failed to remove user from space: %s", spaceID,
				),
			)
			return
		}
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
