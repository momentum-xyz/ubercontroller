package user

import (
	"fmt"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (u *User) OnMessage(msg *posbus.Message) error {
	switch msg.Type() {
	case posbus.MsgTypeSendPosition:
		if err := u.UpdatePosition(msg.AsSendPos()); err != nil {
			return errors.WithMessage(err, "failed to handle: send position")
		}
		return nil
	case posbus.MsgTypeFlatBufferMessage:
		switch msg.AsFlatBufferMessage().MsgType() {
		default:
			return errors.Errorf(
				"unknown flatbuffer message: %d", msg.AsFlatBufferMessage().MsgType(),
			)
		}
	case posbus.MsgTriggerInteraction:
		if err := u.InteractionHandler(msg.AsTriggerInteraction()); err != nil {
			return errors.WithMessage(err, "failed to handle: interaction")
		}
		return nil
	case posbus.MsgTypeRelayToController:
		if err := u.RelayToControllerHandler(msg.AsRelayToController()); err != nil {
			return errors.WithMessage(err, "failed to handle: relay to controller")
		}
		return nil
	case posbus.MsgTypeSwitchWorld:
		if err := u.Teleport(msg.AsSwitchWorld()); err != nil {
			return errors.WithMessage(err, "failed to handle: teleport")
		}
		return nil
	case posbus.MsgTypeSignal:
		if err := u.SignalsHandler(msg.AsSignal().Signal()); err != nil {
			return errors.WithMessage(err, "failed to handle: signal")
		}
		return nil
	case posbus.MsgTypeSetStaticObjectPosition:
		if err := u.UpdateSpacePosition(msg.AsSetStaticObjectPosition()); err != nil {
			return errors.WithMessage(err, "failed to update space position")
		}
		return nil
	case posbus.MsgTypeSetObjectLockState:
		if err := u.LockObject(msg.AsSetObjectLockState()); err != nil {
			return errors.WithMessage(err, "failed to handle: signal")
		}
		return nil
	}

	return errors.Errorf("unknown message: %d", msg.Type())
}

func (u *User) UpdateSpacePosition(msg *posbus.SetStaticObjectPosition) error {
	space, ok := universe.GetNode().GetSpaceFromAllSpaces(msg.ObjectID())
	if !ok {
		return errors.Errorf("space not found: %s", msg.ObjectID())
	}

	if err := space.SetPosition(utils.GetPTR(msg.Position()), true); err != nil {
		return errors.WithMessage(err, "failed to set space position")
	}

	return space.GetWorld().Send(msg.WebsocketMessage(), true)
}

func (u *User) Teleport(msg *posbus.SwitchWorld) error {
	// TODO: teleport function
	//if err := u.SwitchWorld(msg.AsSwitchWorld().World()); err != nil {
	//	u.log.Error(errors.WithMessage(err, "User: OnMessage: failed to switch world"))
	//}
	return nil
}

func (u *User) RelayToControllerHandler(m *posbus.RelayToController) error {
	if m.Topic() == "emoji" {
		// TODO: comes as plugin?
		//u.HandleEmoji(msg.AsRelayToController())
	}
	return nil
}

func (u *User) SignalsHandler(s posbus.Signal) error {
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

	return nil
}

func (u *User) InteractionHandler(m *posbus.TriggerInteraction) error {
	kind := m.Kind()
	targetUUID := m.Target()
	flag := m.Flag()
	label := m.Label()
	u.log.Info(
		"Incoming interaction for user", u.id, "kind:", kind, "target:", targetUUID, "flag:", flag, "label:", label,
	)
	switch kind {
	case posbus.TriggerEnteredSpace:
		space, ok := universe.GetNode().GetSpaceFromAllSpaces(targetUUID)
		if !ok {
			return errors.WithMessage(
				errors.Errorf("space not found: %s", targetUUID), "failed to handle: enter space",
			)
		}
		if err := space.AddUser(u, true); err != nil {
			return errors.WithMessage(
				errors.Errorf("failed to add user to space: %s", targetUUID), "failed to handle: enter space",
			)
		}
		return nil
	case posbus.TriggerLeftSpace:
		space, ok := universe.GetNode().GetSpaceFromAllSpaces(targetUUID)
		if !ok {
			return errors.WithMessage(
				errors.Errorf("space not found: %s", targetUUID), "failed to handle: left space",
			)
		}
		if err := space.RemoveUser(u, true); err != nil {
			return errors.WithMessage(
				errors.Errorf("failed to remove user from space: %s", targetUUID), "failed to handle: left space",
			)
		}
		return nil
	}
	//case posbus.TriggerHighFive:
	//	if err := u.HandleHighFive(m); err != nil {
	//		u.log.Warn(errors.WithMessage(err, "InteractionHandler: trigger high fives: failed to handle high five"))
	//	}
	//case posbus.TriggerStake:
	//	u.HandleStake(m)

	return errors.Errorf("unknown message: %d", kind)
}

func (u *User) LockObject(m *posbus.SetObjectLockState) error {
	id := m.ObjectID()
	state := m.State()
	if space, ok := u.GetWorld().GetSpace(id, true); ok {
		result := space.LockUnityObject(u, state)
		newState := state
		if !result {
			newState = 1 - state
		}
		m.SetLockState(id, newState)
		u.Send(m.WebsocketMessage())
	}
	return nil
}
