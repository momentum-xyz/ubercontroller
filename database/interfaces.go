package database

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/utils/mid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type DB interface {
	GetCommonDB() CommonDB
	GetNodesDB() NodesDB
	GetWorldsDB() WorldsDB
	GetObjectsDB() ObjectsDB
	GetUsersDB() UsersDB
	GetAssets2dDB() Assets2dDB
	GetAssets3dDB() Assets3dDB
	GetPluginsDB() PluginsDB
	GetUserObjectsDB() UserObjectsDB
	GetObjectTypesDB() ObjectTypesDB
	GetUserTypesDB() UserTypesDB
	GetAttributeTypesDB() AttributeTypesDB
	GetNodeAttributesDB() NodeAttributesDB
	GetObjectAttributesDB() ObjectAttributesDB
	GetObjectUserAttributesDB() ObjectUserAttributesDB
	GetUserAttributesDB() UserAttributesDB
	GetUserUserAttributesDB() UserUserAttributesDB
}

type CommonDB interface {
}

type NodesDB interface {
	GetNode(ctx context.Context) (*entry.Node, error)
}

type WorldsDB interface {
	GetWorldIDs(ctx context.Context) ([]mid.ID, error)
	GetWorlds(ctx context.Context) ([]*entry.Object, error)
}

type ObjectsDB interface {
	GetObjectByID(ctx context.Context, objectID mid.ID) (*entry.Object, error)
	GetObjectIDsByParentID(ctx context.Context, parentID mid.ID) ([]mid.ID, error)
	GetObjectsByParentID(ctx context.Context, parentID mid.ID) ([]*entry.Object, error)

	UpsertObject(ctx context.Context, object *entry.Object) error
	UpsertObjects(ctx context.Context, objects []*entry.Object) error

	UpdateObjectParentID(ctx context.Context, objectID mid.ID, parentID mid.ID) error
	UpdateObjectPosition(ctx context.Context, objectID mid.ID, position *cmath.ObjectTransform) error
	UpdateObjectOwnerID(ctx context.Context, objectID, ownerID mid.ID) error
	UpdateObjectAsset2dID(ctx context.Context, objectID mid.ID, asset2dID *mid.ID) error
	UpdateObjectAsset3dID(ctx context.Context, objectID mid.ID, asset3dID *mid.ID) error
	UpdateObjectObjectTypeID(ctx context.Context, objectID, objectTypeID mid.ID) error
	UpdateObjectOptions(ctx context.Context, objectID mid.ID, options *entry.ObjectOptions) error

	RemoveObjectByID(ctx context.Context, objectID mid.ID) error
	RemoveObjectsByIDs(ctx context.Context, objectIDs []mid.ID) error
}

type UsersDB interface {
	GetUserByID(ctx context.Context, userID mid.ID) (*entry.User, error)
	GetUsersByIDs(ctx context.Context, userIDs []mid.ID) ([]*entry.User, error)
	GetUserByWallet(ctx context.Context, wallet string) (*entry.User, error)
	GetUserWalletByUserID(ctx context.Context, userID mid.ID) (*string, error)
	GetUserProfileByUserID(ctx context.Context, userID mid.ID) (*entry.UserProfile, error)

	CheckIsUserExistsByName(ctx context.Context, name string) (bool, error)

	UpsertUser(ctx context.Context, user *entry.User) error
	UpsertUsers(ctx context.Context, user []*entry.User) error

	UpdateUserUserTypeID(ctx context.Context, userID, userTypeID mid.ID) error
	UpdateUserOptions(ctx context.Context, userID mid.ID, options *entry.UserOptions) error
	UpdateUserProfile(ctx context.Context, userID mid.ID, profile *entry.UserProfile) error

	RemoveUserByID(ctx context.Context, userID mid.ID) error
	RemoveUsersByIDs(ctx context.Context, userID []mid.ID) error
}

type Assets2dDB interface {
	GetAssets(ctx context.Context) ([]*entry.Asset2d, error)

	UpsertAsset(ctx context.Context, asset2d *entry.Asset2d) error
	UpsertAssets(ctx context.Context, assets2d []*entry.Asset2d) error

	UpdateAssetMeta(ctx context.Context, asset2dID mid.ID, meta entry.Asset2dMeta) error
	UpdateAssetOptions(ctx context.Context, asset2dID mid.ID, options *entry.Asset2dOptions) error

	RemoveAssetByID(ctx context.Context, asset2dID mid.ID) error
	RemoveAssetsByIDs(ctx context.Context, asset2dIDs []mid.ID) error
}

type Assets3dDB interface {
	GetAssets(ctx context.Context) ([]*entry.Asset3d, error)

	UpsertAsset(ctx context.Context, asset3d *entry.Asset3d) error
	UpsertAssets(ctx context.Context, assets3d []*entry.Asset3d) error

	UpdateAssetMeta(ctx context.Context, asset3dID mid.ID, meta *entry.Asset3dMeta) error
	UpdateAssetOptions(ctx context.Context, asset3dID mid.ID, options *entry.Asset3dOptions) error

	RemoveAssetByID(ctx context.Context, asset3dID mid.ID) error
	RemoveAssetsByIDs(ctx context.Context, asset3dIDs []mid.ID) error
}

type PluginsDB interface {
	GetPlugins(ctx context.Context) ([]*entry.Plugin, error)

	UpsertPlugin(ctx context.Context, plugin *entry.Plugin) error
	UpsertPlugins(ctx context.Context, plugins []*entry.Plugin) error

	UpdatePluginMeta(ctx context.Context, pluginID mid.ID, meta entry.PluginMeta) error
	UpdatePluginOptions(ctx context.Context, pluginID mid.ID, options *entry.PluginOptions) error

	RemovePluginByID(ctx context.Context, pluginID mid.ID) error
	RemovePluginsByIDs(ctx context.Context, pluginIDs []mid.ID) error
}

type UserObjectsDB interface {
	GetUserObjects(ctx context.Context) ([]*entry.UserObject, error)
	GetUserObjectByID(ctx context.Context, userObjectID entry.UserObjectID) (*entry.UserObject, error)
	GetUserObjectsByUserID(ctx context.Context, userID mid.ID) ([]*entry.UserObject, error)
	GetUserObjectsByObjectID(ctx context.Context, objectID mid.ID) ([]*entry.UserObject, error)
	GetUserObjectValueByID(ctx context.Context, userObjectID entry.UserObjectID) (*entry.UserObjectValue, error)

	GetObjectIndirectAdmins(ctx context.Context, objectID mid.ID) ([]*mid.ID, error)
	CheckIsIndirectAdminByID(ctx context.Context, userObjectID entry.UserObjectID) (bool, error)

	UpsertUserObject(
		ctx context.Context, userObjectID entry.UserObjectID,
		modifyFn modify.Fn[entry.UserObjectValue],
	) (*entry.UserObjectValue, error)

	UpdateUserObjectValue(
		ctx context.Context, userObjectID entry.UserObjectID,
		modifyFn modify.Fn[entry.UserObjectValue],
	) (*entry.UserObjectValue, error)

	RemoveUserObjectByID(ctx context.Context, userObjectID entry.UserObjectID) error
	RemoveUserObjectsByIDs(ctx context.Context, userObjectIDs []entry.UserObjectID) error
}

type ObjectTypesDB interface {
	GetObjectTypes(ctx context.Context) ([]*entry.ObjectType, error)

	UpsertObjectType(ctx context.Context, objectType *entry.ObjectType) error
	UpsertObjectTypes(ctx context.Context, objectTypes []*entry.ObjectType) error

	UpdateObjectTypeName(ctx context.Context, objectTypeID mid.ID, name string) error
	UpdateObjectTypeCategoryName(ctx context.Context, objectTypeID mid.ID, categoryName string) error
	UpdateObjectTypeDescription(ctx context.Context, objectTypeID mid.ID, description *string) error
	UpdateObjectTypeOptions(ctx context.Context, objectTypeID mid.ID, options *entry.ObjectOptions) error

	RemoveObjectTypeByID(ctx context.Context, objectTypeID mid.ID) error
	RemoveObjectTypesByIDs(ctx context.Context, objectTypeIDs []mid.ID) error
}

type UserTypesDB interface {
	GetUserTypes(ctx context.Context) ([]*entry.UserType, error)

	UpsertUserType(ctx context.Context, userType *entry.UserType) error
	UpsertUserTypes(ctx context.Context, userTypes []*entry.UserType) error

	UpdateUserTypeName(ctx context.Context, userTypeID mid.ID, name string) error
	UpdateUserTypeDescription(ctx context.Context, userTypeID mid.ID, description string) error
	UpdateUserTypeOptions(ctx context.Context, userTypeID mid.ID, options *entry.UserOptions) error

	RemoveUserTypeByID(ctx context.Context, userTypeID mid.ID) error
	RemoveUserTypesByIDs(ctx context.Context, userTypeIDs []mid.ID) error
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

	RemoveAttributeTypeByID(ctx context.Context, attributeTypeID entry.AttributeTypeID) error
	RemoveAttributeTypesByIDs(ctx context.Context, attributeTypeIDs []entry.AttributeTypeID) error
	RemoveAttributeTypesByName(ctx context.Context, name string) error
	RemoveAttributeTypesByNames(ctx context.Context, names []string) error
	RemoveAttributeTypesByPluginID(ctx context.Context, pluginID mid.ID) error
}

type NodeAttributesDB interface {
	GetNodeAttributes(ctx context.Context) ([]*entry.NodeAttribute, error)
	GetNodeAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) (*entry.NodeAttribute, error)
	GetNodeAttributeValueByAttributeID(ctx context.Context, attributeID entry.AttributeID) (
		*entry.AttributeValue, error,
	)
	GetNodeAttributeOptionsByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) (
		*entry.AttributeOptions, error,
	)

	UpsertNodeAttribute(ctx context.Context, nodeAttribute *entry.NodeAttribute) error
	UpsertNodeAttributes(ctx context.Context, nodeAttributes []*entry.NodeAttribute) error

	UpdateNodeAttributeValue(ctx context.Context, attributeID entry.AttributeID, value *entry.AttributeValue) error
	UpdateNodeAttributeOptions(
		ctx context.Context, attributeID entry.AttributeID, options *entry.AttributeOptions,
	) error

	RemoveNodeAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveNodeAttributesByName(ctx context.Context, name string) error
	RemoveNodeAttributesByNames(ctx context.Context, names []string) error
	RemoveNodeAttributesByPluginID(ctx context.Context, pluginID mid.ID) error
}

type ObjectAttributesDB interface {
	GetObjectAttributes(ctx context.Context) ([]*entry.ObjectAttribute, error)
	GetObjectAttributeByID(ctx context.Context, objectAttributeID entry.ObjectAttributeID) (
		*entry.ObjectAttribute, error,
	)
	GetObjectAttributesByObjectID(ctx context.Context, objectID mid.ID) ([]*entry.ObjectAttribute, error)
	GetObjectAttributesByAttributeID(
		ctx context.Context, attributeID entry.AttributeID,
	) ([]*entry.ObjectAttribute, error)

	UpsertObjectAttribute(ctx context.Context, objectAttribute *entry.ObjectAttribute) error
	UpsertObjectAttributes(ctx context.Context, objectAttributes []*entry.ObjectAttribute) error

	UpdateObjectAttributeValue(
		ctx context.Context, objectAttributeID entry.ObjectAttributeID, value *entry.AttributeValue,
	) error
	UpdateObjectAttributeOptions(
		ctx context.Context, objectAttributeID entry.ObjectAttributeID, options *entry.AttributeOptions,
	) error

	RemoveObjectAttributeByID(ctx context.Context, objectAttributeID entry.ObjectAttributeID) error
	RemoveObjectAttributesByName(ctx context.Context, name string) error
	RemoveObjectAttributesByNames(ctx context.Context, names []string) error
	RemoveObjectAttributesByPluginID(ctx context.Context, pluginID mid.ID) error
	RemoveObjectAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveObjectAttributesByObjectID(ctx context.Context, objectID mid.ID) error
	RemoveObjectAttributesByNameAndObjectID(ctx context.Context, name string, objectID mid.ID) error
	RemoveObjectAttributesByNamesAndObjectID(ctx context.Context, names []string, objectID mid.ID) error
	RemoveObjectAttributesByPluginIDAndObjectID(ctx context.Context, pluginID mid.ID, objectID mid.ID) error
}

type ObjectUserAttributesDB interface {
	GetObjectUserAttributes(ctx context.Context) ([]*entry.ObjectUserAttribute, error)
	GetObjectUserAttributeByID(
		ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
	) (*entry.ObjectUserAttribute, error)
	GetObjectUserAttributePayloadByID(
		ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
	) (*entry.AttributePayload, error)
	GetObjectUserAttributeValueByID(
		ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
	) (*entry.AttributeValue, error)
	GetObjectUserAttributeOptionsByID(
		ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
	) (*entry.AttributeOptions, error)
	GetObjectUserAttributesByObjectID(ctx context.Context, objectID mid.ID) ([]*entry.ObjectUserAttribute, error)
	GetObjectUserAttributesByUserID(ctx context.Context, userID mid.ID) ([]*entry.ObjectUserAttribute, error)
	GetObjectUserAttributesByObjectIDAndUserID(
		ctx context.Context, objectID mid.ID, userID mid.ID,
	) ([]*entry.ObjectUserAttribute, error)
	GetObjectUserAttributesByObjectAttributeID(
		ctx context.Context, objectAttributeID entry.ObjectAttributeID,
	) ([]*entry.ObjectUserAttribute, error)

	GetObjectUserAttributesCount(ctx context.Context) (int64, error)

	UpsertObjectUserAttribute(
		ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
		modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	UpdateObjectUserAttributeOptions(
		ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
		modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)
	UpdateObjectUserAttributeValue(
		ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID,
		modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)

	RemoveObjectUserAttributeByID(ctx context.Context, objectUserAttributeID entry.ObjectUserAttributeID) error
	RemoveObjectUserAttributesByName(ctx context.Context, name string) error
	RemoveObjectUserAttributesByNames(ctx context.Context, names []string) error
	RemoveObjectUserAttributesByPluginID(ctx context.Context, pluginID mid.ID) error
	RemoveObjectUserAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveObjectUserAttributesByObjectID(ctx context.Context, objectID mid.ID) error
	RemoveObjectUserAttributesByNameAndObjectID(ctx context.Context, name string, objectID mid.ID) error
	RemoveObjectUserAttributesByNamesAndObjectID(ctx context.Context, names []string, objectID mid.ID) error
	RemoveObjectUserAttributesByUserID(ctx context.Context, userID mid.ID) error
	RemoveObjectUserAttributesByNameAndUserID(ctx context.Context, name string, userID mid.ID) error
	RemoveObjectUserAttributesByNamesAndUserID(ctx context.Context, names []string, userID mid.ID) error
	RemoveObjectUserAttributesByObjectIDAndUserID(ctx context.Context, objectID mid.ID, userID mid.ID) error
	RemoveObjectUserAttributesByPluginIDAndObjectID(ctx context.Context, pluginID mid.ID, objectID mid.ID) error
	RemoveObjectUserAttributesByObjectAttributeID(ctx context.Context, objectAttributeID entry.ObjectAttributeID) error
	RemoveObjectUserAttributesByPluginIDAndUserID(ctx context.Context, pluginID mid.ID, userID mid.ID) error
	RemoveObjectUserAttributesByUserAttributeID(ctx context.Context, userAttributeID entry.UserAttributeID) error
	RemoveObjectUserAttributesByNameAndObjectIDAndUserID(
		ctx context.Context, name string, objectID mid.ID, userID mid.ID,
	) error
	RemoveObjectUserAttributesByNamesAndObjectIDAndUserID(
		ctx context.Context, names []string, objectID mid.ID, userID mid.ID,
	) error
	RemoveObjectUserAttributesByPluginIDAndObjectIDAndUserID(
		ctx context.Context, pluginID mid.ID, objectID mid.ID, userID mid.ID,
	) error
}

type UserAttributesDB interface {
	GetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error)
	GetUserAttributeByID(ctx context.Context, userAttributeID entry.UserAttributeID) (*entry.UserAttribute, error)
	GetUserAttributesByUserID(ctx context.Context, userID mid.ID) ([]*entry.UserAttribute, error)
	GetUserAttributePayloadByID(ctx context.Context, userAttributeID entry.UserAttributeID) (
		*entry.AttributePayload, error,
	)
	GetUserAttributeValueByID(ctx context.Context, userAttributeID entry.UserAttributeID) (*entry.AttributeValue, error)
	GetUserAttributeOptionsByID(ctx context.Context, userAttributeID entry.UserAttributeID) (
		*entry.AttributeOptions, error,
	)

	GetUserAttributesCount(ctx context.Context) (int64, error)

	UpsertUserAttribute(
		ctx context.Context, userAttributeID entry.UserAttributeID,
		modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	UpdateUserAttributeValue(
		ctx context.Context, userAttributeID entry.UserAttributeID,
		modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
	UpdateUserAttributeOptions(
		ctx context.Context, userAttributeID entry.UserAttributeID,
		modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)

	RemoveUserAttributeByID(ctx context.Context, userAttributeID entry.UserAttributeID) error
	RemoveUserAttributesByName(ctx context.Context, name string) error
	RemoveUserAttributesByNames(ctx context.Context, names []string) error
	RemoveUserAttributesByPluginID(ctx context.Context, pluginID mid.ID) error
	RemoveUserAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveUserAttributesByUserID(ctx context.Context, userID mid.ID) error
	RemoveUserAttributesByNameAndUserID(ctx context.Context, name string, userID mid.ID) error
	RemoveUserAttributesByNamesAndUserID(ctx context.Context, names []string, userID mid.ID) error
	RemoveUserAttributesByPluginIDAndUserID(ctx context.Context, pluginID mid.ID, userID mid.ID) error
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
		ctx context.Context, sourceUserID mid.ID,
	) ([]*entry.UserUserAttribute, error)
	GetUserUserAttributesByTargetUserID(
		ctx context.Context, targetUserID mid.ID,
	) ([]*entry.UserUserAttribute, error)
	GetUserUserAttributesBySourceUserIDAndTargetUserID(
		ctx context.Context, sourceUserID mid.ID, targetUserID mid.ID,
	) ([]*entry.UserUserAttribute, error)

	GetUserUserAttributesCount(ctx context.Context) (int64, error)

	UpsertUserUserAttribute(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
		modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	UpdateUserUserAttributeValue(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
		modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
	UpdateUserUserAttributeOptions(
		ctx context.Context, userUserAttributeID entry.UserUserAttributeID,
		modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)

	RemoveUserUserAttributeByID(ctx context.Context, userUserAttributeID entry.UserUserAttributeID) error
	RemoveUserUserAttributesByName(ctx context.Context, name string) error
	RemoveUserUserAttributesByNames(ctx context.Context, names []string) error
	RemoveUserUserAttributesByPluginID(ctx context.Context, pluginID mid.ID) error
	RemoveUserUserAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveUserUserAttributesBySourceUserID(ctx context.Context, sourceUserID mid.ID) error
	RemoveUserUserAttributesByNameAndSourceUserID(ctx context.Context, name string, sourceUserID mid.ID) error
	RemoveUserUserAttributesByNamesAndSourceUserID(ctx context.Context, names []string, sourceUserID mid.ID) error
	RemoveUserUserAttributesByTargetUserID(ctx context.Context, targetUserID mid.ID) error
	RemoveUserUserAttributesByNameAndTargetUserID(ctx context.Context, name string, targetUserID mid.ID) error
	RemoveUserUserAttributesByNamesAndTargetUserID(ctx context.Context, names []string, targetUserID mid.ID) error
	RemoveUserUserAttributesBySourceUserIDAndTargetUserID(
		ctx context.Context, sourceUserID mid.ID, targetUserID mid.ID,
	) error
	RemoveUserUserAttributesByNameAndSourceUserIDAndTargetUserID(
		ctx context.Context, name string, sourceUserID mid.ID, targetUserID mid.ID,
	) error
	RemoveUserUserAttributesByNamesAndSourceUserIDAndTargetUserID(
		ctx context.Context, names []string, sourceUserID mid.ID, targetUserID mid.ID,
	) error
	RemoveUserUserAttributesByPluginIDAndSourceUserID(
		ctx context.Context, pluginID mid.ID, sourceUserID mid.ID,
	) error
	RemoveUserUserAttributesBySourceUserAttributeID(
		ctx context.Context, sourceUserAttributeID entry.UserAttributeID,
	) error
	RemoveUserUserAttributesByPluginIDAndTargetUserID(
		ctx context.Context, pluginID mid.ID, targetUserID mid.ID,
	) error
	RemoveUserUserAttributesByTargetUserAttributeID(
		ctx context.Context, targetUserAttributeID entry.UserAttributeID,
	) error
	RemoveUserUserAttributesByPluginIDAndSourceUserIDAndTargetUserID(
		ctx context.Context, pluginID mid.ID, sourceUserID mid.ID, targetUserID mid.ID,
	) error
}
