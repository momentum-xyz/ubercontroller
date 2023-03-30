package tree

import (
	"github.com/hashicorp/go-multierror"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type ObjectTemplate struct {
	ObjectID         *umid.UMID           `json:"object_id"`
	ObjectName       *string              `json:"object_name"`
	ObjectTypeID     umid.UMID            `json:"object_type_id"`
	ParentID         umid.UMID            `json:"parent_id"`
	OwnerID          *umid.UMID           `json:"owner_id"`
	Asset2dID        *umid.UMID           `json:"asset_2d_id"`
	Asset3dID        *umid.UMID           `json:"asset_3d_id"`
	Options          *entry.ObjectOptions `json:"options"`
	Position         *cmath.Transform     `json:"position"`
	Label            *string              `json:"label"`
	ObjectAttributes []*entry.Attribute   `json:"object_attributes"`
	Objects          []*ObjectTemplate    `json:"objects"`
	RandomObjects    []*ObjectTemplate    `json:"random_objects"`
}

func AddObjectFromTemplate(objectTemplate *ObjectTemplate, updateDB bool) (umid.UMID, error) {
	node := universe.GetNode()

	// loading
	objectType, ok := node.GetObjectTypes().GetObjectType(objectTemplate.ObjectTypeID)
	if !ok {
		return umid.Nil, errors.Errorf("failed to get object type: %s", objectTemplate.ObjectTypeID)
	}

	parent, ok := node.GetObjectFromAllObjects(objectTemplate.ParentID)
	if !ok {
		return umid.Nil, errors.Errorf("parent object not found: %s", objectTemplate.ParentID)
	}

	// TODO: should be available for admin or owner of parent
	var asset2d universe.Asset2d
	if objectTemplate.Asset2dID != nil {
		asset2d, ok = node.GetAssets2d().GetAsset2d(*objectTemplate.Asset2dID)
		if !ok {
			return umid.Nil, errors.Errorf("asset 2d not found: %s", objectTemplate.Asset2dID)
		}
	}

	// TODO: should be available for admin or owner of parent
	var asset3d universe.Asset3d
	if objectTemplate.Asset3dID != nil {
		asset3d, ok = node.GetAssets3d().GetAsset3d(*objectTemplate.Asset3dID)
		if !ok {
			return umid.Nil, errors.Errorf("asset 3d not found: %s", objectTemplate.Asset3dID)
		}
	}

	objectID := objectTemplate.ObjectID
	ownerID := objectTemplate.OwnerID
	objectName := objectTemplate.ObjectName
	if objectID == nil {
		objectID = utils.GetPTR(umid.New())
	}
	if ownerID == nil {
		ownerID = utils.GetPTR(parent.GetOwnerID())
	}
	if objectName == nil {
		objectName = utils.GetPTR(objectID.String())
	}

	// creating
	object, err := parent.CreateObject(*objectID)
	if err != nil {
		return umid.Nil, errors.WithMessagef(err, "failed to create object: %s", objectID)
	}

	if err := object.SetOwnerID(*ownerID, false); err != nil {
		return umid.Nil, errors.WithMessagef(err, "failed to set owner umid: %s", ownerID)
	}
	if err := object.SetObjectType(objectType, false); err != nil {
		return umid.Nil, errors.WithMessagef(err, "failed to set object type: %s", objectTemplate.ObjectTypeID)
	}
	if asset2d != nil {
		if err := object.SetAsset2D(asset2d, false); err != nil {
			return umid.Nil, errors.WithMessagef(err, "failed to set asset 2d: %s", objectTemplate.Asset2dID)
		}
	}
	if asset3d != nil {
		if err := object.SetAsset3D(asset3d, false); err != nil {
			return umid.Nil, errors.WithMessagef(err, "failed to set asset 3d: %s", objectTemplate.Asset3dID)
		}
	}
	if objectTemplate.Position != nil {
		if err := object.SetTransform(objectTemplate.Position, false); err != nil {
			return umid.Nil, errors.WithMessagef(err, "failed to set position: %+v", objectTemplate.Position)
		}
	}

	if objectTemplate.Options != nil {
		if _, err := object.SetOptions(modify.MergeWith(objectTemplate.Options), false); err != nil {
			return umid.Nil, errors.WithMessage(err, "failed to set options")
		}
	}

	// saving in database
	if updateDB {
		if err := parent.AddObject(object, updateDB); err != nil {
			return umid.Nil, errors.WithMessage(err, "failed to add object")
		}
	}

	// running
	if err := object.Run(); err != nil {
		return umid.Nil, errors.WithMessage(err, "failed to run object")
	}
	if err := parent.UpdateChildrenPosition(true); err != nil {
		return umid.Nil, errors.WithMessage(err, "failed to update children position")
	}

	// adding children
	for i := range objectTemplate.Objects {
		objectTemplate.Objects[i].ParentID = *objectID
		if _, err := AddObjectFromTemplate(objectTemplate.Objects[i], updateDB); err != nil {
			return umid.Nil, errors.WithMessagef(
				err, "failed to add object from template: %+v", objectTemplate.Objects[i],
			)
		}
	}

	// enabling
	object.SetEnabled(true)

	// adding attributes
	if err := object.SetName(*objectName, true); err != nil {
		return umid.Nil, errors.WithMessage(err, "failed to set object name")
	}

	for i := range objectTemplate.ObjectAttributes {
		if _, err := object.GetObjectAttributes().Upsert(
			objectTemplate.ObjectAttributes[i].AttributeID,
			modify.MergeWith(objectTemplate.ObjectAttributes[i].AttributePayload),
			updateDB,
		); err != nil {
			return umid.Nil, errors.WithMessagef(
				err, "failed to upsert object attribute: %+v", objectTemplate.ObjectAttributes[i],
			)
		}
	}

	// updating
	if err := object.Update(true); err != nil {
		return umid.Nil, errors.WithMessage(err, "failed to update object")
	}

	return *objectID, nil
}

func RemoveObjectFromParent(parent, object universe.Object, updateDB bool) (bool, error) {
	if parent == nil {
		return false, errors.Errorf("parent is nil")
	}

	removed, err := parent.RemoveObject(object, true, updateDB)
	if err != nil {
		return false, errors.WithMessagef(err, "failed to remove object from parent: %s", parent.GetID())
	}
	if !removed {
		return false, nil
	}

	var errs *multierror.Error
	if object.GetEnabled() { // we need this check to avoid spam while removing children
		msg := posbus.WSMessage(&posbus.RemoveObjects{Objects: []umid.UMID{object.GetID()}})

		if err := object.GetWorld().Send(msg, true); err != nil {
			errs = multierror.Append(errs, errors.WithMessage(err, "failed to send remove message"))
		}
	}

	if err := parent.UpdateChildrenPosition(true); err != nil {
		errs = multierror.Append(
			errs, errors.WithMessagef(err, "failed to update children position: %s", parent.GetID()),
		)
	}

	if err := object.Stop(); err != nil {
		errs = multierror.Append(errs, errors.WithMessage(err, "failed to stop object"))
	}
	object.SetEnabled(false)

	logic.GetLogger().Infof("Helper: RemoveObjectFromParent: object removed: %s", object.GetID())

	go func() {
		for _, child := range object.GetObjects(false) {
			// prevent spam while removing
			child.SetEnabled(false)

			if _, err := RemoveObjectFromParent(object, child, false); err != nil {
				logic.GetLogger().Error(
					errors.WithMessagef(
						err, "Helper: RemoveObjectFromParent: failed to remove child from object: %s", child.GetID(),
					),
				)
			}
		}
	}()

	return removed, errs.ErrorOrNil()
}
