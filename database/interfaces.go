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
	Assets2dUpdateAssetName(ctx context.Context, asset2dID uuid.UUID, name string) error
	Assets2dUpdateAssetOptions(ctx context.Context, asset2dID uuid.UUID, options *entry.Asset2dOptions) error
}

type Assets3dDB interface {
	Assets3dGetAssets(ctx context.Context) ([]*entry.Asset3d, error)
	Assets3dUpsertAsset(ctx context.Context, asset3d *entry.Asset3d) error
	Assets3dUpsertAssets(ctx context.Context, assets3d []*entry.Asset3d) error
	Assets3dRemoveAssetByID(ctx context.Context, asset3dID uuid.UUID) error
	Assets3dRemoveAssetsByIDs(ctx context.Context, asset3dIDs []uuid.UUID) error
	Assets3dUpdateAssetName(ctx context.Context, asset3dID uuid.UUID, name string) error
	Assets3dUpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, options *entry.Asset3dOptions) error
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

	// QUESTION: will remove attributes for all plugins, is it designed?
	AttributesRemoveAttributeByName(ctx context.Context, name string) error
	// QUESTION: the same as for "AttributesRemoveAttributeByName"
	AttributesRemoveAttributesByNames(ctx context.Context, names []string) error
	AttributesRemoveAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error

	AttributesRemoveAttributeByID(ctx context.Context, attributeID entry.AttributeID) error
	AttributesRemoveAttributesByIDs(ctx context.Context, attributeIDs []entry.AttributeID) error

	AttributesUpdateAttributeName(ctx context.Context, attributeID entry.AttributeID, name string) error

	AttributesUpdateAttributeDescription(
		ctx context.Context, attributeID entry.AttributeID, description *string,
	) error
	AttributesUpdateAttributeOptions(
		ctx context.Context, attributeID entry.AttributeID, options *entry.AttributeOptions,
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
	SpaceAttributesGetSpaceAttributesBySpaceID(ctx context.Context, spaceID uuid.UUID) ([]*entry.SpaceAttribute, error)
	SpaceAttributesUpsertSpaceAttribute(ctx context.Context, spaceAttribute *entry.SpaceAttribute) error
	SpaceAttributesUpsertSpaceAttributes(ctx context.Context, spaceAttributes []*entry.SpaceAttribute) error
	SpaceAttributesRemoveSpaceAttributeByName(ctx context.Context, name string) error
	SpaceAttributesRemoveSpaceAttributesByNames(ctx context.Context, names []string) error
	SpaceAttributesRemoveSpaceAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByPluginIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	SpaceAttributesRemoveSpaceAttributeBySpaceID(ctx context.Context, spaceID uuid.UUID) error
	// QUESTION: will remove attribute for all plugins, is it designed?
	SpaceAttributesRemoveSpaceAttributeByNameAndSpaceID(
		ctx context.Context, attributeName string, spaceID uuid.UUID,
	) error
	// QUESTION: same as for "SpaceAttributesRemoveSpaceAttributeByNameAndSpaceID"
	SpaceAttributesRemoveSpaceAttributeByNamesAndSpaceID(
		ctx context.Context, attributeNames []string, spaceID uuid.UUID,
	) error
	SpaceAttributesRemoveSpaceAttributeByPluginIDAndSpaceID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID,
	) error
	SpaceAttributesRemoveSpaceAttributeByPluginIDAndSpaceIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceID uuid.UUID,
	) error
	SpaceAttributesUpdateSpaceAttributeOptions(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceID uuid.UUID,
		options *entry.AttributeOptions,
	) error
	SpaceAttributesUpdateSpaceAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceID uuid.UUID, value *entry.AttributeValue,
	) error
}

type UserAttributesDB interface {
	UserAttributesGetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error)
	UserAttributesGetUserAttributesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.SpaceAttribute, error)
	UserAttributesUpsertUserAttribute(ctx context.Context, userAttribute *entry.UserAttribute) error
	UserAttributesUpsertUserAttributes(ctx context.Context, userAttributes []*entry.UserAttribute) error
	UserAttributesRemoveUserAttributeByName(ctx context.Context, attributeName string) error
	UserAttributesRemoveUserAttributesByNames(ctx context.Context, attributeNames []string) error
	UserAttributesRemoveUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	UserAttributesRemoveUserAttributeByPluginIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	UserAttributesRemoveUserAttributeByUserID(ctx context.Context, userID uuid.UUID) error
	UserAttributesRemoveUserAttributeByNameAndUserID(
		ctx context.Context, attributeName string, userID uuid.UUID,
	) error
	UserAttributesRemoveUserAttributeByNamesAndUserID(
		ctx context.Context, attributeNames []string, userID uuid.UUID,
	) error
	UserAttributesRemoveUserAttributeByPluginIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, userID uuid.UUID,
	) error
	UserAttributesRemoveUserAttributeByPluginIDAndUserIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string, userID uuid.UUID,
	) error
	UserAttributesUpdateUserAttributeOptions(
		ctx context.Context, pluginID uuid.UUID, attributeName string, userID uuid.UUID,
		options *entry.AttributeOptions,
	) error
	UserAttributesUpdateUserAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, userID uuid.UUID, value *entry.AttributeValue,
	) error
}

type NodeAttributesDB interface {
	NodeAttributesGetNodeAttributes(ctx context.Context) ([]*entry.NodeAttribute, error)
	NodeAttributesUpsertNodeAttribute(ctx context.Context, nodeAttribute *entry.NodeAttribute) error
	NodeAttributesUpsertNodeAttributes(ctx context.Context, nodeAttributes []*entry.NodeAttribute) error
	// QUESTION: for all plugins?
	NodeAttributesRemoveNodeAttributeByName(ctx context.Context, attributeName string) error
	// QUESTION: as above?
	NodeAttributesRemoveNodeAttributesByNames(ctx context.Context, attributeNames []string) error
	NodeAttributesUpdateNodeAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, nodeID uuid.UUID, value *entry.AttributeValue,
	) error
	NodeAttributesRemoveNodeAttributeByPluginIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	NodeAttributesRemoveNodeAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
}

type SpaceUserAttributesDB interface {
	SpaceUserAttributesGetSpaceUserAttributes(ctx context.Context) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributesBySpaceID(
		ctx context.Context, spaceID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributesByUserID(
		ctx context.Context, userID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributesBySpaceIDAndUserID(
		ctx context.Context, spaceID uuid.UUID, userID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesUpsertSpaceUserAttribute(
		ctx context.Context, spaceUserAttribute *entry.SpaceUserAttribute,
	) error
	SpaceUserAttributesUpsertSpaceUserAttributes(
		ctx context.Context, spaceUserAttributes []*entry.SpaceUserAttribute,
	) error
	// QUESTION: for all plugins?
	SpaceUserAttributesRemoveSpaceUserAttributeByName(ctx context.Context, attributeName string) error
	// QUESTION: for all plugins?
	SpaceUserAttributesRemoveSpaceUserAttributesByNames(ctx context.Context, attributeNames []string) error
	SpaceUserAttributesRemoveSpaceUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeBySpaceID(
		ctx context.Context, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceID(
		ctx context.Context, attributeName string, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceID(
		ctx context.Context, attributeNames []string, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByUserID(
		ctx context.Context, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNameAndUserID(
		ctx context.Context, attributeName string, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndUserID(
		ctx context.Context, attributeNames []string, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeBySpaceIDAndUserID(
		ctx context.Context, spaceID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNameAndSpaceIDAndUserID(
		ctx context.Context, attributeName string, spaceID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByNamesAndSpaceIDAndUserID(
		ctx context.Context, attributeNames []string, spaceID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndSpaceID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndSpaceIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndUserIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndSpaceIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesRemoveSpaceUserAttributeByPluginIDAndSpaceIDAndUserIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceID uuid.UUID, userID uuid.UUID,
	) error
	SpaceUserAttributesUpdateSpaceUserAttributeOptions(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceID uuid.UUID, userID uuid.UUID,
		options *entry.AttributeOptions,
	) error
	SpaceUserAttributesUpdateSpaceUserAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, spaceID uuid.UUID, userID uuid.UUID,
		value *entry.AttributeValue,
	) error
}

type UserUserAttributesDB interface {
	UserUserAttributesGetUserUserAttributes(ctx context.Context) ([]*entry.UserUserAttribute, error)
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
		ctx context.Context, userUserAttribute *entry.UserUserAttribute,
	) error
	UserUserAttributesUpsertUserUserAttributes(
		ctx context.Context, userUserAttributes []*entry.UserUserAttribute,
	) error
	UserUserAttributesRemoveUserUserAttributeByName(ctx context.Context, attributeName string) error
	UserUserAttributesRemoveUserUserAttributesByNames(ctx context.Context, attributeNames []string) error
	UserUserAttributesRemoveUserUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string,
	) error
	UserUserAttributesRemoveUserUserAttributeBySourceUserID(
		ctx context.Context, sourceUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserID(
		ctx context.Context, attributeName string, sourceUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserID(
		ctx context.Context, attributeNames []string, sourceUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByTargetUserID(
		ctx context.Context, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNameAndTargetUserID(
		ctx context.Context, attributeName string, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNamesAndTargetUserID(
		ctx context.Context, attributeNames []string, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeBySourceUserIDAndTargetUserID(
		ctx context.Context, sourceUserId uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNameAndSourceUserIDAndTargetUserID(
		ctx context.Context, attributeName string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByNamesAndSourceUserIDAndTargetUserID(
		ctx context.Context, attributeNames []string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndSourceUserID(
		ctx context.Context, pluginID uuid.UUID, sourceUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndSourceUserIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string, sourceUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndTargetUserID(
		ctx context.Context, pluginID uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndTargetUserIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndSourceUserIDAndTargetUserID(
		ctx context.Context, pluginID uuid.UUID, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	UserUserAttributesRemoveUserUserAttributeByPluginIDAndSourceUserIDAndTargetUserIDAndName(
		ctx context.Context, pluginID uuid.UUID, attributeName string, sourceUserID uuid.UUID,
		targetUserID uuid.UUID,
	) error
	UserUserAttributesUpdateUserUserAttributeOptions(
		ctx context.Context, pluginID uuid.UUID, attributeName string, sourceUserID uuid.UUID,
		targetUserID uuid.UUID, options *entry.AttributeOptions,
	) error
	UserUserAttributesUpdateUserUserAttributeValue(
		ctx context.Context, pluginID uuid.UUID, attributeName string, sourceUserID uuid.UUID,
		targetUserID uuid.UUID, value *entry.AttributeValue,
	) error
}
