package universe

import (
	"context"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generics"
)

type Node interface {
	types.IDer
	types.Initializer
	types.RunStopper

	GetWorlds() Worlds
	GetAssets2d() Assets2d
	GetAssets3d() Assets3d
	GetSpaceTypes() SpaceTypes

	Load(ctx context.Context) error
}

type Worlds interface {
	types.Initializer

	GetWorld(worldID uuid.UUID) (World, bool)
	GetWorlds() *generics.SyncMap[uuid.UUID, World]

	Load(ctx context.Context) error
}

type World interface {
	Space
	types.RunStopper
}

type Space interface {
	types.IDer
	types.Initializer

	GetWorld() World

	GetParent() Space
	SetParent(parent Space, updateDB bool) error

	GetOwnerID() uuid.UUID
	SetOwnerID(ownerID uuid.UUID, updateDB bool) error

	GetPosition() cmath.Vec3
	SetPosition(pos cmath.Vec3, updateDB bool) error

	GetOptions() *entry.SpaceOptions
	GetEffectiveOptions() *entry.SpaceOptions
	SetOptions(options *entry.SpaceOptions, updateDB bool) error

	GetAsset2D() Asset2d
	SetAsset2D(asset2d Asset2d, updateDB bool) error

	GetAsset3D() Asset3d
	SetAsset3D(asset3d Asset3d, updateDB bool) error

	GetSpaceType() SpaceType
	SetSpaceType(spaceType SpaceType, updateDB bool) error

	Update(recursive bool) error
	LoadFromEntry(ctx context.Context, entry *entry.Space, recursive bool) error

	GetSpace(spaceID uuid.UUID, recursive bool) (Space, bool)
	GetSpaces(recursive bool) *generics.SyncMap[uuid.UUID, Space]
	AddSpace(space Space, updateDB bool) error
	AddSpaces(spaces []Space, updateDB bool) error
	RemoveSpace(spaceID uuid.UUID, recursive, updateDB bool) (bool, error)
	RemoveSpaces(spaceIDs []uuid.UUID, recursive, updateDB bool) (bool, error)

	GetUser(userID uuid.UUID, recursive bool) (User, bool)
	GetUsers(recursive bool) *generics.SyncMap[uuid.UUID, User]
	AddUser(user User, updateDB bool) error
	RemoveUser(user User, updateDB bool) error

	SendToUser(userID uuid.UUID, msg *websocket.PreparedMessage, recursive bool) error
	SendToUsers(msg *websocket.PreparedMessage, recursive bool) error
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

	Load(ctx context.Context) error
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

	GetOptions() *entry.SpaceOptions
	SetOptions(options *entry.SpaceOptions, updateDB bool) error

	LoadFromEntry(ctx context.Context, entry *entry.SpaceType) error
}

type Assets2d interface {
	types.Initializer

	GetAsset2d(asset2dID uuid.UUID) (Asset2d, bool)

	Load(ctx context.Context) error
}

type Asset2d interface {
	types.IDer
	types.Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	GetOptions() *entry.Asset2dOptions
	SetOptions(options *entry.Asset2dOptions, updateDB bool) error

	LoadFromEntry(ctx context.Context, entry *entry.Asset2d) error
}

type Assets3d interface {
	types.Initializer

	GetAsset3d(asset3dID uuid.UUID) (Asset3d, bool)

	Load(ctx context.Context) error
}

type Asset3d interface {
	types.IDer
	types.Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	GetOptions() *entry.Asset3dOptions
	SetOptions(options *entry.Asset3dOptions, updateDB bool) error

	LoadFromEntry(ctx context.Context, entry *entry.Asset3d) error
}
