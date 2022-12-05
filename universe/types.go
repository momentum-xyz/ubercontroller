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

type Attribute struct {
	Name string
	Key  string
}

var (
	Attributes = struct {
		Node struct {
			GuestUserType  Attribute
			NormalUserType Attribute
			JWTKey         Attribute
		}
		World struct {
			Meta Attribute
		}
		Space struct {
			Name        Attribute
			Description Attribute
		}
		Kusama struct {
			User struct {
				Wallet Attribute
			}
			Challenges Attribute
		}
	}{
		Node: struct {
			GuestUserType  Attribute
			NormalUserType Attribute
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
			JWTKey: Attribute{
				Name: "jwt_key",
				Key:  "secret",
			},
		},
		World: struct {
			Meta Attribute
		}{
			Meta: Attribute{
				Name: "world_meta",
			},
		},
		Space: struct {
			Name        Attribute
			Description Attribute
		}{
			Name: Attribute{
				Name: "name",
				Key:  "name",
			},
			Description: Attribute{
				Name: "description",
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
	}
)
