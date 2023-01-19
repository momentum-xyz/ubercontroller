package world

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) GetAllObjects() map[uuid.UUID]universe.Object {
	w.allSpaces.Mu.RLock()
	defer w.allSpaces.Mu.RUnlock()

	spaces := make(map[uuid.UUID]universe.Object, len(w.allSpaces.Data))

	for id, space := range w.allSpaces.Data {
		spaces[id] = space
	}

	return spaces
}

func (w *World) FilterAllObjects(predicateFn universe.ObjectsFilterPredicateFn) map[uuid.UUID]universe.Object {
	return w.allSpaces.Filter(predicateFn)
}

func (w *World) GetObjectFromAllObjects(spaceID uuid.UUID) (universe.Object, bool) {
	return w.allSpaces.Load(spaceID)
}

func (w *World) AddObjectToAllObjects(space universe.Object) error {
	if space.GetWorld().GetID() != w.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", space.GetWorld().GetID(), w.GetID())
	}

	w.allSpaces.Store(space.GetID(), space)

	return nil
}

func (w *World) RemoveObjectFromAllObjects(space universe.Object) (bool, error) {
	w.allSpaces.Mu.Lock()
	defer w.allSpaces.Mu.Unlock()

	if _, ok := w.allSpaces.Data[space.GetID()]; ok {
		delete(w.allSpaces.Data, space.GetID())

		return true, nil
	}

	return false, nil
}
