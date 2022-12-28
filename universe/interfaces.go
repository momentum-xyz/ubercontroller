package universe

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	influxWrite "github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type IDer interface {
	GetID() uuid.UUID
}

type Enabler interface {
	GetEnabled() bool
	SetEnabled(enabled bool)
}

type Initializer interface {
	Initialize(ctx context.Context) error
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

type SpaceCacher interface {
	GetAllSpaces() map[uuid.UUID]Space
	GetSpaceFromAllSpaces(spaceID uuid.UUID) (Space, bool)
	FilterAllSpaces(predicateFn SpacesFilterPredicateFn) map[uuid.UUID]Space
	AddSpaceToAllSpaces(space Space) error
	RemoveSpaceFromAllSpaces(space Space) (bool, error)
}

type DropCacher interface {
	DropCache()
}

type Node interface {
	Space
	LoadSaver
	APIRegister
	SpaceCacher

	ToSpace() Space

	GetWorlds() Worlds
	GetAssets2d() Assets2d
	GetAssets3d() Assets3d
	GetSpaceTypes() SpaceTypes
	GetUserTypes() UserTypes
	GetAttributeTypes() AttributeTypes
	GetPlugins() Plugins

	GetNodeAttributes() Attributes[entry.AttributeID]
	GetUserAttributes() Attributes[entry.UserAttributeID]           // ignores "updateDB" flag
	GetSpaceUserAttributes() Attributes[entry.SpaceUserAttributeID] // ignores "updateDB" flag

	GetUserUserAttributePayload(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributePayload, bool)
	GetUserUserAttributeValue(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributeValue, bool)
	GetUserUserAttributeOptions(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributeOptions, bool)
	GetUserUserAttributeEffectiveOptions(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributeOptions, bool)

	UpsertUserUserAttribute(
		userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
	) (*entry.AttributePayload, error)

	UpdateUserUserAttributeValue(
		userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
	) (*entry.AttributeValue, error)
	UpdateUserUserAttributeOptions(
		userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
	) (*entry.AttributeOptions, error)

	RemoveUserUserAttribute(userUserAttributeID entry.UserUserAttributeID) (bool, error)

	AddAPIRegister(register APIRegister)

	WriteInfluxPoint(point *influxWrite.Point) error
}

type Worlds interface {
	Initializer
	RunStopper
	LoadSaver
	APIRegister

	CreateWorld(worldID uuid.UUID) (World, error)

	GetWorld(worldID uuid.UUID) (World, bool)
	GetWorlds() map[uuid.UUID]World
	FilterWorlds(predicateFn WorldsFilterPredicateFn) map[uuid.UUID]World
	AddWorld(world World, updateDB bool) error
	AddWorlds(worlds []World, updateDB bool) error
	RemoveWorld(world World, updateDB bool) error
	RemoveWorlds(worlds []World, updateDB bool) error
}

type World interface {
	Space
	LoadSaver
	SpaceCacher

	ToSpace() Space

	GetSettings() *WorldSettings

	GetCalendar() Calendar

	WriteInfluxPoint(point *influxWrite.Point) error
}

type Space interface {
	IDer
	Enabler
	Initializer
	RunStopper
	DropCacher

	CreateSpace(spaceID uuid.UUID) (Space, error)

	GetWorld() World

	GetName() string
	SetName(name string, updateDB bool) error

	GetParent() Space
	SetParent(parent Space, updateDB bool) error

	GetOwnerID() uuid.UUID
	SetOwnerID(ownerID uuid.UUID, updateDB bool) error

	GetPosition() *cmath.SpacePosition
	GetActualPosition() *cmath.SpacePosition
	SetPosition(position *cmath.SpacePosition, updateDB bool) error
	SetActualPosition(pos cmath.SpacePosition, theta float64) error

	GetOptions() *entry.SpaceOptions
	GetEffectiveOptions() *entry.SpaceOptions
	SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) (*entry.SpaceOptions, error)

	GetAsset2D() Asset2d
	SetAsset2D(asset2d Asset2d, updateDB bool) error

	GetAsset3D() Asset3d
	SetAsset3D(asset3d Asset3d, updateDB bool) error

	GetSpaceType() SpaceType
	SetSpaceType(spaceType SpaceType, updateDB bool) error

	GetSpaceAttributes() Attributes[entry.AttributeID]

	GetEntry() *entry.Space
	LoadFromEntry(entry *entry.Space, recursive bool) error

	Update(recursive bool) error
	UpdateChildrenPosition(recursive bool) error

	GetSpace(spaceID uuid.UUID, recursive bool) (Space, bool)
	GetSpaces(recursive bool) map[uuid.UUID]Space
	FilterSpaces(predicateFn SpacesFilterPredicateFn, recursive bool) map[uuid.UUID]Space
	AddSpace(space Space, updateDB bool) error
	AddSpaces(spaces []Space, updateDB bool) error
	RemoveSpace(space Space, recursive, updateDB bool) (bool, error)
	RemoveSpaces(spaces []Space, recursive, updateDB bool) (bool, error)

	GetUser(userID uuid.UUID, recursive bool) (User, bool)
	GetUsers(recursive bool) map[uuid.UUID]User
	AddUser(user User, updateDB bool) error
	RemoveUser(user User, updateDB bool) error

	Send(msg *websocket.PreparedMessage, recursive bool) error

	SendSpawnMessage(sendFn func(msg *websocket.PreparedMessage) error, recursive bool)
	SendAttributes(sendFn func(*websocket.PreparedMessage), recursive bool)
	SendTextures(sendFn func(msg *websocket.PreparedMessage) error, recursive bool)

	LockUnityObject(user User, state uint32) bool
}

type User interface {
	IDer
	Initializer
	RunStopper

	GetWorld() World
	SetWorld(world World)

	GetSpace() Space
	SetSpace(space Space)

	GetUserType() UserType
	SetUserType(userType UserType, updateDB bool) error

	GetProfile() *entry.UserProfile

	GetPosition() cmath.Vec3
	GetRotation() cmath.Vec3
	SetPosition(position cmath.Vec3)

	GetPosBuffer() []byte

	Update() error

	GetSessionID() uuid.UUID
	SetConnection(sessionID uuid.UUID, socketConnection *websocket.Conn) error

	Send(message *websocket.PreparedMessage) error
	SendDirectly(message *websocket.PreparedMessage) error

	ReleaseSendBuffer()

	AddInfluxTags(prefix string, point *influxWrite.Point) *influxWrite.Point
}

type Attributes[K comparable] interface {
	GetPayload(attributeID K) (*entry.AttributePayload, bool)
	GetValue(attributeID K) (*entry.AttributeValue, bool)
	GetOptions(attributeID K) (*entry.AttributeOptions, bool)
	GetEffectiveOptions(attributeID K) (*entry.AttributeOptions, bool)

	Upsert(attributeID K, modifyFn modify.Fn[entry.AttributePayload], updateDB bool) (*entry.AttributePayload, error)

	UpdateValue(attributeID K, modifyFn modify.Fn[entry.AttributeValue], updateDB bool) (*entry.AttributeValue, error)
	UpdateOptions(attributeID K, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool) (*entry.AttributeOptions, error)

	Remove(attributeID K, updateDB bool) (bool, error)

	Len() int
}

type Assets2d interface {
	Initializer
	LoadSaver
	APIRegister

	CreateAsset2d(asset2dID uuid.UUID) (Asset2d, error)

	GetAsset2d(asset2dID uuid.UUID) (Asset2d, bool)
	GetAssets2d() map[uuid.UUID]Asset2d
	FilterAssets2d(predicateFn Assets2dFilterPredicateFn) map[uuid.UUID]Asset2d
	AddAsset2d(asset2d Asset2d, updateDB bool) error
	AddAssets2d(assets2d []Asset2d, updateDB bool) error
	RemoveAsset2d(asset2d Asset2d, updateDB bool) error
	RemoveAssets2d(assets2d []Asset2d, updateDB bool) error
}

type Asset2d interface {
	IDer
	Initializer

	GetMeta() *entry.Asset2dMeta
	SetMeta(meta *entry.Asset2dMeta, updateDB bool) error

	GetOptions() *entry.Asset2dOptions
	SetOptions(modifyFn modify.Fn[entry.Asset2dOptions], updateDB bool) (*entry.Asset2dOptions, error)

	GetEntry() *entry.Asset2d
	LoadFromEntry(entry *entry.Asset2d) error
}

type Assets3d interface {
	Initializer
	LoadSaver
	APIRegister

	CreateAsset3d(asset3dID uuid.UUID) (Asset3d, error)
	GetAsset3d(asset3dID uuid.UUID) (Asset3d, bool)
	GetAssets3d() map[uuid.UUID]Asset3d
	FilterAssets3d(predicateFn Assets3dFilterPredicateFn) map[uuid.UUID]Asset3d
	AddAsset3d(asset3d Asset3d, updateDB bool) error
	AddAssets3d(assets3d []Asset3d, updateDB bool) error
	RemoveAsset3d(asset3d Asset3d, updateDB bool) error
	RemoveAssets3d(assets3d []Asset3d, updateDB bool) error
	RemoveAsset3dByID(asset3dID uuid.UUID, updateDB bool) error
	RemoveAssets3dByIDs(assets3dIDs []uuid.UUID, updateDB bool) error
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

	CreatePlugin(pluginID uuid.UUID) (Plugin, error)

	GetPlugin(pluginID uuid.UUID) (Plugin, bool)
	GetPlugins() map[uuid.UUID]Plugin
	FilterPlugins(predicateFn PluginsFilterPredicateFn) map[uuid.UUID]Plugin
	AddPlugin(plugin Plugin, updateDB bool) error
	AddPlugins(plugins []Plugin, updateDB bool) error
	RemovePlugin(plugin Plugin, updateDB bool) error
	RemovePlugins(plugins []Plugin, updateDB bool) error
}

type Plugin interface {
	IDer
	Initializer

	GetMeta() *entry.PluginMeta
	SetMeta(meta *entry.PluginMeta, updateDB bool) error

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
	RemoveAttributeType(attributeType AttributeType, updateDB bool) error
	RemoveAttributeTypes(attributeTypes []AttributeType, updateDB bool) error
}

type AttributeType interface {
	Initializer

	GetID() entry.AttributeTypeID
	GetName() string
	GetPluginID() uuid.UUID

	GetOptions() *entry.AttributeOptions
	SetOptions(modifyFn modify.Fn[entry.AttributeOptions], updateDB bool) (*entry.AttributeOptions, error)

	GetDescription() *string
	SetDescription(description *string, updateDB bool) error

	GetEntry() *entry.AttributeType
	LoadFromEntry(entry *entry.AttributeType) error
}

type SpaceTypes interface {
	Initializer
	LoadSaver
	APIRegister

	CreateSpaceType(spaceTypeID uuid.UUID) (SpaceType, error)

	GetSpaceType(spaceTypeID uuid.UUID) (SpaceType, bool)
	GetSpaceTypes() map[uuid.UUID]SpaceType
	FilterSpaceTypes(predicateFn SpaceTypesFilterPredicateFn) map[uuid.UUID]SpaceType
	AddSpaceType(spaceType SpaceType, updateDB bool) error
	AddSpaceTypes(spaceTypes []SpaceType, updateDB bool) error
	RemoveSpaceType(spaceType SpaceType, updateDB bool) error
	RemoveSpaceTypes(spaceTypes []SpaceType, updateDB bool) error
}

type SpaceType interface {
	IDer
	Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	GetCategoryName() string
	SetCategoryName(categoryName string, updateDB bool) error

	GetDescription() *string
	SetDescription(description *string, updateDB bool) error

	GetOptions() *entry.SpaceOptions
	SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) (*entry.SpaceOptions, error)

	GetAsset2d() Asset2d
	SetAsset2d(asset2d Asset2d, updateDB bool) error

	GetAsset3d() Asset3d
	SetAsset3d(asset3d Asset3d, updateDB bool) error

	GetEntry() *entry.SpaceType
	LoadFromEntry(entry *entry.SpaceType) error
}

type UserTypes interface {
	Initializer
	LoadSaver
	APIRegister

	CreateUserType(userTypeID uuid.UUID) (UserType, error)

	GetUserType(userTypeID uuid.UUID) (UserType, bool)
	GetUserTypes() map[uuid.UUID]UserType
	FilterUserTypes(predicateFn UserTypesFilterPredicateFn) map[uuid.UUID]UserType
	AddUserType(spaceType UserType, updateDB bool) error
	AddUserTypes(spaceTypes []UserType, updateDB bool) error
	RemoveUserType(spaceType UserType, updateDB bool) error
	RemoveUserTypes(spaceTypes []UserType, updateDB bool) error
}

type UserType interface {
	IDer
	Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	GetDescription() *string
	SetDescription(description *string, updateDB bool) error

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
