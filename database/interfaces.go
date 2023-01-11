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
	GetSpaceAttributesDB() SpaceAttributesDB
	GetSpaceUserAttributesDB() SpaceUserAttributesDB
	GetUserAttributesDB() UserAttributesDB
	GetUserUserAttributesDB() UserUserAttributesDB
	GetUserSpaceDB() UserSpaceDB
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
	GetSpaceAttributes(ctx context.Context) ([]*entry.SpaceAttribute, error)
	GetSpaceAttributesByPluginIDAndName(
		ctx context.Context, pluginID uuid.UUID, name string,
	) ([]*entry.SpaceAttribute, error)
	GetSpaceAttributeByID(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID,
	) (*entry.SpaceAttribute, error)
	GetSpaceAttributesBySpaceID(ctx context.Context, spaceID uuid.UUID) ([]*entry.SpaceAttribute, error)

	UpsertSpaceAttribute(ctx context.Context, spaceAttribute *entry.SpaceAttribute) error
	UpsertSpaceAttributes(ctx context.Context, spaceAttributes []*entry.SpaceAttribute) error

	UpdateSpaceAttributeValue(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID, value *entry.AttributeValue,
	) error
	UpdateSpaceAttributeOptions(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID, options *entry.AttributeOptions,
	) error

	RemoveSpaceAttributeByID(ctx context.Context, spaceAttributeID entry.SpaceAttributeID) error
	RemoveSpaceAttributeByName(ctx context.Context, name string) error
	RemoveSpaceAttributesByNames(ctx context.Context, names []string) error
	RemoveSpaceAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	RemoveSpaceAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveSpaceAttributeBySpaceID(ctx context.Context, spaceID uuid.UUID) error
	RemoveSpaceAttributeByNameAndSpaceID(ctx context.Context, name string, spaceID uuid.UUID) error
	RemoveSpaceAttributeByNamesAndSpaceID(ctx context.Context, names []string, spaceID uuid.UUID) error
	RemoveSpaceAttributeByPluginIDAndSpaceID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID,
	) error
}

type SpaceUserAttributesDB interface {
	GetSpaceUserAttributes(ctx context.Context) ([]*entry.SpaceUserAttribute, error)
	GetSpaceUserAttributeByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) (*entry.SpaceUserAttribute, error)
	GetSpaceUserAttributePayloadByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) (*entry.AttributePayload, error)
	GetSpaceUserAttributeValueByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) (*entry.AttributeValue, error)
	GetSpaceUserAttributeOptionsByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) (*entry.AttributeOptions, error)
	GetSpaceUserAttributesBySpaceID(
		ctx context.Context, spaceID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	GetSpaceUserAttributesByUserID(
		ctx context.Context, userID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	GetSpaceUserAttributesBySpaceIDAndUserID(
		ctx context.Context, spaceID uuid.UUID, userID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)
	GetSpaceUserAttributesByPluginIDAndNameAndSpaceID(
		ctx context.Context, pluginID uuid.UUID, name string, spaceID uuid.UUID,
	) ([]*entry.SpaceUserAttribute, error)

	GetSpaceUserAttributesCount(ctx context.Context) (int64, error)

	UpsertSpaceUserAttribute(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
		modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	UpdateSpaceUserAttributeOptions(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
		modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)
	UpdateSpaceUserAttributeValue(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)

	RemoveSpaceUserAttributeByID(
		ctx context.Context, spaceUserAttributeID entry.SpaceUserAttributeID,
	) error
	RemoveSpaceUserAttributeByName(ctx context.Context, name string) error
	RemoveSpaceUserAttributesByNames(ctx context.Context, names []string) error
	RemoveSpaceUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	RemoveSpaceUserAttributeByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) error
	RemoveSpaceUserAttributeBySpaceID(ctx context.Context, spaceID uuid.UUID) error
	RemoveSpaceUserAttributeByNameAndSpaceID(
		ctx context.Context, name string, spaceID uuid.UUID,
	) error
	RemoveSpaceUserAttributeByNamesAndSpaceID(
		ctx context.Context, names []string, spaceID uuid.UUID,
	) error
	RemoveSpaceUserAttributeByUserID(ctx context.Context, userID uuid.UUID) error
	RemoveSpaceUserAttributeByNameAndUserID(
		ctx context.Context, name string, userID uuid.UUID,
	) error
	RemoveSpaceUserAttributeByNamesAndUserID(
		ctx context.Context, names []string, userID uuid.UUID,
	) error
	RemoveSpaceUserAttributeBySpaceIDAndUserID(
		ctx context.Context, spaceID uuid.UUID, userID uuid.UUID,
	) error
	RemoveSpaceUserAttributeByNameAndSpaceIDAndUserID(
		ctx context.Context, name string, spaceID uuid.UUID, userID uuid.UUID,
	) error
	RemoveSpaceUserAttributeByNamesAndSpaceIDAndUserID(
		ctx context.Context, names []string, spaceID uuid.UUID, userID uuid.UUID,
	) error
	RemoveSpaceUserAttributeByPluginIDAndSpaceID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID,
	) error
	RemoveSpaceUserAttributeBySpaceAttributeID(
		ctx context.Context, spaceAttributeID entry.SpaceAttributeID,
	) error
	RemoveSpaceUserAttributeByPluginIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, userID uuid.UUID,
	) error
	RemoveSpaceUserAttributeByUserAttributeID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) error
	RemoveSpaceUserAttributeByPluginIDAndSpaceIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, spaceID uuid.UUID, userID uuid.UUID,
	) error
}

type UserAttributesDB interface {
	GetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error)
	GetUserAttributesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserAttribute, error)
	GetUserAttributeByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.UserAttribute, error)
	GetUserAttributePayloadByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.AttributePayload, error)
	GetUserAttributeValueByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.AttributeValue, error)
	GetUserAttributeOptionsByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) (*entry.AttributeOptions, error)

	GetUserAttributesCount(ctx context.Context) (int64, error)

	UpsertUserAttribute(
		ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	UpdateUserAttributeValue(
		ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
	UpdateUserAttributeOptions(
		ctx context.Context, userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)

	RemoveUserAttributeByID(
		ctx context.Context, userAttributeID entry.UserAttributeID,
	) error
	RemoveUserAttributeByName(ctx context.Context, name string) error
	RemoveUserAttributesByNames(ctx context.Context, names []string) error
	RemoveUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	RemoveUserAttributeByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) error
	RemoveUserAttributeByUserID(ctx context.Context, userID uuid.UUID) error
	RemoveUserAttributeByNameAndUserID(
		ctx context.Context, name string, userID uuid.UUID,
	) error
	RemoveUserAttributeByNamesAndUserID(
		ctx context.Context, names []string, userID uuid.UUID,
	) error
	RemoveUserAttributeByPluginIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, userID uuid.UUID,
	) error
}

type UserUserAttributesDB interface {
	GetUserUserAttributes(ctx context.Context) ([]*entry.UserUserAttribute, error)
	GetUserUserAttributeByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) (*entry.UserUserAttribute, error)
	GetUserUserAttributePayloadByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) (*entry.AttributePayload, error)
	GetUserUserAttributeValueByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) (*entry.AttributeValue, error)
	GetUserUserAttributeOptionsByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) (*entry.AttributeOptions, error)
	GetUserUserAttributesBySourceUserID(
		ctx context.Context, sourceUserID uuid.UUID,
	) ([]*entry.UserUserAttribute, error)
	GetUserUserAttributesByTargetUserID(
		ctx context.Context, targetUserID uuid.UUID,
	) ([]*entry.UserUserAttribute, error)
	GetUserUserAttributesBySourceUserIDAndTargetUserID(
		ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) ([]*entry.UserUserAttribute, error)

	GetUserUserAttributesCount(ctx context.Context) (int64, error)

	UpsertUserUserAttribute(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	UpdateUserUserAttributeValue(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
	UpdateUserUserAttributeOptions(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)

	RemoveUserUserAttributeByID(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
	) error
	RemoveUserUserAttributeByName(ctx context.Context, name string) error
	RemoveUserUserAttributesByNames(ctx context.Context, names []string) error
	RemoveUserUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	RemoveUserUserAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveUserUserAttributeBySourceUserID(
		ctx context.Context, sourceUserID uuid.UUID,
	) error
	RemoveUserUserAttributeByNameAndSourceUserID(
		ctx context.Context, name string, sourceUserID uuid.UUID,
	) error
	RemoveUserUserAttributeByNamesAndSourceUserID(
		ctx context.Context, names []string, sourceUserID uuid.UUID,
	) error
	RemoveUserUserAttributeByTargetUserID(
		ctx context.Context, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributeByNameAndTargetUserID(
		ctx context.Context, name string, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributeByNamesAndTargetUserID(
		ctx context.Context, names []string, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributeBySourceUserIDAndTargetUserID(
		ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributeByNameAndSourceUserIDAndTargetUserID(
		ctx context.Context, name string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributeByNamesAndSourceUserIDAndTargetUserID(
		ctx context.Context, names []string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributeByPluginIDAndSourceUserID(
		ctx context.Context, pluginID uuid.UUID, sourceUserID uuid.UUID,
	) error
	RemoveUserUserAttributeBySourceUserAttributeID(
		ctx context.Context, sourceUserAttributeID entry.UserAttributeID,
	) error
	RemoveUserUserAttributeByPluginIDAndTargetUserID(
		ctx context.Context, pluginID uuid.UUID, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributeByTargetUserAttributeID(
		ctx context.Context, targetUserAttributeID entry.UserAttributeID,
	) error
	RemoveUserUserAttributeByPluginIDAndSourceUserIDAndTargetUserID(
		ctx context.Context, pluginID uuid.UUID, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
}

type UserSpaceDB interface {
	GetUserSpaces(ctx context.Context) ([]*entry.UserSpace, error)
	GetUserSpaceByID(ctx context.Context, userSpaceID entry.UserSpaceID) (*entry.UserSpace, error)
	GetUserSpacesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserSpace, error)
	GetUserSpacesBySpaceID(ctx context.Context, spaceID uuid.UUID) ([]*entry.UserSpace, error)
	GetValueByID(ctx context.Context, userSpaceID entry.UserSpaceID) (*entry.UserSpaceValue, error)

	GetSpaceIndirectAdmins(ctx context.Context, spaceID uuid.UUID) ([]*uuid.UUID, error)
	CheckIsUserIndirectSpaceAdmin(ctx context.Context, userID, spaceID uuid.UUID) (bool, error)

	UpsertUserSpace(ctx context.Context, userSpace *entry.UserSpace) error
	UpsertUserSpaces(ctx context.Context, userSpaces []*entry.UserSpace) error

	UpdateValueByID(
		ctx context.Context, userSpaceID entry.UserSpaceID, modifyFn modify.Fn[entry.UserSpaceValue],
	) (*entry.UserSpaceValue, error)

	RemoveUserSpace(ctx context.Context, userSpace *entry.UserSpace) error
	RemoveUserSpaces(ctx context.Context, userSpaces []*entry.UserSpace) error
}
