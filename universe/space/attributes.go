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
	Space      *Space
	Attributes map[entry.AttributeID]*entry.AttributePayload
}

func NewAttributes(space *Space) *Attributes {
	return &Attributes{
		Space:      space,
		Attributes: make(map[entry.AttributeID]*entry.AttributePayload),
	}
}

func (a *Attributes) GetPayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	a.Space.mu.RLock()
	defer a.Space.mu.RUnlock()

	if payload, ok := a.Attributes[attributeID]; ok {
		return payload, true
	}
	return nil, false
}

func (a *Attributes) GetValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	a.Space.mu.RLock()
	defer a.Space.mu.RUnlock()

	if payload, ok := a.Attributes[attributeID]; ok && payload != nil {
		return payload.Value, true
	}
	return nil, false
}

func (a *Attributes) GetOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	a.Space.mu.RLock()
	defer a.Space.mu.RUnlock()

	if payload, ok := a.Attributes[attributeID]; ok && payload != nil {
		return payload.Options, true
	}
	return nil, false
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
		a.Space.log.Error(
			errors.WithMessagef(
				err,
				"Space: Attributes: GetEffectiveOptions: failed to merge options: %s: %+v",
				a.Space.GetID(), attributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (a *Attributes) Upsert(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	a.Space.mu.Lock()
	defer a.Space.mu.Unlock()

	payload, err := modifyFn(a.Attributes[attributeID])
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify payload")
	}

	if updateDB {
		if err := a.Space.db.SpaceAttributesUpsertSpaceAttribute(
			a.Space.ctx, entry.NewSpaceAttribute(entry.NewSpaceAttributeID(attributeID, a.Space.GetID()), payload),
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to upsert space attribute")
		}
	}

	a.Attributes[attributeID] = payload

	if a.Space.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go a.Space.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return payload, nil
}

func (a *Attributes) UpdateValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	a.Space.mu.Lock()
	defer a.Space.mu.Unlock()

	payload, ok := a.Attributes[attributeID]
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
		if err := a.Space.db.SpaceAttributesUpdateSpaceAttributeValue(
			a.Space.ctx, entry.NewSpaceAttributeID(attributeID, a.Space.GetID()), value,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update space attribute value")
		}
	}

	payload.Value = value
	a.Attributes[attributeID] = payload

	if a.Space.GetEnabled() {
		go a.Space.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return value, nil
}

func (a *Attributes) UpdateOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	a.Space.mu.Lock()
	defer a.Space.mu.Unlock()

	payload, ok := a.Attributes[attributeID]
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
		if err := a.Space.db.SpaceAttributesUpdateSpaceAttributeOptions(
			a.Space.ctx, entry.NewSpaceAttributeID(attributeID, a.Space.GetID()), options,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update space attribute options")
		}
	}

	payload.Options = options
	a.Attributes[attributeID] = payload

	if a.Space.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go a.Space.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return options, nil
}

func (a *Attributes) Remove(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	effectiveOptions, ok := a.GetEffectiveOptions(attributeID)
	if !ok {
		return false, nil
	}

	a.Space.mu.Lock()
	defer a.Space.mu.Unlock()

	if _, ok := a.Attributes[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := a.Space.db.SpaceAttributesRemoveSpaceAttributeByID(
			a.Space.ctx, entry.NewSpaceAttributeID(attributeID, a.Space.GetID()),
		); err != nil {
			return false, errors.WithMessagef(err, "failed to remove space attribute")
		}
	}

	delete(a.Attributes, attributeID)

	if a.Space.GetEnabled() {
		go a.Space.onSpaceAttributeChanged(universe.RemovedAttributeChangeType, attributeID, nil, effectiveOptions)
	}

	return true, nil
}

func (a *Attributes) Len() int {
	a.Space.mu.RLock()
	defer a.Space.mu.RUnlock()

	return len(a.Attributes)
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
			// QUESTION: why our "attribute_type.attribute_name" is not a foreign key in database?
			s.log.Warnf(
				"Space: loadSpaceAttributes: failed to get attribute effective options: %+v", entry.SpaceAttributeID,
			)
			continue
		}
		autoOption, err := unity.GetOptionAutoOption(entry.AttributeID, effectiveOptions)
		if err != nil {
			return errors.WithMessagef(err, "failed to get option auto option: %+v", entry)
		}
		s.UpdateAutoTextureMap(autoOption, entry.Value)
	}

	s.log.Debugf("Space attributes loaded: %s: %d", s.GetID(), s.GetSpaceAttributes().Len())

	return nil
}
