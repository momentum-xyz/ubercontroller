package universe

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type IDer interface {
	GetID() uuid.UUID
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

type Node interface {
	IDer
	Initializer
	RunStopper
	LoadSaver
	APIRegister

	GetWorlds() Worlds
	GetAssets2d() Assets2d
	GetAssets3d() Assets3d
	GetSpaceTypes() SpaceTypes

	AddAPIRegister(register APIRegister)
}

type Worlds interface {
	Initializer
	RunStopper
	LoadSaver
	APIRegister

	NewWorld(worldID uuid.UUID) (World, error)

	GetWorld(worldID uuid.UUID) (World, bool)
	GetWorlds() map[uuid.UUID]World
	AddWorld(world World, updateDB bool) error
	AddWorlds(worlds []World, updateDB bool) error
	RemoveWorld(world World, updateDB bool) error
	RemoveWorlds(worlds []World, updateDB bool) error
}

type World interface {
	Space
	RunStopper
	LoadSaver
	APIRegister
}

type Space interface {
	IDer
	Initializer

	NewSpace(spaceID uuid.UUID) (Space, error)

	GetWorld() World

	GetParent() Space
	SetParent(parent Space, updateDB bool) error

	GetOwnerID() uuid.UUID
	SetOwnerID(ownerID uuid.UUID, updateDB bool) error

	GetPosition() *cmath.Vec3
	SetPosition(position *cmath.Vec3, updateDB bool) error

	GetOptions() *entry.SpaceOptions
	SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) error

	GetEffectiveOptions() *entry.SpaceOptions

	GetAsset2D() Asset2d
	SetAsset2D(asset2d Asset2d, updateDB bool) error

	GetAsset3D() Asset3d
	SetAsset3D(asset3d Asset3d, updateDB bool) error

	GetSpaceType() SpaceType
	SetSpaceType(spaceType SpaceType, updateDB bool) error

	GetEntry() *entry.Space
	LoadFromEntry(entry *entry.Space, recursive bool) error

	Update(recursive bool) error

	GetSpace(spaceID uuid.UUID, recursive bool) (Space, bool)
	GetSpaces(recursive bool) map[uuid.UUID]Space
	AddSpace(space Space, updateDB bool) error
	AddSpaces(spaces []Space, updateDB bool) error
	RemoveSpace(space Space, recursive, updateDB bool) (bool, error)
	RemoveSpaces(spaces []Space, recursive, updateDB bool) (bool, error)

	GetUser(userID uuid.UUID, recursive bool) (User, bool)
	GetUsers(recursive bool) map[uuid.UUID]User
	AddUser(user User, updateDB bool) error
	RemoveUser(user User, updateDB bool) error

	SendToUser(userID uuid.UUID, msg *websocket.PreparedMessage, recursive bool) error
	Broadcast(msg *websocket.PreparedMessage, recursive bool) error
}

type User interface {
	IDer
	Initializer
	RunStopper
	APIRegister

	GetWorld() World
	SetWorld(world World, updateDB bool) error

	GetSpace() Space
	SetSpace(space Space, updateDB bool) error
}

type SpaceTypes interface {
	Initializer
	LoadSaver
	APIRegister

	NewSpaceType(spaceTypeID uuid.UUID) (SpaceType, error)

	GetSpaceType(spaceTypeID uuid.UUID) (SpaceType, bool)
	GetSpaceTypes() map[uuid.UUID]SpaceType
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
	SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) error

	GetAsset2d() Asset2d
	SetAsset2d(asset2d Asset2d, updateDB bool) error

	GetAsset3d() Asset3d
	SetAsset3d(asset3d Asset3d, updateDB bool) error

	GetEntry() *entry.SpaceType
	LoadFromEntry(entry *entry.SpaceType) error
}

type Assets2d interface {
	Initializer
	LoadSaver
	APIRegister

	NewAsset2d(asset2dID uuid.UUID) (Asset2d, error)

	GetAsset2d(asset2dID uuid.UUID) (Asset2d, bool)
	GetAssets2d() map[uuid.UUID]Asset2d
	AddAsset2d(asset2d Asset2d, updateDB bool) error
	AddAssets2d(assets2d []Asset2d, updateDB bool) error
	RemoveAsset2d(asset2d Asset2d, updateDB bool) error
	RemoveAssets2d(assets2d []Asset2d, updateDB bool) error
}

type Asset2d interface {
	IDer
	Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	GetOptions() *entry.Asset2dOptions
	SetOptions(modifyFn modify.Fn[entry.Asset2dOptions], updateDB bool) error

	GetEntry() *entry.Asset2d
	LoadFromEntry(entry *entry.Asset2d) error
}

type Assets3d interface {
	Initializer
	LoadSaver
	APIRegister

	NewAsset3d(asset3dID uuid.UUID) (Asset3d, error)

	GetAsset3d(asset3dID uuid.UUID) (Asset3d, bool)
	GetAssets3d() map[uuid.UUID]Asset3d
	AddAsset3d(asset3d Asset3d, updateDB bool) error
	AddAssets3d(assets3d []Asset3d, updateDB bool) error
	RemoveAsset3d(asset3d Asset3d, updateDB bool) error
	RemoveAssets3d(assets3d []Asset3d, updateDB bool) error
}

type Asset3d interface {
	IDer
	Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	GetOptions() *entry.Asset3dOptions
	SetOptions(modifyFn modify.Fn[entry.Asset3dOptions], updateDB bool) error

	GetEntry() *entry.Asset3d
	LoadFromEntry(entry *entry.Asset3d) error
}
