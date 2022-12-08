package universe

import (
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

type WorldSettings struct {
	Kind       string               `db:"kind" json:"kind"`
	Spaces     map[string]uuid.UUID `db:"spaces" json:"spaces"`
	Attributes map[string]uuid.UUID `db:"spaces" json:"attributes"`
	SpaceTypes map[string]uuid.UUID `db:"space_types" json:"space_types"`
	Effects    map[string]uuid.UUID `db:"effects" json:"effects"`
}

type Attribute struct {
	Name string
	Key  string
}

var (
	Attributes = struct {
		Node struct {
			GuestUserType  Attribute
			NormalUserType Attribute
			WorldTemplate  Attribute
			JWTKey         Attribute
		}
		World struct {
			Meta     Attribute
			Settings Attribute
		}
		Space struct {
			Name          Attribute
			Description   Attribute
			NewsFeedItems Attribute
		}
		Kusama struct {
			User struct {
				Wallet Attribute
			}
			Challenges Attribute
		}
		User struct {
			HighFive Attribute
		}
	}{
		Node: struct {
			GuestUserType  Attribute
			NormalUserType Attribute
			WorldTemplate  Attribute
			JWTKey         Attribute
		}{
			GuestUserType: Attribute{
				Name: "node_settings",
				Key:  "guest_user_type",
			},
			NormalUserType: Attribute{
				Name: "node_settings",
				Key:  "normal_user_type",
			},
			WorldTemplate: Attribute{
				Name: "world_template",
			},
			JWTKey: Attribute{
				Name: "jwt_key",
				Key:  "secret",
			},
		},
		World: struct {
			Meta     Attribute
			Settings Attribute
		}{
			Meta: Attribute{
				Name: "world_meta",
			},
			Settings: Attribute{
				Name: "world_settings",
			},
		},
		Space: struct {
			Name          Attribute
			Description   Attribute
			NewsFeedItems Attribute
		}{
			Name: Attribute{
				Name: "name",
				Key:  "name",
			},
			Description: Attribute{
				Name: "description",
			},
			NewsFeedItems: Attribute{
				Name: "news_feed",
				Key:  "items",
			},
		},
		Kusama: struct {
			User struct {
				Wallet Attribute
			}
			Challenges Attribute
		}{
			User: struct {
				Wallet Attribute
			}{
				Wallet: Attribute{
					Name: "wallet",
					Key:  "wallet",
				},
			},
			Challenges: Attribute{
				Name: "challenge_store",
				Key:  "challenges",
			},
		},
		User: struct {
			HighFive Attribute
		}{
			HighFive: Attribute{
				Name: "high_five",
				Key:  "counter",
			},
		},
	}
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
