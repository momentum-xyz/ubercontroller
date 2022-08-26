package universe

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/momentum-xyz/controller/pkg/cmath"
	"github.com/momentum-xyz/controller/types"
	"github.com/momentum-xyz/controller/types/generics"
)

type Node interface {
}

type World interface {
	Space
}

type Space interface {
	types.IDer
	types.Initializer

	GetWorld() World

	GetRoot() Space
	GetParent() Space
	SetParent(parent Space, updateDB bool) error

	GetTheta() float64
	SetTheta(theta float64, updateDB bool) error
	GetPosition() cmath.Vec3
	SetPosition(pos cmath.Vec3, updateDB bool) error

	Load(recursive bool) error

	GetOptions() *SpaceOptionsEntry
	SetOptions(options *SpaceOptionsEntry, updateDB bool) error

	GetSpace(spaceID uuid.UUID, recursive bool) (Space, bool)
	GetSpaces(recursive bool) *generics.SyncMap[uuid.UUID, Space]
	AttachSpace(space Space, updateDB bool) error
	AttachSpaces(spaces []Space, updateDB bool) error
	DetachSpace(spaceID uuid.UUID, recursive, updateDB bool) (bool, error)
	DetachSpaces(spaceIDs []uuid.UUID, recursive, updateDB bool) (bool, error)

	GetOwner() User
	SetOwner(owner User, updateDB bool) error

	GetUser(userID uuid.UUID, recursive bool) (User, bool)
	GetUsers(recursive bool) *generics.SyncMap[uuid.UUID, User]
	AttachUser(user User, updateDB bool) error
	DetachUser(user User, updateDB bool) error

	SendToUser(userID uuid.UUID, msg *websocket.PreparedMessage, recursive bool) error
	SendToUsers(msg *websocket.PreparedMessage, recursive bool) error

	//GetSpaceAttributes() *SpaceAttributes
	//GetUserSpaceAttributes() *UserSpaceAttributes
}

type User interface {
	types.IDer

	GetWorld() World
	SetWorld(world World, updateDB bool) error

	GetSpace() Space
	SetSpace(space Space, updateDB bool) error
}

type SpaceTypes interface {
	Load() error

	Get(spaceTypeID uuid.UUID) (SpaceTypes, bool)
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

	Load() error

	GetOptions() *SpaceOptionsEntry
	SetOptions(options *SpaceOptionsEntry, updateDB bool) error
}

type Assets2D interface {
	types.Initializer

	Get(asset2DID uuid.UUID) (Assets2D, bool)

	Load() error
}

type Asset2D interface {
	types.IDer
	types.Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	Load() error

	GetOptions() *SpaceAsset2DOptionsEntry
	SetOptions(options *SpaceAsset2DOptionsEntry, updateDB bool) error
}

type Assets3D interface {
	types.Initializer

	Get(asset3DID uuid.UUID) (Asset3D, bool)

	Load() error
}

type Asset3D interface {
	types.IDer
	types.Initializer

	GetName() string
	SetName(name string, updateDB bool) error

	Load() error

	GetOptions() *SpaceAsset3DOptionsEntry
	SetOptions(options *SpaceAsset3DOptionsEntry, updateDB bool) error
}
