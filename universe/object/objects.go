package object

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

const (
	chanIsClosed = -0x3FFFFFFFFFFFFFFF
)

func (s *Object) CreateObject(objectID uuid.UUID) (universe.Object, error) {
	object := NewObject(objectID, s.db, s.GetWorld())

	if err := object.Initialize(s.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize object: %s", objectID)
	}
	if err := s.AddObject(object, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add object %s to object: %s", objectID, s.GetID())
	}

	return object, nil
}

func (s *Object) FilterObjects(predicateFn universe.ObjectsFilterPredicateFn, recursive bool) map[uuid.UUID]universe.Object {
	objects := s.Children.Filter(predicateFn)

	if !recursive {
		return objects
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		for id, object := range child.FilterObjects(predicateFn, recursive) {
			objects[id] = object
		}
	}

	return objects
}

func (s *Object) GetObject(objectID uuid.UUID, recursive bool) (universe.Object, bool) {
	object, ok := s.Children.Load(objectID)
	if ok {
		return object, true
	}

	if !recursive {
		return nil, false
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		object, ok := child.GetObject(objectID, recursive)
		if ok {
			return object, true
		}
	}

	return nil, false
}

// GetObjects return map with all nested children if recursive is true,
// otherwise the method return map with children dependent only to current object.
func (s *Object) GetObjects(recursive bool) map[uuid.UUID]universe.Object {
	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	objects := make(map[uuid.UUID]universe.Object, len(s.Children.Data))

	for id, child := range s.Children.Data {
		objects[id] = child
		if !recursive {
			continue
		}

		for id, child := range child.GetObjects(recursive) {
			objects[id] = child
		}
	}

	return objects
}

func (s *Object) AddObject(object universe.Object, updateDB bool) error {
	s.Children.Mu.Lock()
	defer s.Children.Mu.Unlock()

	if object == s {
		return errors.Errorf("object can't be part of itself")
	} else if object.GetWorld().GetID() != s.GetWorld().GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", object.GetWorld().GetID(), s.GetWorld().GetID())
	}

	if err := object.SetParent(s, false); err != nil {
		return errors.WithMessagef(err, "failed to set parent %s to object: %s", s.GetID(), object.GetID())
	}

	if updateDB {
		if err := s.db.GetObjectsDB().UpsertObject(s.ctx, object.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.Children.Data[object.GetID()] = object

	return universe.GetNode().AddObjectToAllObjects(object)
}

// TODO: think about rollaback on error
func (s *Object) AddObjects(objects []universe.Object, updateDB bool) error {
	s.Children.Mu.Lock()
	defer s.Children.Mu.Unlock()

	for i := range objects {
		if objects[i] == s {
			return errors.Errorf("object can't be part of itself")
		} else if objects[i].GetWorld().GetID() != s.GetWorld().GetID() {
			return errors.Errorf(
				"object %s: worlds mismatch: %s != %s", objects[i].GetID(), objects[i].GetWorld().GetID(),
				s.GetWorld().GetID(),
			)
		}

		if err := objects[i].SetParent(s, false); err != nil {
			return errors.WithMessagef(err, "failed to set parent %s to object: %s", s.GetID(), objects[i].GetID())
		}
	}

	if updateDB {
		entries := make([]*entry.Object, len(objects))
		for i := range objects {
			entries[i] = objects[i].GetEntry()
		}
		if err := s.db.GetObjectsDB().UpsertObjects(s.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	node := universe.GetNode()
	for i := range objects {
		s.Children.Data[objects[i].GetID()] = objects[i]

		if err := node.AddObjectToAllObjects(objects[i]); err != nil {
			return errors.WithMessagef(err, "failed to add object to all objects: %s", objects[i].GetID())
		}
	}

	return nil
}

// TODO: think about rollback on error
func (s *Object) RemoveObject(object universe.Object, recursive, updateDB bool) (bool, error) {
	if object.GetWorld().GetID() != s.GetWorld().GetID() {
		return false, errors.Errorf("worlds mismatch: %s != %s", object.GetWorld().GetID(), s.GetWorld().GetID())
	}

	removed, err := func() (bool, error) {
		s.Children.Mu.Lock()
		defer s.Children.Mu.Unlock()

		if _, ok := s.Children.Data[object.GetID()]; !ok {
			return false, nil
		}

		if _, err := s.DoRemoveObject(object, updateDB); err != nil {
			return false, errors.WithMessage(err, "failed to do remove object")
		}

		delete(s.Children.Data, object.GetID())

		return true, nil
	}()
	if err != nil {
		return false, err
	}
	if removed {
		return removed, nil
	}

	if !recursive {
		return false, nil
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		removed, err := child.RemoveObject(object, recursive, updateDB)
		if err != nil {
			return false, errors.WithMessagef(err, "failed to remove object: %s", object.GetID())
		}
		if removed {
			return true, nil
		}
	}

	return false, nil
}

func (s *Object) DoRemoveObject(object universe.Object, updateDB bool) (bool, error) {
	if err := object.SetParent(nil, false); err != nil {
		return false, errors.WithMessage(err, "failed to set parent to nil")
	}

	if updateDB {
		if err := s.db.GetObjectsDB().RemoveObjectByID(s.ctx, object.GetID()); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	return universe.GetNode().RemoveObjectFromAllObjects(object)
}

// RemoveObjects return true in first value if all objects with object ids was removed.
// TODO: rethink and optimize/reimplement
func (s *Object) RemoveObjects(objects []universe.Object, recursive, updateDB bool) (bool, error) {
	res := true

	var errs *multierror.Error
	for i := range objects {
		removed, err := s.RemoveObject(objects[i], recursive, updateDB)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to remove object: %s", objects[i].GetID()))
		}
		if !removed {
			res = false
		}
	}

	return res, errs.ErrorOrNil()
}
