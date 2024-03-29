package world

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) GetAllObjects() map[umid.UMID]universe.Object {
	w.allObjects.Mu.RLock()
	defer w.allObjects.Mu.RUnlock()

	objects := make(map[umid.UMID]universe.Object, len(w.allObjects.Data))
	for id, object := range w.allObjects.Data {
		objects[id] = object
	}

	return objects
}

func (w *World) FilterAllObjects(predicateFn universe.ObjectsFilterPredicateFn) map[umid.UMID]universe.Object {
	return w.allObjects.Filter(predicateFn)
}

func (w *World) GetObjectFromAllObjects(objectID umid.UMID) (universe.Object, bool) {
	return w.allObjects.Load(objectID)
}

func (w *World) AddObjectToAllObjects(object universe.Object) error {
	if object.GetWorld().GetID() != w.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", object.GetWorld().GetID(), w.GetID())
	}

	w.allObjects.Store(object.GetID(), object)

	return nil
}

func (w *World) RemoveObjectFromAllObjects(object universe.Object) (bool, error) {
	w.allObjects.Mu.Lock()
	defer w.allObjects.Mu.Unlock()

	if _, ok := w.allObjects.Data[object.GetID()]; ok {
		delete(w.allObjects.Data, object.GetID())

		return true, nil
	}

	return false, nil
}
