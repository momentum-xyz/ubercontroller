package world

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) FilterAllSpaces(predicateFn universe.SpacesFilterPredicateFn) map[uuid.UUID]universe.Space {
	return w.allSpaces.Filter(predicateFn)
}

func (w *World) GetSpaceFromAllSpaces(spaceID uuid.UUID) (universe.Space, bool) {
	return w.allSpaces.Load(spaceID)
}

func (w *World) AddSpaceToAllSpaces(space universe.Space) error {
	if space.GetWorld().GetID() != w.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", space.GetWorld().GetID(), w.GetID())
	}

	w.allSpaces.Store(space.GetID(), space)

	return nil
}

func (w *World) RemoveSpaceFromAllSpaces(space universe.Space) (bool, error) {
	w.allSpaces.Mu.Lock()
	defer w.allSpaces.Mu.Unlock()

	if _, ok := w.allSpaces.Data[space.GetID()]; ok {
		delete(w.allSpaces.Data, space.GetID())

		return true, nil
	}

	return false, nil
}
