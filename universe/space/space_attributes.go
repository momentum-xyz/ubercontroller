package space

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
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

func (s *Space) GetSpaceAttributePayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
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
			options[spaceAttributeID] = payload.Options
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
	s.spaceAttributes.Mu.Lock()
	defer s.spaceAttributes.Mu.Unlock()

	payload, ok := s.spaceAttributes.Data[attributeID]
	if !ok {
		payload = (*entry.AttributePayload)(nil)
	}

	payload, err := modifyFn(payload)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify attribute payload")
	}

	spaceAttribute := entry.NewSpaceAttribute(entry.NewSpaceAttributeID(attributeID, s.GetID()), payload)
	if updateDB {
		if err := s.db.SpaceAttributesUpsertSpaceAttribute(s.ctx, spaceAttribute); err != nil {
			return nil, errors.WithMessage(err, "failed to upsert space attribute")
		}
	}

	s.spaceAttributes.Data[attributeID] = payload

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

	return options, nil
}

func (s *Space) RemoveSpaceAttribute(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	s.spaceAttributes.Mu.Lock()
	defer s.spaceAttributes.Mu.Unlock()

	if _, ok := s.spaceAttributes.Data[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := s.db.SpaceAttributesRemoveSpaceAttributeByID(
			s.ctx, entry.NewSpaceAttributeID(attributeID, s.GetID()),
		); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(s.spaceAttributes.Data, attributeID)

	return true, nil
}

// TODO: optimize
func (s *Space) RemoveSpaceAttributes(attributeIDs []entry.AttributeID, updateDB bool) (bool, error) {
	res := true
	var errs *multierror.Error
	for i := range attributeIDs {
		removed, err := s.RemoveSpaceAttribute(attributeIDs[i], updateDB)
		if err != nil {
			errs = multierror.Append(errs,
				errors.WithMessagef(err, "failed to remove space attribute: %+v", attributeIDs[i]),
			)
		}
		if !removed {
			res = false
		}
	}
	return res, errs.ErrorOrNil()
}

func (s *Space) CheckIfRendered(instance *entry.SpaceAttribute) {
	attr, ok := universe.GetNode().GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(instance.AttributeID))
	if !ok {
		return
	}
	var opts entry.AttributeOptions
	//if instance.Options != nil {
	//	opts = *instance.Options
	//} else {
	//	opts = entry.AttributeOptions{}
	//}

	if attr.GetOptions() != nil {
		opts = *attr.GetOptions()
	} else {
		return
	}
	//utils.MergePTRs(&opts, attr.GetOptions())
	//if opts == nil {
	//	return
	//}
	if v, ok := opts["render_type"]; ok && v.(string) == "texture" {
		if c, ok := (map[string]any)(*instance.Value)["render_hash"]; ok {
			s.renderTextureAttr[attr.GetName()] = c.(string)
		}
	}
	s.textMsg.Store(message.GetBuilder().SetObjectTextures(s.id, s.renderTextureAttr))
}

func (s *Space) loadSpaceAttributes() error {
	entries, err := s.db.SpaceAttributesGetSpaceAttributesBySpaceID(s.ctx, s.id)
	if err != nil {
		return errors.WithMessage(err, "failed to get space spaceAttributes")
	}

	for _, instance := range entries {
		s.CheckIfRendered(instance)
		if _, err := s.UpsertSpaceAttribute(
			instance.AttributeID, modify.MergeWith(instance.AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert space attribute: %+v", instance.AttributeID)
		}
	}

	return nil
}
