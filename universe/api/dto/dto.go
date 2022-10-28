package dto

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
)

type Plugins map[uuid.UUID]string

type PluginsMeta map[uuid.UUID]*PluginMeta

type PluginMeta entry.PluginMeta

type PluginsOptions map[uuid.UUID]*PluginOptions

type PluginOptions entry.PluginOptions

type SpaceOptions map[uuid.UUID]*entry.SpaceOptions

type SpaceSubOptions map[uuid.UUID]any

type SpaceAttributes map[uuid.UUID]*entry.SpaceAttribute

type SpaceSubAttributes map[uuid.UUID]any

type SpaceEffectiveOptions SpaceOptions

type SpaceEffectiveSubOptions map[uuid.UUID]any

type Assets3d map[uuid.UUID]*entry.Asset3d

type Profile struct {
	Bio         *string `json:"bio,omitempty"`
	Location    *string `json:"location,omitempty"`
	AvatarHash  *string `json:"avatarHash,omitempty"`
	ProfileLink *string `json:"profileLink,omitempty"`
	OnBoarded   *bool   `json:"onBoarded,omitempty"`
	ImageHash   string  `json:"imageHash"`
}

type User struct {
	ID          string          `json:"id"`
	UserTypeID  string          `json:"userTypeId"`
	Wallet      *string         `json:"wallet"`
	Name        string          `json:"name"`
	Email       *string         `json:"email,omitempty"`
	Description *string         `json:"description"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   *string         `json:"updatedAt"`
	IsNodeAdmin bool            `json:"isNodeAdmin"`
	Status      *UserStatusType `json:"status,omitempty"`
	Profile     Profile         `json:"profile"`
}

type Tile struct {
	ID            string        `json:"id"`
	Hash          string        `json:"hash"`
	SpaceID       string        `json:"spaceId"`
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
	ID      string `json:"id"`
	Code    string `json:"code"`
	Hash    string `json:"hash"`
	Name    string `json:"name"`
	SpaceID string `json:"spaceId"`
	Order   int    `json:"order"`
}

type Event struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	HostedBy    string  `json:"hosted_by"`
	ImageHash   *string `json:"image_hash,omitempty"`
	WebLink     *string `json:"web_link,omitempty"`
	SpaceID     string  `json:"spaceId"`
	SpaceName   string  `json:"spaceName"`
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
	SpaceID string `json:"spaceId"`
	Name    string `json:"name"`
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

type SpaceInfo struct {
	ID          string    `json:"id"`
	ParentID    *string   `json:"parentId,omitempty"`
	SpaceType   SpaceType `json:"spaceType"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"ownerId"`
	OwnerName   string    `json:"ownerName"`
	CreateAt    string    `json:"createAt"`
	UpdatedAt   string    `json:"updatedAt"`
	IsPrivate   bool      `json:"isPrivate"`
	IsAdmin     bool      `json:"isAdmin"`
}

type Space struct {
	SpaceInfo
	UITypeID string `json:"uiTypeId"`
	IsMember bool   `json:"isMember"`
	IsOwner  bool   `json:"isOwner"`
	Metadata *struct {
		KusamaMetadata *struct {
			OperatorID *string `json:"operator_id,omitempty"`
		}
	} `json:"metadata,omitempty"`
}

type SubSpace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        SpaceType `json:"type"`
	SubSpaces   []struct {
		ID           string    `json:"id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		Type         SpaceType `json:"type"`
		HasSubSpaces bool      `json:"hasSubSpaces"`
	} `json:"subSpaces"`
}

type SpaceAncestor struct {
	SpaceID   string `json:"spaceId"`
	SpaceName string `json:"spaceName"`
}

type SpaceMember struct {
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
	SpaceID         *string                    `json:"spaceId,omitempty"`
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
	SpaceName        *string                   `json:"spaceName,omitempty"`
}

type Validator struct {
	ID                string  `json:"id"`
	ParentID          string  `json:"parentId"`
	SpaceTypeID       string  `json:"spaceTypeId"`
	OperatorSpaceID   *string `json:"operatorSpaceId,omitempty"`
	UITypeID          string  `json:"uiTypeId"`
	OperatorSpaceName string  `json:"operatorSpaceName"`
	Name              string  `json:"name"`
	IsFavorited       bool    `json:"isFavorited"`
	Metadata          struct {
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
	SpaceID string `json:"spaceId"`
	Amount  int    `json:"amount"`
}

type Plugin struct {
	Name      string  `json:"name"`
	Title     string  `json:"title"`
	SubTitle  *string `json:"subTitle,omitempty"`
	ScriptURL string  `json:"scriptUrl"`
	IconName  *string `json:"iconName,omitempty"`
}
