package universe

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
)

type AttributeChangeType string

const (
	InvalidAttributeChangeType AttributeChangeType = ""
	ChangedAttributeChangeType AttributeChangeType = "attribute_changed"
	RemovedAttributeChangeType AttributeChangeType = "attribute_removed"
)

type SpaceFilterPredicateFn func(spaceID uuid.UUID, space Space) bool
type WorldsFilterPredicateFn func(worldID uuid.UUID, world World) bool
type SpacesFilterPredicateFn func(spaceID uuid.UUID, space Space) bool
type Assets2dFilterPredicateFn func(asset2dID uuid.UUID, asset2d Asset2d) bool
type Assets3dFilterPredicateFn func(asset3dID uuid.UUID, asset3d Asset3d) bool
type PluginsFilterPredicateFn func(pluginID uuid.UUID, plugin Plugin) bool
type AttributeTypesFilterPredicateFn func(attributeTypeID entry.AttributeTypeID, attributeType AttributeType) bool
type SpaceTypesFilterPredicateFn func(spaceTypeID uuid.UUID, spaceType SpaceType) bool
type UserTypesFilterPredicateFn func(userTypeID uuid.UUID, userType UserType) bool

type WorldSettings struct {
	Kind       string               `db:"kind" json:"kind"`
	Spaces     map[string]uuid.UUID `db:"spaces" json:"spaces"`
	Attributes map[string]uuid.UUID `db:"spaces" json:"attributes"`
	SpaceTypes map[string]uuid.UUID `db:"space_types" json:"space_types"`
	Effects    map[string]uuid.UUID `db:"effects" json:"effects"`
}

type ReservedAttribute struct {
	Name string
	Key  string
}

var (
	ReservedAttributes = struct {
		Node struct {
			GuestUserType  ReservedAttribute
			NormalUserType ReservedAttribute
			WorldTemplate  ReservedAttribute
			JWTKey         ReservedAttribute
		}
		World struct {
			Meta     ReservedAttribute
			Settings ReservedAttribute
		}
		Space struct {
			Name          ReservedAttribute
			Description   ReservedAttribute
			NewsFeedItems ReservedAttribute
			Events        ReservedAttribute
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
			GuestUserType  ReservedAttribute
			NormalUserType ReservedAttribute
			WorldTemplate  ReservedAttribute
			JWTKey         ReservedAttribute
		}{
			GuestUserType: ReservedAttribute{
				Name: "node_settings",
				Key:  "guest_user_type",
			},
			NormalUserType: ReservedAttribute{
				Name: "node_settings",
				Key:  "normal_user_type",
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
			Meta     ReservedAttribute
			Settings ReservedAttribute
		}{
			Meta: ReservedAttribute{
				Name: "world_meta",
			},
			Settings: ReservedAttribute{
				Name: "world_settings",
			},
		},
		Space: struct {
			Name          ReservedAttribute
			Description   ReservedAttribute
			NewsFeedItems ReservedAttribute
			Events        ReservedAttribute
		}{
			Name: ReservedAttribute{
				Name: "name",
				Key:  "name",
			},
			Description: ReservedAttribute{
				Name: "description",
			},
			NewsFeedItems: ReservedAttribute{
				Name: "news_feed",
				Key:  "items",
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
