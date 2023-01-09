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
	SpaceTemplate
}

func AddWorldFromTemplate(worldTemplate *WorldTemplate, updateDB bool) (uuid.UUID, error) {
	node := universe.GetNode()

	// loading
	worldSpaceType, ok := node.GetSpaceTypes().GetSpaceType(worldTemplate.SpaceTypeID)
	if !ok {
		return uuid.Nil, errors.Errorf("failed to get world space type: %s", worldTemplate.SpaceTypeID)
	}

	worldID := worldTemplate.SpaceID
	if worldID == nil {
		worldID = utils.GetPTR(uuid.New())
	}
	worldName := worldTemplate.SpaceName
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
	if err := world.SetSpaceType(worldSpaceType, false); err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to set space type: %s", worldTemplate.SpaceTypeID)
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
	spaceLabelToID := make(map[string]uuid.UUID)
	for i := range worldTemplate.Spaces {
		worldTemplate.Spaces[i].ParentID = *worldID
		spaceID, err := AddSpaceFromTemplate(worldTemplate.Spaces[i], updateDB)
		if err != nil {
			return uuid.Nil, errors.WithMessagef(
				err, "failed to add space from template: %+v", worldTemplate.Spaces[i],
			)
		}

		if worldTemplate.Spaces[i].Label != nil {
			spaceLabelToID[*worldTemplate.Spaces[i].Label] = spaceID
		}
	}

	// enabling
	world.SetEnabled(true)

	// adding attributes
	if err := world.SetName(*worldName, true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to set world name")
	}

	worldTemplate.SpaceAttributes = append(
		worldTemplate.SpaceAttributes,
		entry.NewAttribute(
			entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.World.Settings.Name),
			entry.NewAttributePayload(
				&entry.AttributeValue{
					"kind":        "basic",
					"spaces":      spaceLabelToID,
					"attributes":  map[string]any{},
					"space_types": map[string]any{},
					"effects":     map[string]any{},
				},
				nil,
			),
		),
	)

	for i := range worldTemplate.SpaceAttributes {
		if _, err := world.GetSpaceAttributes().Upsert(
			worldTemplate.SpaceAttributes[i].AttributeID,
			modify.MergeWith(worldTemplate.SpaceAttributes[i].AttributePayload),
			updateDB,
		); err != nil {
			return uuid.Nil, errors.WithMessagef(
				err, "failed to upsert world space attribute: %+v", worldTemplate.SpaceAttributes[i],
			)
		}
	}

	// updating
	if err := world.Update(true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to update world")
	}

	return *worldID, nil
}
