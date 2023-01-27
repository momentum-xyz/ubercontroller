package user

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func (u *User) OnMessage(msg *posbus.Message) error {
	switch msg.Type() {
	case posbus.MsgTypeSendPosition:
		if err := u.UpdatePosition(msg.AsSendPos()); err != nil {
			return errors.WithMessage(err, "failed to handle: send position")
		}
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
	case posbus.MsgTypeRelayToController:
		if err := u.RelayToControllerHandler(msg.AsRelayToController()); err != nil {
			return errors.WithMessage(err, "failed to handle: relay to controller")
		}
	case posbus.MsgTypeSwitchWorld:
		if err := u.Teleport(msg.AsSwitchWorld()); err != nil {
			return errors.WithMessage(err, "failed to handle: teleport")
		}
	case posbus.MsgTypeSignal:
		if err := u.SignalsHandler(msg.AsSignal().Signal()); err != nil {
			return errors.WithMessage(err, "failed to handle: signal")
		}
	case posbus.MsgTypeSetStaticObjectPosition:
		if err := u.UpdateSpacePosition(msg.AsSetStaticObjectPosition()); err != nil {
			return errors.WithMessage(err, "failed to update space position")
		}
	case posbus.MsgTypeSetObjectLockState:
		if err := u.LockObject(msg.AsSetObjectLockState()); err != nil {
			return errors.WithMessage(err, "failed to set object lock state")
		}
	default:
		return errors.Errorf("unknown message: %d", msg.Type())
	}

	return nil
}

func (u *User) UpdateSpacePosition(msg *posbus.SetStaticObjectPosition) error {
	space, ok := universe.GetNode().GetSpaceFromAllSpaces(msg.ObjectID())
	if !ok {
		return errors.Errorf("space not found: %s", msg.ObjectID())
	}
	return space.SetPosition(utils.GetPTR(msg.Position()), true)
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
		u.log.Infof("SKYBOX: Got SignalReady from %+v\n", u.GetID().String())
		sm := u.world.TempGetSkybox()
		if sm != nil {
			u.log.Infof("SKYBOX: Sending texture to  %+v\n", u.GetID().String())
			u.SendDirectly(sm)
		}

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
	u.log.Infof(
		"Incoming interaction for user: %s, kind: %d, target: %s, flag: %d, label: %s",
		u.GetID(), kind, targetUUID, flag, label,
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
	case posbus.TriggerHighFive:
		if err := u.HandleHighFive(m); err != nil {
			u.log.Warn(errors.WithMessage(err, "InteractionHandler: trigger high fives: failed to handle high five"))
		}
		return nil
	}
	// case posbus.TriggerStake:
	// 	u.HandleStake(m)

	return errors.Errorf("unknown message: %d", kind)
}

func (u *User) LockObject(msg *posbus.SetObjectLockState) error {
	id := msg.ObjectID()
	state := msg.State()

	space, ok := u.GetWorld().GetSpaceFromAllSpaces(id)
	if !ok {
		return errors.Errorf("space not found: %s", id)
	}

	result := space.LockUnityObject(u, state)
	newState := state
	if !result {
		newState = 1 - state
	}

	msg.SetLockState(id, newState)

	return u.GetWorld().Send(msg.WebsocketMessage(), true)
}

func (u *User) HandleHighFive(m *posbus.TriggerInteraction) error {
	targetID := m.Target()
	if targetID == u.GetID() {
		return errors.New("high-five yourself not permitted")
	}

	world := u.GetWorld()
	target, ok := world.GetUser(targetID, false)
	if !ok {
		u.Send(
			posbus.NewSimpleNotificationMsg(
				posbus.DestinationReact, posbus.NotificationTextMessage, 0, "Target user not found",
			).WebsocketMessage(),
		)
		return errors.Errorf("failed to get target: %s", targetID)
	}

	var uName string
	var tName string
	uProfile := u.GetProfile()
	tProfile := target.GetProfile()
	if uProfile != nil && uProfile.Name != nil {
		uName = *uProfile.Name
	}
	if tProfile != nil && tProfile.Name != nil {
		tName = *tProfile.Name
	}

	high5Msg := struct {
		SenderID   string `json:"senderId"`
		ReceiverID string `json:"receiverId"`
		Message    string `json:"message"`
	}{
		SenderID:   u.GetID().String(),
		ReceiverID: targetID.String(),
		Message:    uName,
	}
	high5Data, err := json.Marshal(&high5Msg)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal high5 message")
	}

	u.Send(
		posbus.NewSimpleNotificationMsg(
			posbus.DestinationReact, posbus.NotificationHighFive, 0, tName,
		).WebsocketMessage(),
	)
	target.Send(posbus.NewRelayToReactMsg("high5", high5Data).WebsocketMessage())

	effectsEmitterID := world.GetSettings().Spaces["effects_emitter"]
	effect := posbus.NewTriggerTransitionalBridgingEffectsOnPositionMsg(1)
	effect.SetEffect(0, effectsEmitterID, u.GetPosition(), target.GetPosition(), 1001)
	u.GetWorld().Send(effect.WebsocketMessage(), false)

	go u.SendHighFiveStats(target)

	return nil
}
