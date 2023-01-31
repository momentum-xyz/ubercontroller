package seed

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedUserTypes(node universe.Node) error {
	type item struct {
		id           uuid.UUID
		userTypeName string
		description  *string
		options      *entry.UserOptions
	}

	items := []*item{
		{
			id:           uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			userTypeName: "Deity",
			description:  utils.GetPTR("They rule the world"),
			options: &entry.UserOptions{
				IsGuest: utils.GetPTR(false),
			},
		},
		{
			id:           uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			userTypeName: "User",
			description:  utils.GetPTR("Momentum user"),
			options: &entry.UserOptions{
				IsGuest: utils.GetPTR(false),
			},
		},
	}

	for _, item := range items {
		userType, err := node.GetUserTypes().CreateUserType(item.id)
		if err != nil {
			return errors.WithMessagef(err, "failed to create user type: %s", item.id)
		}

		if err := userType.SetName(item.userTypeName, false); err != nil {
			return errors.WithMessagef(err, "failed to set user type name: %s %s", item.id, item.userTypeName)
		}

		if err := userType.SetDescription(item.description, false); err != nil {
			return errors.WithMessagef(err, "failed to set user type description: %s", item.id)
		}

		_, err = userType.SetOptions(modify.MergeWith(item.options), false)
		if err != nil {
			return errors.WithMessagef(err, "failed to set user type options: %s", item.id)
		}
	}

	return nil
}
