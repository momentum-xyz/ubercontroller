package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/universe"
)

type DB interface {
	CommonDB
	NodesDB
	WorldsDB
	SpacesDB
	UsersDB
	Assets2dDB
	Assets3dDB
	SpaceTypesDB
}

type CommonDB interface {
}

type NodesDB interface {
}

type WorldsDB interface {
	WorldsGetWorlds(ctx context.Context) ([]universe.SpaceEntry, error)
}

type SpacesDB interface {
	SpacesGetSpacesByParentID(ctx context.Context, parentID uuid.UUID) ([]universe.SpaceEntry, error)
}

type UsersDB interface {
}

type Assets2dDB interface {
}

type Assets3dDB interface {
}

type SpaceTypesDB interface {
}
