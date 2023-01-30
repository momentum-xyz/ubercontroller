package seed

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

// TODO Add to seed.go
func seedUsers(node universe.Node) error {
	type item struct {
		id         uuid.UUID
		userTypeID uuid.UUID
		profile    entry.UserProfile
		options    *entry.UserOptions
	}

	items := []*item{
		{
			id:         uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			userTypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			profile:    entry.UserProfile{},
			options:    nil,
		},
	}

	for _, item := range items {
		fmt.Println(item)
		//node.AddUser()
	}

	return nil
}
