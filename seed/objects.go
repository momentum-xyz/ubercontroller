package seed

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
)

func seedObjects(node universe.Node) error {
	type item struct {
		id          mid.ID
		spaceTypeID mid.ID
		ownerID     mid.ID
		parentID    mid.ID
		asset2dID   mid.ID
		asset3dID   *mid.ID
		options     *entry.ObjectOptions
		position    *map[string]any
	}

	return nil
}
