package universe

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type AttributeChangeType string

const (
	InvalidAttributeChangeType AttributeChangeType = ""
	ChangedAttributeChangeType AttributeChangeType = "attribute_changed"
	RemovedAttributeChangeType AttributeChangeType = "attribute_removed"
)

type ObjectFilterPredicateFn func(objectID umid.UMID, object Object) bool
type WorldsFilterPredicateFn func(worldID umid.UMID, world World) bool
type UsersFilterPredicateFn func(userID umid.UMID, user User) bool
type ObjectsFilterPredicateFn func(objectID umid.UMID, object Object) bool
type Assets2dFilterPredicateFn func(asset2dID umid.UMID, asset2d Asset2d) bool
type Assets3dFilterPredicateFn func(asset3dUserID AssetUserIDPair, asset3d Asset3d) bool
type PluginsFilterPredicateFn func(pluginID umid.UMID, plugin Plugin) bool
type AttributeTypesFilterPredicateFn func(attributeTypeID entry.AttributeTypeID, attributeType AttributeType) bool
type ObjectTypesFilterPredicateFn func(objectTypeID umid.UMID, objectType ObjectType) bool
type UserTypesFilterPredicateFn func(userTypeID umid.UMID, userType UserType) bool

type WorldSettings struct {
	Kind        string               `db:"kind" json:"kind"`
	Objects     map[string]umid.UMID `db:"objects" json:"objects"`
	Attributes  map[string]umid.UMID `db:"attributes" json:"attributes"`
	ObjectTypes map[string]umid.UMID `db:"object_types" json:"object_types"`
	Effects     map[string]umid.UMID `db:"effects" json:"effects"`
}

type AssetUserIDPair struct {
	AssetID umid.UMID
	UserID  umid.UMID
}

type SortType string

const (
	ASC  SortType = "ASC"
	DESC SortType = "DESC"
)

type ReservedAttribute struct {
	Name string
	Key  string
}

var (
	ReservedAttributes = struct {
		Node struct {
			GuestUserType    ReservedAttribute
			NormalUserType   ReservedAttribute
			PortalObjectType ReservedAttribute
			WorldTemplate    ReservedAttribute
			JWTKey           ReservedAttribute
		}
		World struct {
			Meta                ReservedAttribute
			Settings            ReservedAttribute
			TeleportDestination ReservedAttribute
		}
		Object struct {
			Name           ReservedAttribute
			Description    ReservedAttribute
			WorldAvatar    ReservedAttribute
			WebsiteLink    ReservedAttribute
			NewsFeedItems  ReservedAttribute
			PortalDockFace ReservedAttribute
			Events         ReservedAttribute
		}
		Kusama struct {
			User struct {
				Wallet ReservedAttribute
			}
			Challenges ReservedAttribute
		}
		User struct {
			HighFive ReservedAttribute
		}
	}{
		Node: struct {
			GuestUserType    ReservedAttribute
			NormalUserType   ReservedAttribute
			PortalObjectType ReservedAttribute
			WorldTemplate    ReservedAttribute
			JWTKey           ReservedAttribute
		}{
			GuestUserType: ReservedAttribute{
				Name: "node_settings",
				Key:  "guest_user_type",
			},
			NormalUserType: ReservedAttribute{
				Name: "node_settings",
				Key:  "normal_user_type",
			},
			PortalObjectType: ReservedAttribute{
				Name: "node_settings",
				Key:  "docking_hub_object_type",
			},
			WorldTemplate: ReservedAttribute{
				Name: "world_template",
			},
			JWTKey: ReservedAttribute{
				Name: "jwt_key",
				Key:  "secret",
			},
		},
		World: struct {
			Meta                ReservedAttribute
			Settings            ReservedAttribute
			TeleportDestination ReservedAttribute
		}{
			Meta: ReservedAttribute{
				Name: "world_meta",
			},
			Settings: ReservedAttribute{
				Name: "world_settings",
			},
			TeleportDestination: ReservedAttribute{
				Name: "teleport",
				Key:  "DestinationWorldID",
			},
		},
		Object: struct {
			Name           ReservedAttribute
			Description    ReservedAttribute
			WorldAvatar    ReservedAttribute
			WebsiteLink    ReservedAttribute
			NewsFeedItems  ReservedAttribute
			PortalDockFace ReservedAttribute
			Events         ReservedAttribute
		}{
			Name: ReservedAttribute{
				Name: "name",
				Key:  "name",
			},
			Description: ReservedAttribute{
				Name: "description",
				Key:  "description",
			},
			WorldAvatar: ReservedAttribute{
				Name: "world_avatar",
				Key:  "render_hash",
			},
			WebsiteLink: ReservedAttribute{
				Name: "website_link",
				Key:  "website_link",
			},
			NewsFeedItems: ReservedAttribute{
				Name: "news_feed",
				Key:  "items",
			},
			PortalDockFace: ReservedAttribute{
				Name: "dock_face",
				Key:  "render_hash",
			},
			Events: ReservedAttribute{
				Name: "events",
				Key:  "",
			},
		},
		Kusama: struct {
			User struct {
				Wallet ReservedAttribute
			}
			Challenges ReservedAttribute
		}{
			User: struct {
				Wallet ReservedAttribute
			}{
				Wallet: ReservedAttribute{
					Name: "wallet",
					Key:  "wallet",
				},
			},
			Challenges: ReservedAttribute{
				Name: "challenge_store",
				Key:  "challenges",
			},
		},
		User: struct {
			HighFive ReservedAttribute
		}{
			HighFive: ReservedAttribute{
				Name: "high_five",
				Key:  "counter",
			},
		},
	}
)
