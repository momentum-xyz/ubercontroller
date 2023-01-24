package object

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/unity"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.ObjectAttributes = (*objectAttributes)(nil)

type objectAttributes struct {
	object *Object
	data   map[entry.AttributeID]*entry.AttributePayload
}

func newObjectAttributes(object *Object) *objectAttributes {
	return &objectAttributes{
		object: object,
		data:   make(map[entry.AttributeID]*entry.AttributePayload),
	}
}

func (sa *objectAttributes) Load() error {
	entries, err := sa.object.db.GetObjectAttributesDB().GetObjectAttributesByObjectID(sa.object.ctx, sa.object.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get object attributes")
	}

	for i := range entries {
		if _, err := sa.Upsert(
			entries[i].AttributeID, modify.MergeWith(entries[i].AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert object attribute: %+v", entries[i].AttributeID)
		}

		effectiveOptions, ok := sa.GetEffectiveOptions(entries[i].AttributeID)
		if !ok {
			// QUESTION: why our "attribute_type.attribute_name" is not a foreign key in database?
			sa.object.log.Warnf(
				"Object attributes: Load: failed to get object attribute effective options: %+v",
				entries[i].ObjectAttributeID,
			)
			continue
		}
		autoOption, err := unity.GetOptionAutoOption(entries[i].AttributeID, effectiveOptions)
		if err != nil {
			return errors.WithMessagef(err, "failed to get option auto option: %+v", entries[i])
		}
		sa.object.UpdateAutoTextureMap(autoOption, entries[i].Value)
	}

	sa.object.log.Debugf("Object attributes loaded: %s: %d", sa.object.GetID(), sa.Len())

	return nil
}

func (sa *objectAttributes) Save() error {
	sa.object.Mu.RLock()
	defer sa.object.Mu.RUnlock()

	attributes := make([]*entry.ObjectAttribute, 0, len(sa.data))
	for id, payload := range sa.data {
		attributes = append(attributes, entry.NewObjectAttribute(entry.NewObjectAttributeID(id, sa.object.GetID()), payload))
	}

	if err := sa.object.db.GetObjectAttributesDB().UpsertObjectAttributes(sa.object.ctx, attributes); err != nil {
		return errors.WithMessage(err, "failed to upsert object attributes")
	}

	return nil
}

func (sa *objectAttributes) GetPayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	sa.object.Mu.RLock()
	defer sa.object.Mu.RUnlock()

	if payload, ok := sa.data[attributeID]; ok {
		return payload, true
	}
	return nil, false
}

func (sa *objectAttributes) GetValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	sa.object.Mu.RLock()
	defer sa.object.Mu.RUnlock()

	if payload, ok := sa.data[attributeID]; ok && payload != nil {
		return payload.Value, true
	}
	return nil, false
}

func (sa *objectAttributes) GetOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	sa.object.Mu.RLock()
	defer sa.object.Mu.RUnlock()

	if payload, ok := sa.data[attributeID]; ok && payload != nil {
		return payload.Options, true
	}
	return nil, false
}

func (sa *objectAttributes) GetEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
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
		sa.object.log.Error(
			errors.WithMessagef(
				err,
				"Object attributes: GetEffectiveOptions: failed to merge options: %s: %+v",
				sa.object.GetID(), attributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (sa *objectAttributes) Upsert(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	sa.object.Mu.Lock()
	defer sa.object.Mu.Unlock()

	payload, err := modifyFn(sa.data[attributeID])
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify payload")
	}

	if updateDB {
		if err := sa.object.db.GetObjectAttributesDB().UpsertObjectAttribute(
			sa.object.ctx, entry.NewObjectAttribute(entry.NewObjectAttributeID(attributeID, sa.object.GetID()), payload),
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to upsert object attribute")
		}
	}

	sa.data[attributeID] = payload

	if sa.object.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go sa.object.onObjectAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return payload, nil
}

func (sa *objectAttributes) UpdateValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	sa.object.Mu.Lock()
	defer sa.object.Mu.Unlock()

	payload, ok := sa.data[attributeID]
	if !ok {
		return nil, errors.Errorf("object attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify value")
	}

	if updateDB {
		if err := sa.object.db.GetObjectAttributesDB().UpdateObjectAttributeValue(
			sa.object.ctx, entry.NewObjectAttributeID(attributeID, sa.object.GetID()), value,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update object attribute value")
		}
	}

	payload.Value = value
	sa.data[attributeID] = payload

	if sa.object.GetEnabled() {
		go sa.object.onObjectAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return value, nil
}

func (sa *objectAttributes) UpdateOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	sa.object.Mu.Lock()
	defer sa.object.Mu.Unlock()

	payload, ok := sa.data[attributeID]
	if !ok {
		return nil, errors.Errorf("object attribute not found")
	}
	if payload == nil {
		payload = entry.NewAttributePayload(nil, nil)
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify options")
	}

	if updateDB {
		if err := sa.object.db.GetObjectAttributesDB().UpdateObjectAttributeOptions(
			sa.object.ctx, entry.NewObjectAttributeID(attributeID, sa.object.GetID()), options,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update object attribute options")
		}
	}

	payload.Options = options
	sa.data[attributeID] = payload

	if sa.object.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go sa.object.onObjectAttributeChanged(universe.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return options, nil
}

func (sa *objectAttributes) Remove(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	effectiveOptions, ok := sa.GetEffectiveOptions(attributeID)
	if !ok {
		return false, nil
	}

	sa.object.Mu.Lock()
	defer sa.object.Mu.Unlock()

	if _, ok := sa.data[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := sa.object.db.GetObjectAttributesDB().RemoveObjectAttributeByID(
			sa.object.ctx, entry.NewObjectAttributeID(attributeID, sa.object.GetID()),
		); err != nil {
			return false, errors.WithMessagef(err, "failed to remove object attribute")
		}
	}

	delete(sa.data, attributeID)

	if sa.object.GetEnabled() {
		go sa.object.onObjectAttributeChanged(universe.RemovedAttributeChangeType, attributeID, nil, effectiveOptions)
	}

	return true, nil
}

func (sa *objectAttributes) Len() int {
	sa.object.Mu.RLock()
	defer sa.object.Mu.RUnlock()

	return len(sa.data)
}

func (s *Object) onObjectAttributeChanged(
	changeType universe.AttributeChangeType, attributeID entry.AttributeID,
	value *entry.AttributeValue, effectiveOptions *entry.AttributeOptions,
) {
	go s.calendarOnObjectAttributeChanged(changeType, attributeID, value, effectiveOptions)

	if effectiveOptions == nil {
		options, ok := s.GetObjectAttributes().GetEffectiveOptions(attributeID)
		if !ok {
			s.log.Error(
				errors.Errorf(
					"Object: onObjectAttributeChanged: failed to get object attribute effective options: %+v",
					attributeID,
				),
			)
			return
		}
		effectiveOptions = options
	}

	go func() {
		if err := s.posBusAutoOnObjecteAttributeChanged(changeType, attributeID, value, effectiveOptions); err != nil {
			s.log.Error(
				errors.WithMessagef(
					err, "Object: onObjectAttributeChanged: failed to handle posbus auto: %s: %+v",
					s.GetID(), attributeID,
				),
			)
		}
	}()

	go func() {
		if err := s.unityAutoOnObjectAttributeChanged(changeType, attributeID, value, effectiveOptions); err != nil {
			s.log.Error(
				errors.WithMessagef(
					err, "Object: onObjectAttributeChanged: failed to handle unity auto: %s: %+v",
					s.GetID(), attributeID,
				),
			)
		}
	}()
}

func (s *Object) calendarOnObjectAttributeChanged(
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
