package space

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/unity"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.Attributes[entry.AttributeID] = (*spaceAttributes)(nil)

type spaceAttributes struct {
	space *Space
	data  map[entry.AttributeID]*entry.AttributePayload
}

func newSpaceAttributes(space *Space) *spaceAttributes {
	return &spaceAttributes{
		space: space,
		data:  make(map[entry.AttributeID]*entry.AttributePayload),
	}
}

func (sa *spaceAttributes) GetPayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	sa.space.Mu.RLock()
	defer sa.space.Mu.RUnlock()

	if payload, ok := sa.data[attributeID]; ok {
		return payload, true
	}
	return nil, false
}

func (sa *spaceAttributes) GetValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	sa.space.Mu.RLock()
	defer sa.space.Mu.RUnlock()

	if payload, ok := sa.data[attributeID]; ok && payload != nil {
		return payload.Value, true
	}
	return nil, false
}

func (sa *spaceAttributes) GetOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	sa.space.Mu.RLock()
	defer sa.space.Mu.RUnlock()

	if payload, ok := sa.data[attributeID]; ok && payload != nil {
		return payload.Options, true
	}
	return nil, false
}

func (sa *spaceAttributes) GetEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attributeType, ok := universe.GetNode().GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(attributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := sa.GetOptions(attributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		sa.space.log.Error(
			errors.WithMessagef(
				err,
				"Space attributes: GetEffectiveOptions: failed to merge options: %s: %+v",
				sa.space.GetID(), attributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (sa *spaceAttributes) Upsert(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	sa.space.Mu.Lock()
	defer sa.space.Mu.Unlock()

	payload, err := modifyFn(sa.data[attributeID])
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify payload")
	}

	if updateDB {
		if err := sa.space.db.GetSpaceAttributesDB().UpsertSpaceAttribute(
			sa.space.ctx, entry.NewSpaceAttribute(entry.NewSpaceAttributeID(attributeID, sa.space.GetID()), payload),
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to upsert space attribute")
		}
	}

	sa.data[attributeID] = payload

	if sa.space.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go sa.space.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return payload, nil
}

func (sa *spaceAttributes) UpdateValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	sa.space.Mu.Lock()
	defer sa.space.Mu.Unlock()

	payload, ok := sa.data[attributeID]
	if !ok {
		return nil, errors.Errorf("space attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify value")
	}

	if updateDB {
		if err := sa.space.db.GetSpaceAttributesDB().UpdateSpaceAttributeValue(
			sa.space.ctx, entry.NewSpaceAttributeID(attributeID, sa.space.GetID()), value,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update space attribute value")
		}
	}

	payload.Value = value
	sa.data[attributeID] = payload

	if sa.space.GetEnabled() {
		go sa.space.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return value, nil
}

func (sa *spaceAttributes) UpdateOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	sa.space.Mu.Lock()
	defer sa.space.Mu.Unlock()

	payload, ok := sa.data[attributeID]
	if !ok {
		return nil, errors.Errorf("space attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify options")
	}

	if updateDB {
		if err := sa.space.db.GetSpaceAttributesDB().UpdateSpaceAttributeOptions(
			sa.space.ctx, entry.NewSpaceAttributeID(attributeID, sa.space.GetID()), options,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update space attribute options")
		}
	}

	payload.Options = options
	sa.data[attributeID] = payload

	if sa.space.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go sa.space.onSpaceAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return options, nil
}

func (sa *spaceAttributes) Remove(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	effectiveOptions, ok := sa.GetEffectiveOptions(attributeID)
	if !ok {
		return false, nil
	}

	sa.space.Mu.Lock()
	defer sa.space.Mu.Unlock()

	if _, ok := sa.data[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := sa.space.db.GetSpaceAttributesDB().RemoveSpaceAttributeByID(
			sa.space.ctx, entry.NewSpaceAttributeID(attributeID, sa.space.GetID()),
		); err != nil {
			return false, errors.WithMessagef(err, "failed to remove space attribute")
		}
	}

	delete(sa.data, attributeID)

	if sa.space.GetEnabled() {
		go sa.space.onSpaceAttributeChanged(universe.RemovedAttributeChangeType, attributeID, nil, effectiveOptions)
	}

	return true, nil
}

func (sa *spaceAttributes) Len() int {
	sa.space.Mu.RLock()
	defer sa.space.Mu.RUnlock()

	return len(sa.data)
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
					"Space: onSpaceAttributeChanged: failed to get space attribute effective options: %+v",
					attributeID,
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
	entries, err := s.db.GetSpaceAttributesDB().GetSpaceAttributesBySpaceID(s.ctx, s.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get space attributes")
	}

	attributes := s.GetSpaceAttributes()
	for i := range entries {
		if _, err := attributes.Upsert(
			entries[i].AttributeID, modify.MergeWith(entries[i].AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert space attribute: %+v", entries[i].AttributeID)
		}

		effectiveOptions, ok := attributes.GetEffectiveOptions(entries[i].AttributeID)
		if !ok {
			// QUESTION: why our "attribute_type.attribute_name" is not a foreign key in database?
			s.log.Warnf(
				"Space: loadSpaceAttributes: failed to get space attribute effective options: %+v",
				entries[i].SpaceAttributeID,
			)
			continue
		}
		autoOption, err := unity.GetOptionAutoOption(entries[i].AttributeID, effectiveOptions)
		if err != nil {
			return errors.WithMessagef(err, "failed to get option auto option: %+v", entries[i])
		}
		s.UpdateAutoTextureMap(autoOption, entries[i].Value)
	}

	s.log.Debugf("Space attributes loaded: %s: %d", s.GetID(), attributes.Len())

	return nil
}
