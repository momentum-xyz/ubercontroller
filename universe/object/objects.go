package object

import (
	"github.com/hashicorp/go-multierror"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

const (
	chanIsClosed = -0x3FFFFFFFFFFFFFFF
)

func (o *Object) CreateObject(objectID mid.ID) (universe.Object, error) {
	object := NewObject(objectID, o.db, o.GetWorld())

	if err := object.Initialize(o.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize object: %s", objectID)
	}
	if err := o.AddObject(object, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add object %s to object: %s", objectID, o.GetID())
	}

	return object, nil
}

func (o *Object) FilterObjects(
	predicateFn universe.ObjectsFilterPredicateFn, recursive bool,
) map[mid.ID]universe.Object {
	objects := o.Children.Filter(predicateFn)

	if !recursive {
		return objects
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		for id, object := range child.FilterObjects(predicateFn, true) {
			objects[id] = object
		}
	}

	return objects
}

func (o *Object) GetObject(objectID mid.ID, recursive bool) (universe.Object, bool) {
	object, ok := o.Children.Load(objectID)
	if ok {
		return object, true
	}

	if !recursive {
		return nil, false
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		object, ok := child.GetObject(objectID, true)
		if ok {
			return object, true
		}
	}

	return nil, false
}

// GetObjects return map with all nested children if recursive is true,
// otherwise the method return map with children dependent only to current object.
func (o *Object) GetObjects(recursive bool) map[mid.ID]universe.Object {
	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	objects := make(map[mid.ID]universe.Object, len(o.Children.Data))
	for id, child := range o.Children.Data {
		objects[id] = child

		if !recursive {
			continue
		}

		for id, child := range child.GetObjects(true) {
			objects[id] = child
		}
	}

	return objects
}

// TODO: think about rollaback
func (o *Object) AddObject(object universe.Object, updateDB bool) error {
	o.Children.Mu.Lock()
	defer o.Children.Mu.Unlock()

	if object == o {
		return errors.Errorf("object can't be part of itself")
	} else if object.GetWorld().GetID() != o.GetWorld().GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", object.GetWorld().GetID(), o.GetWorld().GetID())
	}

	if err := object.SetParent(o, false); err != nil {
		return errors.WithMessagef(err, "failed to set parent %s to object %s", o.GetID(), object.GetID())
	}

	if updateDB {
		if err := object.Save(); err != nil {
			return errors.WithMessage(err, "failed to save object")
		}
	}

	o.Children.Data[object.GetID()] = object

	return universe.GetNode().AddObjectToAllObjects(object)
}

// TODO: optimize
func (o *Object) AddObjects(objects []universe.Object, updateDB bool) error {
	for _, object := range objects {
		if err := o.AddObject(object, updateDB); err != nil {
			return errors.WithMessagef(err, "failed to add object: %s", object.GetID())
		}
	}
	return nil
}

// TODO: think about rollback on error
func (o *Object) RemoveObject(object universe.Object, recursive, updateDB bool) (bool, error) {
	if object.GetWorld().GetID() != o.GetWorld().GetID() {
		return false, errors.Errorf("worlds mismatch: %s != %s", object.GetWorld().GetID(), o.GetWorld().GetID())
	}

	removed, err := func() (bool, error) {
		o.Children.Mu.Lock()
		defer o.Children.Mu.Unlock()

		if _, ok := o.Children.Data[object.GetID()]; !ok {
			return false, nil
		}

		if _, err := o.DoRemoveObject(object, updateDB); err != nil {
			return false, errors.WithMessage(err, "failed to do remove object")
		}

		delete(o.Children.Data, object.GetID())

		return true, nil
	}()
	if err != nil {
		return false, err
	}
	if removed {
		return true, nil
	}

	if !recursive {
		return false, nil
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		removed, err := child.RemoveObject(object, true, updateDB)
		if err != nil {
			return false, errors.WithMessagef(err, "failed to remove object: %s", object.GetID())
		}
		if removed {
			return true, nil
		}
	}

	return false, nil
}

func (o *Object) DoRemoveObject(object universe.Object, updateDB bool) (bool, error) {
	if err := object.SetParent(nil, false); err != nil {
		return false, errors.WithMessage(err, "failed to set parent to nil")
	}

	if updateDB {
		if err := o.db.GetObjectsDB().RemoveObjectByID(o.ctx, object.GetID()); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	return universe.GetNode().RemoveObjectFromAllObjects(object)
}

// TODO: rethink and optimize/reimplement
func (o *Object) RemoveObjects(objects []universe.Object, recursive, updateDB bool) (bool, error) {
	var res bool
	var errs *multierror.Error
	for _, object := range objects {
		removed, err := o.RemoveObject(object, recursive, updateDB)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to remove object: %s", object.GetID()))
			continue
		}
		if removed {
			res = true
		}
	}

	return res, errs.ErrorOrNil()
}
