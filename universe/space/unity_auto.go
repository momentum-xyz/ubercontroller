package space

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/unity"
)

func (s *Space) unityAutoOnSpaceAttributeChanged(
	changeType universe.AttributeChangeType,
	attributeID entry.AttributeID,
	value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
	s.log.Infof("attribute Unuty Auto processing for %+v %+v", s.GetID(), attributeID)
	autoOption, err := unity.GetOptionAutoOption(attributeID, effectiveOptions)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto option: %+v", attributeID)
	}
	if autoOption == nil {
		return nil
	}

	s.log.Infof("unity-auto stage3 for %+v %+v", s.GetID(), attributeID)

	hash, err := unity.PrerenderAutoValue(s.ctx, autoOption, value)
	if err != nil {
		return errors.WithMessagef(err, "prerendering error: %+v", attributeID)
	}

	s.log.Infof("unity-auto stage4 for %+v %+v %+v", s.GetID(), attributeID, hash)

	//dirty hack to set auto_render_hash value without triggering processing again
	// TODO: fix it properly later
	if hash != nil && hash.Hash != "" {
		return func() error {
			s.spaceAttributes.space.Mu.Lock()
			defer s.spaceAttributes.space.Mu.Unlock()

			(*value)["auto_render_hash"] = hash.Hash

			if err := s.db.GetSpaceAttributesDB().UpdateSpaceAttributeValue(
				s.ctx, entry.NewSpaceAttributeID(attributeID, s.GetID()), value,
			); err != nil {
				return errors.WithMessage(err, "failed to update db")
			}

			return nil
		}()
	}
	s.SendUnityAutoAttributeMessage(
		autoOption, value, func(m *websocket.PreparedMessage) error { return s.GetWorld().Send(m, false) },
	)
	return nil
}

func (s *Space) SendUnityAutoAttributeMessage(
	option *entry.UnityAutoAttributeOption,
	value *entry.AttributeValue,
	send func(*websocket.PreparedMessage) error,
) {
	msg := s.UpdateAutoTextureMap(option, value)
	if msg != nil {
		send(msg)
	}
	return
}

func (s *Space) UpdateAutoTextureMap(
	option *entry.UnityAutoAttributeOption, value *entry.AttributeValue,
) *websocket.PreparedMessage {
	if option == nil || value == nil {
		return nil
	}

	var msg *websocket.PreparedMessage
	switch option.SlotType {
	case entry.UnitySlotTypeNumber:
		v, ok := (*value)[option.ValueField]
		if !ok {
			return nil
		}
		val, ok := v.(int)
		if !ok {
			return nil
		}
		sendMap := map[string]int32{option.SlotName: int32(val)}
		msg = message.GetBuilder().SetObjectAttributes(s.GetID(), sendMap)
	case entry.UnitySlotTypeString:
		v, ok := (*value)[option.ValueField]
		if !ok {
			return nil
		}
		val, ok := v.(string)
		if !ok {
			return nil
		}

		sendMap := map[string]string{option.SlotName: val}
		msg = message.GetBuilder().SetObjectStrings(s.GetID(), sendMap)
	case entry.UnitySlotTypeTexture:
		valField := "auto_render_hash"
		if option.ContentType == "image" {
			valField = "render_hash"
		}
		v, ok := (*value)[valField]
		if !ok {
			return nil
		}
		val, ok := v.(string)
		if !ok {
			return nil
		}

		s.renderTextureMap.Store(option.SlotName, val)
		func() {
			s.renderTextureMap.Mu.RLock()
			defer s.renderTextureMap.Mu.RUnlock()

			s.textMsg.Store(message.GetBuilder().SetObjectTextures(s.GetID(), s.renderTextureMap.Data))
		}()

		sendMap := map[string]string{option.SlotName: val}
		msg = message.GetBuilder().SetObjectTextures(s.GetID(), sendMap)

	}
	return msg
}
