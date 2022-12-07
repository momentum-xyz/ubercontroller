package database

import (
	"context"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type DB interface {
	CommonDB
	NodesDB
	WorldsDB
	SpacesDB
	UsersDB
	Assets2dDB
	Assets3dDB
	PluginsDB
	SpaceTypesDB
	UserTypesDB
	AttributeTypesDB
	NodeAttributesDB
	SpaceAttributesDB
	SpaceUserAttributesDB
	UserAttributesDB
	UserUserAttributesDB
	UserSpaceDB
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
	SpacesUpdateSpacePosition(ctx context.Context, spaceID uuid.UUID, position *cmath.SpacePosition) error
	SpacesUpdateSpaceOwnerID(ctx context.Context, spaceID, ownerID uuid.UUID) error
	SpacesUpdateSpaceAsset2dID(ctx context.Context, spaceID uuid.UUID, asset2dID *uuid.UUID) error
	SpacesUpdateSpaceAsset3dID(ctx context.Context, spaceID uuid.UUID, asset3dID *uuid.UUID) error
	SpacesUpdateSpaceSpaceTypeID(ctx context.Context, spaceID, spaceTypeID uuid.UUID) error
	SpacesUpdateSpaceOptions(ctx context.Context, spaceID uuid.UUID, options *entry.SpaceOptions) error
}

type UsersDB interface {
	UsersGetUserByID(ctx context.Context, userID uuid.UUID) (*entry.User, error)
	UsersGetUserByWallet(ctx context.Context, wallet string) (*entry.User, error)
	UsersGetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*entry.UserProfile, error)

	UsersUpsertUser(ctx context.Context, user *entry.User) error
	UsersUpsertUsers(ctx context.Context, user []*entry.User) error

	UsersRemoveUsersByIDs(ctx context.Context, userID []uuid.UUID) error
	UsersRemoveUserByID(ctx context.Context, userID uuid.UUID) error

	UsersUpdateUserUserTypeID(ctx context.Context, userID, userTypeID uuid.UUID) error
	UsersUpdateUserOptions(ctx context.Context, userID uuid.UUID, options *entry.UserOptions) error
	UsersUpdateUserProfile(ctx context.Context, userID uuid.UUID, profile *entry.UserProfile) error
}

type Assets2dDB interface {
	Assets2dGetAssets(ctx context.Context) ([]*entry.Asset2d, error)

	Assets2dUpsertAsset(ctx context.Context, asset2d *entry.Asset2d) error
	Assets2dUpsertAssets(ctx context.Context, assets2d []*entry.Asset2d) error

	Assets2dRemoveAssetByID(ctx context.Context, asset2dID uuid.UUID) error
	Assets2dRemoveAssetsByIDs(ctx context.Context, asset2dIDs []uuid.UUID) error

	Assets2dUpdateAssetMeta(ctx context.Context, asset2dID uuid.UUID, meta *entry.Asset2dMeta) error
	Assets2dUpdateAssetOptions(ctx context.Context, asset2dID uuid.UUID, options *entry.Asset2dOptions) error
}

type Assets3dDB interface {
	Assets3dGetAssets(ctx context.Context) ([]*entry.Asset3d, error)

	Assets3dUpsertAsset(ctx context.Context, asset3d *entry.Asset3d) error
	Assets3dUpsertAssets(ctx context.Context, assets3d []*entry.Asset3d) error

	Assets3dRemoveAssetByID(ctx context.Context, asset3dID uuid.UUID) error
	Assets3dRemoveAssetsByIDs(ctx context.Context, asset3dIDs []uuid.UUID) error

	Assets3dUpdateAssetMeta(ctx context.Context, asset3dID uuid.UUID, meta *entry.Asset3dMeta) error
	Assets3dUpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, options *entry.Asset3dOptions) error
}

type PluginsDB interface {
	PluginsGetPlugins(ctx context.Context) ([]*entry.Plugin, error)

	PluginsUpsertPlugin(ctx context.Context, plugin *entry.Plugin) error
	PluginsUpsertPlugins(ctx context.Context, plugins []*entry.Plugin) error

	PluginsRemovePluginByID(ctx context.Context, pluginID uuid.UUID) error
	PluginsRemovePluginsByIDs(ctx context.Context, pluginIDs []uuid.UUID) error

	PluginsUpdatePluginMeta(ctx context.Context, pluginID uuid.UUID, meta *entry.PluginMeta) error
	PluginsUpdatePluginOptions(
		ctx context.Context, pluginID uuid.UUID, options *entry.PluginOptions,
	) error
}

type SpaceTypesDB interface {
	SpaceTypesGetSpaceTypes(ctx context.Context) ([]*entry.SpaceType, error)

	SpaceTypesUpsertSpaceType(ctx context.Context, spaceType *entry.SpaceType) error
	SpaceTypesUpsertSpaceTypes(ctx context.Context, spaceTypes []*entry.SpaceType) error

	SpaceTypesRemoveSpaceTypeByID(ctx context.Context, spaceTypeID uuid.UUID) error
	SpaceTypesRemoveSpaceTypesByIDs(ctx context.Context, spaceTypeIDs []uuid.UUID) error

	SpaceTypesUpdateSpaceTypeName(ctx context.Context, spaceTypeID uuid.UUID, name string) error
	SpaceTypesUpdateSpaceTypeCategoryName(ctx context.Context, spaceTypeID uuid.UUID, categoryName string) error
	SpaceTypesUpdateSpaceTypeDescription(ctx context.Context, spaceTypeID uuid.UUID, description *string) error
	SpaceTypesUpdateSpaceTypeOptions(ctx context.Context, spaceTypeID uuid.UUID, options *entry.SpaceOptions) error
}

type UserTypesDB interface {
	UserTypesGetUserTypes(ctx context.Context) ([]*entry.UserType, error)

	UserTypesUpsertUserType(ctx context.Context, userType *entry.UserType) error
	UserTypesUpsertUserTypes(ctx context.Context, userTypes []*entry.UserType) error

	UserTypesRemoveUserTypeByID(ctx context.Context, userTypeID uuid.UUID) error
	UserTypesRemoveUserTypesByIDs(ctx context.Context, userTypeIDs []uuid.UUID) error

	UserTypesUpdateUserTypeName(ctx context.Context, userTypeID uuid.UUID, name string) error
	UserTypesUpdateUserTypeDescription(ctx context.Context, userTypeID uuid.UUID, description *string) error
	UserTypesUpdateUserTypeOptions(ctx context.Context, userTypeID uuid.UUID, options *entry.UserOptions) error
}

type AttributeTypesDB interface {
	AttributeTypesGetAttributeTypes(ctx context.Context) ([]*entry.AttributeType, error)

	AttributeTypesUpsertAttributeType(ctx context.Context, attributeType *entry.AttributeType) error
	AttributeTypesUpsertAttributeTypes(ctx context.Context, attributeTypes []*entry.AttributeType) error

	AttributeTypesRemoveAttributeTypeByName(ctx context.Context, name string) error
	AttributeTypesRemoveAttributeTypesByNames(ctx context.Context, names []string) error
	AttributeTypesRemoveAttributeTypesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	AttributeTypesRemoveAttributeTypeByID(ctx context.Context, attributeTypeID entry.AttributeTypeID) error
	AttributeTypesRemoveAttributeTypesByIDs(ctx context.Context, attributeTypeIDs []entry.AttributeTypeID) error

	AttributeTypesUpdateAttributeTypeName(ctx context.Context, attributeTypeID entry.AttributeTypeID, name string) error
	AttributeTypesUpdateAttributeTypeDescription(
		ctx context.Context, attributeTypeID entry.AttributeTypeID, description *string,
	) error
	AttributeTypesUpdateAttributeTypeOptions(
		ctx context.Context, attributeTypeID entry.AttributeTypeID, options *entry.AttributeOptions,
	) error
}

type NodeAttributesDB interface {
	NodeAttributesGetNodeAttributes(ctx context.Context) ([]*entry.NodeAttribute, error)
	NodeAttributesGetNodeAttributeByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) (*entry.NodeAttribute, error)
	NodeAttributesGetNodeAttributeValueByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) (*entry.AttributeValue, error)
	NodeAttributesGetNodeAttributeOptionsByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) (*entry.AttributeOptions, error)

	NodeAttributesUpsertNodeAttribute(ctx context.Context, nodeAttribute *entry.NodeAttribute) error
	NodeAttributesUpsertNodeAttributes(ctx context.Context, nodeAttributes []*entry.NodeAttribute) error

	NodeAttributesRemoveNodeAttributeByName(ctx context.Context, name string) error
	NodeAttributesRemoveNodeAttributesByNames(ctx context.Context, names []string) error
	NodeAttributesRemoveNodeAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	NodeAttributesRemoveNodeAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error

	NodeAttributesUpdateNodeAttributeValue(
		ctx context.Context, attributeID entry.AttributeID, value *entry.AttributeValue,
	) error
	NodeAttributesUpdateNodeAttributeOptions(
		ctx context.Context, attributeID entry.AttributeID, options *entry.AttributeOptions,
	) error
}

type SpaceAttributesDB interface {
	SpaceAttributesGetSpaceAttributes(ctx context.Context) ([]*entry.SpaceAttribute, error)
	SpaceAttributesGetSpaceAttributeByID(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID,
	) (*entry.SpaceAttribute, error)
	SpaceAttributesGetSpaceAttributesBySpaceID(ctx context.Context, spaceID uuid.UUID) ([]*entry.SpaceAttribute, error)

	SpaceAttributesUpsertSpaceAttribute(ctx context.Context, spaceAttribute *entry.SpaceAttribute) error
	SpaceAttributesUpsertSpaceAttributes(ctx context.Context, spaceAttributes []*entry.SpaceAttribute) error

	SpaceAttributesRemoveSpaceAttributeByName(ctx context.Context, name string) error
	SpaceAttributesRemoveSpaceAttributesByNames(ctx context.Context, names []string) error
	SpaceAttributesRemoveSpaceAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	SpaceAttributesRemoveSpaceAttributeBySpaceID(ctx context.Context, spaceID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByNameAndSpaceID(ctx context.Context, name string, spaceID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByNamesAndSpaceID(ctx context.Context, names []string, spaceID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByID(ctx context.Context, spaceAttributeID entry.SpaceAttributeID) error
	SpaceAttributesRemoveSpaceAttributeByPluginIDAndSpaceID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID,
	) error

	SpaceAttributesUpdateSpaceAttributeValue(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID, value *entry.AttributeValue,
	) error
	SpaceAttributesUpdateSpaceAttributeOptions(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID, options *entry.AttributeOptions,
	) error
}

type SpaceUserAttributesDB interface {
	SpaceUserAttributesGetSpaceUserAttributes(ctx context.Context) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributeByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) (*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributeValueByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) (*entry.AttributeValue, error)
	SpaceUserAttributesGetSpaceUserAttributeOptionsByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) (*entry.AttributeOptions, error)
	SpaceUserAttributesGetSpaceUserAttributesBySpaceID(
		ctx context.Context, spaceID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributesByUserID(
		ctx context.Context, userID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributesBySpaceIDAndUserID(
		ctx context.Context, spaceID uuid.UUID, userID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributesByPluginIDAndNameAndSpaceID(
		ctx context.Context, pluginID uuid.UUID, name string, spaceID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)

	SpaceUserAttributesUpsertSpaceUserAttribute(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
		modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.SpaceUserAttribute, error)

	SpaceUserAttributesRemoveSpaceUserAttributeByName(ctx context.Context, name string) error
	SpaceUserAttributesRemoveSpaceUserAttributesByNames(ctx context.Context, names []string) error
	SpaceUserAttributesRemoveSpaceUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	SpaceUserAttributesRemoveSpaceUserAttributeByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeBySpaceID(ctx context.Context, spaceID uuid.UUID) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceID(
		ctx context.Context, name string, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceID(
		ctx context.Context, names []string, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByUserID(ctx context.Context, userID uuid.UUID) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNameAndUserID(
		ctx context.Context, name string, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndUserID(
		ctx context.Context, names []string, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeBySpaceIDAndUserID(
		ctx context.Context, spaceID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceIDAndUserID(
		ctx context.Context, name string, spaceID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceIDAndUserID(
		ctx context.Context, names []string, spaceID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndSpaceID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeBySpaceAttributeID(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByUserAttributeID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndSpaceIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) error

	SpaceUserAttributesUpdateSpaceUserAttributeOptions(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
		modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)
	SpaceUserAttributesUpdateSpaceUserAttributeValue(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
}

type UserAttributesDB interface {
	UserAttributesGetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error)
	UserAttributesGetUserAttributesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserAttribute, error)
	UserAttributesGetUserAttributeByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.UserAttribute, error)
	UserAttributesGetUserAttributeValueByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.AttributeValue, error)
	UserAttributesGetUserAttributeOptionsByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.AttributeOptions, error)

	UserAttributesUpsertUserAttribute(
		ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.UserAttribute, error)

	UserAttributesRemoveUserAttributeByName(ctx context.Context, name string) error
	UserAttributesRemoveUserAttributesByNames(ctx context.Context, names []string) error
	UserAttributesRemoveUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	UserAttributesRemoveUserAttributeByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) error
	UserAttributesRemoveUserAttributeByUserID(ctx context.Context, userID uuid.UUID) error
	UserAttributesRemoveUserAttributeByNameAndUserID(
		ctx context.Context, name string, userID uuid.UUID,
	) error
	UserAttributesRemoveUserAttributeByNamesAndUserID(
		ctx context.Context, names []string, userID uuid.UUID,
	) error
	UserAttributesRemoveUserAttributeByPluginIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, userID uuid.UUID,
	) error
	UserAttributesRemoveUserAttributeByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) error

	UserAttributesUpdateUserAttributeValue(
		ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
	UserAttributesUpdateUserAttributeOptions(
		ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)
}

type UserUserAttributesDB interface {
	UserUserAttributesGetUserUserAttributes(ctx context.Context) ([]*entry.UserUserAttribute, error)
	UserUserAttributesGetUserUserAttributeByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) (*entry.UserUserAttribute, error)
	UserUserAttributesGetUserUserAttributeValueByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) (*entry.AttributeValue, error)
	UserUserAttributesGetUserUserAttributeOptionsByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) (*entry.AttributeOptions, error)
	UserUserAttributesGetUserUserAttributesBySourceUserID(
		ctx context.Context, sourceUserID uuid.UUID,
	) ([]*entry.UserUserAttribute, error)
	UserUserAttributesGetUserUserAttributesByTargetUserID(
		ctx context.Context, targetUserID uuid.UUID,
	) ([]*entry.UserUserAttribute, error)
	UserUserAttributesGetUserUserAttributesBySourceUserIDAndTargetUserID(
		ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) ([]*entry.UserUserAttribute, error)

	UserUserAttributesUpsertUserUserAttribute(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.UserUserAttribute, error)

	UserUserAttributesRemoveUserUserAttributeByName(ctx context.Context, name string) error
	UserUserAttributesRemoveUserUserAttributesByNames(ctx context.Context, names []string) error
	UserUserAttributesRemoveUserUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	UserUserAttributesRemoveUserUserAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	UserUserAttributesRemoveUserUserAttributeBySourceUserID(
		ctx context.Context, sourceUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserID(
		ctx context.Context, name string, sourceUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserID(
		ctx context.Context, names []string, sourceUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByTargetUserID(
		ctx context.Context, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNameAndTargetUserID(
		ctx context.Context, name string, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNamesAndTargetUserID(
		ctx context.Context, names []string, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeBySourceUserIDAndTargetUserID(
		ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserIDAndTargetUserID(
		ctx context.Context, name string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserIDAndTargetUserID(
		ctx context.Context, names []string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndSourceUserID(
		ctx context.Context, pluginID uuid.UUID, sourceUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeBySourceUserAttributeID(
		ctx context.Context, sourceUserAttributeID entry.UserAttributeID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndTargetUserID(
		ctx context.Context, pluginID uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByTargetUserAttributeID(
		ctx context.Context, targetUserAttributeID entry.UserAttributeID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndSourceUserIDAndTargetUserID(
		ctx context.Context, pluginID uuid.UUID, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) error

	UserUserAttributesUpdateUserUserAttributeValue(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
	UserUserAttributesUpdateUserUserAttributeOptions(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)
}

type UserSpaceDB interface {
	UserSpaceGetUserSpaces(ctx context.Context) ([]*entry.UserSpace, error)

	UserSpaceGetUserSpacesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserSpace, error)
	UserSpaceGetUserSpacesBySpaceID(ctx context.Context, spaceID uuid.UUID) ([]*entry.UserSpace, error)
	UserSpaceGetUserSpaceByUserAndSpaceIDs(ctx context.Context, userSpaceID entry.UserSpaceID) (*entry.UserSpace, error)

	UserSpaceGetValueByUserAndSpaceIDs(ctx context.Context, userSpaceID entry.UserSpaceID) (*entry.UserSpaceValue, error)

	UserSpaceGetIndirectAdmins(ctx context.Context, spaceID uuid.UUID) ([]*uuid.UUID, error)

	UserSpaceUpdateValueByUserAndSpaceIDs(ctx context.Context, userSpaceID entry.UserSpaceID, modifyFn modify.Fn[entry.UserSpaceValue]) (*entry.UserSpaceValue, error)
}
