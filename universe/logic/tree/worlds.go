package tree

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type WorldTemplate struct {
	ObjectTemplate
}

func addWorldFromTemplate(worldTemplate *WorldTemplate, updateDB bool) (umid.UMID, error) {
	node := universe.GetNode()

	// loading
	worldObjectType, ok := node.GetObjectTypes().GetObjectType(worldTemplate.ObjectTypeID)
	if !ok {
		return umid.Nil, errors.Errorf("failed to get world object type: %s", worldTemplate.ObjectTypeID)
	}

	worldID := worldTemplate.ObjectID
	if worldID == nil {
		worldID = utils.GetPTR(umid.New())
	}
	worldName := worldTemplate.ObjectName
	if worldName == nil {
		worldName = utils.GetPTR(worldID.String())
	}

	// creating
	world, err := node.GetWorlds().CreateWorld(*worldID)
	if err != nil {
		return umid.Nil, errors.WithMessagef(err, "failed to create world: %s", worldID)
	}

	if err := world.SetOwnerID(*worldTemplate.OwnerID, false); err != nil {
		return umid.Nil, errors.WithMessagef(err, "failed to set owner: %s", worldTemplate.OwnerID)
	}
	if err := world.SetObjectType(worldObjectType, false); err != nil {
		return umid.Nil, errors.WithMessagef(err, "failed to set object type: %s", worldTemplate.ObjectTypeID)
	}

	// saving in database
	if updateDB {
		if err := node.GetWorlds().AddWorld(world, updateDB); err != nil {
			return umid.Nil, errors.WithMessage(err, "failed to add world")
		}
	}

	// running
	if err := world.Run(); err != nil {
		return umid.Nil, errors.WithMessage(err, "failed to run world")
	}

	// adding children
	objectLabelToID := make(map[string]umid.UMID)
	if len(worldTemplate.RandomObjects) > 0 {
		randomObject := worldTemplate.RandomObjects[rand.Intn(len(worldTemplate.RandomObjects))]
		worldTemplate.Objects = append(worldTemplate.Objects, randomObject)
		for i := range worldTemplate.Objects {
			worldTemplate.Objects[i].ParentID = *worldID
			objectID, err := AddObjectFromTemplate(worldTemplate.Objects[i], updateDB)
			if err != nil {
				return umid.Nil, errors.WithMessagef(
					err, "failed to add object from template: %+v", worldTemplate.Objects[i],
				)
			}

			if worldTemplate.Objects[i].Label != nil {
				objectLabelToID[*worldTemplate.Objects[i].Label] = objectID
			}
		}
	}

	// enabling
	world.SetEnabled(true)

	// adding attributes
	if err := world.SetName(*worldName, true); err != nil {
		return umid.Nil, errors.WithMessage(err, "failed to set world name")
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
			return umid.Nil, errors.WithMessagef(
				err, "failed to upsert world object attribute: %+v", worldTemplate.ObjectAttributes[i],
			)
		}
	}

	// updating
	if err := world.Update(true); err != nil {
		return umid.Nil, errors.WithMessage(err, "failed to update world")
	}

	return *worldID, nil
}

// This func wraps create world function with add activity
func AddWorldFromTemplate(worldTemplate *WorldTemplate, updateDB bool) (umid.UMID, error) {
	id, err := addWorldFromTemplate(worldTemplate, updateDB)
	if id != umid.Nil && err == nil {
		err := addNewWorldCreatedActivity(id)
		if err != nil {
			return id, errors.WithMessage(err, "failed to AddNewOdysseyActivity")
		}
	}

	return id, err
}

func addNewWorldCreatedActivity(id umid.UMID) error {
	node := universe.GetNode()
	a, err := node.GetActivities().CreateActivity(umid.New())
	if err != nil {
		return errors.WithMessage(err, "failed to CreateActivity")
	}

	world, ok := node.GetWorlds().GetWorld(id)
	if !ok {
		return errors.New("world not found by id:" + id.String())
	}

	if err := a.SetObjectID(id, true); err != nil {
		return errors.WithMessage(err, "failed to set object ID")
	}

	if err := a.SetUserID(world.GetOwnerID(), true); err != nil {
		return errors.WithMessage(err, "failed to set user ID")
	}

	if err := a.SetCreatedAt(time.Now(), true); err != nil {
		return errors.WithMessage(err, "failed to set created_at")
	}

	aType := entry.ActivityWorldCreate
	if err := a.SetType(&aType, true); err != nil {
		return errors.WithMessage(err, "failed to set activity type")
	}

	modifyFn := func(current *entry.ActivityData) (*entry.ActivityData, error) {
		if current == nil {
			current = &entry.ActivityData{}
		}

		//current.Position = &position
		//current.Hash = &inBody.Hash
		//current.Description = &inBody.Description

		return current, nil
	}

	_, err = a.SetData(modifyFn, true)
	if err != nil {
		return errors.New("failed to set activity data")
	}

	if err := a.GetActivities().Inject(a); err != nil {
		return errors.New("failed to inject activity")
	}

	return nil
}
