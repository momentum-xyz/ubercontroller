package space

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func (s *Space) GetSpace(spaceID uuid.UUID, recursive bool) (universe.Space, bool) {
	space, ok := s.Children.Load(spaceID)
	if ok {
		return space, true
	}

	if !recursive {
		return nil, false
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		space, ok := child.GetSpace(spaceID, recursive)
		if ok {
			return space, true
		}
	}

	return nil, false
}

// GetSpaces return map with all nested children if recursive is true,
// otherwise the method return map with children dependent only to current space.
func (s *Space) GetSpaces(recursive bool) map[uuid.UUID]universe.Space {
	spaces := make(map[uuid.UUID]universe.Space)

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for id, child := range s.Children.Data {
		spaces[id] = child
		if !recursive {
			continue
		}

		for id, child := range child.GetSpaces(recursive) {
			spaces[id] = child
		}
	}

	return spaces
}

func (s *Space) AddSpace(space universe.Space, updateDB bool) error {
	s.Children.Mu.Lock()
	defer s.Children.Mu.Unlock()

	if err := space.SetParent(s, false); err != nil {
		return errors.WithMessagef(err, "failed to set parent to space: %s", space.GetID())
	}

	if updateDB {
		if err := s.db.SpacesUpsertSpace(s.ctx, space.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.Children.Data[space.GetID()] = space

	return nil
}

func (s *Space) AddSpaces(spaces []universe.Space, updateDB bool) error {
	s.Children.Mu.Lock()
	defer s.Children.Mu.Unlock()

	for i := range spaces {
		if err := spaces[i].SetParent(s, false); err != nil {
			return errors.WithMessagef(err, "failed to set parent to space: %s", spaces[i].GetID())
		}
	}

	if updateDB {
		entries := make([]*entry.Space, len(spaces))
		for i := range spaces {
			entries[i] = spaces[i].GetEntry()
		}
		if err := s.db.SpacesUpsertSpaces(s.ctx, entries); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	for i := range spaces {
		s.Children.Data[spaces[i].GetID()] = spaces[i]
	}

	return nil
}

func (s *Space) RemoveSpace(space universe.Space, recursive, updateDB bool) (bool, error) {
	s.Children.Mu.Lock()
	if _, ok := s.Children.Data[space.GetID()]; ok {
		defer s.Children.Mu.Unlock()

		// TODO: move to morgue
		if err := space.SetParent(nil, false); err != nil {
			return false, errors.WithMessage(err, "failed to set parent")
		}

		if updateDB {
			if err := s.db.SpacesRemoveSpaceByID(s.ctx, space.GetID()); err != nil {
				return false, errors.WithMessage(err, "failed to remove space by id")
			}
		}

		delete(s.Children.Data, space.GetID())

		return true, nil
	}
	s.Children.Mu.Unlock()

	if !recursive {
		return false, nil
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		removed, err := child.RemoveSpace(space, recursive, updateDB)
		if err != nil {
			return false, errors.WithMessage(err, "failed to remove space")
		}
		if removed {
			return true, nil
		}
	}

	return false, nil
}

// RemoveSpaces return true in first value if any of spaces with space ids was removed.
func (s *Space) RemoveSpaces(spaces []universe.Space, recursive, updateDB bool) (bool, error) {
	var res bool
	group, _ := errgroup.WithContext(s.ctx)

	for _, space := range spaces {
		space := space

		group.Go(func() error {
			removed, err := s.RemoveSpace(space, recursive, updateDB)
			if err != nil {
				return errors.WithMessagef(err, "failed to remove space: %s", space.GetID())
			}
			if removed {
				res = true
			}
			return nil
		})
	}

	return res, group.Wait()
}
