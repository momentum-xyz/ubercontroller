package universe

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/momentum-xyz/controller/pkg/cmath"
	"github.com/momentum-xyz/controller/types"
	"github.com/momentum-xyz/controller/types/generics"
)

type Node interface {
	types.IDer
	types.Initializer
	types.RunStopper

	Load() error
}

type Worlds interface {
	types.Initializer

	GetWorld(worldID uuid.UUID) (World, bool)
	GetWorlds() *generics.SyncMap[uuid.UUID, World]

	Load() error
}

type World interface {
	Space
	types.RunStopper

	Load() error
}

type Space interface {
	types.IDer
	types.Initializer

	GetWorld() World

	GetParent() Space
	SetParent(parent Space, updateDB bool) error

	GetTheta() float64
	SetTheta(theta float64, updateDB bool) error
	GetPosition() cmath.Vec3
	SetPosition(pos cmath.Vec3, updateDB bool) error

	Update(recursive bool) error
	LoadFromEntry(entry *SpaceEntry) error

	GetAsset2D() Asset2D
	GetAsset3D() Asset3D
	SetAsset2D(asset2d Asset2D, updateDB bool) error
	SetAsset3D(asset3d Asset3D, updateDB bool) error

	GetSpaceType() SpaceType
	SetSpaceType(spaceType SpaceType, updateDB bool) error

	GetOptions() *SpaceOptionsEntry
	SetOptions(options *SpaceOptionsEntry, updateDB bool) error

	GetSpace(spaceID uuid.UUID, recursive bool) (Space, bool)
	GetSpaces(recursive bool) *generics.SyncMap[uuid.UUID, Space]
	AddSpace(space Space, updateDB bool) error
	AddSpaces(spaces []Space, updateDB bool) error
	RemoveSpace(spaceID uuid.UUID, recursive, updateDB bool) (bool, error)
	RemoveSpaces(spaceIDs []uuid.UUID, recursive, updateDB bool) (bool, error)

	GetOwner() User
	SetOwner(owner User, updateDB bool) error

	GetUser(userID uuid.UUID, recursive bool) (User, bool)
	GetUsers(recursive bool) *generics.SyncMap[uuid.UUID, User]
	AddUser(user User, updateDB bool) error
	RemoveUser(user User, updateDB bool) error

	SendToUser(userID uuid.UUID, msg *websocket.PreparedMessage, recursive bool) error
	SendToUsers(msg *websocket.PreparedMessage, recursive bool) error

	//GetSpaceAttributes() *SpaceAttributes
	//GetUserSpaceAttributes() *UserSpaceAttributes
}

type User interface {
	types.IDer
	types.Initializer
	types.RunStopper

	GetWorld() World
	SetWorld(world World, updateDB bool) error

	GetSpace() Space
	SetSpace(space Space, updateDB bool) error
}

type SpaceTypes interface {
	types.Initializer

	GetSpaceType(spaceTypeID uuid.UUID) (SpaceType, bool)

	Load() error
}

type SpaceType interface {
	types.IDer
	types.Initializer

	GetName() string
	SetName(name string, updateDB bool) error
	GetCategoryName() string
	SetCategoryName(categoryName string, updateDB bool) error
	GetDescription() *string
	SetDescription(description *string, updateDB bool) error

	LoadFromEntry(entry *SpaceTypeEntry) error

	GetOptions() *SpaceOptionsEntry
	SetOptions(options *SpaceOptionsEntry, updateDB bool) error
}

type Assets2D interface {
	types.Initializer

	GetAsset2D(asset2DID uuid.UUID) (Asset2D, bool)

	Load() error
}

type Asset2D interface {
	types.IDer
	types.Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	LoadFromEntry(entry *Asset2DEntry) error

	GetOptions() *Asset2DOptionsEntry
	SetOptions(options *Asset2DOptionsEntry, updateDB bool) error
}

type Assets3D interface {
	types.Initializer

	GetAsset3D(asset3DID uuid.UUID) (Asset3D, bool)

	Load() error
}

type Asset3D interface {
	types.IDer
	types.Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	LoadFromEntry(entry *Asset3DEntry) error

	GetOptions() *Asset3DOptionsEntry
	SetOptions(options *Asset3DOptionsEntry, updateDB bool) error
}
