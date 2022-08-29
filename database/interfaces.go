package database

import (
	"context"

	"github.com/google/uuid"
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
	WorldsGetWorldIDs(ctx context.Context) ([]uuid.UUID, error)
}

type SpacesDB interface {
}

type UsersDB interface {
}

type Assets2dDB interface {
}

type Assets3dDB interface {
}

type SpaceTypesDB interface {
}
