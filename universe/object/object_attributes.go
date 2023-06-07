package object

import (
	"context"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common/slot"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
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

func (oa *objectAttributes) GetAll() map[entry.AttributeID]*entry.AttributePayload {
	oa.object.Mu.RLock()
	defer oa.object.Mu.RUnlock()

	attributes := make(map[entry.AttributeID]*entry.AttributePayload, len(oa.data))
	for id, payload := range oa.data {
		attributes[id] = payload
	}

	return attributes
}

func (oa *objectAttributes) Load() error {
	oa.object.log.Debugf("Loading object attributes: %s...", oa.object.GetID())

	entries, err := oa.object.db.GetObjectAttributesDB().GetObjectAttributesByObjectID(oa.object.ctx, oa.object.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get object attributes")
	}

	for _, oaEntry := range entries {
		if _, err := oa.Upsert(
			oaEntry.AttributeID, modify.MergeWith(oaEntry.AttributePayload), false,
		); err != nil {
			return errors.WithMessagef(err, "failed to upsert object attribute: %+v", oaEntry.AttributeID)
		}

		effectiveOptions, ok := oa.GetEffectiveOptions(oaEntry.AttributeID)
		if !ok {
			// QUESTION: why our "attribute_type.attribute_name" is not a foreign key in database?
			oa.object.log.Warnf(
				"Object attributes: Load: failed to get object attribute effective options: %+v",
				oaEntry.ObjectAttributeID,
			)
			continue
		}
		autoOption, err := slot.GetOptionAutoOption(oaEntry.AttributeID, effectiveOptions)
		if err != nil {
			return errors.WithMessagef(err, "failed to get option auto option: %+v", oaEntry)
		}
		oa.object.UpdateAutoTextureMap(autoOption, oaEntry.Value)
	}

	oa.object.log.Debugf("Object attributes loaded: %s: %d", oa.object.GetID(), oa.Len())

	return nil
}

func (oa *objectAttributes) Save() error {
	oa.object.log.Debugf("Saving object attributes: %s...", oa.object.GetID())

	oa.object.Mu.RLock()
	defer oa.object.Mu.RUnlock()

	attributes := make([]*entry.ObjectAttribute, 0, len(oa.data))
	for id, payload := range oa.data {
		attributes = append(
			attributes, entry.NewObjectAttribute(entry.NewObjectAttributeID(id, oa.object.GetID()), payload),
		)
	}

	if err := oa.object.db.GetObjectAttributesDB().UpsertObjectAttributes(oa.object.ctx, attributes); err != nil {
		return errors.WithMessage(err, "failed to upsert object attributes")
	}

	oa.object.log.Debugf("Object attributes saved: %s: %d", oa.object.GetID(), len(oa.data))

	return nil
}

func (oa *objectAttributes) GetPayload(attributeID entry.AttributeID) (*entry.AttributePayload, bool) {
	oa.object.Mu.RLock()
	defer oa.object.Mu.RUnlock()

	if payload, ok := oa.data[attributeID]; ok {
		return payload, true
	}
	return nil, false
}

func (oa *objectAttributes) GetValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	oa.object.Mu.RLock()
	defer oa.object.Mu.RUnlock()

	if payload, ok := oa.data[attributeID]; ok && payload != nil {
		return payload.Value, true
	}
	return nil, false
}

func (oa *objectAttributes) GetOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	oa.object.Mu.RLock()
	defer oa.object.Mu.RUnlock()

	if payload, ok := oa.data[attributeID]; ok && payload != nil {
		return payload.Options, true
	}
	return nil, false
}

func (oa *objectAttributes) GetEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attributeType, ok := universe.GetNode().GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(attributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := oa.GetOptions(attributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		oa.object.log.Error(
			errors.WithMessagef(
				err,
				"Object attributes: GetEffectiveOptions: failed to merge options: %s: %+v",
				oa.object.GetID(), attributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (oa *objectAttributes) Upsert(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool,
) (*entry.AttributePayload, error) {
	oa.object.Mu.Lock()
	defer oa.object.Mu.Unlock()

	payload, err := modifyFn(oa.data[attributeID])
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to modify payload")
	}

	if updateDB {
		if err := oa.object.db.GetObjectAttributesDB().UpsertObjectAttribute(
			oa.object.ctx,
			entry.NewObjectAttribute(entry.NewObjectAttributeID(attributeID, oa.object.GetID()), payload),
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to upsert object attribute")
		}
	}

	oa.data[attributeID] = payload

	if oa.object.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go oa.object.onObjectAttributeChanged(posbus.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return payload, nil
}

func (oa *objectAttributes) UpdateValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) (*entry.AttributeValue, error) {
	oa.object.Mu.Lock()
	defer oa.object.Mu.Unlock()

	payload, ok := oa.data[attributeID]
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
		if err := oa.object.db.GetObjectAttributesDB().UpdateObjectAttributeValue(
			oa.object.ctx, entry.NewObjectAttributeID(attributeID, oa.object.GetID()), value,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update object attribute value")
		}
	}

	payload.Value = value
	oa.data[attributeID] = payload

	if oa.object.GetEnabled() {
		go oa.object.onObjectAttributeChanged(posbus.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return value, nil
}

func (oa *objectAttributes) UpdateOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) (*entry.AttributeOptions, error) {
	oa.object.Mu.Lock()
	defer oa.object.Mu.Unlock()

	payload, ok := oa.data[attributeID]
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
		if err := oa.object.db.GetObjectAttributesDB().UpdateObjectAttributeOptions(
			oa.object.ctx, entry.NewObjectAttributeID(attributeID, oa.object.GetID()), options,
		); err != nil {
			return nil, errors.WithMessagef(err, "failed to update object attribute options")
		}
	}

	payload.Options = options
	oa.data[attributeID] = payload

	if oa.object.GetEnabled() {
		var value *entry.AttributeValue
		if payload != nil {
			value = payload.Value
		}
		go oa.object.onObjectAttributeChanged(posbus.ChangedAttributeChangeType, attributeID, value, nil)
	}

	return options, nil
}

func (oa *objectAttributes) Remove(attributeID entry.AttributeID, updateDB bool) (bool, error) {
	effectiveOptions, ok := oa.GetEffectiveOptions(attributeID)
	if !ok {
		return false, nil
	}

	oa.object.Mu.Lock()
	defer oa.object.Mu.Unlock()

	if _, ok := oa.data[attributeID]; !ok {
		return false, nil
	}

	if updateDB {
		if err := oa.object.db.GetObjectAttributesDB().RemoveObjectAttributeByID(
			oa.object.ctx, entry.NewObjectAttributeID(attributeID, oa.object.GetID()),
		); err != nil {
			return false, errors.WithMessagef(err, "failed to remove object attribute")
		}
	}

	delete(oa.data, attributeID)

	if oa.object.GetEnabled() {
		go oa.object.onObjectAttributeChanged(posbus.RemovedAttributeChangeType, attributeID, nil, effectiveOptions)
	}

	return true, nil
}

func (oa *objectAttributes) Len() int {
	oa.object.Mu.RLock()
	defer oa.object.Mu.RUnlock()

	return len(oa.data)
}

func (o *Object) onObjectAttributeChanged(
	changeType posbus.AttributeChangeType, attributeID entry.AttributeID,
	value *entry.AttributeValue, effectiveOptions *entry.AttributeOptions,
) {
	go o.calendarOnObjectAttributeChanged(changeType, attributeID, value, effectiveOptions)

	if effectiveOptions == nil {
		options, ok := o.GetObjectAttributes().GetEffectiveOptions(attributeID)
		if !ok {
			o.log.Error(
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
		if err := o.posBusAutoOnObjecteAttributeChanged(changeType, attributeID, value, effectiveOptions); err != nil {
			o.log.Error(
				errors.WithMessagef(
					err, "Object: onObjectAttributeChanged: failed to handle posbus auto: %s: %+v",
					o.GetID(), attributeID,
				),
			)
		}
	}()

	go func() {
		if err := o.renderAutoOnObjectAttributeChanged(changeType, attributeID, value, effectiveOptions); err != nil {
			o.log.Error(
				errors.WithMessagef(
					err, "Object: onObjectAttributeChanged: failed to handle slot auto: %s: %+v",
					o.GetID(), attributeID,
				),
			)
		}
	}()
}

func (o *Object) calendarOnObjectAttributeChanged(
	changeType posbus.AttributeChangeType, attributeID entry.AttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
	world := o.GetWorld()
	if world == nil {
		return nil
	}

	switch changeType {
	case posbus.ChangedAttributeChangeType:
		world.GetCalendar().OnAttributeUpsert(attributeID, value)
	case posbus.RemovedAttributeChangeType:
		world.GetCalendar().OnAttributeRemove(attributeID)
	default:
		return errors.Errorf("unsupported change type: %s", changeType)
	}

	return nil
}

// AttributePermissionsAuthorizer
func (oa *objectAttributes) GetUserRoles(
	ctx context.Context,
	attrType entry.AttributeType,
	targetID entry.AttributeID,
	userID umid.UMID,
) ([]entry.PermissionsRoleType, error) {
	var roles []entry.PermissionsRoleType
	// owner is always considered an admin, TODO: add this to check function
	if oa.object.GetOwnerID() == userID {
		roles = append(roles, entry.PermissionAdmin)
	} else { // we have to lookup through the db user tree
		userObjectID := entry.NewUserObjectID(userID, oa.object.GetID())
		isAdmin, err := oa.object.db.GetUserObjectsDB().CheckIsIndirectAdminByID(ctx, userObjectID)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to check admin status")
		}
		if isAdmin {
			roles = append(roles, entry.PermissionAdmin)
		}
	}
	return roles, nil
}
