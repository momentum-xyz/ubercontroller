package helper

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/pkg/errors"
)

type ObjectTemplate struct {
	entry.Object
	ObjectName       *string            `json:"object_name"`
	Label            *string            `json:"label"`
	ObjectAttributes []*entry.Attribute `json:"object_attributes"`
	Children         []*ObjectTemplate  `json:"children"`
}

// TODO: think about rollback
func AddObjectFromTemplate(objectTemplate *ObjectTemplate, updateDB bool) (universe.Object, error) {
	// loading
	parent, ok := universe.GetNode().GetObjectFromAllObjects(objectTemplate.ParentID)
	if !ok {
		return nil, errors.Errorf("parent object not found: %s", objectTemplate.ParentID)
	}

	// creating
	object, err := createObjectFromTemplate(parent, objectTemplate)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create object from template: %+v", objectTemplate)
	}

	// saving in database
	if updateDB {
		if err := parent.AddObject(object, true); err != nil {
			return object, errors.WithMessagef(
				err, "failed to add object %s to parent %s", object.GetID(), parent.GetID(),
			)
		}
	}

	// updating
	if err := object.Update(true); err != nil {
		return object, errors.WithMessage(err, "failed to update object")
	}
	if err := parent.UpdateChildrenPosition(true); err != nil {
		return object, errors.WithMessage(err, "failed to update children position")
	}

	return object, nil
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
		removeMsg := posbus.NewRemoveStaticObjectsMsg(1)
		removeMsg.SetObject(0, object.GetID())
		if err := object.GetWorld().Send(removeMsg.WebsocketMessage(), true); err != nil {
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

	common.GetLogger().Infof("Helper: RemoveObjectFromParent: object removed: %s", object.GetID())

	go func() {
		for _, child := range object.GetObjects(false) {
			// prevent spam while removing
			child.SetEnabled(false)

			if _, err := RemoveObjectFromParent(object, child, false); err != nil {
				common.GetLogger().Error(
					errors.WithMessagef(
						err, "Helper: RemoveObjectFromParent: failed to remove child: %s", child.GetID(),
					),
				)
			}
		}
	}()

	return removed, errs.ErrorOrNil()
}

func CalcObjectSpawnPosition(parentID, userID uuid.UUID) (*cmath.SpacePosition, error) {
	parent, ok := universe.GetNode().GetObjectFromAllObjects(parentID)
	if !ok {
		return nil, errors.Errorf("object parent not found: %s", parentID)
	}

	var position *cmath.SpacePosition
	effectiveOptions := parent.GetEffectiveOptions()
	if effectiveOptions == nil || len(effectiveOptions.ChildPlacements) == 0 {
		world := parent.GetWorld()
		if world != nil {
			user, ok := world.GetUser(userID, true)
			if ok {
				fmt.Printf("User rotation: %v", user.GetRotation())
				//distance := float32(10)
				position = &cmath.SpacePosition{
					// TODO: recalc based on euler angles, not lookat: Location: cmath.Add(user.GetPosition(), cmath.MultiplyN(user.GetRotation(), distance)),
					Location: user.GetPosition(),
					Rotation: cmath.Vec3{},
					Scale:    cmath.Vec3{X: 1, Y: 1, Z: 1},
				}
			}
		}
	}

	return position, nil
}

// createObjectFromTemplate creates in-memory ready for use object with children from template
func createObjectFromTemplate(parent universe.Object, objectTemplate *ObjectTemplate) (universe.Object, error) {
	// TODO: think about rollback
	if parent == nil {
		return nil, errors.Errorf("parent is nil")
	}

	node := universe.GetNode()

	// loading
	objectType, ok := node.GetObjectTypes().GetObjectType(objectTemplate.ObjectTypeID)
	if !ok {
		return nil, errors.Errorf("failed to get object type: %s", objectTemplate.ObjectTypeID)
	}
	// TODO: should be available for admin or owner of parent
	var asset2d universe.Asset2d
	if objectTemplate.Asset2dID != nil {
		asset2d, ok = node.GetAssets2d().GetAsset2d(*objectTemplate.Asset2dID)
		if !ok {
			return nil, errors.Errorf("asset 2d not found: %s", objectTemplate.Asset2dID)
		}
	}
	// TODO: should be available for admin or owner of parent
	var asset3d universe.Asset3d
	if objectTemplate.Asset3dID != nil {
		asset3d, ok = node.GetAssets3d().GetAsset3d(*objectTemplate.Asset3dID)
		if !ok {
			return nil, errors.Errorf("asset 3d not found: %s", objectTemplate.Asset3dID)
		}
	}

	ownerID := objectTemplate.OwnerID
	objectID := objectTemplate.ObjectID
	objectName := objectTemplate.ObjectName
	if ownerID == uuid.Nil {
		ownerID = parent.GetOwnerID()
	}
	if objectID == uuid.Nil {
		objectID = uuid.New()
	}
	if objectName == nil {
		objectName = utils.GetPTR(objectID.String())
	}

	// creating
	object, err := parent.CreateObject(objectID)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create object: %s", objectID)
	}

	if err := object.SetOwnerID(ownerID, false); err != nil {
		return object, errors.WithMessagef(err, "failed to set owner id: %s", ownerID)
	}
	if err := object.SetObjectType(objectType, false); err != nil {
		return object, errors.WithMessagef(err, "failed to set object type: %s", objectTemplate.ObjectTypeID)
	}
	if asset2d != nil {
		if err := object.SetAsset2D(asset2d, false); err != nil {
			return object, errors.WithMessagef(err, "failed to set asset 2d: %s", objectTemplate.Asset2dID)
		}
	}
	if asset3d != nil {
		if err := object.SetAsset3D(asset3d, false); err != nil {
			return object, errors.WithMessagef(err, "failed to set asset 3d: %s", objectTemplate.Asset3dID)
		}
	}
	if objectTemplate.Position != nil {
		if err := object.SetPosition(objectTemplate.Position, false); err != nil {
			return object, errors.WithMessagef(err, "failed to set position: %+v", objectTemplate.Position)
		}
	}

	// running
	if err := object.Run(); err != nil {
		return object, errors.WithMessage(err, "failed to run object")
	}

	// adding children
	for _, childTemplate := range objectTemplate.Children {
		if _, err := createObjectFromTemplate(object, childTemplate); err != nil {
			return object, errors.WithMessagef(err, "failed to create child from template: %+v", childTemplate)
		}
	}

	// enabling
	object.SetEnabled(true)

	// adding attributes
	if err := object.SetName(*objectName, false); err != nil {
		return object, errors.WithMessage(err, "failed to set object name")
	}

	for _, attribute := range objectTemplate.ObjectAttributes {
		if _, err := object.GetObjectAttributes().Upsert(
			attribute.AttributeID, modify.MergeWith(attribute.AttributePayload), false,
		); err != nil {
			return object, errors.WithMessagef(err, "failed to upsert object attribute: %+v", attribute.AttributeID)
		}
	}

	return object, nil
}
