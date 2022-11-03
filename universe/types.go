package universe

import "github.com/google/uuid"

type SpaceFilterPredicateFn func(spaceID uuid.UUID, space Space) bool

type Attribute struct {
	Name string
	Key  string
}

var (
	Attributes = struct {
		Node struct {
			Settings Attribute
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
		}
	}{
		Node: struct {
			Settings Attribute
		}{
			Settings: Attribute{
				Name: "node_settings",
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
			},
			Description: Attribute{
				Name: "description",
			},
		},
		Kusama: struct {
			User struct {
				Wallet Attribute
			}
		}{
			User: struct {
				Wallet Attribute
			}{
				Wallet: Attribute{
					Name: "wallet",
					Key:  "wallet",
				},
			},
		},
	}
)
