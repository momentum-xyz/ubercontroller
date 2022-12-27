package space

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/unity"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.Attributes[entry.AttributeID] = (*Attributes)(nil)

type Attributes struct {
	space      *Space
	attributes map[entry.AttributeID]*entry.AttributePayload
}

func NewAttributes(space *Space) *Attributes {
	return &Attributes{
		space:      space,
		attributes: make(map[entry.AttributeID]*entry.AttributePayload),
	}
}

func (a *Attributes) GetPayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	a.space.mu.RLock()
	defer a.space.mu.RUnlock()

	if payload, ok := a.attributes[attributeID]; ok {
		return payload, true
	}

	return nil, false
}

func (a *Attributes) GetValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	payload, ok := a.GetPayload(attributeID)
	if !ok {
		return nil, false
	}
	if payload == nil {
		return nil, true
	}
	return payload.Value, true
}

func (a *Attributes) GetOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	payload, ok := a.GetPayload(attributeID)
	if !ok {
		return nil, false
	}
	if payload == nil {
		return nil, true
	}
	return payload.Options, true
}

func (a *Attributes) GetEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attributeType, ok := universe.GetNode().GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(attributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := a.GetOptions(attributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		a.space.log.Error(
			errors.WithMessagef(
				err,
				"Space: Attributes: GetEffectiveOptions: failed to merge options: %s: %+v",
				a.space.GetID(), attributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (a *Attributes) Upsert(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	a.space.mu.Lock()
	defer a.space.mu.Unlock()

	payload, err := modifyFn(a.attributes[attributeID])
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify payload")
	}

	if updateDB {
		if err := a.space.db.SpaceAttributesUpsertSpaceAttribute(
			a.space.ctx, entry.NewSpaceAttribute(entry.NewSpaceAttributeID(attributeID, a.space.GetID()), payload),
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to upsert space attribute")
		}
	}

	a.attributes[attributeID] = payload

	if a.space.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go a.space.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return payload, nil
}

func (a *Attributes) UpdateValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	a.space.mu.Lock()
	defer a.space.mu.Unlock()

	payload, ok := a.attributes[attributeID]
	if !ok {
		return nil, errors.Errorf("attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify value")
	}

	if updateDB {
		if err := a.space.db.SpaceAttributesUpdateSpaceAttributeValue(
			a.space.ctx, entry.NewSpaceAttributeID(attributeID, a.space.GetID()), value,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update space attribute value")
		}
	}

	payload.Value = value
	a.attributes[attributeID] = payload

	if a.space.GetEnabled() {
		go a.space.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return value, nil
}

func (a *Attributes) UpdateOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	a.space.mu.Lock()
	defer a.space.mu.Unlock()

	payload, ok := a.attributes[attributeID]
	if !ok {
		return nil, errors.Errorf("attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify options")
	}

	if updateDB {
		if err := a.space.db.SpaceAttributesUpdateSpaceAttributeOptions(
			a.space.ctx, entry.NewSpaceAttributeID(attributeID, a.space.GetID()), options,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update space attribute options")
		}
	}

	payload.Options = options
	a.attributes[attributeID] = payload

	if a.space.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go a.space.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return options, nil
}

func (a *Attributes) Remove(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	effectiveOptions, ok := a.GetEffectiveOptions(attributeID)
	if !ok {
		return false, nil
	}

	a.space.mu.Lock()
	defer a.space.mu.Unlock()

	if _, ok := a.attributes[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := a.space.db.SpaceAttributesRemoveSpaceAttributeByID(
			a.space.ctx, entry.NewSpaceAttributeID(attributeID, a.space.GetID()),
		); err != nil {
			return false, errors.WithMessagef(err, "failed to remove space attribute")
		}
	}

	delete(a.attributes, attributeID)

	if a.space.GetEnabled() {
		go a.space.onSpaceAttributeChanged(universe.RemovedAttributeChangeType, attributeID, nil, effectiveOptions)
	}

	return true, nil
}

func (a *Attributes) Len() int {
	a.space.mu.RLock()
	defer a.space.mu.RUnlock()

	return len(a.attributes)
}

func (s *Space) onSpaceAttributeChanged(
	changeType universe.AttributeChangeType, attributeID entry.AttributeID,
	value *entry.AttributeValue, effectiveOptions *entry.AttributeOptions,
) {
	go s.calendarOnSpaceAttributeChanged(changeType, attributeID, value, effectiveOptions)

	if effectiveOptions == nil {
		options, ok := s.GetSpaceAttributes().GetEffectiveOptions(attributeID)
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
					err, "Space: onSpaceAttributeChanged: failed to handle posbus auto: %s: %+v",
					s.GetID(), attributeID,
				),
			)
		}
	}()

	go func() {
		if err := s.unityAutoOnSpaceAttributeChanged(changeType, attributeID, value, effectiveOptions); err != nil {
			s.log.Error(
				errors.WithMessagef(
					err, "Space: onSpaceAttributeChanged: failed to handle unity auto: %s: %+v",
					s.GetID(), attributeID,
				),
			)
		}
	}()
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

func (s *Space) loadSpaceAttributes() error {
	entries, err := s.db.SpaceAttributesGetSpaceAttributesBySpaceID(s.ctx, s.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get space attributes")
	}

	attributes := s.GetSpaceAttributes()
	for i := range entries {
		entry := entries[i]

		if _, err := attributes.Upsert(
			entry.AttributeID, modify.MergeWith(entry.AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert space attribute: %+v", entry.AttributeID)
		}

		effectiveOptions, ok := attributes.GetEffectiveOptions(entry.AttributeID)
		if !ok {
			continue
		}
		autoOption, err := unity.GetOptionAutoOption(entry.AttributeID, effectiveOptions)
		if err != nil {
			return errors.WithMessagef(err, "failed to get option auto option: %+v", entry)
		}
		s.UpdateAutoTextureMap(autoOption, entry.Value)
	}

	s.log.Debugf("Space attributes loaded: %s: %d", s.GetID(), s.Attributes.Len())

	return nil
}
