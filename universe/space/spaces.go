package space

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

const (
	chanIsClosed = -0x3FFFFFFFFFFFFFFF
)

func (s *Space) CreateSpace(spaceID uuid.UUID) (universe.Space, error) {
	space := NewSpace(spaceID, s.db, s.GetWorld())

	if err := space.Initialize(s.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize space: %s", spaceID)
	}
	if err := s.AddSpace(space, false); err != nil {
		return nil, errors.WithMessagef(err, "failed to add space %s to space: %s", spaceID, s.GetID())
	}

	return space, nil
}

func (s *Space) FilterSpaces(
	predicateFn universe.SpacesFilterPredicateFn, recursive bool,
) map[uuid.UUID]universe.Space {
	spaces := s.Children.Filter(predicateFn)

	if !recursive {
		return spaces
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		for id, space := range child.FilterSpaces(predicateFn, recursive) {
			spaces[id] = space
		}
	}

	return spaces
}

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
	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	spaces := make(map[uuid.UUID]universe.Space, len(s.Children.Data))

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

	if space == s {
		return errors.Errorf("space can't be part of itself")
	} else if space.GetWorld().GetID() != s.GetWorld().GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", space.GetWorld().GetID(), s.GetWorld().GetID())
	}

	if err := space.SetParent(s, false); err != nil {
		return errors.WithMessagef(err, "failed to set parent %s to space: %s", s.GetID(), space.GetID())
	}

	if updateDB {
		if err := s.db.SpacesUpsertSpace(s.ctx, space.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.Children.Data[space.GetID()] = space

	return universe.GetNode().AddSpaceToAllSpaces(space)
}

// TODO: think about rollaback on error
func (s *Space) AddSpaces(spaces []universe.Space, updateDB bool) error {
	s.Children.Mu.Lock()
	defer s.Children.Mu.Unlock()

	for i := range spaces {
		if spaces[i] == s {
			return errors.Errorf("space can't be part of itself")
		} else if spaces[i].GetWorld().GetID() != s.GetWorld().GetID() {
			return errors.Errorf(
				"space %s: worlds mismatch: %s != %s", spaces[i].GetID(), spaces[i].GetWorld().GetID(),
				s.GetWorld().GetID(),
			)
		}

		if err := spaces[i].SetParent(s, false); err != nil {
			return errors.WithMessagef(err, "failed to set parent %s to space: %s", s.GetID(), spaces[i].GetID())
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

	node := universe.GetNode()
	for i := range spaces {
		s.Children.Data[spaces[i].GetID()] = spaces[i]

		if err := node.AddSpaceToAllSpaces(spaces[i]); err != nil {
			return errors.WithMessagef(err, "failed to add space to all spaces: %s", spaces[i].GetID())
		}
	}

	return nil
}

// TODO: think about rollback on error
func (s *Space) RemoveSpace(space universe.Space, recursive, updateDB bool) (bool, error) {
	if space.GetWorld().GetID() != s.GetWorld().GetID() {
		return false, errors.Errorf("worlds mismatch: %s != %s", space.GetWorld().GetID(), s.GetWorld().GetID())
	}

	removed, err := func() (bool, error) {
		s.Children.Mu.Lock()
		defer s.Children.Mu.Unlock()

		if _, ok := s.Children.Data[space.GetID()]; !ok {
			return false, nil
		}

		if _, err := s.DoRemoveSpace(space, updateDB); err != nil {
			return false, errors.WithMessage(err, "failed to do remove space")
		}

		delete(s.Children.Data, space.GetID())

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
		removed, err := child.RemoveSpace(space, recursive, updateDB)
		if err != nil {
			return false, errors.WithMessagef(err, "failed to remove space: %s", space.GetID())
		}
		if removed {
			return true, nil
		}
	}

	return false, nil
}

func (s *Space) DoRemoveSpace(space universe.Space, updateDB bool) (bool, error) {
	if err := space.SetParent(nil, false); err != nil {
		return false, errors.WithMessage(err, "failed to set parent to nil")
	}

	if updateDB {
		if err := s.db.SpacesRemoveSpaceByID(s.ctx, space.GetID()); err != nil {
			return false, errors.WithMessage(err, "failed to update db")
		}
	}

	// we need this check to avoid spam while removing children
	if space.GetEnabled() {
		removeMsg := posbus.NewRemoveStaticObjectsMsg(1)
		removeMsg.SetObject(0, space.GetID())
		if err := space.GetWorld().Send(removeMsg.WebsocketMessage(), true); err != nil {
			s.log.Warn(
				errors.WithMessagef(
					err, "Helper: RemoveSpaceFromParent: failed to send remove message: %s", space.GetID(),
				),
			)
		}
	}

	return universe.GetNode().RemoveSpaceFromAllSpaces(space)
}

// RemoveSpaces return true in first value if all spaces with space ids was removed.
func (s *Space) RemoveSpaces(spaces []universe.Space, recursive, updateDB bool) (bool, error) {
	// TODO: optimize
	res := true

	var errs *multierror.Error
	for i := range spaces {
		removed, err := s.RemoveSpace(spaces[i], recursive, updateDB)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to remove space: %s", spaces[i].GetID()))
		}
		if !removed {
			res = false
		}
	}

	return res, errs.ErrorOrNil()
}
