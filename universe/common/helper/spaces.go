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

type SpaceTemplate struct {
	SpaceID         *uuid.UUID           `json:"space_id"`
	SpaceName       *string              `json:"space_name"`
	SpaceTypeID     uuid.UUID            `json:"space_type_id"`
	ParentID        uuid.UUID            `json:"parent_id"`
	OwnerID         *uuid.UUID           `json:"owner_id"`
	Asset2dID       *uuid.UUID           `json:"asset_2d_id"`
	Asset3dID       *uuid.UUID           `json:"asset_3d_id"`
	Options         *entry.SpaceOptions  `json:"options"`
	Position        *cmath.SpacePosition `json:"position"`
	Label           *string              `json:"label"`
	SpaceAttributes []*entry.Attribute   `json:"space_attributes"`
	Spaces          []*SpaceTemplate     `json:"spaces"`
	RandomSpaces    []*SpaceTemplate     `json:"random_spaces"`
}

func AddSpaceFromTemplate(spaceTemplate *SpaceTemplate, updateDB bool) (uuid.UUID, error) {
	node := universe.GetNode()

	// loading
	spaceType, ok := node.GetSpaceTypes().GetSpaceType(spaceTemplate.SpaceTypeID)
	if !ok {
		return uuid.Nil, errors.Errorf("failed to get space type: %s", spaceTemplate.SpaceTypeID)
	}

	parent, ok := node.GetSpaceFromAllSpaces(spaceTemplate.ParentID)
	if !ok {
		return uuid.Nil, errors.Errorf("parent space not found: %s", spaceTemplate.ParentID)
	}

	// TODO: should be available for admin or owner of parent
	var asset2d universe.Asset2d
	if spaceTemplate.Asset2dID != nil {
		asset2d, ok = node.GetAssets2d().GetAsset2d(*spaceTemplate.Asset2dID)
		if !ok {
			return uuid.Nil, errors.Errorf("asset 2d not found: %s", spaceTemplate.Asset2dID)
		}
	}

	// TODO: should be available for admin or owner of parent
	var asset3d universe.Asset3d
	if spaceTemplate.Asset3dID != nil {
		asset3d, ok = node.GetAssets3d().GetAsset3d(*spaceTemplate.Asset3dID)
		if !ok {
			return uuid.Nil, errors.Errorf("asset 3d not found: %s", spaceTemplate.Asset3dID)
		}
	}

	spaceID := spaceTemplate.SpaceID
	ownerID := spaceTemplate.OwnerID
	spaceName := spaceTemplate.SpaceName
	if spaceID == nil {
		spaceID = utils.GetPTR(uuid.New())
	}
	if ownerID == nil {
		ownerID = utils.GetPTR(parent.GetOwnerID())
	}
	if spaceName == nil {
		spaceName = utils.GetPTR(spaceID.String())
	}

	// creating
	space, err := parent.CreateSpace(*spaceID)
	if err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to create space: %s", spaceID)
	}

	if err := space.SetOwnerID(*ownerID, false); err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to set owner id: %s", ownerID)
	}
	if err := space.SetSpaceType(spaceType, false); err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to set space type: %s", spaceTemplate.SpaceTypeID)
	}
	if asset2d != nil {
		if err := space.SetAsset2D(asset2d, false); err != nil {
			return uuid.Nil, errors.WithMessagef(err, "failed to set asset 2d: %s", spaceTemplate.Asset2dID)
		}
	}
	if asset3d != nil {
		if err := space.SetAsset3D(asset3d, false); err != nil {
			return uuid.Nil, errors.WithMessagef(err, "failed to set asset 3d: %s", spaceTemplate.Asset3dID)
		}
	}
	if spaceTemplate.Position != nil {
		if err := space.SetPosition(spaceTemplate.Position, false); err != nil {
			return uuid.Nil, errors.WithMessagef(err, "failed to set position: %+v", spaceTemplate.Position)
		}
	}

	if spaceTemplate.Options != nil {
		if _, err := space.SetOptions(modify.MergeWith(spaceTemplate.Options), false); err != nil {
			return uuid.Nil, errors.WithMessage(err, "failed to set options")
		}
	}

	// saving in database
	if updateDB {
		if err := parent.AddSpace(space, updateDB); err != nil {
			return uuid.Nil, errors.WithMessage(err, "failed to add space")
		}
	}

	// running
	if err := space.Run(); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to run space")
	}
	if err := parent.UpdateChildrenPosition(true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to update children position")
	}

	// adding children
	for i := range spaceTemplate.Spaces {
		spaceTemplate.Spaces[i].ParentID = *spaceID
		if _, err := AddSpaceFromTemplate(spaceTemplate.Spaces[i], updateDB); err != nil {
			return uuid.Nil, errors.WithMessagef(
				err, "failed to add space from template: %+v", spaceTemplate.Spaces[i],
			)
		}
	}

	// enabling
	space.SetEnabled(true)

	// adding attributes
	if err := space.SetName(*spaceName, true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to set space name")
	}

	for i := range spaceTemplate.SpaceAttributes {
		if _, err := space.GetSpaceAttributes().Upsert(
			spaceTemplate.SpaceAttributes[i].AttributeID,
			modify.MergeWith(spaceTemplate.SpaceAttributes[i].AttributePayload),
			updateDB,
		); err != nil {
			return uuid.Nil, errors.WithMessagef(err, "failed to upsert space attribute: %+v", spaceTemplate.SpaceAttributes[i])
		}
	}

	// updating
	if err := space.Update(true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to update space")
	}

	return *spaceID, nil
}

func RemoveSpaceFromParent(parent, space universe.Space, updateDB bool) (bool, error) {
	if parent == nil {
		return false, errors.Errorf("parent is nil")
	}

	removed, err := parent.RemoveSpace(space, true, updateDB)
	if err != nil {
		return false, errors.WithMessagef(err, "failed to remove space from parent: %s", parent.GetID())
	}
	if !removed {
		return false, nil
	}

	var errs *multierror.Error
	if space.GetEnabled() { // we need this check to avoid spam while removing children
		removeMsg := posbus.NewRemoveStaticObjectsMsg(1)
		removeMsg.SetObject(0, space.GetID())
		if err := space.GetWorld().Send(removeMsg.WebsocketMessage(), true); err != nil {
			errs = multierror.Append(errs, errors.WithMessage(err, "failed to send remove message"))
		}
	}

	if err := parent.UpdateChildrenPosition(true); err != nil {
		errs = multierror.Append(
			errs, errors.WithMessagef(err, "failed to update children position: %s", parent.GetID()),
		)
	}

	if err := space.Stop(); err != nil {
		errs = multierror.Append(errs, errors.WithMessage(err, "failed to stop space"))
	}
	space.SetEnabled(false)

	common.GetLogger().Infof("Helper: RemoveSpaceFromParent: space removed: %s", space.GetID())

	go func() {
		for _, child := range space.GetSpaces(false) {
			// prevent spam while removing
			child.SetEnabled(false)

			if _, err := RemoveSpaceFromParent(space, child, false); err != nil {
				common.GetLogger().Error(
					errors.WithMessagef(
						err, "Helper: RemoveSpaceFromParent: failed to remove child from space: %s", child.GetID(),
					),
				)
			}
		}
	}()

	return removed, errs.ErrorOrNil()
}

func CalcSpaceSpawnPosition(parentID, userID uuid.UUID) (*cmath.SpacePosition, error) {
	parent, ok := universe.GetNode().GetSpaceFromAllSpaces(parentID)
	if !ok {
		return nil, errors.Errorf("space parent not found: %s", parentID)
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
