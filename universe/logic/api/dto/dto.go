//go:generate go run gen/mus.go
package dto

import (
	"math/big"
	"time"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
)

type ExploreOptions []ExploreOption

type SearchOptions []ExploreOption

type UserSearchResults []UserSearchResult

type Activities map[umid.UMID]Activity

type Plugins map[umid.UMID]string

type PluginsMeta map[umid.UMID]PluginMeta

type PluginMeta entry.PluginMeta

type PluginsOptions map[umid.UMID]PluginOptions

type PluginOptions *entry.PluginOptions

type ObjectOptions *entry.ObjectOptions

type ObjectSubOptions map[string]any

type ObjectAttributes map[umid.UMID]*entry.ObjectAttribute

type ObjectAttributeValues map[umid.UMID]*entry.AttributeValue

type ObjectSubAttributes map[string]any

type UserSubAttributes map[string]any

type Asset2dMeta entry.Asset2dMeta

type Asset2dOptions *entry.Asset2dOptions

type Assets3dOptions map[umid.UMID]Asset3dOptions

type Asset3dOptions *entry.Asset3dOptions

type Assets3dMeta map[umid.UMID]Asset3dMeta

type Asset3dMeta *entry.Asset3dMeta

type ExploreOption struct {
	ID          umid.UMID `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
}

type Activity struct {
	ActivityID      umid.UMID           `json:"activity_id"`
	UserID          umid.UMID           `json:"user_id"`
	ObjectID        umid.UMID           `json:"object_id"`
	Type            *entry.ActivityType `json:"type"`
	Data            *entry.ActivityData `json:"data"`
	AvatarHash      *string             `json:"avatar_hash"`
	WorldName       string              `json:"world_name"`
	WorldAvatarHash *string             `json:"world_avatar_hash"`
	UserName        *string             `json:"user_name"`
	CreatedAt       time.Time           `json:"created_at"`
}

type RecentWorld struct {
	ID          umid.UMID `json:"id"`
	OwnerID     umid.UMID `json:"owner_id"`
	OwnerName   *string   `json:"owner_name"`
	Description *string   `json:"description"`
	StakeTotal  *string   `json:"stake_total,omitempty"`
	Name        *string   `json:"name"`
	WebsiteLink *string   `json:"website_link"`
	AvatarHash  *string   `json:"avatarHash"`
}

type WorldDetails struct {
	ID                 umid.UMID     `json:"id"`
	OwnerID            umid.UMID     `json:"owner_id"`
	OwnerName          *string       `json:"owner_name"`
	Description        *string       `json:"description"`
	StakeTotal         *string       `json:"stake_total,omitempty"`
	Name               *string       `json:"name"`
	CreatedAt          string        `json:"createdAt,omitempty"`
	UpdatedAt          string        `json:"updatedAt,omitempty"`
	AvatarHash         *string       `json:"avatarHash"`
	WebsiteLink        *string       `json:"website_link"`
	WorldStakers       []WorldStaker `json:"stakers"`
	LastStakingComment *string       `json:"last_staking_comment"`
}

type WorldNFTMeta struct {
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	ExternalURL string          `json:"external_url"`
	Attributes  []NFTAttributes `json:"attributes"`
}

type NFTAttributes struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type WorldStaker struct {
	UserID     umid.UMID `json:"user_id"`
	Name       *string   `json:"name"`
	Stake      *string   `json:"stake,omitempty"`
	AvatarHash *string   `json:"avatarHash"`
}

type TopStaker struct {
	UserID     umid.UMID `json:"user_id"`
	Name       *string   `json:"name"`
	TotalStake *big.Int  `json:"total_stake,omitempty"`
	StakeCount *uint8    `json:"stake_count,omitempty"`
	AvatarHash *string   `json:"avatarHash"`
}

type Profile struct {
	Bio         *string `json:"bio,omitempty"`
	Location    *string `json:"location,omitempty"`
	AvatarHash  *string `json:"avatarHash,omitempty"`
	ProfileLink *string `json:"profileLink,omitempty"`
}

type JWTToken struct {
	Subject      *string `json:"subject,omitempty"`
	Issuer       *string `json:"issuer,omitempty"`
	ExpiresAt    *string `json:"expiresAt,omitempty"`
	IssuedAt     *string `json:"issuedAt,omitempty"`
	SignedString *string `json:"signedString,omitempty"`
}

type HashResponse struct {
	Hash string `json:"hash"`
}

type User struct {
	ID         umid.UMID `json:"id"`
	UserTypeID umid.UMID `json:"userTypeId"`
	Name       string    `json:"name"`
	Wallet     *string   `json:"wallet,omitempty"`
	Profile    Profile   `json:"profile"`
	JWTToken   *string   `json:"token,omitempty"`
	CreatedAt  string    `json:"createdAt"`
	UpdatedAt  string    `json:"updatedAt,omitempty"`
	IsGuest    bool      `json:"isGuest"`
}

type UserSearchResult struct {
	ID      umid.UMID `json:"id"`
	Name    *string   `json:"name,omitempty"`
	Wallet  *string   `json:"wallet,omitempty"`
	Profile Profile   `json:"profile,omitempty"`
}

type RecentUser struct {
	ID      umid.UMID `json:"id"`
	Name    *string   `json:"name,omitempty"`
	Wallet  *string   `json:"wallet,omitempty"`
	Profile Profile   `json:"profile,omitempty"`
}

type Object struct {
	OwnerID      string          `json:"owner_id"`
	ParentID     string          `json:"parent_id"`
	ObjectTypeID string          `json:"object_type_id"`
	Asset2dID    string          `json:"asset_2d_id"`
	Asset3dID    string          `json:"asset_3d_id"`
	Transform    cmath.Transform `json:"transform"`
}

type Member struct {
	ObjectID   *umid.UMID `json:"object_id"`
	UserID     *umid.UMID `json:"user_id"`
	Name       *string    `json:"name"`
	AvatarHash *string    `json:"avatar_hash"`
	Role       *string    `json:"role"`
}

type OwnedWorld struct {
	ID          umid.UMID `json:"id"`
	OwnerID     umid.UMID `json:"owner_id"`
	OwnerName   *string   `json:"owner_name"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description"`
	WebsiteLink *string   `json:"website_link"`
	AvatarHash  *string   `json:"avatarHash"`
}

type StakedWorld struct {
	ID          umid.UMID `json:"id"`
	OwnerID     umid.UMID `json:"owner_id"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description"`
	WebsiteLink *string   `json:"website_link"`
	AvatarHash  *string   `json:"avatarHash"`
}

type Asset2d struct {
	Meta    Asset2dMeta    `json:"meta"`
	Options Asset2dOptions `json:"options"`
}

type Asset3d struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Name      string      `json:"name,omitempty"`
	Meta      Asset3dMeta `json:"meta,omitempty"`
	Private   bool        `json:"is_private"`
	CreatedAt string      `json:"createdAt,omitempty"`
	UpdatedAt string      `json:"updatedAt,omitempty"`
}

type Tile struct {
	ID            string        `json:"id"`
	Hash          string        `json:"hash"`
	ObjectID      string        `json:"objectId"`
	UITypeID      string        `json:"uiTypeId"`
	OwnerID       string        `json:"owner_id"`
	UpdatedAt     string        `json:"updatedAt"`
	PermanentType PermanentType `json:"permanentType"`
	Edited        int           `json:"edited"`
	Render        uint8         `json:"render"`
	Column        int           `json:"column"`
	Row           int           `json:"row"`
	Type          TileType      `json:"type"`
	Content       TileContent   `json:"content"`
}

type TileContent struct {
	Text  *string `json:"text"`
	Title *string `json:"title"`
	Type  *string `json:"type"`
	URL   *string `json:"url"`
}

type Emoji struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Hash     string `json:"hash"`
	Name     string `json:"name"`
	ObjectID string `json:"objectId"`
	Order    int    `json:"order"`
}

type Event struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	HostedBy    string  `json:"hosted_by"`
	ImageHash   *string `json:"image_hash,omitempty"`
	WebLink     *string `json:"web_link,omitempty"`
	ObjectID    string  `json:"objectId"`
	ObjectName  string  `json:"objectName"`
	Start       string  `json:"start"`
	End         string  `json:"end"`
	Created     string  `json:"created"`
	Modified    string  `json:"modified"`
	Attendees   []User  `json:"attendees"`
}

type EventForm struct {
	Start       string  `json:"start"`
	End         string  `json:"end"`
	Title       string  `json:"title"`
	HostedBy    string  `json:"hosted_by"`
	WebLink     *string `json:"web_link"`
	Description string  `json:"description"`
	//Image *File `json:"image,omitempty"` QUESTION: what is it "File"?
}

type Favorite struct {
	ObjectID string `json:"objectId"`
	Name     string `json:"name"`
}

type Miro struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ViewLink    string `json:"viewLink"`
	AccessLink  string `json:"accessLink"`
}

type GoogleDrive struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Broadcast struct {
	URL             string              `json:"url"`
	YoutubeURL      string              `json:"youtubeUrl"`
	BroadcastStatus BroadcastStatusType `json:"broadcastStatus"`
}

type Magic struct {
	ID   string    `json:"id"`
	Key  string    `json:"key"`
	Type MagicType `json:"type"`
	Data struct {
		ID       string  `json:"id"`
		EventID  *string `json:"eventId,omitempty"`
		Position *any    `json:"position,omitempty"`
	} `json:"data"`
	Expired  string `json:"expired"`
	UpdateAt string `json:"update_at"`
	CratedAt string `json:"cratedAt"`
}

type ObjectInfo struct {
	ID          string     `json:"id"`
	ParentID    *string    `json:"parentId,omitempty"`
	ObjectType  ObjectType `json:"objectType"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OwnerID     string     `json:"ownerId"`
	OwnerName   string     `json:"ownerName"`
	CreateAt    string     `json:"createAt"`
	UpdatedAt   string     `json:"updatedAt"`
	IsPrivate   bool       `json:"isPrivate"`
	IsAdmin     bool       `json:"isAdmin"`
}

type ObjectAncestor struct {
	ObjectID   string `json:"objectId"`
	ObjectName string `json:"objectName"`
}

type ObjectMember struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	IsAdmin  bool   `json:"isAdmin"`
}

type StageModeUser struct {
	UserID string `json:"userId"`
	Flag   int    `json:"flag"`
}

type TokenInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Token struct {
	ID              string                     `json:"id"`
	Name            string                     `json:"name"`
	ContractAddress *string                    `json:"contractAddress,omitempty"`
	TokenType       *TokenType                 `json:"tokenType,omitempty"`
	CreatedAt       *string                    `json:"createdAt,omitempty"`
	UpdatedAt       *string                    `json:"updatedAt,omitempty"`
	Status          *TokenRuleReviewStatusType `json:"status,omitempty"`
	WorldID         *string                    `json:"worldId,omitempty"`
	ObjectID        *string                    `json:"objectId,omitempty"`
}

type TokenRule struct {
	ID               string                    `json:"id"`
	Status           TokenRuleReviewStatusType `json:"status"`
	CreatedAt        *string                   `json:"createdAt,omitempty"`
	TokenGroupUserID *string                   `json:"tokenGroupUserId,omitempty"`
	Name             string                    `json:"name"`
	UpdatedAt        *string                   `json:"updatedAt,omitempty"`
	MinBalance       int                       `json:"minBalance"`
	Network          TokenNetworkType          `json:"network"`
	ContractAddress  string                    `json:"contractAddress"`
	TokenType        TokenType                 `json:"tokenType"`
	UserName         *string                   `json:"userName,omitempty"`
	UserID           *string                   `json:"userId,omitempty"`
	ObjectName       *string                   `json:"objectName,omitempty"`
}

type Validator struct {
	ID                 string  `json:"id"`
	ParentID           string  `json:"parentId"`
	ObjectTypeID       string  `json:"objectTypeId"`
	OperatorObjectID   *string `json:"operatorObjectId,omitempty"`
	UITypeID           string  `json:"uiTypeId"`
	OperatorObjectName string  `json:"operatorObjectName"`
	Name               string  `json:"name"`
	IsFavorited        bool    `json:"isFavorited"`
	Metadata           struct {
		KusamaMetadata KusamaMetaData `json:"kusama_metadata"`
	} `json:"metadata"`
}

type KusamaMetaData struct {
	ValidatorID     string `json:"validator_id"`
	ValidatorReward int    `json:"validator_reward"`
	ValidatorInfo   struct {
		AccountID string `json:"account_id"`
		Entity    struct {
			Name      string `json:"name"`
			AccountID string `json:"accountId"`
		} `json:"entity"`
		Commission              any    `json:"commission"` // QUESTION: originally here is "string | number", really?
		OwnStake                any    `json:"ownStake"`   // QUESTION: originally here is "string | number", really?
		Status                  string `json:"status"`
		TotalStake              string `json:"totalStake"`
		ValidatorAccountDetails struct {
			Name string `json:"name"`
		} `json:"validatorAccountDetails"`
	} `json:"validator_info"`
}

type VibeAmount struct {
	ObjectID string `json:"objectId"`
	Amount   int    `json:"amount"`
}

type Plugin struct {
	Name      string  `json:"name"`
	Title     string  `json:"title"`
	SubTitle  *string `json:"subTitle,omitempty"`
	ScriptURL string  `json:"scriptUrl"`
	IconName  *string `json:"iconName,omitempty"`
}

type Stake struct {
	ObjectID     umid.UMID `json:"object_id"`
	Name         string    `json:"name"`
	WalletID     string    `json:"wallet_id"`
	BlockchainID umid.UMID `json:"blockchain_id"`
	Amount       string    `json:"amount,omitempty"`
	Reward       string    `json:"reward"`
	LastComment  string    `json:"last_comment"`
	AvatarHash   string    `json:"avatar_hash"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type WalletInfo struct {
	WalletID       string    `json:"wallet_id"`
	ContractID     string    `json:"contract_id"`
	Balance        string    `json:"balance"`
	BlockchainName string    `json:"blockchain_name"`
	UpdatedAt      time.Time `json:"updated_at"`
	Reward         string    `json:"reward"`
	Transferable   string    `json:"transferable"`
	Staked         string    `json:"staked"`
	Unbonding      string    `json:"unbonding"`
}
