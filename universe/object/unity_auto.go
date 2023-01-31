package object

import (
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
		return func() error {
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
		v, ok := (*value)[option.ValueField]
		if !ok {
			return nil
		}
		val, ok := v.(string)
		if !ok {
			return nil
		}

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

			o.textMsg.Store(message.GetBuilder().SetObjectTextures(o.GetID(), o.renderTextureMap.Data))
		}()

		sendMap := map[string]string{option.SlotName: val}
		msg = message.GetBuilder().SetObjectTextures(o.GetID(), sendMap)

	}
	return msg
}
