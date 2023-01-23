package object

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/unity"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.ObjectAttributes = (*spaceAttributes)(nil)

type spaceAttributes struct {
	space *Object
	data  map[entry.AttributeID]*entry.AttributePayload
}

func newSpaceAttributes(space *Object) *spaceAttributes {
	return &spaceAttributes{
		space: space,
		data:  make(map[entry.AttributeID]*entry.AttributePayload),
	}
}

func (sa *spaceAttributes) Load() error {
	entries, err := sa.space.db.GetObjectAttributesDB().GetObjectAttributesByObjectID(sa.space.ctx, sa.space.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get space attributes")
	}

	for i := range entries {
		if _, err := sa.Upsert(
			entries[i].AttributeID, modify.MergeWith(entries[i].AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert space attribute: %+v", entries[i].AttributeID)
		}

		effectiveOptions, ok := sa.GetEffectiveOptions(entries[i].AttributeID)
		if !ok {
			// QUESTION: why our "attribute_type.attribute_name" is not a foreign key in database?
			sa.space.log.Warnf(
				"Object: loadSpaceAttributes: failed to get space attribute effective options: %+v",
				entries[i].ObjectAttributeID,
			)
			continue
		}
		autoOption, err := unity.GetOptionAutoOption(entries[i].AttributeID, effectiveOptions)
		if err != nil {
			return errors.WithMessagef(err, "failed to get option auto option: %+v", entries[i])
		}
		sa.space.UpdateAutoTextureMap(autoOption, entries[i].Value)
	}

	sa.space.log.Debugf("Object attributes loaded: %s: %d", sa.space.GetID(), sa.Len())

	return nil
}

func (sa *spaceAttributes) Save() error {
	sa.space.Mu.RLock()
	defer sa.space.Mu.RUnlock()

	attributes := make([]*entry.ObjectAttribute, 0, len(sa.data))

	for id, payload := range sa.data {
		attributes = append(attributes, entry.NewObjectAttribute(entry.NewObjectAttributeID(id, sa.space.GetID()), payload))
	}

	if err := sa.space.db.GetObjectAttributesDB().UpsertObjectAttributes(sa.space.ctx, attributes); err != nil {
		return errors.WithMessage(err, "failed to upsert object attributes")
	}

	return nil
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
				"Object attributes: GetEffectiveOptions: failed to merge options: %s: %+v",
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
		if err := sa.space.db.GetObjectAttributesDB().UpsertObjectAttribute(
			sa.space.ctx, entry.NewObjectAttribute(entry.NewObjectAttributeID(attributeID, sa.space.GetID()), payload),
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
		if err := sa.space.db.GetObjectAttributesDB().UpdateObjectAttributeValue(
			sa.space.ctx, entry.NewObjectAttributeID(attributeID, sa.space.GetID()), value,
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
		if err := sa.space.db.GetObjectAttributesDB().UpdateObjectAttributeOptions(
			sa.space.ctx, entry.NewObjectAttributeID(attributeID, sa.space.GetID()), options,
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
		if err := sa.space.db.GetObjectAttributesDB().RemoveObjectAttributeByID(
			sa.space.ctx, entry.NewObjectAttributeID(attributeID, sa.space.GetID()),
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

func (s *Object) onSpaceAttributeChanged(
	changeType universe.AttributeChangeType, attributeID entry.AttributeID,
	value *entry.AttributeValue, effectiveOptions *entry.AttributeOptions,
) {
	go s.calendarOnSpaceAttributeChanged(changeType, attributeID, value, effectiveOptions)

	if effectiveOptions == nil {
		options, ok := s.GetObjectAttributes().GetEffectiveOptions(attributeID)
		if !ok {
			s.log.Error(
				errors.Errorf(
					"Object: onSpaceAttributeChanged: failed to get space attribute effective options: %+v",
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
					err, "Object: onSpaceAttributeChanged: failed to handle posbus auto: %s: %+v",
					s.GetID(), attributeID,
				),
			)
		}
	}()

	go func() {
		if err := s.unityAutoOnSpaceAttributeChanged(changeType, attributeID, value, effectiveOptions); err != nil {
			s.log.Error(
				errors.WithMessagef(
					err, "Object: onSpaceAttributeChanged: failed to handle unity auto: %s: %+v",
					s.GetID(), attributeID,
				),
			)
		}
	}()
}

func (s *Object) calendarOnSpaceAttributeChanged(
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
