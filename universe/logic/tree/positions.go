package tree

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe"
)

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
