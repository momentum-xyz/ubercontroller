package seed

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

// TODO Add to seed.go
func seedUsers(ctx context.Context, node universe.Node, db database.DB) error {
	type item struct {
		id         uuid.UUID
		userTypeID uuid.UUID
		profile    entry.UserProfile
		options    *entry.UserOptions
	}

	// 1 Odin (Deity)
	// 2 Node Admin (User)
	items := []*item{
		{
			id:         uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			userTypeID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			profile:    entry.UserProfile{},
			options:    nil,
		},
	}

	for _, item := range items {

		entry := &entry.User{
			UserID:     item.id,
			UserTypeID: item.userTypeID,
			Profile:    entry.UserProfile{},
			Auth:       map[string]any{},
			Options:    nil,
		}

		if err := db.GetUsersDB().UpsertUser(ctx, entry); err != nil {
			return errors.WithMessage(err, "failed to upsert user")
		}

		// TODO Make it work
		//userItem := user.NewUser(item.id, db)
		//if err := userItem.Load(); err != nil {
		//	return errors.WithMessage(err, "failed to load user")
		//}
		//
		//if err := node.AddUser(userItem, true); err != nil {
		//	return errors.WithMessage(err, "failed to add user to node")
		//}
	}

	return nil
}
