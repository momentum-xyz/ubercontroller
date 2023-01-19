package node

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) GetUserSpaceValue(userSpaceID entry.UserObjectID) (*entry.UserObjectValue, bool) {
	value, err := n.db.GetUserObjectDB().GetValueByID(n.ctx, userSpaceID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (n *Node) UpdateUserSpaceValue(
	userSpaceID entry.UserObjectID, modifyFn modify.Fn[entry.UserObjectValue],
) (*entry.UserObjectValue, error) {
	value, err := n.db.GetUserObjectDB().UpdateValueByID(n.ctx, userSpaceID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update space user attribute value")
	}

	return value, nil
}
