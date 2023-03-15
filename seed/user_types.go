package seed

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedUserTypes(ctx context.Context, node universe.Node, db database.DB) error {
	type item struct {
		id           mid.ID
		userTypeName string
		description  *string
		options      *entry.UserOptions
	}

	items := []*item{
		{
			id:           mid.MustParse("00000000-0000-0000-0000-000000000002"),
			userTypeName: "Deity",
			description:  utils.GetPTR("They rule the world"),
			options: &entry.UserOptions{
				IsGuest: utils.GetPTR(false),
			},
		},
		{
			id:           mid.MustParse(normalUserTypeID),
			userTypeName: "User",
			description:  utils.GetPTR("Momentum user"),
			options: &entry.UserOptions{
				IsGuest: utils.GetPTR(false),
			},
		},
		{
			id:           mid.MustParse(guestUserTypeID),
			userTypeName: "Temporary User",
			description:  utils.GetPTR("Temporary Momentum user"),
			options: &entry.UserOptions{
				IsGuest: utils.GetPTR(true),
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

		if err := userType.SetDescription(*item.description, false); err != nil {
			return errors.WithMessagef(err, "failed to set user type description: %s", item.id)
		}

		_, err = userType.SetOptions(modify.MergeWith(item.options), false)
		if err != nil {
			return errors.WithMessagef(err, "failed to set user type options: %s", item.id)
		}

		if err := db.GetUserTypesDB().UpsertUserType(ctx, userType.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to upsert user_type")
		}
	}

	return nil
}
