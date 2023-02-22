package object

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common/unity"
)

func (o *Object) unityAutoOnObjectAttributeChanged(
	changeType universe.AttributeChangeType,
	attributeID entry.AttributeID,
	value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
	o.log.Infof("attribute Unuty Auto processing for %+v %+v", o.GetID(), attributeID)
	autoOption, err := unity.GetOptionAutoOption(attributeID, effectiveOptions)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto option: %+v", attributeID)
	}
	if autoOption == nil {
		return nil
	}

	o.log.Infof("unity-auto stage3 for %+v %+v", o.GetID(), attributeID)

	hash, err := unity.PrerenderAutoValue(o.ctx, autoOption, value)
	if err != nil {
		return errors.WithMessagef(err, "prerendering error: %+v", attributeID)
	}

	o.log.Infof("unity-auto stage4 for %+v %+v %+v", o.GetID(), attributeID, hash)

	//dirty hack to set auto_render_hash value without triggering processing again
	// TODO: fix it properly later
	if hash != nil && hash.Hash != "" {
		func() error {
			o.objectAttributes.object.Mu.Lock()
			defer o.objectAttributes.object.Mu.Unlock()

			(*value)["auto_render_hash"] = hash.Hash

			if err := o.db.GetObjectAttributesDB().UpdateObjectAttributeValue(
				o.ctx, entry.NewObjectAttributeID(attributeID, o.GetID()), value,
			); err != nil {
				return errors.WithMessage(err, "failed to update db")
			}

			return nil
		}()
	}
	o.SendUnityAutoAttributeMessage(
		autoOption, value, func(m *websocket.PreparedMessage) error { return o.GetWorld().Send(m, false) },
	)
	return nil
}

func (o *Object) SendUnityAutoAttributeMessage(
	option *entry.UnityAutoAttributeOption,
	value *entry.AttributeValue,
	send func(*websocket.PreparedMessage) error,
) {
	msg := o.UpdateAutoTextureMap(option, value)
	if msg != nil {
		send(msg)
	}
	return
}

func (o *Object) UpdateAutoTextureMap(
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
		msg = message.GetBuilder().SetObjectAttributes(o.GetID(), sendMap)
	case entry.UnitySlotTypeString:
		o.log.Infof("Processing String Slot %+v for %+v \n", option.SlotName, o.GetID())
		v, ok := (*value)[option.ValueField]
		if !ok {
			return nil
		}
		val, ok := v.(string)
		if !ok {
			return nil
		}

		o.log.Infof("Setting String Slot %+v for %+v to  %+v  \n", option.SlotName, o.GetID(), val)

		o.renderStringMap.Store(option.SlotName, val)
		func() {
			o.renderStringMap.Mu.RLock()
			defer o.renderStringMap.Mu.RUnlock()
			o.stringMsg.Store(message.GetBuilder().SetObjectStrings(o.GetID(), o.renderStringMap.Data))
		}()

		sendMap := map[string]string{option.SlotName: val}
		msg = message.GetBuilder().SetObjectStrings(o.GetID(), sendMap)
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

		o.renderTextureMap.Store(option.SlotName, val)
		func() {
			o.renderTextureMap.Mu.RLock()
			defer o.renderTextureMap.Mu.RUnlock()

			o.textureMsg.Store(message.GetBuilder().SetObjectTextures(o.GetID(), o.renderTextureMap.Data))
			ot := o.GetObjectType()
			if ot != nil && option.SlotName == "skybox_custom" && ot.GetAsset3d() != nil && ot.GetAsset3d().GetID() == uuid.MustParse("313a597a-8b9a-47a7-9908-52bdc7a21a3e") {
				o.log.Infof(
					"unity-auto stage7b :Setting skybox texture for %+v to %+v | %+v from %+v\n", o.world.GetID(), val,
					option,
					o.GetID(),
				)
				skyBoxTextureMap := map[string]string{option.SlotName: val}
				o.world.TempSetSkybox(message.GetBuilder().SetObjectTextures(o.GetID(), skyBoxTextureMap))
			}
		}()

		sendMap := map[string]string{option.SlotName: val}
		msg = message.GetBuilder().SetObjectTextures(o.GetID(), sendMap)

	}
	return msg
}
