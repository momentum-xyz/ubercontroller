package object

import (
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common/slot"
)

func (o *Object) renderAutoOnObjectAttributeChanged(
	changeType posbus.AttributeChangeType,
	attributeID entry.AttributeID,
	value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
	o.log.Infof("attribute render auto processing for %+v %+v", o.GetID(), attributeID)
	autoOption, err := slot.GetOptionAutoOption(attributeID, effectiveOptions)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto option: %+v", attributeID)
	}
	if autoOption == nil {
		return nil
	}

	o.log.Infof("render-auto stage3 for %+v %+v", o.GetID(), attributeID)

	hash, err := slot.PrerenderAutoValue(o.ctx, autoOption, value)
	if err != nil {
		return errors.WithMessagef(err, "prerendering error: %+v", attributeID)
	}

	o.log.Infof("render-auto stage4 for %+v %+v %+v", o.GetID(), attributeID, hash)

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
	o.SendRenderAutoAttributeMessage(
		autoOption, value, func(m *websocket.PreparedMessage) error { return o.GetWorld().Send(m, false) },
	)
	return nil
}

func (o *Object) SendRenderAutoAttributeMessage(
	option *entry.RenderAutoAttributeOption,
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
	option *entry.RenderAutoAttributeOption, value *entry.AttributeValue,
) *websocket.PreparedMessage {
	if option == nil || value == nil {
		return nil
	}

	var data interface{}
	switch option.SlotType {
	case entry.SlotTypeNumber:
		v, ok := (*value)[option.ValueField]
		if !ok {
			return nil
		}
		val, ok := v.(int)
		if !ok {
			return nil
		}
		data = int32(val)
	case entry.SlotTypeString:
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

		data = val
	case entry.SlotTypeTexture:
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
		data = val
	case entry.SlotTypeAudio:
		// TODO: really need a more flexible render_auto system,
		// audio requires multiple fields, so for now, pass through the whole value:
		data = map[string]any(*value)
	}

	// store to the full list and update message (one which send on spawn)
	func() {
		o.renderDataMap.Mu.RLock()
		defer o.renderDataMap.Mu.RUnlock()
		slots, ok := o.renderDataMap.Data[option.SlotType]
		if !ok {
			slots = utils.GetPTR(make(posbus.StringAnyMap))
			o.renderDataMap.Data[option.SlotType] = slots
		}
		(*slots)[option.SlotName] = data
		msg := posbus.WSMessage(&posbus.ObjectData{ID: o.GetID(), Entries: o.renderDataMap.Data})
		o.dataMsg.Store(msg)
	}()

	// prepare message for this atomic update
	//FIXME

	msg := posbus.ObjectData{ID: o.GetID(), Entries: map[entry.SlotType]*posbus.StringAnyMap{option.SlotType: utils.GetPTR(posbus.StringAnyMap{option.SlotName: data})}}
	return posbus.WSMessage(&msg)
}
