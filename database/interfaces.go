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
	UserTypesDB
	AttributesDB
	PluginsDB
	SpaceAttributesDB
	UserAttributesDB
	NodeAttributesDB
	UserUserAttributesDB
	SpaceUserAttributesDB
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
	UsersGetUserByID(ctx context.Context, userID uuid.UUID) (*entry.User, error)
	UsersUpsertUser(ctx context.Context, user *entry.User) error
	UsersUpsertUsers(ctx context.Context, user []*entry.User) error
	UsersRemoveUsersByIDs(ctx context.Context, userID []uuid.UUID) error
	UsersRemoveUserByID(ctx context.Context, userID uuid.UUID) error
	UsersUpdateUserUserTypeID(ctx context.Context, userID, UserTypeID uuid.UUID) error
	UsersUpdateUserOptions(ctx context.Context, userID uuid.UUID, options *entry.UserOptions) error
	UsersUpdateUserProfile(ctx context.Context, userID uuid.UUID, options *entry.UserProfile) error
}

type Assets2dDB interface {
	Assets2dGetAssets(ctx context.Context) ([]*entry.Asset2d, error)
	Assets2dUpsertAsset(ctx context.Context, asset2d *entry.Asset2d) error
	Assets2dUpsertAssets(ctx context.Context, assets2d []*entry.Asset2d) error
	Assets2dRemoveAssetByID(ctx context.Context, asset2dID uuid.UUID) error
	Assets2dRemoveAssetsByIDs(ctx context.Context, asset2dIDs []uuid.UUID) error
	Assets2dUpdateAssetName(ctx context.Context, asset2dID uuid.UUID, name string) error
	Assets2dUpdateAssetOptions(ctx context.Context, asset2dID uuid.UUID, asset2dOptions *entry.Asset2dOptions) error
}

type Assets3dDB interface {
	Assets3dGetAssets(ctx context.Context) ([]*entry.Asset3d, error)
	Assets3dUpsertAsset(ctx context.Context, asset3d *entry.Asset3d) error
	Assets3dUpsertAssets(ctx context.Context, assets3d []*entry.Asset3d) error
	Assets3dRemoveAssetByID(ctx context.Context, asset3dID uuid.UUID) error
	Assets3dRemoveAssetsByIDs(ctx context.Context, asset3dIDs []uuid.UUID) error
	Assets3dUpdateAssetName(ctx context.Context, asset3dID uuid.UUID, name string) error
	Assets3dUpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, asset3dOptions *entry.Asset3dOptions) error
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

type AttributesDB interface {
	AttributesGetAttributes(ctx context.Context) ([]*entry.Attribute, error)
	AttributesUpsertAttribute(ctx context.Context, attribute *entry.Attribute) error
	AttributesUpsertAttributes(ctx context.Context, attributes []*entry.Attribute) error

	AttributesRemoveAttributeByName(ctx context.Context, attributeName string) error
	AttributesRemoveAttributesByNames(ctx context.Context, attributeNames []string) error
	AttributesRemoveAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error

	AttributesRemoveAttributeByID(ctx context.Context, attributeId entry.AttributeID) error
	AttributesRemoveAttributesByIDs(ctx context.Context, attributeIds []entry.AttributeID) error

	AttributesUpdateAttributeName(ctx context.Context, attributeId entry.AttributeID, newName string) error

	AttributesUpdateAttributeDescription(
		ctx context.Context, attributeId entry.AttributeID, description *string,
	) error
	AttributesUpdateAttributeOptions(
		ctx context.Context, attributeId entry.AttributeID, attributeOptions *entry.AttributeOptions,
	) error
}

type PluginsDB interface {
	PluginsGetPlugins(ctx context.Context) ([]*entry.Plugin, error)
	PluginsUpsertPlugin(ctx context.Context, plugin *entry.Plugin) error
	PluginsUpsertPlugins(ctx context.Context, plugins []*entry.Plugin) error
	PluginsRemovePluginByID(ctx context.Context, pluginID uuid.UUID) error
	PluginsRemovePluginsByIDs(ctx context.Context, pluginIDs []uuid.UUID) error
	PluginsUpdatePluginName(ctx context.Context, pluginID uuid.UUID, name string) error
	PluginsUpdatePluginDescription(ctx context.Context, pluginID uuid.UUID, description *string) error
	PluginsUpdatePluginOptions(ctx context.Context, pluginID uuid.UUID, options *entry.PluginOptions) error
}

type SpaceAttributesDB interface {
	SpaceAttributesGetSpaceAttributes(ctx context.Context) ([]*entry.SpaceAttribute, error)
	SpaceAttributesGetSpaceAttributesBySpaceId(ctx context.Context, spaceid uuid.UUID) ([]*entry.SpaceAttribute, error)
	SpaceAttributesUpsertSpaceAttribute(ctx context.Context, spaceAttribute *entry.SpaceAttribute) error
	SpaceAttributesUpsertSpaceAttributes(ctx context.Context, spaceAttributes []*entry.SpaceAttribute) error
	SpaceAttributesRemoveSpaceAttributeByName(ctx context.Context, attributeName string) error
	SpaceAttributesRemoveSpaceAttributesByNames(ctx context.Context, attributeNames []string) error
	SpaceAttributesRemoveSpaceAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByPluginIdAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	SpaceAttributesRemoveSpaceAttributeBySpaceId(ctx context.Context, spaceID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByNameAndSpaceId(
		ctx context.Context, attributeName string, spaceID uuid.UUID,
	) error
	SpaceAttributesRemoveSpaceAttributeByNamesAndSpaceId(
		ctx context.Context, attributeNames []string, spaceID uuid.UUID,
	) error
	SpaceAttributesRemoveSpaceAttributeByPluginIdAndSpaceId(
		ctx context.Context, pluginId uuid.UUID, spaceID uuid.UUID,
	) error
	SpaceAttributesRemoveSpaceAttributeByPluginIdAndSpaceIdAndName(
		ctx context.Context, pluginId uuid.UUID, attributeName string, spaceID uuid.UUID,
	) error
	SpaceAttributesUpdateSpaceAttributeOptions(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceId uuid.UUID,
		options *entry.AttributeOptions,
	) error
	SpaceAttributesUpdateSpaceAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceId uuid.UUID, value *entry.AttributeValue,
	) error
}

type UserAttributesDB interface {
	UserAttributesGetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error)
	UserAttributesGetUserAttributesByUserId(ctx context.Context, userId uuid.UUID) ([]*entry.SpaceAttribute, error)
	UserAttributesUpsertUserAttribute(ctx context.Context, userAttribute *entry.UserAttribute) error
	UserAttributesUpsertUserAttributes(ctx context.Context, userAttributes []*entry.UserAttribute) error
	UserAttributesRemoveUserAttributeByName(ctx context.Context, attributeName string) error
	UserAttributesRemoveUserAttributesByNames(ctx context.Context, attributeNames []string) error
	UserAttributesRemoveUserAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error
	UserAttributesRemoveUserAttributeByPluginIdAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	UserAttributesRemoveUserAttributeByUserId(ctx context.Context, userID uuid.UUID) error
	UserAttributesRemoveUserAttributeByNameAndUserId(
		ctx context.Context, attributeName string, userID uuid.UUID,
	) error
	UserAttributesRemoveUserAttributeByNamesAndUserId(
		ctx context.Context, attributeNames []string, userID uuid.UUID,
	) error
	UserAttributesRemoveUserAttributeByPluginIdAndUserId(
		ctx context.Context, pluginId uuid.UUID, userID uuid.UUID,
	) error
	UserAttributesRemoveUserAttributeByPluginIdAndUserIdAndName(
		ctx context.Context, pluginId uuid.UUID, attributeName string, userID uuid.UUID,
	) error
	UserAttributesUpdateUserAttributeOptions(
		ctx context.Context, pluginID uuid.UUID, attributeName string, userId uuid.UUID,
		options *entry.AttributeOptions,
	) error
	UserAttributesUpdateUserAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, userId uuid.UUID, value *entry.AttributeValue,
	) error
}

type NodeAttributesDB interface {
	NodeAttributesGetNodeAttributes(ctx context.Context) ([]*entry.NodeAttribute, error)
	NodeAttributesUpsertNodeAttribute(ctx context.Context, nodeAttribute *entry.NodeAttribute) error
	NodeAttributesUpsertNodeAttributes(ctx context.Context, nodeAttributes []*entry.NodeAttribute) error
	NodeAttributesRemoveNodeAttributeByName(ctx context.Context, attributeName string) error
	NodeAttributesRemoveNodeAttributesByNames(ctx context.Context, attributeNames []string) error
	NodeAttributesUpdateNodeAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, nodeId uuid.UUID, value *entry.AttributeValue,
	) error
	NodeAttributesRemoveNodeAttributeByPluginIdAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	NodeAttributesRemoveNodeAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error
}

type SpaceUserAttributesDB interface {
	SpaceUserAttributesGetSpaceUserAttributes(ctx context.Context) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributesBySpaceId(
		ctx context.Context, spaceId uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributesByUserId(
		ctx context.Context, userId uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributesBySpaceIdAndUserId(
		ctx context.Context, spaceId uuid.UUID, userId uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesUpsertSpaceUserAttribute(
		ctx context.Context, spaceUserAttribute *entry.SpaceUserAttribute,
	) error
	SpaceUserAttributesUpsertSpaceUserAttributes(
		ctx context.Context, spaceUserAttributes []*entry.SpaceUserAttribute,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByName(ctx context.Context, attributeName string) error
	SpaceUserAttributesRemoveSpaceUserAttributesByNames(ctx context.Context, attributeNames []string) error
	SpaceUserAttributesRemoveSpaceUserAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeBySpaceId(
		ctx context.Context, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceId(
		ctx context.Context, attributeName string, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceId(
		ctx context.Context, attributeNames []string, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByUserId(
		ctx context.Context, userId uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNameAndUserId(
		ctx context.Context, attributeName string, userId uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndUserId(
		ctx context.Context, attributeNames []string, userId uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeBySpaceIdAndUserId(
		ctx context.Context, spaceID uuid.UUID, userId uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceIdAndUserId(
		ctx context.Context, attributeName string, spaceID uuid.UUID, userId uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceIdAndUserId(
		ctx context.Context, attributeNames []string, spaceID uuid.UUID, userId uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndSpaceId(
		ctx context.Context, pluginId uuid.UUID, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndSpaceIdAndName(
		ctx context.Context, pluginId uuid.UUID, attributeName string, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndUserId(
		ctx context.Context, pluginId uuid.UUID, userId uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndUserIdAndName(
		ctx context.Context, pluginId uuid.UUID, attributeName string, userId uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndSpaceIdAndUserId(
		ctx context.Context, pluginId uuid.UUID, spaceID uuid.UUID, userId uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIdAndSpaceIdAndUserIdAndName(
		ctx context.Context, pluginId uuid.UUID, attributeName string, spaceID uuid.UUID, userId uuid.UUID,
	) error
	SpaceUserAttributesUpdateSpaceUserAttributeOptions(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceId uuid.UUID, userId uuid.UUID,
		options *entry.AttributeOptions,
	) error
	SpaceUserAttributesUpdateSpaceUserAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceId uuid.UUID, userId uuid.UUID,
		value *entry.AttributeValue,
	) error
}

type UserUserAttributesDB interface {
	UserUserAttributesGetUserUserAttributes(ctx context.Context) ([]*entry.UserUserAttribute, error)
	UserUserAttributesGetUserUserAttributesBySourceUserId(
		ctx context.Context, sourceUserId uuid.UUID,
	) ([]*entry.UserUserAttribute, error)
	UserUserAttributesGetUserUserAttributesByTargetUserId(
		ctx context.Context, targetTargetUserId uuid.UUID,
	) ([]*entry.UserUserAttribute, error)
	UserUserAttributesGetUserUserAttributesBySourceUserIdAndTargetUserId(
		ctx context.Context, sourceUserId uuid.UUID, targetTargetUserId uuid.UUID,
	) ([]*entry.UserUserAttribute, error)
	UserUserAttributesUpsertUserUserAttribute(
		ctx context.Context, userUserAttribute *entry.UserUserAttribute,
	) error
	UserUserAttributesUpsertUserUserAttributes(
		ctx context.Context, userUserAttributes []*entry.UserUserAttribute,
	) error
	UserUserAttributesRemoveUserUserAttributeByName(ctx context.Context, attributeName string) error
	UserUserAttributesRemoveUserUserAttributesByNames(ctx context.Context, attributeNames []string) error
	UserUserAttributesRemoveUserUserAttributesByPluginId(ctx context.Context, pluginID uuid.UUID) error
	UserUserAttributesRemoveUserUserAttributeByPluginIdAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	UserUserAttributesRemoveUserUserAttributeBySourceUserId(
		ctx context.Context, sourceUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserId(
		ctx context.Context, attributeName string, sourceUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserId(
		ctx context.Context, attributeNames []string, sourceUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByTargetUserId(
		ctx context.Context, targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNameAndTargetUserId(
		ctx context.Context, attributeName string, targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNamesAndTargetUserId(
		ctx context.Context, attributeNames []string, targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeBySourceUserIdAndTargetUserId(
		ctx context.Context, sourceUserId uuid.UUID, targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserIdAndTargetUserId(
		ctx context.Context, attributeName string, sourceUserId uuid.UUID, targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserIdAndTargetUserId(
		ctx context.Context, attributeNames []string, sourceUserId uuid.UUID, targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIdAndSourceUserId(
		ctx context.Context, pluginId uuid.UUID, sourceUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIdAndSourceUserIdAndName(
		ctx context.Context, pluginId uuid.UUID, attributeName string, sourceUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIdAndTargetUserId(
		ctx context.Context, pluginId uuid.UUID, targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIdAndTargetUserIdAndName(
		ctx context.Context, pluginId uuid.UUID, attributeName string, targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIdAndSourceUserIdAndTargetUserId(
		ctx context.Context, pluginId uuid.UUID, sourceUserId uuid.UUID, targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIdAndSourceUserIdAndTargetUserIdAndName(
		ctx context.Context, pluginId uuid.UUID, attributeName string, sourceUserId uuid.UUID,
		targetTargetUserId uuid.UUID,
	) error
	UserUserAttributesUpdateUserUserAttributeOptions(
		ctx context.Context, pluginID uuid.UUID, attributeName string, sourceUserId uuid.UUID,
		targetTargetUserId uuid.UUID,
		options *entry.AttributeOptions,
	) error
	UserUserAttributesUpdateUserUserAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, sourceUserId uuid.UUID,
		targetTargetUserId uuid.UUID,
		value *entry.AttributeValue,
	) error
}
