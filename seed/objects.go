package seed

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func seedObjects(node universe.Node) error {
	type item struct {
		id          uuid.UUID
		spaceTypeID uuid.UUID
		ownerID     uuid.UUID
		parentID    uuid.UUID
		asset2dID   uuid.UUID
		asset3dID   *uuid.UUID
		options     *entry.ObjectOptions
		position    *map[string]any
	}

	return nil
}
