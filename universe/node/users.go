package node

import (
	"context"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// Create new users on this node.
func (n *Node) CreateUsers(ctx context.Context, users ...*entry.User) error {
	if err := n.db.GetUsersDB().UpsertUsers(ctx, users); err != nil {
		return errors.Wrap(err, "create users")
	}
	// TODO: add 'returning` to the query, so we can give back user objects
	// with their filled in database defaults and ID.
	return nil
}

func (n *Node) Filter(predicateFn func(userID umid.UMID, user universe.User) bool) (
	map[umid.UMID]universe.User, error,
) {
	data := make(map[umid.UMID]universe.User)
	userTypeID, err := common.GetNormalUserTypeID()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get normal user type id")
	}

	users, err := n.db.GetUsersDB().GetAllUsers(n.ctx, userTypeID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get all user entries")
	}

	n.Mu.RLock()
	defer n.Mu.RUnlock()

	for _, v := range users {
		loadedUser, _ := n.LoadUser(v.UserID)
		userID := loadedUser.GetID()

		if predicateFn(loadedUser.GetID(), loadedUser) {
			data[userID] = loadedUser
		}
	}

	return data, nil
}
