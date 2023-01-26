package helper

import (
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/pkg/errors"
)

type WorldTemplate struct {
	ObjectTemplate
}

func AddWorldFromTemplate(worldTemplate *WorldTemplate, updateDB bool) (universe.World, error) {
	node := universe.GetNode()

	// loading
	worldObjectType, ok := node.GetObjectTypes().GetObjectType(worldTemplate.ObjectTypeID)
	if !ok {
		return nil, errors.Errorf("failed to get world object type: %s", worldTemplate.ObjectTypeID)
	}

	worldID := worldTemplate.ObjectID
	worldName := worldTemplate.ObjectName
	if worldID == nil {
		worldID = utils.GetPTR(uuid.New())
	}
	if worldName == nil {
		worldName = utils.GetPTR(worldID.String())
	}

	// creating
	world, err := node.GetWorlds().CreateWorld(*worldID)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create world: %s", worldID)
	}

	if err := world.SetOwnerID(*worldTemplate.OwnerID, false); err != nil {
		return world, errors.WithMessagef(err, "failed to set owner: %s", worldTemplate.OwnerID)
	}
	if err := world.SetObjectType(worldObjectType, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to set object type: %s", worldTemplate.ObjectTypeID)
	}

	// running
	if err := world.Run(); err != nil {
		return world, errors.WithMessage(err, "failed to run world")
	}

	// adding children
	labelToChild := make(map[string]universe.Object)
	for _, childTemplate := range worldTemplate.Children {
		child, err := createObjectFromTemplate(world.ToObject(), childTemplate)
		if err != nil {
			return world, errors.WithMessagef(err, "failed to create child from template: %+v", childTemplate)
		}

		if childTemplate.Label != nil {
			labelToChild[*childTemplate.Label] = child
		}
	}

	// enabling
	world.SetEnabled(true)

	// adding attributes
	if err := world.SetName(*worldName, false); err != nil {
		return world, errors.WithMessage(err, "failed to set world name")
	}

	worldTemplate.ObjectAttributes = append(
		worldTemplate.ObjectAttributes,
		entry.NewAttribute(
			entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.World.Settings.Name),
			entry.NewAttributePayload(
				&entry.AttributeValue{
					"kind":         "basic",
					"objects":      labelToChild,
					"attributes":   map[string]any{},
					"object_types": map[string]any{},
					"effects":      map[string]any{},
				},
				nil,
			),
		),
	)

	for _, attribute := range worldTemplate.ObjectAttributes {
		if _, err := world.GetObjectAttributes().Upsert(
			attribute.AttributeID, modify.MergeWith(attribute.AttributePayload), false,
		); err != nil {
			return world, errors.WithMessagef(err, "failed to upsert world attribute: %+v", attribute.AttributeID)
		}
	}

	if updateDB {
		if err := node.GetWorlds().AddWorld(world, true); err != nil {
			return world, errors.WithMessagef(err, "failed to add to worlds: %s", world.GetID())
		}
	}

	// updating
	if err := world.Update(true); err != nil {
		return world, errors.WithMessage(err, "failed to update world")
	}
	if err := world.UpdateChildrenPosition(true); err != nil {
		return world, errors.WithMessage(err, "failed to update children position")
	}

	return world, nil
}
