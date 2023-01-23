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

func AddWorldFromTemplate(worldTemplate *WorldTemplate, updateDB bool) (uuid.UUID, error) {
	node := universe.GetNode()

	// loading
	worldObjectType, ok := node.GetObjectTypes().GetObjectType(worldTemplate.ObjectTypeID)
	if !ok {
		return uuid.Nil, errors.Errorf("failed to get world object type: %s", worldTemplate.ObjectTypeID)
	}

	worldID := worldTemplate.ObjectID
	if worldID == nil {
		worldID = utils.GetPTR(uuid.New())
	}
	worldName := worldTemplate.ObjectName
	if worldName == nil {
		worldName = utils.GetPTR(worldID.String())
	}

	// creating
	world, err := node.GetWorlds().CreateWorld(*worldID)
	if err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to create world: %s", worldID)
	}

	if err := world.SetOwnerID(*worldTemplate.OwnerID, false); err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to set owner: %s", worldTemplate.OwnerID)
	}
	if err := world.SetObjectType(worldObjectType, false); err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to set space type: %s", worldTemplate.ObjectTypeID)
	}

	// saving in database
	if updateDB {
		if err := node.GetWorlds().AddWorld(world, updateDB); err != nil {
			return uuid.Nil, errors.WithMessage(err, "failed to add world")
		}
	}

	// running
	if err := world.Run(); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to run world")
	}

	// adding children
	objectLabelToID := make(map[string]uuid.UUID)
	for i := range worldTemplate.Children {
		worldTemplate.Children[i].ParentID = *worldID
		objectID, err := AddObjectFromTemplate(worldTemplate.Children[i], updateDB)
		if err != nil {
			return uuid.Nil, errors.WithMessagef(
				err, "failed to add object from template: %+v", worldTemplate.Children[i],
			)
		}

		if worldTemplate.Children[i].Label != nil {
			objectLabelToID[*worldTemplate.Children[i].Label] = objectID
		}
	}

	// enabling
	world.SetEnabled(true)

	// adding attributes
	if err := world.SetName(*worldName, true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to set world name")
	}

	worldTemplate.ObjectAttributes = append(
		worldTemplate.ObjectAttributes,
		entry.NewAttribute(
			entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.World.Settings.Name),
			entry.NewAttributePayload(
				&entry.AttributeValue{
					"kind":         "basic",
					"objects":      objectLabelToID,
					"attributes":   map[string]any{},
					"object_types": map[string]any{},
					"effects":      map[string]any{},
				},
				nil,
			),
		),
	)

	for i := range worldTemplate.ObjectAttributes {
		if _, err := world.GetObjectAttributes().Upsert(
			worldTemplate.ObjectAttributes[i].AttributeID,
			modify.MergeWith(worldTemplate.ObjectAttributes[i].AttributePayload),
			updateDB,
		); err != nil {
			return uuid.Nil, errors.WithMessagef(
				err, "failed to upsert world space attribute: %+v", worldTemplate.ObjectAttributes[i],
			)
		}
	}

	// updating
	if err := world.Update(true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to update world")
	}

	return *worldID, nil
}
