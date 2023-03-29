package universe

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	influxWrite "github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type IDer interface {
	GetID() umid.UMID
}

type Initializer interface {
	Initialize(ctx context.Context) error
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

	CreateUsers(ctx context.Context, users ...*entry.User) error // TODO: refactor, place Users next to Nodes in a universe

	AddAPIRegister(register APIRegister)

	WriteInfluxPoint(point *influxWrite.Point) error
}

type Worlds interface {
	Initializer
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

	GetParent() Object
	SetParent(parent Object, updateDB bool) error

	GetOwnerID() umid.UMID
	SetOwnerID(ownerID umid.UMID, updateDB bool) error

	GetTransform() *cmath.ObjectTransform
	GetActualTransform() *cmath.ObjectTransform
	SetTransform(position *cmath.ObjectTransform, updateDB bool) error
	SetActualTransform(pos cmath.ObjectTransform, theta float64) error

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

	LockUnityObject(user User, state uint32) bool
}

type User interface {
	IDer
	Initializer
	RunStopper

	GetWorld() World
	SetWorld(world World)

	GetObject() Object
	SetObject(object Object)

	GetUserType() UserType
	SetUserType(userType UserType, updateDB bool) error

	GetProfile() *entry.UserProfile

	GetTransform() *cmath.UserTransform
	SetTransform(cmath.UserTransform)

	GetPosition() cmath.Vec3
	GetRotation() cmath.Vec3
	SetPosition(position cmath.Vec3)

	//GetPosBuffer() []byte
	GetLastPosTime() int64

	Update() error
	ReleaseSendBuffer()
	LockSendBuffer()

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

type Attributes[ID comparable] interface {
	GetPayload(attributeID ID) (*entry.AttributePayload, bool)
	GetValue(attributeID ID) (*entry.AttributeValue, bool)
	GetOptions(attributeID ID) (*entry.AttributeOptions, bool)
	GetEffectiveOptions(attributeID ID) (*entry.AttributeOptions, bool)

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
	Initializer
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
	Initializer

	GetMeta() entry.Asset2dMeta
	SetMeta(meta entry.Asset2dMeta, updateDB bool) error

	GetOptions() *entry.Asset2dOptions
	SetOptions(modifyFn modify.Fn[entry.Asset2dOptions], updateDB bool) (*entry.Asset2dOptions, error)

	GetEntry() *entry.Asset2d
	LoadFromEntry(entry *entry.Asset2d) error
}

type Assets3d interface {
	Initializer
	LoadSaver
	APIRegister

	CreateAsset3d(asset3dID umid.UMID) (Asset3d, error)
	GetAsset3d(asset3dID umid.UMID) (Asset3d, bool)
	GetAssets3d() map[umid.UMID]Asset3d
	FilterAssets3d(predicateFn Assets3dFilterPredicateFn) map[umid.UMID]Asset3d
	AddAsset3d(asset3d Asset3d, updateDB bool) error
	AddAssets3d(assets3d []Asset3d, updateDB bool) error
	RemoveAsset3d(asset3d Asset3d, updateDB bool) (bool, error)
	RemoveAssets3d(assets3d []Asset3d, updateDB bool) (bool, error)
	RemoveAsset3dByID(assets3dID umid.UMID, updateDB bool) (bool, error)
	RemoveAssets3dByIDs(assets3dIDs []umid.UMID, updateDB bool) (bool, error)
}

type Asset3d interface {
	IDer
	Initializer

	GetMeta() *entry.Asset3dMeta
	SetMeta(meta *entry.Asset3dMeta, updateDB bool) error

	GetOptions() *entry.Asset3dOptions
	SetOptions(modifyFn modify.Fn[entry.Asset3dOptions], updateDB bool) (*entry.Asset3dOptions, error)

	GetEntry() *entry.Asset3d
	LoadFromEntry(entry *entry.Asset3d) error
}

type Plugins interface {
	Initializer
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
	Initializer

	GetMeta() entry.PluginMeta
	SetMeta(meta entry.PluginMeta, updateDB bool) error

	GetOptions() *entry.PluginOptions
	SetOptions(modifyFn modify.Fn[entry.PluginOptions], updateDB bool) (*entry.PluginOptions, error)

	GetEntry() *entry.Plugin
	LoadFromEntry(entry *entry.Plugin) error
}

type AttributeTypes interface {
	Initializer
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
	Initializer

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
	Initializer
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
	Initializer

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
	Initializer
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
	Initializer

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
	Initializer
	RunStopper

	OnAttributeUpsert(attributeID entry.AttributeID, value any)
	OnAttributeRemove(attributeID entry.AttributeID)
}
