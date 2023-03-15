package node

import (
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
)

func (n *Node) detectSpawnWorld(userId mid.ID) universe.World {
	// TODO: implement. Temporary, just first world from the list
	wid := mid.MustParse("d83670c7-a120-47a4-892d-f9ec75604f74")
	if world, ok := n.worlds.GetWorld(wid); ok != false {
		return world
	}
	return nil
}
