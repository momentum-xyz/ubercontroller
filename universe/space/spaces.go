package space

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/generics"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func (s *Space) GetSpace(spaceID uuid.UUID, recursive bool) (universe.Space, bool) {
	space, ok := s.children.Load(spaceID)
	if ok {
		return space, true
	}

	if !recursive {
		return nil, false
	}

	s.children.Mu.RLock()
	defer s.children.Mu.RUnlock()

	for _, child := range s.children.Data {
		space, ok := child.GetSpace(spaceID, recursive)
		if ok {
			return space, true
		}
	}

	return nil, false
}

// GetSpaces return new sync map with all nested children if recursive is true,
// otherwise the method return existing sync map with children dependent only to current space.
func (s *Space) GetSpaces(recursive bool) *generics.SyncMap[uuid.UUID, universe.Space] {
	if !recursive {
		return s.children
	}

	spaces := generics.NewSyncMap[uuid.UUID, universe.Space]()

	s.children.Mu.RLock()
	defer s.children.Mu.RUnlock()

	// maybe we will need lock here in future
	for id, child := range s.children.Data {
		spaces.Data[id] = child

		for id, child := range child.GetSpaces(recursive).Data {
			spaces.Data[id] = child
		}
	}

	return spaces
}

func (s *Space) AddSpace(space universe.Space, updateDB bool) error {
	s.children.Mu.Lock()
	defer s.children.Mu.Unlock()

	if err := space.SetParent(s, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set parent to space: %s", space.GetID())
	}
	s.children.Data[space.GetID()] = space

	return nil
}

func (s *Space) AddSpaces(spaces []universe.Space, updateDB bool) error {
	var errs *multierror.Error
	for i := range spaces {
		if err := s.AddSpace(spaces[i], updateDB); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to add space: %s", spaces[i].GetID()))
		}
	}
	return errs.ErrorOrNil()
}

func (s *Space) RemoveSpace(spaceID uuid.UUID, recursive, updateDB bool) (bool, error) {
	s.children.Mu.Lock()
	space, ok := s.children.Data[spaceID]
	if ok {
		defer s.children.Mu.Unlock()

		if err := space.SetParent(nil, updateDB); err != nil {
			return false, errors.WithMessagef(err, "failed to set parent to space: %s", spaceID)
		}
		delete(s.children.Data, spaceID)

		return true, nil
	}
	s.children.Mu.Unlock()

	if !recursive {
		return true, nil
	}

	s.children.Mu.RLock()
	defer s.children.Mu.RUnlock()

	for _, child := range s.children.Data {
		removed, err := child.RemoveSpace(spaceID, recursive, updateDB)
		if err != nil {
			return false, errors.WithMessagef(err, "failed to remove space: %s", spaceID)
		}
		if removed {
			return true, nil
		}
	}

	return false, nil
}

// RemoveSpaces return true in first value if any of spaces with space ids was removed.
func (s *Space) RemoveSpaces(spaceIDs []uuid.UUID, recursive, updateDB bool) (bool, error) {
	var res bool
	for i := range spaceIDs {
		removed, err := s.RemoveSpace(spaceIDs[i], recursive, updateDB)
		if err != nil {
			return false, errors.WithMessagef(err, "failed to remove space: %s", spaceIDs[i])
		}
		if removed {
			res = true
		}
	}
	return res, nil
}
