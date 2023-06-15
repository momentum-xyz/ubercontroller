package universe

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	influxWrite "github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type IDer interface {
	GetID() umid.UMID
}

type Initializer interface {
	Initialize(ctx types.NodeContext) error
}

type Enabler interface {
	GetEnabled() bool
	SetEnabled(enabled bool)
}

type Runner interface {
	Run() error
}

type Stopper interface {
	Stop() error
}

type RunStopper interface {
	Runner
	Stopper
}

type Loader interface {
	Load() error
}

type Saver interface {
	Save() error
}

type LoadSaver interface {
	Loader
	Saver
}

type APIRegister interface {
	RegisterAPI(r *gin.Engine)
}

type ObjectsCacher interface {
	GetAllObjects() map[umid.UMID]Object
	GetObjectFromAllObjects(objectID umid.UMID) (Object, bool)
	FilterAllObjects(predicateFn ObjectsFilterPredicateFn) map[umid.UMID]Object
	AddObjectToAllObjects(object Object) error
	RemoveObjectFromAllObjects(object Object) (bool, error)
}

type CacheInvalidator interface {
	InvalidateCache()
}

type Node interface {
	Object
	Loader
	APIRegister
	ObjectsCacher

	GetDB() database.DB

	GetConfig() *config.Config
	GetLogger() *zap.SugaredLogger

	ToObject() Object

	GetWorlds() Worlds
	GetAssets2d() Assets2d
	GetAssets3d() Assets3d
	GetObjectTypes() ObjectTypes
	GetUserTypes() UserTypes
	GetAttributeTypes() AttributeTypes
	GetPlugins() Plugins

	GetUserObjects() UserObjects

	GetNodeAttributes() NodeAttributes
	GetUserAttributes() UserAttributes
	GetUserUserAttributes() UserUserAttributes
	GetObjectUserAttributes() ObjectUserAttributes

	CreateUsers(
		ctx context.Context, users ...*entry.User,
	) error // TODO: refactor, place Users next to Nodes in a universe

	AddAPIRegister(register APIRegister)

	WriteInfluxPoint(point *influxWrite.Point) error
	LoadUser(userID umid.UMID) (User, error)
}

type Worlds interface {
	RunStopper
	LoadSaver
	APIRegister

	CreateWorld(worldID umid.UMID) (World, error)

	GetWorld(worldID umid.UMID) (World, bool)
	GetWorlds() map[umid.UMID]World

	FilterWorlds(predicateFn WorldsFilterPredicateFn) map[umid.UMID]World

	AddWorld(world World, updateDB bool) error
	AddWorlds(worlds []World, updateDB bool) error

	RemoveWorld(world World, updateDB bool) (bool, error)
	RemoveWorlds(worlds []World, updateDB bool) (bool, error)
}

type World interface {
	Object
	Loader
	ObjectsCacher

	ToObject() Object

	GetSettings() *WorldSettings

	GetTotalStake() uint8

	GetWorldAvatar() string
	GetWebsiteLink() string

	GetCalendar() Calendar

	WriteInfluxPoint(point *influxWrite.Point) error

	TempSetSkybox(msg *websocket.PreparedMessage)
	TempGetSkybox() *websocket.PreparedMessage
}

type Object interface {
	IDer
	Enabler
	Initializer
	RunStopper
	Saver
	CacheInvalidator

	GetWorld() World

	GetName() string
	SetName(name string, updateDB bool) error

	GetDescription() string

	GetParent() Object
	SetParent(parent Object, updateDB bool) error

	GetOwnerID() umid.UMID
	SetOwnerID(ownerID umid.UMID, updateDB bool) error

	GetTransform() *cmath.Transform
	GetActualTransform() *cmath.Transform
	SetTransform(position *cmath.Transform, updateDB bool) error
	SetActualTransform(pos cmath.Transform, theta float64) error

	GetOptions() *entry.ObjectOptions
	GetEffectiveOptions() *entry.ObjectOptions
	SetOptions(modifyFn modify.Fn[entry.ObjectOptions], updateDB bool) (*entry.ObjectOptions, error)

	GetAsset2D() Asset2d
	SetAsset2D(asset2d Asset2d, updateDB bool) error

	GetAsset3D() Asset3d
	SetAsset3D(asset3d Asset3d, updateDB bool) error

	GetObjectType() ObjectType
	SetObjectType(objectType ObjectType, updateDB bool) error

	GetObjectAttributes() ObjectAttributes

	GetEntry() *entry.Object
	LoadFromEntry(entry *entry.Object, recursive bool) error

	Update(recursive bool) error
	UpdateChildrenPosition(recursive bool) error

	CreateObject(objectID umid.UMID) (Object, error)
	GetObject(objectID umid.UMID, recursive bool) (Object, bool)
	GetObjects(recursive bool) map[umid.UMID]Object
	FilterObjects(predicateFn ObjectsFilterPredicateFn, recursive bool) map[umid.UMID]Object
	AddObject(object Object, updateDB bool) error
	AddObjects(objects []Object, updateDB bool) error
	RemoveObject(object Object, recursive, updateDB bool) (bool, error)
	RemoveObjects(objects []Object, recursive, updateDB bool) (bool, error)

	GetUser(userID umid.UMID, recursive bool) (User, bool)
	GetUsers(recursive bool) map[umid.UMID]User
	AddUser(user User, updateDB bool) error
	RemoveUser(user User, updateDB bool) (bool, error)

	Send(msg *websocket.PreparedMessage, recursive bool) error

	SendSpawnMessage(sendFn func(msg *websocket.PreparedMessage) error, recursive bool)
	SendAttributes(sendFn func(*websocket.PreparedMessage), recursive bool)
	SendAllAutoAttributes(sendFn func(msg *websocket.PreparedMessage) error, recursive bool)

	LockUIObject(user User, state uint32) bool

	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

type User interface {
	IDer
	RunStopper

	GetWorld() World
	SetWorld(world World)

	GetObject() Object
	SetObject(object Object)

	GetUserType() UserType
	SetUserType(userType UserType, updateDB bool) error

	GetProfile() *entry.UserProfile

	GetTransform() *cmath.TransformNoScale
	SetTransform(cmath.TransformNoScale)

	GetPosition() cmath.Vec3
	GetRotation() cmath.Vec3
	SetPosition(position cmath.Vec3)

	//GetPosBuffer() []byte
	GetLastPosTime() int64
	GetLastSendPosTime() int64
	SetLastSendPosTime(int64)

	Update() error
	ReleaseSendBuffer()
	LockSendBuffer()

	IsTemporaryUser() (bool, error)
	SetOfflineTimer() (bool, error)
	DeleteTemporaryUser(uid umid.UMID) error

	GetSessionID() umid.UMID
	SetConnection(sessionID umid.UMID, socketConnection *websocket.Conn) error

	Send(message *websocket.PreparedMessage) error
	SendDirectly(message *websocket.PreparedMessage) error

	AddInfluxTags(prefix string, point *influxWrite.Point) *influxWrite.Point
	GetUserDefinition() *posbus.UserData
}

// UserObjects ignores "updateDB" flag
type UserObjects interface {
	GetValue(userObjectID entry.UserObjectID) (*entry.UserObjectValue, bool)

	GetObjectIndirectAdmins(objectID umid.UMID) ([]*umid.UMID, bool)
	CheckIsIndirectAdmin(userObjectID entry.UserObjectID) (bool, error)

	Upsert(
		userObjectID entry.UserObjectID, modifyFn modify.Fn[entry.UserObjectValue], updateDB bool,
	) (*entry.UserObjectValue, error)

	UpdateValue(
		userObjectID entry.UserObjectID, modifyFn modify.Fn[entry.UserObjectValue], updateDB bool,
	) (*entry.UserObjectValue, error)

	Remove(userObjectID entry.UserObjectID, updateDB bool) (bool, error)
	RemoveMany(userObjectIDs []entry.UserObjectID, updateDB bool) (bool, error)
}

type AttributeOptionsGetter[ID comparable] interface {
	// Get the options set directly on this object.
	GetOptions(attributeID ID) (*entry.AttributeOptions, bool)
	// Get the merged options of this object and its parent type.
	GetEffectiveOptions(attributeID ID) (*entry.AttributeOptions, bool)
}

type AttributeUserRoleGetter[T comparable] interface {
	// Retrieve roles a user has on an plugin attribute.
	GetUserRoles(
		ctx context.Context,
		attrType entry.AttributeType,
		targetID T,
		userID umid.UMID,
	) ([]entry.PermissionsRoleType, error)
}

type Attributes[ID comparable] interface {
	AttributeUserRoleGetter[ID]
	AttributeOptionsGetter[ID]

	GetPayload(attributeID ID) (*entry.AttributePayload, bool)
	GetValue(attributeID ID) (*entry.AttributeValue, bool)

	Upsert(attributeID ID, modifyFn modify.Fn[entry.AttributePayload], updateDB bool) (*entry.AttributePayload, error)

	UpdateValue(attributeID ID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool) (*entry.AttributeValue, error)
	UpdateOptions(attributeID ID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool) (
		*entry.AttributeOptions, error,
	)

	Remove(attributeID ID, updateDB bool) (bool, error)
}

type NodeAttributes interface {
	LoadSaver
	Attributes[entry.AttributeID]

	GetAll() map[entry.AttributeID]*entry.AttributePayload

	Len() int
}

type ObjectAttributes interface {
	LoadSaver
	Attributes[entry.AttributeID]

	GetAll() map[entry.AttributeID]*entry.AttributePayload

	Len() int
}

// UserAttributes ignores "updateDB" flag
type UserAttributes interface {
	Attributes[entry.UserAttributeID]
}

// UserUserAttributes ignores "updateDB" flag
type UserUserAttributes interface {
	Attributes[entry.UserUserAttributeID]
}

// ObjectUserAttributes ignores "updateDB" flag
type ObjectUserAttributes interface {
	Attributes[entry.ObjectUserAttributeID]
}

type Assets2d interface {
	LoadSaver
	APIRegister

	CreateAsset2d(asset2dID umid.UMID) (Asset2d, error)
	GetAsset2d(asset2dID umid.UMID) (Asset2d, bool)
	GetAssets2d() map[umid.UMID]Asset2d
	FilterAssets2d(predicateFn Assets2dFilterPredicateFn) map[umid.UMID]Asset2d
	AddAsset2d(asset2d Asset2d, updateDB bool) error
	AddAssets2d(assets2d []Asset2d, updateDB bool) error
	RemoveAsset2d(asset2d Asset2d, updateDB bool) (bool, error)
	RemoveAssets2d(assets2d []Asset2d, updateDB bool) (bool, error)
}

type Asset2d interface {
	IDer

	GetMeta() entry.Asset2dMeta
	SetMeta(meta entry.Asset2dMeta, updateDB bool) error

	GetOptions() *entry.Asset2dOptions
	SetOptions(modifyFn modify.Fn[entry.Asset2dOptions], updateDB bool) (*entry.Asset2dOptions, error)

	GetEntry() *entry.Asset2d
	LoadFromEntry(entry *entry.Asset2d) error
}

type Activity interface {
	IDer
	Initializer

	GetData() *entry.ActivityData
	SetData(modifyFn modify.Fn[entry.ActivityData], updateDB bool) (*entry.ActivityData, error)

	GetType() *entry.ActivityType
	SetType(activityType *entry.ActivityType, updateDB bool) error

	GetObjectID() umid.UMID
	SetObjectID(objectID umid.UMID, updateDB bool) error

	GetUserID() umid.UMID
	SetUserID(userID umid.UMID, updateDB bool) error

	GetEntry() *entry.Activity
	LoadFromEntry(entry *entry.Activity) error

	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time, updateDB bool) error
}

type Activities interface {
	Initializer
	LoadSaver

	CreateActivity(activityID umid.UMID) (Activity, error)

	GetActivity(activityID umid.UMID) (Activity, bool)
	GetActivities() map[umid.UMID]Activity

	GetPaginatedActivitiesByObjectID(objectID *umid.UMID, page int, pageSize int) ([]Activity, int)
	GetActivitiesByUserID(userID umid.UMID) map[umid.UMID]Activity

	AddActivity(activity Activity, updateDB bool) error
	AddActivities(activities []Activity, updateDB bool) error

	RemoveActivity(activity Activity, updateDB bool) (bool, error)
	RemoveActivities(activities2d []Activity, updateDB bool) (bool, error)
}

type Assets3d interface {
	LoadSaver
	APIRegister

	// Create new instance if doesn't exist, returns the existing/created asset3d and bool isCreated
	CreateAsset3d(assetID umid.UMID) (Asset3d, error, bool)
	CreateUserAsset3d(assetID umid.UMID, userID umid.UMID, isPrivate bool) (UserAsset3d, error)

	GetAsset3d(assetID umid.UMID) (Asset3d, bool)
	GetUserAsset3d(assetID umid.UMID, userID umid.UMID) (UserAsset3d, bool)

	GetAssets3d() map[umid.UMID]Asset3d
	GetUserAssets3d() map[AssetUserIDPair]UserAsset3d

	FilterUserAssets3d(predicateFn Assets3dFilterPredicateFn) map[AssetUserIDPair]UserAsset3d

	AddAsset3d(asset3d Asset3d, updateDB bool) error
	AddUserAsset3d(asset3d UserAsset3d, updateDB bool) error

	RemoveUserAsset3dByID(assets3dID AssetUserIDPair, updateDB bool) (bool, error)
}

type Asset3d interface {
	IDer
	GetMeta() *entry.Asset3dMeta
	SetMeta(meta *entry.Asset3dMeta, updateDB bool) error

	GetOptions() *entry.Asset3dOptions
	SetOptions(modifyFn modify.Fn[entry.Asset3dOptions], updateDB bool) (*entry.Asset3dOptions, error)

	GetEntry() *entry.Asset3d
	LoadFromEntry(entry *entry.Asset3d) error
}

type UserAsset3d interface {
	GetAssetUserIDPair() AssetUserIDPair
	GetAssetID() umid.UMID
	GetUserID() umid.UMID

	GetAsset3d() *Asset3d

	GetMeta() *entry.Asset3dMeta
	SetMeta(meta *entry.Asset3dMeta, updateDB bool) error

	IsPrivate() bool
	SetIsPrivate(isPrivate bool, updateDB bool) error

	GetEntry() *entry.UserAsset3d
	LoadFromEntry(entry *entry.UserAsset3d) error
}

type Plugins interface {
	LoadSaver
	APIRegister

	CreatePlugin(pluginID umid.UMID) (Plugin, error)
	GetPlugin(pluginID umid.UMID) (Plugin, bool)
	GetPlugins() map[umid.UMID]Plugin
	FilterPlugins(predicateFn PluginsFilterPredicateFn) map[umid.UMID]Plugin
	AddPlugin(plugin Plugin, updateDB bool) error
	AddPlugins(plugins []Plugin, updateDB bool) error
	RemovePlugin(plugin Plugin, updateDB bool) (bool, error)
	RemovePlugins(plugins []Plugin, updateDB bool) (bool, error)
}

type Plugin interface {
	IDer

	GetMeta() entry.PluginMeta
	SetMeta(meta entry.PluginMeta, updateDB bool) error

	GetOptions() *entry.PluginOptions
	SetOptions(modifyFn modify.Fn[entry.PluginOptions], updateDB bool) (*entry.PluginOptions, error)

	GetEntry() *entry.Plugin
	LoadFromEntry(entry *entry.Plugin) error
}

type AttributeTypes interface {
	LoadSaver
	APIRegister

	CreateAttributeType(attributeTypeID entry.AttributeTypeID) (AttributeType, error)
	GetAttributeType(attributeTypeID entry.AttributeTypeID) (AttributeType, bool)
	GetAttributeTypes() map[entry.AttributeTypeID]AttributeType
	FilterAttributeTypes(predicateFn AttributeTypesFilterPredicateFn) map[entry.AttributeTypeID]AttributeType
	AddAttributeType(attributeType AttributeType, updateDB bool) error
	AddAttributeTypes(attributeTypes []AttributeType, updateDB bool) error
	RemoveAttributeType(attributeType AttributeType, updateDB bool) (bool, error)
	RemoveAttributeTypes(attributeTypes []AttributeType, updateDB bool) (bool, error)
}

type AttributeType interface {
	GetID() entry.AttributeTypeID
	GetName() string
	GetPluginID() umid.UMID

	GetOptions() *entry.AttributeOptions
	SetOptions(modifyFn modify.Fn[entry.AttributeOptions], updateDB bool) (*entry.AttributeOptions, error)

	GetDescription() *string
	SetDescription(description *string, updateDB bool) error

	GetEntry() *entry.AttributeType
	LoadFromEntry(entry *entry.AttributeType) error
}

type ObjectTypes interface {
	LoadSaver
	APIRegister

	CreateObjectType(objectTypeID umid.UMID) (ObjectType, error)
	GetObjectType(objectTypeID umid.UMID) (ObjectType, bool)
	GetObjectTypes() map[umid.UMID]ObjectType
	FilterObjectTypes(predicateFn ObjectTypesFilterPredicateFn) map[umid.UMID]ObjectType
	AddObjectType(objectType ObjectType, updateDB bool) error
	AddObjectTypes(objectTypes []ObjectType, updateDB bool) error
	RemoveObjectType(objectType ObjectType, updateDB bool) (bool, error)
	RemoveObjectTypes(objectTypes []ObjectType, updateDB bool) (bool, error)
}

type ObjectType interface {
	IDer

	GetName() string
	SetName(name string, updateDB bool) error

	GetCategoryName() string
	SetCategoryName(categoryName string, updateDB bool) error

	GetDescription() *string
	SetDescription(description *string, updateDB bool) error

	GetOptions() *entry.ObjectOptions
	SetOptions(modifyFn modify.Fn[entry.ObjectOptions], updateDB bool) (*entry.ObjectOptions, error)

	GetAsset2d() Asset2d
	SetAsset2d(asset2d Asset2d, updateDB bool) error

	GetAsset3d() Asset3d
	SetAsset3d(asset3d Asset3d, updateDB bool) error

	GetEntry() *entry.ObjectType
	LoadFromEntry(entry *entry.ObjectType) error
}

type UserTypes interface {
	LoadSaver
	APIRegister

	CreateUserType(userTypeID umid.UMID) (UserType, error)
	GetUserType(userTypeID umid.UMID) (UserType, bool)
	GetUserTypes() map[umid.UMID]UserType
	FilterUserTypes(predicateFn UserTypesFilterPredicateFn) map[umid.UMID]UserType
	AddUserType(userType UserType, updateDB bool) error
	AddUserTypes(userTypes []UserType, updateDB bool) error
	RemoveUserType(userType UserType, updateDB bool) (bool, error)
	RemoveUserTypes(userTypes []UserType, updateDB bool) (bool, error)
}

type UserType interface {
	IDer

	GetName() string
	SetName(name string, updateDB bool) error

	GetDescription() string
	SetDescription(description string, updateDB bool) error

	GetOptions() *entry.UserOptions
	SetOptions(modifyFn modify.Fn[entry.UserOptions], updateDB bool) (*entry.UserOptions, error)

	GetEntry() *entry.UserType
	LoadFromEntry(entry *entry.UserType) error
}

type Calendar interface {
	RunStopper

	OnAttributeUpsert(attributeID entry.AttributeID, value any)
	OnAttributeRemove(attributeID entry.AttributeID)
}
