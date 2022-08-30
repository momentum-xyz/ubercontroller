package database

import (
	"context"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

type DB interface {
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
	WorldsGetWorlds(ctx context.Context) ([]entry.Space, error)
}

type SpacesDB interface {
	SpacesGetSpacesByParentID(ctx context.Context, parentID uuid.UUID) ([]entry.Space, error)
	SpacesUpdateSpaceParentID(ctx context.Context, spaceID, parentID uuid.UUID) error
	SpacesUpdateSpacePosition(ctx context.Context, spaceID uuid.UUID, position cmath.Vec3) error
	SpacesUpdateSpaceOwnerID(ctx context.Context, spaceID, ownerID uuid.UUID) error
	SpacesUpdateSpaceAsset2dID(ctx context.Context, spaceID, asset2dID uuid.UUID) error
	SpacesUpdateSpaceAsset3dID(ctx context.Context, spaceID, asset3dID uuid.UUID) error
	SpacesUpdateSpaceSpaceTypeID(ctx context.Context, spaceID, spaceTypeID uuid.UUID) error
	SpacesUpdateSpaceOptions(ctx context.Context, spaceID uuid.UUID, options *entry.SpaceOptions) error
}

type UsersDB interface {
}

type Assets2dDB interface {
	Assets2dUpdateAssetName(ctx context.Context, asset2dID uuid.UUID, name string) error
	Assets2dUpdateAssetOptions(ctx context.Context, asset2dID uuid.UUID, asset2dOptions *entry.Asset2dOptions) error
}

type Assets3dDB interface {
	Assets3dUpdateAssetName(ctx context.Context, asset3dID uuid.UUID, name string) error
	Assets3dUpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, asset3dOptions *entry.Asset3dOptions) error
}

type SpaceTypesDB interface {
	SpaceTypesUpdateSpaceTypeName(ctx context.Context, spaceTypeID uuid.UUID, name string) error
	SpaceTypesUpdateSpaceTypeCategoryName(ctx context.Context, spaceTypeID uuid.UUID, categoryName string) error
	SpaceTypesUpdateSpaceTypeDescription(ctx context.Context, spaceTypeID uuid.UUID, description *string) error
	SpaceTypesUpdateSpaceTypeOptions(ctx context.Context, spaceTypeID uuid.UUID, options *entry.SpaceOptions) error
}
