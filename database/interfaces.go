package database

import (
	"context"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
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
	NodesGetNode(ctx context.Context) (*entry.Node, error)
}

type WorldsDB interface {
	WorldsGetWorldIDs(ctx context.Context) ([]uuid.UUID, error)
	WorldsGetWorlds(ctx context.Context) ([]*entry.Space, error)
}

type SpacesDB interface {
	SpacesGetSpaceByID(ctx context.Context, spaceID uuid.UUID) (*entry.Space, error)
	SpacesGetSpaceIDsByParentID(ctx context.Context, parentID uuid.UUID) ([]uuid.UUID, error)
	SpacesGetSpacesByParentID(ctx context.Context, parentID uuid.UUID) ([]*entry.Space, error)
	SpacesUpsertSpace(ctx context.Context, space *entry.Space) error
	SpacesUpsertSpaces(ctx context.Context, spaces []*entry.Space) error
	SpacesRemoveSpaceByID(ctx context.Context, spaceID uuid.UUID) error
	SpacesRemoveSpacesByIDs(ctx context.Context, spaceIDs []uuid.UUID) error
	SpacesUpdateSpaceParentID(ctx context.Context, spaceID uuid.UUID, parentID uuid.UUID) error
	SpacesUpdateSpacePosition(ctx context.Context, spaceID uuid.UUID, position *cmath.Vec3) error
	SpacesUpdateSpaceOwnerID(ctx context.Context, spaceID, ownerID uuid.UUID) error
	SpacesUpdateSpaceAsset2dID(ctx context.Context, spaceID uuid.UUID, asset2dID *uuid.UUID) error
	SpacesUpdateSpaceAsset3dID(ctx context.Context, spaceID uuid.UUID, asset3dID *uuid.UUID) error
	SpacesUpdateSpaceSpaceTypeID(ctx context.Context, spaceID, spaceTypeID uuid.UUID) error
	SpacesUpdateSpaceOptions(ctx context.Context, spaceID uuid.UUID, options *entry.SpaceOptions) error
}

type UsersDB interface {
}

type Assets2dDB interface {
	Assets2dGetAssets(ctx context.Context) ([]*entry.Asset2d, error)
	Assets2dUpsetAsset(ctx context.Context, asset2d *entry.Asset2d) error
	Assets2dUpsetAssets(ctx context.Context, assets2d []*entry.Asset2d) error
	Assets2dRemoveAssetByID(ctx context.Context, asset2dID uuid.UUID) error
	Assets2dRemoveAssetsByIDs(ctx context.Context, asset2dIDs []uuid.UUID) error
	Assets2dUpdateAssetName(ctx context.Context, asset2dID uuid.UUID, name string) error
	Assets2dUpdateAssetOptions(ctx context.Context, asset2dID uuid.UUID, asset2dOptions *entry.Asset2dOptions) error
}

type Assets3dDB interface {
	Assets3dGetAssets(ctx context.Context) ([]*entry.Asset3d, error)
	Assets3dUpsetAsset(ctx context.Context, asset3d *entry.Asset3d) error
	Assets3dUpsetAssets(ctx context.Context, assets3d []*entry.Asset3d) error
	Assets3dRemoveAssetByID(ctx context.Context, asset3dID uuid.UUID) error
	Assets3dRemoveAssetsByIDs(ctx context.Context, asset3dIDs []uuid.UUID) error
	Assets3dUpdateAssetName(ctx context.Context, asset3dID uuid.UUID, name string) error
	Assets3dUpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, asset3dOptions *entry.Asset3dOptions) error
}

type SpaceTypesDB interface {
	SpaceTypesGetSpaceTypes(ctx context.Context) ([]*entry.SpaceType, error)
	SpaceTypesUpsetSpaceType(ctx context.Context, spaceType *entry.SpaceType) error
	SpaceTypesUpsetSpaceTypes(ctx context.Context, spaceTypes []*entry.SpaceType) error
	SpaceTypesRemoveSpaceTypeByID(ctx context.Context, spaceTypeID uuid.UUID) error
	SpaceTypesRemoveSpaceTypesByIDs(ctx context.Context, spaceTypeIDs []uuid.UUID) error
	SpaceTypesUpdateSpaceTypeName(ctx context.Context, spaceTypeID uuid.UUID, name string) error
	SpaceTypesUpdateSpaceTypeCategoryName(ctx context.Context, spaceTypeID uuid.UUID, categoryName string) error
	SpaceTypesUpdateSpaceTypeDescription(ctx context.Context, spaceTypeID uuid.UUID, description *string) error
	SpaceTypesUpdateSpaceTypeOptions(ctx context.Context, spaceTypeID uuid.UUID, options *entry.SpaceOptions) error
}
