package database

import (
	"context"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type DB interface {
	CommonDB
	GetNodesDB() NodesDB
	GetWorldsDB() WorldsDB
	GetSpacesDB() SpacesDB
	GetUsersDB() UsersDB
	GetAssets2dDB() Assets2dDB
	GetAssets3dDB() Assets3dDB
	GetPluginsDB() PluginsDB
	GetSpaceTypesDB() SpaceTypesDB
	GetUserTypesDB() UserTypesDB
	GetAttributeTypesDB() AttributeTypesDB
	GetNodeAttributesDB() NodeAttributesDB
	SpaceAttributesDB
	SpaceUserAttributesDB
	UserAttributesDB
	UserUserAttributesDB
	UserSpaceDB
}

type CommonDB interface {
}

type NodesDB interface {
	GetNode(ctx context.Context) (*entry.Node, error)
}

type WorldsDB interface {
	GetWorldIDs(ctx context.Context) ([]uuid.UUID, error)
	GetWorlds(ctx context.Context) ([]*entry.Space, error)
}

type SpacesDB interface {
	GetSpaceByID(ctx context.Context, spaceID uuid.UUID) (*entry.Space, error)
	GetSpaceIDsByParentID(ctx context.Context, parentID uuid.UUID) ([]uuid.UUID, error)
	GetSpacesByParentID(ctx context.Context, parentID uuid.UUID) ([]*entry.Space, error)

	UpsertSpace(ctx context.Context, space *entry.Space) error
	UpsertSpaces(ctx context.Context, spaces []*entry.Space) error

	UpdateSpaceParentID(ctx context.Context, spaceID uuid.UUID, parentID uuid.UUID) error
	UpdateSpacePosition(ctx context.Context, spaceID uuid.UUID, position *cmath.SpacePosition) error
	UpdateSpaceOwnerID(ctx context.Context, spaceID, ownerID uuid.UUID) error
	UpdateSpaceAsset2dID(ctx context.Context, spaceID uuid.UUID, asset2dID *uuid.UUID) error
	UpdateSpaceAsset3dID(ctx context.Context, spaceID uuid.UUID, asset3dID *uuid.UUID) error
	UpdateSpaceSpaceTypeID(ctx context.Context, spaceID, spaceTypeID uuid.UUID) error
	UpdateSpaceOptions(ctx context.Context, spaceID uuid.UUID, options *entry.SpaceOptions) error

	RemoveSpaceByID(ctx context.Context, spaceID uuid.UUID) error
	RemoveSpacesByIDs(ctx context.Context, spaceIDs []uuid.UUID) error
}

type UsersDB interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entry.User, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entry.User, error)
	GetUserByWallet(ctx context.Context, wallet string) (*entry.User, error)
	GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*entry.UserProfile, error)

	UpsertUser(ctx context.Context, user *entry.User) error
	UpsertUsers(ctx context.Context, user []*entry.User) error

	UpdateUserUserTypeID(ctx context.Context, userID, userTypeID uuid.UUID) error
	UpdateUserOptions(ctx context.Context, userID uuid.UUID, options *entry.UserOptions) error
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, profile *entry.UserProfile) error

	RemoveUsersByIDs(ctx context.Context, userID []uuid.UUID) error
	RemoveUserByID(ctx context.Context, userID uuid.UUID) error
}

type Assets2dDB interface {
	GetAssets(ctx context.Context) ([]*entry.Asset2d, error)

	UpsertAsset(ctx context.Context, asset2d *entry.Asset2d) error
	UpsertAssets(ctx context.Context, assets2d []*entry.Asset2d) error

	UpdateAssetMeta(ctx context.Context, asset2dID uuid.UUID, meta *entry.Asset2dMeta) error
	UpdateAssetOptions(ctx context.Context, asset2dID uuid.UUID, options *entry.Asset2dOptions) error

	RemoveAssetByID(ctx context.Context, asset2dID uuid.UUID) error
	RemoveAssetsByIDs(ctx context.Context, asset2dIDs []uuid.UUID) error
}

type Assets3dDB interface {
	GetAssets(ctx context.Context) ([]*entry.Asset3d, error)

	UpsertAsset(ctx context.Context, asset3d *entry.Asset3d) error
	UpsertAssets(ctx context.Context, assets3d []*entry.Asset3d) error

	UpdateAssetMeta(ctx context.Context, asset3dID uuid.UUID, meta *entry.Asset3dMeta) error
	UpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, options *entry.Asset3dOptions) error

	RemoveAssetByID(ctx context.Context, asset3dID uuid.UUID) error
	RemoveAssetsByIDs(ctx context.Context, asset3dIDs []uuid.UUID) error
}

type PluginsDB interface {
	GetPlugins(ctx context.Context) ([]*entry.Plugin, error)

	UpsertPlugin(ctx context.Context, plugin *entry.Plugin) error
	UpsertPlugins(ctx context.Context, plugins []*entry.Plugin) error

	UpdatePluginMeta(ctx context.Context, pluginID uuid.UUID, meta *entry.PluginMeta) error
	UpdatePluginOptions(
		ctx context.Context, pluginID uuid.UUID, options *entry.PluginOptions,
	) error

	RemovePluginByID(ctx context.Context, pluginID uuid.UUID) error
	RemovePluginsByIDs(ctx context.Context, pluginIDs []uuid.UUID) error
}

type SpaceTypesDB interface {
	GetSpaceTypes(ctx context.Context) ([]*entry.SpaceType, error)

	UpsertSpaceType(ctx context.Context, spaceType *entry.SpaceType) error
	UpsertSpaceTypes(ctx context.Context, spaceTypes []*entry.SpaceType) error

	UpdateSpaceTypeName(ctx context.Context, spaceTypeID uuid.UUID, name string) error
	UpdateSpaceTypeCategoryName(ctx context.Context, spaceTypeID uuid.UUID, categoryName string) error
	UpdateSpaceTypeDescription(ctx context.Context, spaceTypeID uuid.UUID, description *string) error
	UpdateSpaceTypeOptions(ctx context.Context, spaceTypeID uuid.UUID, options *entry.SpaceOptions) error

	RemoveSpaceTypeByID(ctx context.Context, spaceTypeID uuid.UUID) error
	RemoveSpaceTypesByIDs(ctx context.Context, spaceTypeIDs []uuid.UUID) error
}

type UserTypesDB interface {
	GetUserTypes(ctx context.Context) ([]*entry.UserType, error)

	UpsertUserType(ctx context.Context, userType *entry.UserType) error
	UpsertUserTypes(ctx context.Context, userTypes []*entry.UserType) error

	UpdateUserTypeName(ctx context.Context, userTypeID uuid.UUID, name string) error
	UpdateUserTypeDescription(ctx context.Context, userTypeID uuid.UUID, description *string) error
	UpdateUserTypeOptions(ctx context.Context, userTypeID uuid.UUID, options *entry.UserOptions) error

	RemoveUserTypeByID(ctx context.Context, userTypeID uuid.UUID) error
	RemoveUserTypesByIDs(ctx context.Context, userTypeIDs []uuid.UUID) error
}

type AttributeTypesDB interface {
	GetAttributeTypes(ctx context.Context) ([]*entry.AttributeType, error)

	UpsertAttributeType(ctx context.Context, attributeType *entry.AttributeType) error
	UpsertAttributeTypes(ctx context.Context, attributeTypes []*entry.AttributeType) error

	UpdateAttributeTypeName(ctx context.Context, attributeTypeID entry.AttributeTypeID, name string) error
	UpdateAttributeTypeDescription(
		ctx context.Context, attributeTypeID entry.AttributeTypeID, description *string,
	) error
	UpdateAttributeTypeOptions(
		ctx context.Context, attributeTypeID entry.AttributeTypeID, options *entry.AttributeOptions,
	) error

	RemoveAttributeTypeByName(ctx context.Context, name string) error
	RemoveAttributeTypesByNames(ctx context.Context, names []string) error
	RemoveAttributeTypesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	RemoveAttributeTypeByID(ctx context.Context, attributeTypeID entry.AttributeTypeID) error
	RemoveAttributeTypesByIDs(ctx context.Context, attributeTypeIDs []entry.AttributeTypeID) error
}

type NodeAttributesDB interface {
	GetNodeAttributes(ctx context.Context) ([]*entry.NodeAttribute, error)
	GetNodeAttributeByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) (*entry.NodeAttribute, error)
	GetNodeAttributeValueByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) (*entry.AttributeValue, error)
	GetNodeAttributeOptionsByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) (*entry.AttributeOptions, error)

	UpsertNodeAttribute(ctx context.Context, nodeAttribute *entry.NodeAttribute) error
	UpsertNodeAttributes(ctx context.Context, nodeAttributes []*entry.NodeAttribute) error

	UpdateNodeAttributeValue(
		ctx context.Context, attributeID entry.AttributeID, value *entry.AttributeValue,
	) error
	UpdateNodeAttributeOptions(
		ctx context.Context, attributeID entry.AttributeID, options *entry.AttributeOptions,
	) error

	RemoveNodeAttributeByName(ctx context.Context, name string) error
	RemoveNodeAttributesByNames(ctx context.Context, names []string) error
	RemoveNodeAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveNodeAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
}

type SpaceAttributesDB interface {
	SpaceAttributesGetSpaceAttributes(ctx context.Context) ([]*entry.SpaceAttribute, error)
	SpaceAttributesGetSpaceAttributesByPluginIDAndName(
		ctx context.Context, pluginID uuid.UUID, name string,
	) ([]*entry.SpaceAttribute, error)
	SpaceAttributesGetSpaceAttributeByID(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID,
	) (*entry.SpaceAttribute, error)
	SpaceAttributesGetSpaceAttributesBySpaceID(ctx context.Context, spaceID uuid.UUID) ([]*entry.SpaceAttribute, error)

	SpaceAttributesUpsertSpaceAttribute(ctx context.Context, spaceAttribute *entry.SpaceAttribute) error
	SpaceAttributesUpsertSpaceAttributes(ctx context.Context, spaceAttributes []*entry.SpaceAttribute) error

	SpaceAttributesUpdateSpaceAttributeValue(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID, value *entry.AttributeValue,
	) error
	SpaceAttributesUpdateSpaceAttributeOptions(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID, options *entry.AttributeOptions,
	) error

	SpaceAttributesRemoveSpaceAttributeByID(ctx context.Context, spaceAttributeID entry.SpaceAttributeID) error
	SpaceAttributesRemoveSpaceAttributeByName(ctx context.Context, name string) error
	SpaceAttributesRemoveSpaceAttributesByNames(ctx context.Context, names []string) error
	SpaceAttributesRemoveSpaceAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	SpaceAttributesRemoveSpaceAttributeBySpaceID(ctx context.Context, spaceID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByNameAndSpaceID(ctx context.Context, name string, spaceID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByNamesAndSpaceID(ctx context.Context, names []string, spaceID uuid.UUID) error
	SpaceAttributesRemoveSpaceAttributeByPluginIDAndSpaceID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID,
	) error
}

type SpaceUserAttributesDB interface {
	SpaceUserAttributesGetSpaceUserAttributes(ctx context.Context) ([]*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributeByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) (*entry.SpaceUserAttribute, error)
	SpaceUserAttributesGetSpaceUserAttributePayloadByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) (*entry.AttributePayload, error)
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

	SpaceUserAttributesGetSpaceUserAttributesCount(ctx context.Context) (int64, error)

	SpaceUserAttributesUpsertSpaceUserAttribute(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
		modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	SpaceUserAttributesUpdateSpaceUserAttributeOptions(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
		modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)
	SpaceUserAttributesUpdateSpaceUserAttributeValue(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)

	SpaceUserAttributesRemoveSpaceUserAttributeByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) error
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
}

type UserAttributesDB interface {
	UserAttributesGetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error)
	UserAttributesGetUserAttributesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserAttribute, error)
	UserAttributesGetUserAttributeByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.UserAttribute, error)
	UserAttributesGetUserAttributePayloadByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.AttributePayload, error)
	UserAttributesGetUserAttributeValueByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.AttributeValue, error)
	UserAttributesGetUserAttributeOptionsByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.AttributeOptions, error)

	UserAttributesGetUserAttributesCount(ctx context.Context) (int64, error)

	UserAttributesUpsertUserAttribute(
		ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	UserAttributesUpdateUserAttributeValue(
		ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
	UserAttributesUpdateUserAttributeOptions(
		ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)

	UserAttributesRemoveUserAttributeByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) error
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
}

type UserUserAttributesDB interface {
	UserUserAttributesGetUserUserAttributes(ctx context.Context) ([]*entry.UserUserAttribute, error)
	UserUserAttributesGetUserUserAttributeByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) (*entry.UserUserAttribute, error)
	UserUserAttributesGetUserUserAttributePayloadByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) (*entry.AttributePayload, error)
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

	UserUserAttributesGetUserUserAttributesCount(ctx context.Context) (int64, error)

	UserUserAttributesUpsertUserUserAttribute(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	UserUserAttributesUpdateUserUserAttributeValue(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
	UserUserAttributesUpdateUserUserAttributeOptions(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)

	UserUserAttributesRemoveUserUserAttributeByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) error
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
}

type UserSpaceDB interface {
	UserSpaceGetUserSpaces(ctx context.Context) ([]*entry.UserSpace, error)
	UserSpaceGetUserSpaceByID(ctx context.Context, userSpaceID entry.UserSpaceID) (*entry.UserSpace, error)
	UserSpaceGetUserSpacesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserSpace, error)
	UserSpaceGetUserSpacesBySpaceID(ctx context.Context, spaceID uuid.UUID) ([]*entry.UserSpace, error)
	UserSpaceGetValueByID(ctx context.Context, userSpaceID entry.UserSpaceID) (*entry.UserSpaceValue, error)

	UserSpaceGetSpaceIndirectAdmins(ctx context.Context, spaceID uuid.UUID) ([]*uuid.UUID, error)
	UserSpaceCheckIsUserIndirectSpaceAdmin(ctx context.Context, userID, spaceID uuid.UUID) (bool, error)

	UserSpacesUpsertUserSpace(ctx context.Context, userSpace *entry.UserSpace) error
	UserSpacesUpsertUserSpaces(ctx context.Context, userSpaces []*entry.UserSpace) error

	UserSpaceUpdateValueByID(
		ctx context.Context, userSpaceID entry.UserSpaceID, modifyFn modify.Fn[entry.UserSpaceValue],
	) (*entry.UserSpaceValue, error)

	UserSpaceRemoveUserSpace(ctx context.Context, userSpace *entry.UserSpace) error
	UserSpaceRemoveUserSpaces(ctx context.Context, userSpaces []*entry.UserSpace) error
}
