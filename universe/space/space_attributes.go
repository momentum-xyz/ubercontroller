package space

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/posbus"
	"github.com/momentum-xyz/ubercontroller/universe/common/unity"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (s *Space) GetSpaceAttributeValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	payload, ok := s.GetSpaceAttributePayload(attributeID)
	if !ok {
		return nil, false
	}
	if payload == nil {
		return nil, true
	}
	return payload.Value, true
}

func (s *Space) GetSpaceAttributeOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	payload, ok := s.GetSpaceAttributePayload(attributeID)
	if !ok {
		return nil, false
	}
	if payload == nil {
		return nil, true
	}
	return payload.Options, true
}

func (s *Space) GetSpaceAttributeEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attributeType, ok := universe.GetNode().GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(attributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := s.GetSpaceAttributeOptions(attributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		s.log.Error(
			errors.WithMessagef(
				err, "Space: GetSpaceAttributeEffectiveOptions: failed to merge options: %+v", attributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (s *Space) GetSpaceAttributePayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	if s.id.String() == "0e347f9e-bba9-48c1-a3f9-4258ca230481" {
		fmt.Println("***> GetSpaceAttributePayload", attributeID.PluginID.String(), attributeID.Name)

		for k, v := range s.spaceAttributes.Data {
			if k.Name == "news_feed" {
				fmt.Println("***> ", k.PluginID.String(), k.Name, "...", v.Options)
			} else {
				fmt.Println("***> ", k.PluginID.String(), k.Name, v.Value, v.Options)
			}
		}
	}
	return s.spaceAttributes.Load(attributeID)
}

func (s *Space) GetSpaceAttributesValue(recursive bool) map[entry.SpaceAttributeID]*entry.AttributeValue {
	s.spaceAttributes.Mu.RLock()
	values := make(map[entry.SpaceAttributeID]*entry.AttributeValue, len(s.spaceAttributes.Data))
	for attributeID, payload := range s.spaceAttributes.Data {
		spaceAttributeID := entry.NewSpaceAttributeID(attributeID, s.GetID())
		if payload == nil {
			values[spaceAttributeID] = nil
			continue
		}
		values[spaceAttributeID] = payload.Value
	}
	s.spaceAttributes.Mu.RUnlock()

	if !recursive {
		return values
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		for spaceAttributeID, value := range child.GetSpaceAttributesValue(recursive) {
			values[spaceAttributeID] = value
		}
	}

	return values
}

func (s *Space) GetSpaceAttributesOptions(recursive bool) map[entry.SpaceAttributeID]*entry.AttributeOptions {
	s.spaceAttributes.Mu.RLock()
	options := make(map[entry.SpaceAttributeID]*entry.AttributeOptions, len(s.spaceAttributes.Data))
	for attributeID, payload := range s.spaceAttributes.Data {
		spaceAttributeID := entry.NewSpaceAttributeID(attributeID, s.GetID())
		if payload == nil {
			options[spaceAttributeID] = nil
			continue
		}
		options[spaceAttributeID] = payload.Options
	}
	s.spaceAttributes.Mu.RUnlock()

	if !recursive {
		return options
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		for spaceAttributeID, opt := range child.GetSpaceAttributesOptions(recursive) {
			options[spaceAttributeID] = opt
		}
	}

	return options
}

func (s *Space) GetSpaceAttributesPayload(recursive bool) map[entry.SpaceAttributeID]*entry.AttributePayload {
	s.spaceAttributes.Mu.RLock()
	payloads := make(map[entry.SpaceAttributeID]*entry.AttributePayload, len(s.spaceAttributes.Data))
	for attributeID, payload := range s.spaceAttributes.Data {
		payloads[entry.NewSpaceAttributeID(attributeID, s.GetID())] = payload
	}
	s.spaceAttributes.Mu.RUnlock()

	if !recursive {
		return payloads
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		for spaceAttributeID, payload := range child.GetSpaceAttributesPayload(recursive) {
			payloads[spaceAttributeID] = payload
		}
	}

	return payloads
}

func (s *Space) UpsertSpaceAttribute(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.SpaceAttribute, error) {
	if s.id.String() == "0e347f9e-bba9-48c1-a3f9-4258ca230481" {
		fmt.Println("***>> UpsertSpaceAttribute call for:", attributeID.PluginID.String(), attributeID.Name)
	}
	s.spaceAttributes.Mu.Lock()
	defer s.spaceAttributes.Mu.Unlock()

	payload, ok := s.spaceAttributes.Data[attributeID]
	if !ok {
		payload = nil
	}

	payload, err := modifyFn(payload)
	if err != nil {
		if s.id.String() == "0e347f9e-bba9-48c1-a3f9-4258ca230481" {
			fmt.Println("***>> UpsertSpaceAttribute error: failed to modify attribute payload")
		}
		return nil, errors.WithMessage(err, "failed to modify attribute payload")
	}

	spaceAttribute := entry.NewSpaceAttribute(entry.NewSpaceAttributeID(attributeID, s.GetID()), payload)
	if updateDB {
		if err := s.db.SpaceAttributesUpsertSpaceAttribute(s.ctx, spaceAttribute); err != nil {
			return nil, errors.WithMessage(err, "failed to upsert space attribute")
		}
	}

	s.spaceAttributes.Data[attributeID] = payload
	if s.id.String() == "0e347f9e-bba9-48c1-a3f9-4258ca230481" {
		if attributeID.Name == "news_feed" {
			fmt.Println("***>> UpsertSpaceAttribute ", attributeID.PluginID.String(), attributeID.Name, "...")
		} else {
			fmt.Println("***>> UpsertSpaceAttribute ", attributeID.PluginID.String(), attributeID.Name, payload)
		}
	}

	var value *entry.AttributeValue
	if payload != nil {
		value = payload.Value
	}
	if s.GetEnabled() {
		go s.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return spaceAttribute, nil
}

func (s *Space) UpdateSpaceAttributeValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	s.spaceAttributes.Mu.Lock()
	defer s.spaceAttributes.Mu.Unlock()

	payload, ok := s.spaceAttributes.Data[attributeID]
	if !ok {
		return nil, errors.Errorf("space attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify value")
	}

	if updateDB {
		if err := s.db.SpaceAttributesUpdateSpaceAttributeValue(
			s.ctx, entry.NewSpaceAttributeID(attributeID, s.GetID()), value,
		); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Value = value
	s.spaceAttributes.Data[attributeID] = payload

	go s.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)

	return value, nil
}

func (s *Space) UpdateSpaceAttributeOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	s.spaceAttributes.Mu.Lock()
	defer s.spaceAttributes.Mu.Unlock()

	payload, ok := s.spaceAttributes.Data[attributeID]
	if !ok {
		return nil, errors.Errorf("space attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := s.db.SpaceAttributesUpdateSpaceAttributeOptions(
			s.ctx, entry.NewSpaceAttributeID(attributeID, s.GetID()), options,
		); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Options = options
	s.spaceAttributes.Data[attributeID] = payload

	var value *entry.AttributeValue
	if payload != nil {
		value = payload.Value
	}
	if s.GetEnabled() {
		go s.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return options, nil
}

func (s *Space) RemoveSpaceAttribute(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	if _, ok := s.spaceAttributes.Load(attributeID); !ok {
		return false, nil
	}

	attributeEffectiveOptions, attributeEffectiveOptionsOK := s.GetSpaceAttributeEffectiveOptions(attributeID)

	s.spaceAttributes.Mu.Lock()
	defer s.spaceAttributes.Mu.Unlock()

	if updateDB {
		if err := s.db.SpaceAttributesRemoveSpaceAttributeByID(
			s.ctx, entry.NewSpaceAttributeID(attributeID, s.GetID()),
		); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(s.spaceAttributes.Data, attributeID)

	if s.GetEnabled() {
		go func() {
			if !attributeEffectiveOptionsOK {
				s.log.Error(
					errors.Errorf(
						"Space: RemoveSpaceAttribute: failed to get space attribute effective options",
					),
				)
				return
			}
			s.onSpaceAttributeChanged(universe.RemovedAttributeChangeType, attributeID, nil, attributeEffectiveOptions)
		}()
	}

	return true, nil
}

// TODO: optimize
func (s *Space) RemoveSpaceAttributes(attributeIDs []entry.AttributeID, updateDB bool) (bool, error) {
	res := true

	var errs *multierror.Error
	for i := range attributeIDs {
		removed, err := s.RemoveSpaceAttribute(attributeIDs[i], updateDB)
		if err != nil {
			errs = multierror.Append(
				errs,
				errors.WithMessagef(err, "failed to remove space attribute: %+v", attributeIDs[i]),
			)
		}
		if !removed {
			res = false
		}
	}

	return res, errs.ErrorOrNil()
}

func (s *Space) loadSpaceAttributes() error {
	entries, err := s.db.SpaceAttributesGetSpaceAttributesBySpaceID(s.ctx, s.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get space attributes")
	}

	for i := range entries {
		entry := entries[i]

		if _, err := s.UpsertSpaceAttribute(
			entry.AttributeID, modify.MergeWith(entry.AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert space attribute: %+v", entry.AttributeID)
		}

		effectiveOptions, ok := s.GetSpaceAttributeEffectiveOptions(entry.AttributeID)
		if !ok {
			continue
		}
		autoOption, err := unity.GetOptionAutoOption(entry.AttributeID, effectiveOptions)
		if err != nil {
			return errors.WithMessagef(err, "failed to get option auto option: %+v", entry)
		}
		s.UpdateAutoTextureMap(autoOption, entry.Value)
	}

	s.log.Debugf("Space attributes loaded: %s: %d", s.GetID(), s.spaceAttributes.Len())

	return nil
}

func (s *Space) onSpaceAttributeChanged(
	changeType universe.AttributeChangeType, attributeID entry.AttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) {
	go s.calendarOnSpaceAttributeChanged(changeType, attributeID, value, effectiveOptions)

	if effectiveOptions == nil {
		options, ok := s.GetSpaceAttributeEffectiveOptions(attributeID)
		if !ok {
			s.log.Error(
				errors.Errorf(
					"Space: onSpaceAttributeChanged: failed to get attribute effective options: %+v", attributeID,
				),
			)
			return
		}
		effectiveOptions = options
	}

	go func() {
		if err := s.posBusAutoOnSpaceAttributeChanged(changeType, attributeID, value, effectiveOptions); err != nil {
			s.log.Error(
				errors.WithMessagef(
					err, "Space: onSpaceAttributeChanged: failed to handle pos bus auto: %+v", attributeID,
				),
			)
		}
	}()

	go func() {
		if err := s.unityAutoOnSpaceAttributeChanged(changeType, attributeID, value, effectiveOptions); err != nil {
			s.log.Error(
				errors.WithMessagef(
					err, "Space: onSpaceAttributeChanged: failed to handle unity auto: %+v", attributeID,
				),
			)
		}
	}()
}

func (s *Space) posBusAutoOnSpaceAttributeChanged(
	changeType universe.AttributeChangeType, attributeID entry.AttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
	autoOption, err := posbus.GetOptionAutoOption(effectiveOptions)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto option: %+v", attributeID)
	}
	autoMessage, err := posbus.GetOptionAutoMessage(autoOption, changeType, attributeID, value)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto message: %+v", attributeID)
	}
	if autoMessage == nil {
		return nil
	}

	var errs *multierror.Error
	for i := range autoOption.Scope {
		switch autoOption.Scope[i] {
		case entry.WorldPosBusAutoScopeAttributeOption:
			world := s.GetWorld()
			if world == nil {
				errs = multierror.Append(
					err, errors.Errorf("failed to get world: %s", autoOption.Scope[i]),
				)
				continue
			}
			if err := world.Send(autoMessage, true); err != nil {
				errs = multierror.Append(
					errs, errors.WithMessagef(
						err, "failed to send message: %s", autoOption.Scope[i],
					),
				)
			}
		case entry.SpacePosBusAutoScopeAttributeOption:
			if err := s.Send(autoMessage, false); err != nil {
				errs = multierror.Append(
					errs, errors.WithMessagef(
						err, "failed to send message: %s", autoOption.Scope[i],
					),
				)
			}
		default:
			errs = multierror.Append(
				errs, errors.Errorf(
					"scope type in not supported: %s", autoOption.Scope[i],
				),
			)
		}
	}

	return errs.ErrorOrNil()
}

func (s *Space) calendarOnSpaceAttributeChanged(
	changeType universe.AttributeChangeType, attributeID entry.AttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
	world := s.GetWorld()
	if world == nil {
		return nil
	}

	switch changeType {
	case universe.ChangedAttributeChangeType:
		world.GetCalendar().OnAttributeUpsert(attributeID, value)
	case universe.RemovedAttributeChangeType:
		world.GetCalendar().OnAttributeRemove(attributeID)
	default:
		return errors.Errorf("unsupported change type: %s", changeType)
	}

	return nil
}

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
			s.spaceAttributes.Mu.Lock()
			defer s.spaceAttributes.Mu.Unlock()

			(*value)["auto_render_hash"] = hash.Hash

			if err := s.db.SpaceAttributesUpdateSpaceAttributeValue(
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
