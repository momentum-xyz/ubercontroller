package node

import (
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

func (n *Node) detectSpawnWorld(userId umid.UMID) universe.World {
	// TODO: implement. Temporary, just first world from the list
	wid := umid.MustParse("d83670c7-a120-47a4-892d-f9ec75604f74")
	if world, ok := n.worlds.GetWorld(wid); ok != false {
		return world
	}
	return nil
}
