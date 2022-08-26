package space

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/controller/types/generics"
	"github.com/momentum-xyz/controller/universe"
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

// GetSpaces returns new sync map with all nested children if recursive is true,
// otherwise the method returns existing sync map with children dependent only to current space.
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

func (s *Space) AttachSpace(space universe.Space, updateDB bool) error {
	s.children.Mu.Lock()
	defer s.children.Mu.Unlock()

	if err := space.SetParent(s, updateDB); err != nil {
		return errors.WithMessagef(err, "failed to set parent to space: %s", space.GetID())
	}
	s.children.Data[space.GetID()] = space

	return nil
}

func (s *Space) AttachSpaces(spaces []universe.Space, updateDB bool) error {
	var errs *multierror.Error
	for i := range spaces {
		if err := s.AttachSpace(spaces[i], updateDB); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to attach space: %s", spaces[i].GetID()))
		}
	}
	return errs.ErrorOrNil()
}

func (s *Space) DetachSpace(spaceID uuid.UUID, recursive, updateDB bool) (bool, error) {
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
		detached, err := child.DetachSpace(spaceID, recursive, updateDB)
		if err != nil {
			return false, errors.WithMessagef(err, "failed to detach space: %s", spaceID)
		}
		if detached {
			return true, nil
		}
	}

	return false, nil
}

// DetachSpaces returns true in first value if any of spaces with space ids was detached.
func (s *Space) DetachSpaces(spaceIDs []uuid.UUID, recursive, updateDB bool) (bool, error) {
	var res bool
	for i := range spaceIDs {
		detached, err := s.DetachSpace(spaceIDs[i], recursive, updateDB)
		if err != nil {
			return false, errors.WithMessagef(err, "failed to detach space: %s", spaceIDs[i])
		}
		if detached {
			res = true
		}
	}
	return res, nil
}
