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
	return payload.Value, true
}

func (s *Space) GetSpaceAttributeOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	payload, ok := s.GetSpaceAttributePayload(attributeID)
	if !ok {
		return nil, false
	}
	return payload.Options, true
}

func (s *Space) GetSpaceAttributePayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	return s.attributes.Load(attributeID)
}

func (s *Space) GetSpaceAttributesValue(recursive bool) map[entry.SpaceAttributeID]*entry.AttributeValue {
	s.attributes.Mu.RLock()
	values := make(map[entry.SpaceAttributeID]*entry.AttributeValue, len(s.attributes.Data))
	for attributeID, payload := range s.attributes.Data {
		values[entry.NewSpaceAttributeID(attributeID, s.GetID())] = payload.Value
	}
	s.attributes.Mu.RUnlock()

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
	s.attributes.Mu.RLock()
	options := make(map[entry.SpaceAttributeID]*entry.AttributeOptions, len(s.attributes.Data))
	for attributeID, payload := range s.attributes.Data {
		options[entry.NewSpaceAttributeID(attributeID, s.GetID())] = payload.Options
	}
	s.attributes.Mu.RUnlock()

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
	s.attributes.Mu.RLock()
	payloads := make(map[entry.SpaceAttributeID]*entry.AttributePayload, len(s.attributes.Data))
	for attributeID, payload := range s.attributes.Data {
		payloads[entry.NewSpaceAttributeID(attributeID, s.GetID())] = payload
	}
	s.attributes.Mu.RUnlock()

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
	spaceAttribute *entry.SpaceAttribute, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) error {
	if spaceAttribute.SpaceID != s.GetID() {
		return errors.Errorf("space id mismatch: %s != %s", spaceAttribute.SpaceID, s.GetID())
	}

	s.attributes.Mu.Lock()
	defer s.attributes.Mu.Unlock()

	curPayload, ok := s.attributes.Data[spaceAttribute.AttributeID]
	if !ok {
		curPayload = (*entry.AttributePayload)(nil)
	}

	payload, err := modifyFn(curPayload)
	if err != nil {
		return errors.WithMessage(err, "failed to modify attribute payload")
	}

	if updateDB {
		if err := s.db.SpaceAttributesUpsertSpaceAttribute(
			s.ctx, entry.NewSpaceAttribute(spaceAttribute.SpaceAttributeID, payload),
		); err != nil {
			return errors.WithMessage(err, "failed to upsert space attribute")
		}
	}

	spaceAttribute.AttributePayload = payload
	s.attributes.Data[spaceAttribute.AttributeID] = spaceAttribute.AttributePayload

	return nil
}

func (s *Space) UpdateSpaceAttributeValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) error {
	s.attributes.Mu.Lock()
	defer s.attributes.Mu.Unlock()

	payload, ok := s.attributes.Data[attributeID]
	if !ok {
		return errors.Errorf("space attribute not found: %+v", attributeID)
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return errors.WithMessage(err, "failed to modify value")
	}

	if updateDB {
		if err := s.db.SpaceAttributesUpdateSpaceAttributeValue(
			s.ctx, entry.NewSpaceAttributeID(attributeID, s.GetID()), value,
		); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Value = value

	return nil
}

func (s *Space) UpdateSpaceAttributeOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) error {
	s.attributes.Mu.Lock()
	defer s.attributes.Mu.Unlock()

	payload, ok := s.attributes.Data[attributeID]
	if !ok {
		return errors.Errorf("space attribute not found: %+v", attributeID)
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := s.db.SpaceAttributesUpdateSpaceAttributeOptions(
			s.ctx, entry.NewSpaceAttributeID(attributeID, s.GetID()), options,
		); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	payload.Options = options

	return nil
}

func (s *Space) RemoveSpaceAttribute(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	s.attributes.Mu.Lock()
	defer s.attributes.Mu.Unlock()

	if _, ok := s.attributes.Data[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := s.db.SpaceAttributesRemoveSpaceAttributeByID(
			s.ctx, entry.NewSpaceAttributeID(attributeID, s.GetID()),
		); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	delete(s.attributes.Data, attributeID)

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
		return errors.WithMessage(err, "failed to get space attributes")
	}

	attributeTypes := universe.GetNode().GetAttributeTypes()
	for _, instance := range entries {
		if _, ok := attributeTypes.GetAttributeType(entry.AttributeTypeID(instance.AttributeID)); ok {
			s.CheckIfRendered(instance)
			if err := s.UpsertSpaceAttribute(
				instance, modify.MergeWith(instance.AttributePayload), false,
			); err != nil {
				return errors.WithMessagef(err, "failed to upsert space attribute: %+v", instance.AttributeID)
			}
		}
	}

	return nil
}
