package tree

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe"
)

// CalcOjbectSpawnPosition calculate the initial transform for a new object.
// src is an optional Transform to use for default values.
func CalcObjectSpawnPosition(parentID, userID umid.UMID, src *cmath.Transform) (*cmath.Transform, error) {
	parent, ok := universe.GetNode().GetObjectFromAllObjects(parentID)
	if !ok {
		return nil, errors.Errorf("object parent not found: %s", parentID)
	}

	var position *cmath.Transform
	effectiveOptions := parent.GetEffectiveOptions()
	if effectiveOptions == nil || len(effectiveOptions.ChildPlacements) == 0 {
		world := parent.GetWorld()
		if world != nil {
			user, ok := world.GetUser(userID, true)
			if ok {
				// TODO: allow caller (API endpoint) to optionally specifiy a position, rotation and scale.
				pos := user.GetPosition()
				rot := cmath.Vec3{}
				scl := cmath.Vec3{X: 1, Y: 1, Z: 1}
				if src != nil {
					rot = src.Rotation
					scl = src.Scale
				}
				position = &cmath.Transform{
					Position: pos,
					Rotation: rot,
					Scale:    scl,
				}
			}
		}
	}

	return position, nil
}
