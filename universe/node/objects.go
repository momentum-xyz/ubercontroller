package node

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func (n *Node) CreateObject(objectID umid.UMID) (universe.Object, error) {
	return nil, errors.Errorf("not permitted for node")
}

func (n *Node) SetParent(parent universe.Object, updateDB bool) error {
	return errors.Errorf("not permitted for node")
}

func (n *Node) GetAllObjects() map[umid.UMID]universe.Object {
	objects := map[umid.UMID]universe.Object{
		n.GetID(): n.ToObject(),
	}

	for _, world := range n.GetWorlds().GetWorlds() {
		for objectID, object := range world.GetAllObjects() {
			objects[objectID] = object
		}
	}

	return objects
}

func (n *Node) FilterAllObjects(predicateFn universe.ObjectsFilterPredicateFn) map[umid.UMID]universe.Object {
	objects := make(map[umid.UMID]universe.Object)
	if predicateFn(n.GetID(), n.ToObject()) {
		objects[n.GetID()] = n.ToObject()
	}

	for _, world := range n.GetWorlds().GetWorlds() {
		for objectID, object := range world.FilterAllObjects(predicateFn) {
			objects[objectID] = object
		}
	}

	return objects
}

func (n *Node) GetObjectFromAllObjects(objectID umid.UMID) (universe.Object, bool) {
	if objectID == n.GetID() {
		return n.ToObject(), true
	}

	world, ok := n.objectIDToWorld.Load(objectID)
	if !ok {
		return nil, false
	}

	return world.GetObjectFromAllObjects(objectID)
}

func (n *Node) AddObjectToAllObjects(object universe.Object) error {
	if object.GetID() == n.GetID() {
		return errors.Errorf("not permitted for node")
	}

	world := object.GetWorld()
	if world == nil {
		return errors.Errorf("object has nil world")
	}

	if err := world.AddObjectToAllObjects(object); err != nil {
		return errors.WithMessage(err, "failed to add object to world all objects")
	}

	n.objectIDToWorld.Store(object.GetID(), world)

	return nil
}

func (n *Node) RemoveObjectFromAllObjects(object universe.Object) (bool, error) {
	if object.GetID() == n.GetID() {
		return false, errors.Errorf("not permitted for node")
	}

	n.objectIDToWorld.Mu.Lock()
	defer n.objectIDToWorld.Mu.Unlock()

	if world, ok := n.objectIDToWorld.Data[object.GetID()]; ok {
		delete(n.objectIDToWorld.Data, object.GetID())

		return world.RemoveObjectFromAllObjects(object)
	}

	return false, nil
}
