// TODO: deal with this mess
package dto

type Asset3dType int8

const (
	AddressableAssetType Asset3dType = iota
	GLTFAsset3dType
	BasicAsset3dType
)

type FlyToMeType string

const (
	FlyToMeTrigger FlyToMeType = "fly-to-me"
)

type UserStatusType string

const (
	UnknownUserStatusType      UserStatusType = ""
	OnlineUserStatusType       UserStatusType = "online"
	DoNotDisturbUserStatusType UserStatusType = "dnd"
	AwayUserStatusType         UserStatusType = "away"
	InvisibleUserStatusType    UserStatusType = "invisible"
)

type PermanentType string

const (
	UnknownPermanentType        PermanentType = ""
	NonePermanentType           PermanentType = "none"
	PosterPermanentType         PermanentType = "poster"
	MemePermanentType           PermanentType = "meme"
	LogoPermanentType           PermanentType = "logo"
	DescriptionPermanentType    PermanentType = "description"
	VideoPermanentPermanentType PermanentType = "video"
	NamePermanentType           PermanentType = "name"
	ProblemPermanentType        PermanentType = "problem"
	SolutionPermanentType       PermanentType = "solution"
	ThirdPermanentType          PermanentType = "third"
)

type TileType string

const (
	UnknownTileType TileType = ""
	TextTileType    TileType = "TILE_TYPE_TEXT"
	MediaTileType   TileType = "TILE_TYPE_MEDIA"
	VideoTileType   TileType = "TILE_TYPE_VIDEO"
)

type BroadcastStatusType string

const (
	UnknownBroadcastStatusType    BroadcastStatusType = ""
	ForceSmallBroadcastStatusType BroadcastStatusType = "force_small"
	PlaySmallBroadcastStatusType  BroadcastStatusType = "play_small"
	ForceLargeBroadcastStatusType BroadcastStatusType = "force_large"
	PlayLargeBroadcastStatusType  BroadcastStatusType = "play_large"
	PlayBroadcastStatusType       BroadcastStatusType = "play"
	StopBroadcastStatusType       BroadcastStatusType = "stop"
)

type MagicType string

const (
	UnknownMagicType     MagicType = ""
	OpenObjectMagicType  MagicType = "open_object"
	JoinMeetingMagicType MagicType = "join_meeting"
	FlyMagicType         MagicType = "fly"
	EventMagicType       MagicType = "event"
)

type StageModeStatusType string

const (
	UnknownStageModeStatusType   StageModeStatusType = ""
	InitiatedStageModeStatusType StageModeStatusType = "initiated"
	StoppedStageModeStatusType   StageModeStatusType = "stopped"
)

type ObjectType string

const (
	UnknownObjectType    ObjectType = ""
	WorldObjectType      ObjectType = "world"
	ProgramObjectType    ObjectType = "program"
	ChallengeObjectType  ObjectType = "challenge"
	ProjectObjectType    ObjectType = "project"
	GrabATableObjectType ObjectType = "grab-a-table"
)

type StageModeRequestType string

const (
	UnknownStageModeRequestType StageModeRequestType = ""
	RequestStageModeRequestType StageModeRequestType = "request"
	InviteStageModeRequestType  StageModeRequestType = "invite"
	AcceptStageModeRequestType  StageModeRequestType = "accept"
	DeclineStageModeRequestType StageModeRequestType = "decline"
)

type TokenType string

const (
	UnknownTokenType TokenType = ""
	ERC20TokenType   TokenType = "ERC20"
	ERC721TokenType  TokenType = "ERC721"
	ERC1155TokenType TokenType = "ERC1155"
)

type TokenRuleReviewStatusType string

const (
	UnknownTokenRuleReviewStatusType  TokenRuleReviewStatusType = ""
	ApprovedTokenRuleReviewStatusType TokenRuleReviewStatusType = "approved"
	DeniedTokenRuleReviewStatusType   TokenRuleReviewStatusType = "denied"
)

type TokenNetworkType string

const (
	UnknownTokenNetworkType  TokenNetworkType = ""
	MoonbeamTokenNetworkType TokenNetworkType = "moonbeam"
	EthereumTokenNetworkType TokenNetworkType = "eth_mainnet"
)
