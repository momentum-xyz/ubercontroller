package world

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) GetSpace(spaceID uuid.UUID, recursive bool) (universe.Space, bool) {
	return w.Space.GetSpace(spaceID, false)
}

func (w *World) GetSpaces(recursive bool) map[uuid.UUID]universe.Space {
	return w.Space.GetSpaces(false)
}

func (w *World) RemoveSpace(space universe.Space, recursive, updateDB bool) (bool, error) {
	return w.Space.RemoveSpace(space, false, updateDB)
}

func (w *World) RemoveSpaces(spaces []universe.Space, recursive, updateDB bool) (bool, error) {
	return w.Space.RemoveSpaces(spaces, false, updateDB)
}
