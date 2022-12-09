package node

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) GetUserSpaceValue(userSpaceID entry.UserSpaceID) (*entry.UserSpaceValue, bool) {
	value, err := n.db.UserSpaceGetValueByID(n.ctx, userSpaceID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (n *Node) UpdateUserSpaceValue(
	userSpaceID entry.UserSpaceID, modifyFn modify.Fn[entry.UserSpaceValue],
) (*entry.UserSpaceValue, error) {
	value, err := n.db.UserSpaceUpdateValueByID(n.ctx, userSpaceID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update space user attribute value")
	}

	return value, nil
}
