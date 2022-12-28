package node

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (n *Node) CreateSpace(spaceID uuid.UUID) (universe.Space, error) {
	return nil, errors.Errorf("not permitted for node")
}

func (n *Node) SetParent(parent universe.Space, updateDB bool) error {
	return errors.Errorf("not permitted for node")
}

func (n *Node) GetAllSpaces() map[uuid.UUID]universe.Space {
	spaces := map[uuid.UUID]universe.Space{
		n.GetID(): n,
	}

	for _, world := range n.GetWorlds().GetWorlds() {
		for spaceID, space := range world.GetAllSpaces() {
			spaces[spaceID] = space
		}
	}

	return spaces
}

func (n *Node) FilterAllSpaces(predicateFn universe.SpacesFilterPredicateFn) map[uuid.UUID]universe.Space {
	spaces := make(map[uuid.UUID]universe.Space)
	if predicateFn(n.GetID(), n) {
		spaces[n.GetID()] = n
	}

	for _, world := range n.GetWorlds().GetWorlds() {
		for spaceID, space := range world.FilterAllSpaces(predicateFn) {
			spaces[spaceID] = space
		}
	}

	return spaces
}

func (n *Node) GetSpaceFromAllSpaces(spaceID uuid.UUID) (universe.Space, bool) {
	if spaceID == n.GetID() {
		return n, true
	}

	world, ok := n.spaceIDToWorld.Load(spaceID)
	if !ok {
		return nil, false
	}

	return world.GetSpaceFromAllSpaces(spaceID)
}

func (n *Node) AddSpaceToAllSpaces(space universe.Space) error {
	if space.GetID() == n.GetID() {
		return errors.Errorf("not permitted for node")
	}

	world := space.GetWorld()
	if world == nil {
		return errors.Errorf("space has nil world")
	}

	if err := world.AddSpaceToAllSpaces(space); err != nil {
		return errors.WithMessage(err, "failed to add space to world all spaces")
	}

	n.spaceIDToWorld.Store(space.GetID(), world)

	return nil
}

func (n *Node) RemoveSpaceFromAllSpaces(space universe.Space) (bool, error) {
	if space.GetID() == n.GetID() {
		return false, errors.Errorf("not permitted for node")
	}

	n.spaceIDToWorld.Mu.Lock()
	defer n.spaceIDToWorld.Mu.Unlock()

	if world, ok := n.spaceIDToWorld.Data[space.GetID()]; ok {
		delete(n.spaceIDToWorld.Data, space.GetID())

		return world.RemoveSpaceFromAllSpaces(space)
	}

	return false, nil
}
