package database

import (
	"context"

	"github.com/google/uuid"

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
	GetWorldIDs(ctx context.Context) ([]uuid.UUID, error)
	GetWorlds(ctx context.Context) ([]*entry.Object, error)
}

type ObjectsDB interface {
	GetObjectByID(ctx context.Context, objectID uuid.UUID) (*entry.Object, error)
	GetObjectIDsByParentID(ctx context.Context, parentID uuid.UUID) ([]uuid.UUID, error)
	GetObjectsByParentID(ctx context.Context, parentID uuid.UUID) ([]*entry.Object, error)

	UpsertObject(ctx context.Context, object *entry.Object) error
	UpsertObjects(ctx context.Context, objects []*entry.Object) error

	UpdateObjectParentID(ctx context.Context, objectID uuid.UUID, parentID uuid.UUID) error
	UpdateObjectPosition(ctx context.Context, objectID uuid.UUID, position *cmath.ObjectPosition) error
	UpdateObjectOwnerID(ctx context.Context, objectID, ownerID uuid.UUID) error
	UpdateObjectAsset2dID(ctx context.Context, objectID uuid.UUID, asset2dID *uuid.UUID) error
	UpdateObjectAsset3dID(ctx context.Context, objectID uuid.UUID, asset3dID *uuid.UUID) error
	UpdateObjectObjectTypeID(ctx context.Context, objectID, objectTypeID uuid.UUID) error
	UpdateObjectOptions(ctx context.Context, objectID uuid.UUID, options *entry.ObjectOptions) error

	RemoveObjectByID(ctx context.Context, objectID uuid.UUID) error
	RemoveObjectsByIDs(ctx context.Context, objectIDs []uuid.UUID) error
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

	RemoveUserByID(ctx context.Context, userID uuid.UUID) error
	RemoveUsersByIDs(ctx context.Context, userID []uuid.UUID) error
}

type Assets2dDB interface {
	GetAssets(ctx context.Context) ([]*entry.Asset2d, error)

	UpsertAsset(ctx context.Context, asset2d *entry.Asset2d) error
	UpsertAssets(ctx context.Context, assets2d []*entry.Asset2d) error

	UpdateAssetMeta(ctx context.Context, asset2dID uuid.UUID, meta entry.Asset2dMeta) error
	UpdateAssetOptions(ctx context.Context, asset2dID uuid.UUID, options *entry.Asset2dOptions) error

	RemoveAssetByID(ctx context.Context, asset2dID uuid.UUID) error
	RemoveAssetsByIDs(ctx context.Context, asset2dIDs []uuid.UUID) error
}

type Assets3dDB interface {
	GetAssets(ctx context.Context) ([]*entry.Asset3d, error)

	UpsertAsset(ctx context.Context, asset3d *entry.Asset3d) error
	UpsertAssets(ctx context.Context, assets3d []*entry.Asset3d) error

	UpdateAssetMeta(ctx context.Context, asset3dID uuid.UUID, meta entry.Asset3dMeta) error
	UpdateAssetOptions(ctx context.Context, asset3dID uuid.UUID, options *entry.Asset3dOptions) error

	RemoveAssetByID(ctx context.Context, asset3dID uuid.UUID) error
	RemoveAssetsByIDs(ctx context.Context, asset3dIDs []uuid.UUID) error
}

type PluginsDB interface {
	GetPlugins(ctx context.Context) ([]*entry.Plugin, error)

	UpsertPlugin(ctx context.Context, plugin *entry.Plugin) error
	UpsertPlugins(ctx context.Context, plugins []*entry.Plugin) error

	UpdatePluginMeta(ctx context.Context, pluginID uuid.UUID, meta entry.PluginMeta) error
	UpdatePluginOptions(ctx context.Context, pluginID uuid.UUID, options *entry.PluginOptions) error

	RemovePluginByID(ctx context.Context, pluginID uuid.UUID) error
	RemovePluginsByIDs(ctx context.Context, pluginIDs []uuid.UUID) error
}

type UserObjectsDB interface {
	GetUserObjects(ctx context.Context) ([]*entry.UserObject, error)
	GetUserObjectByID(ctx context.Context, userObjectID entry.UserObjectID) (*entry.UserObject, error)
	GetUserObjectsByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserObject, error)
	GetUserObjectsByObjectID(ctx context.Context, objectID uuid.UUID) ([]*entry.UserObject, error)
	GetUserObjectValueByID(ctx context.Context, userObjectID entry.UserObjectID) (*entry.UserObjectValue, error)

	GetObjectIndirectAdmins(ctx context.Context, objectID uuid.UUID) ([]*uuid.UUID, error)
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

	UpdateObjectTypeName(ctx context.Context, objectTypeID uuid.UUID, name string) error
	UpdateObjectTypeCategoryName(ctx context.Context, objectTypeID uuid.UUID, categoryName string) error
	UpdateObjectTypeDescription(ctx context.Context, objectTypeID uuid.UUID, description *string) error
	UpdateObjectTypeOptions(ctx context.Context, objectTypeID uuid.UUID, options *entry.ObjectOptions) error

	RemoveObjectTypeByID(ctx context.Context, objectTypeID uuid.UUID) error
	RemoveObjectTypesByIDs(ctx context.Context, objectTypeIDs []uuid.UUID) error
}

type UserTypesDB interface {
	GetUserTypes(ctx context.Context) ([]*entry.UserType, error)

	UpsertUserType(ctx context.Context, userType *entry.UserType) error
	UpsertUserTypes(ctx context.Context, userTypes []*entry.UserType) error

	UpdateUserTypeName(ctx context.Context, userTypeID uuid.UUID, name string) error
	UpdateUserTypeDescription(ctx context.Context, userTypeID uuid.UUID, description string) error
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

	RemoveAttributeTypeByID(ctx context.Context, attributeTypeID entry.AttributeTypeID) error
	RemoveAttributeTypesByIDs(ctx context.Context, attributeTypeIDs []entry.AttributeTypeID) error
	RemoveAttributeTypesByName(ctx context.Context, name string) error
	RemoveAttributeTypesByNames(ctx context.Context, names []string) error
	RemoveAttributeTypesByPluginID(ctx context.Context, pluginID uuid.UUID) error
}

type NodeAttributesDB interface {
	GetNodeAttributes(ctx context.Context) ([]*entry.NodeAttribute, error)
	GetNodeAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) (*entry.NodeAttribute, error)
	GetNodeAttributeValueByAttributeID(ctx context.Context, attributeID entry.AttributeID) (*entry.AttributeValue, error)
	GetNodeAttributeOptionsByAttributeID(
		ctx context.Context, attributeID entry.AttributeID) (*entry.AttributeOptions, error,
	)

	UpsertNodeAttribute(ctx context.Context, nodeAttribute *entry.NodeAttribute) error
	UpsertNodeAttributes(ctx context.Context, nodeAttributes []*entry.NodeAttribute) error

	UpdateNodeAttributeValue(ctx context.Context, attributeID entry.AttributeID, value *entry.AttributeValue) error
	UpdateNodeAttributeOptions(ctx context.Context, attributeID entry.AttributeID, options *entry.AttributeOptions) error

	RemoveNodeAttributeByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveNodeAttributesByName(ctx context.Context, name string) error
	RemoveNodeAttributesByNames(ctx context.Context, names []string) error
	RemoveNodeAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
}

type ObjectAttributesDB interface {
	GetObjectAttributes(ctx context.Context) ([]*entry.ObjectAttribute, error)
	GetObjectAttributeByID(ctx context.Context, objectAttributeID entry.ObjectAttributeID) (*entry.ObjectAttribute, error)
	GetObjectAttributesByObjectID(ctx context.Context, objectID uuid.UUID) ([]*entry.ObjectAttribute, error)
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
	RemoveObjectAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	RemoveObjectAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveObjectAttributesByObjectID(ctx context.Context, objectID uuid.UUID) error
	RemoveObjectAttributesByNameAndObjectID(ctx context.Context, name string, objectID uuid.UUID) error
	RemoveObjectAttributesByNamesAndObjectID(ctx context.Context, names []string, objectID uuid.UUID) error
	RemoveObjectAttributesByPluginIDAndObjectID(ctx context.Context, pluginID uuid.UUID, objectID uuid.UUID) error
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
	GetObjectUserAttributesByObjectID(ctx context.Context, objectID uuid.UUID) ([]*entry.ObjectUserAttribute, error)
	GetObjectUserAttributesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.ObjectUserAttribute, error)
	GetObjectUserAttributesByObjectIDAndUserID(
		ctx context.Context, objectID uuid.UUID, userID uuid.UUID,
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
	RemoveObjectUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	RemoveObjectUserAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveObjectUserAttributesByObjectID(ctx context.Context, objectID uuid.UUID) error
	RemoveObjectUserAttributesByNameAndObjectID(ctx context.Context, name string, objectID uuid.UUID) error
	RemoveObjectUserAttributesByNamesAndObjectID(ctx context.Context, names []string, objectID uuid.UUID) error
	RemoveObjectUserAttributesByUserID(ctx context.Context, userID uuid.UUID) error
	RemoveObjectUserAttributesByNameAndUserID(ctx context.Context, name string, userID uuid.UUID) error
	RemoveObjectUserAttributesByNamesAndUserID(ctx context.Context, names []string, userID uuid.UUID) error
	RemoveObjectUserAttributesByObjectIDAndUserID(ctx context.Context, objectID uuid.UUID, userID uuid.UUID) error
	RemoveObjectUserAttributesByPluginIDAndObjectID(ctx context.Context, pluginID uuid.UUID, objectID uuid.UUID) error
	RemoveObjectUserAttributesByObjectAttributeID(ctx context.Context, objectAttributeID entry.ObjectAttributeID) error
	RemoveObjectUserAttributesByPluginIDAndUserID(ctx context.Context, pluginID uuid.UUID, userID uuid.UUID) error
	RemoveObjectUserAttributesByUserAttributeID(ctx context.Context, userAttributeID entry.UserAttributeID) error
	RemoveObjectUserAttributesByNameAndObjectIDAndUserID(
		ctx context.Context, name string, objectID uuid.UUID, userID uuid.UUID,
	) error
	RemoveObjectUserAttributesByNamesAndObjectIDAndUserID(
		ctx context.Context, names []string, objectID uuid.UUID, userID uuid.UUID,
	) error
	RemoveObjectUserAttributesByPluginIDAndObjectIDAndUserID(
		ctx context.Context, pluginID uuid.UUID, objectID uuid.UUID, userID uuid.UUID,
	) error
}

type UserAttributesDB interface {
	GetUserAttributes(ctx context.Context) ([]*entry.UserAttribute, error)
	GetUserAttributeByID(ctx context.Context, userAttributeID entry.UserAttributeID) (*entry.UserAttribute, error)
	GetUserAttributesByUserID(ctx context.Context, userID uuid.UUID) ([]*entry.UserAttribute, error)
	GetUserAttributePayloadByID(ctx context.Context, userAttributeID entry.UserAttributeID) (*entry.AttributePayload, error)
	GetUserAttributeValueByID(ctx context.Context, userAttributeID entry.UserAttributeID) (*entry.AttributeValue, error)
	GetUserAttributeOptionsByID(ctx context.Context, userAttributeID entry.UserAttributeID) (*entry.AttributeOptions, error)

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
	RemoveUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	RemoveUserAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveUserAttributesByUserID(ctx context.Context, userID uuid.UUID) error
	RemoveUserAttributesByNameAndUserID(ctx context.Context, name string, userID uuid.UUID) error
	RemoveUserAttributesByNamesAndUserID(ctx context.Context, names []string, userID uuid.UUID) error
	RemoveUserAttributesByPluginIDAndUserID(ctx context.Context, pluginID uuid.UUID, userID uuid.UUID) error
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
	RemoveUserUserAttributesByPluginID(ctx context.Context, pluginID uuid.UUID) error
	RemoveUserUserAttributesByAttributeID(ctx context.Context, attributeID entry.AttributeID) error
	RemoveUserUserAttributesBySourceUserID(ctx context.Context, sourceUserID uuid.UUID) error
	RemoveUserUserAttributesByNameAndSourceUserID(ctx context.Context, name string, sourceUserID uuid.UUID) error
	RemoveUserUserAttributesByNamesAndSourceUserID(ctx context.Context, names []string, sourceUserID uuid.UUID) error
	RemoveUserUserAttributesByTargetUserID(ctx context.Context, targetUserID uuid.UUID) error
	RemoveUserUserAttributesByNameAndTargetUserID(ctx context.Context, name string, targetUserID uuid.UUID) error
	RemoveUserUserAttributesByNamesAndTargetUserID(ctx context.Context, names []string, targetUserID uuid.UUID) error
	RemoveUserUserAttributesBySourceUserIDAndTargetUserID(
		ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributesByNameAndSourceUserIDAndTargetUserID(
		ctx context.Context, name string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributesByNamesAndSourceUserIDAndTargetUserID(
		ctx context.Context, names []string, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributesByPluginIDAndSourceUserID(
		ctx context.Context, pluginID uuid.UUID, sourceUserID uuid.UUID,
	) error
	RemoveUserUserAttributesBySourceUserAttributeID(
		ctx context.Context, sourceUserAttributeID entry.UserAttributeID,
	) error
	RemoveUserUserAttributesByPluginIDAndTargetUserID(
		ctx context.Context, pluginID uuid.UUID, targetUserID uuid.UUID,
	) error
	RemoveUserUserAttributesByTargetUserAttributeID(
		ctx context.Context, targetUserAttributeID entry.UserAttributeID,
	) error
	RemoveUserUserAttributesByPluginIDAndSourceUserIDAndTargetUserID(
		ctx context.Context, pluginID uuid.UUID, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	) error
}
