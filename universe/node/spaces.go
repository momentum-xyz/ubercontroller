package node

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (n *Node) CreateObject(spaceID uuid.UUID) (universe.Object, error) {
	return nil, errors.Errorf("not permitted for node")
}

func (n *Node) SetParent(parent universe.Object, updateDB bool) error {
	return errors.Errorf("not permitted for node")
}

func (n *Node) GetAllObjects() map[uuid.UUID]universe.Object {
	spaces := map[uuid.UUID]universe.Object{
		n.GetID(): n.ToObject(),
	}

	for _, world := range n.GetWorlds().GetWorlds() {
		for spaceID, space := range world.GetAllObjects() {
			spaces[spaceID] = space
		}
	}

	return spaces
}

func (n *Node) FilterAllObjects(predicateFn universe.ObjectsFilterPredicateFn) map[uuid.UUID]universe.Object {
	spaces := make(map[uuid.UUID]universe.Object)
	if predicateFn(n.GetID(), n.ToObject()) {
		spaces[n.GetID()] = n.ToObject()
	}

	for _, world := range n.GetWorlds().GetWorlds() {
		for spaceID, space := range world.FilterAllObjects(predicateFn) {
			spaces[spaceID] = space
		}
	}

	return spaces
}

func (n *Node) GetObjectFromAllObjects(spaceID uuid.UUID) (universe.Object, bool) {
	if spaceID == n.GetID() {
		return n.ToObject(), true
	}

	world, ok := n.spaceIDToWorld.Load(spaceID)
	if !ok {
		return nil, false
	}

	return world.GetObjectFromAllObjects(spaceID)
}

func (n *Node) AddObjectToAllObjects(space universe.Object) error {
	if space.GetID() == n.GetID() {
		return errors.Errorf("not permitted for node")
	}

	world := space.GetWorld()
	if world == nil {
		return errors.Errorf("space has nil world")
	}

	if err := world.AddObjectToAllObjects(space); err != nil {
		return errors.WithMessage(err, "failed to add space to world all spaces")
	}

	n.spaceIDToWorld.Store(space.GetID(), world)

	return nil
}

func (n *Node) RemoveObjectFromAllObjects(space universe.Object) (bool, error) {
	if space.GetID() == n.GetID() {
		return false, errors.Errorf("not permitted for node")
	}

	n.spaceIDToWorld.Mu.Lock()
	defer n.spaceIDToWorld.Mu.Unlock()

	if world, ok := n.spaceIDToWorld.Data[space.GetID()]; ok {
		delete(n.spaceIDToWorld.Data, space.GetID())

		return world.RemoveObjectFromAllObjects(space)
	}

	return false, nil
}
