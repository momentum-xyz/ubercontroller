package seed

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

func seedObjects(node universe.Node) error {
	type item struct {
		id          umid.UMID
		spaceTypeID umid.UMID
		ownerID     umid.UMID
		parentID    umid.UMID
		asset2dID   umid.UMID
		asset3dID   *umid.UMID
		options     *entry.ObjectOptions
		transform   *map[string]any
	}

	return nil
}
