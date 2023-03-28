package node

import (
	"context"
	"fmt"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/user"
)

func (n *Node) LoadUser(userID umid.UMID) (universe.User, error) {
	user := user.NewUser(userID, n.db)
	if err := user.Initialize(n.ctx); err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize user: %s", userID)
	}

	if err := user.Load(); err != nil {
		return nil, errors.WithMessagef(err, "failed to load user: %s", userID)
	}

	fmt.Printf("%+v\n", user.GetPosition())
	user.SetPosition(cmath.Vec3{X: 50, Y: 50, Z: 150})
	fmt.Printf("%+v\n", user.GetPosition())
	return user, nil
}

// Create new users on this node.
func (n *Node) CreateUsers(ctx context.Context, users ...*entry.User) error {
	if err := n.db.GetUsersDB().UpsertUsers(ctx, users); err != nil {
		return errors.Wrap(err, "create users")
	}
	// TODO: add 'returning` to the query, so we can give back user objects
	// with their filled in database defaults and ID.
	return nil
}
