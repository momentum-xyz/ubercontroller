package node

import (
	"github.com/jackc/pgx/v4"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.UserObjects = (*userObjects)(nil)

type userObjects struct {
	node *Node
}

func newUserObjects(node *Node) *userObjects {
	return &userObjects{
		node: node,
	}
}

func (uo *userObjects) GetValue(userObjectID entry.UserObjectID) (*entry.UserObjectValue, bool) {
	value, err := uo.node.db.GetUserObjectsDB().GetUserObjectValueByID(uo.node.ctx, userObjectID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (uo *userObjects) GetObjectIndirectAdmins(objectID umid.UMID) ([]*umid.UMID, bool) {
	admins, err := uo.node.db.GetUserObjectsDB().GetObjectIndirectAdmins(uo.node.ctx, objectID)
	if err != nil {
		return nil, false
	}
	return admins, true
}

func (uo *userObjects) CheckIsIndirectAdmin(userObjectID entry.UserObjectID) (bool, error) {
	isAdmin, err := uo.node.db.GetUserObjectsDB().CheckIsIndirectAdminByID(uo.node.ctx, userObjectID)
	if err != nil {
		return false, errors.WithMessage(err, "failed to check is indirect admin by umid")
	}
	return isAdmin, nil
}

func (uo *userObjects) Upsert(
	userObjectID entry.UserObjectID, modifyFn modify.Fn[entry.UserObjectValue], updateDB bool,
) (*entry.UserObjectValue, error) {
	value, err := uo.node.db.GetUserObjectsDB().UpsertUserObject(uo.node.ctx, userObjectID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upsert user object")
	}
	return value, nil
}

func (uo *userObjects) UpdateValue(
	userObjectID entry.UserObjectID, modifyFn modify.Fn[entry.UserObjectValue], updateDB bool,
) (*entry.UserObjectValue, error) {
	value, err := uo.node.db.GetUserObjectsDB().UpdateUserObjectValue(uo.node.ctx, userObjectID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user object value")
	}
	return value, nil
}

func (uo *userObjects) Remove(userObjectID entry.UserObjectID, updateDB bool) (bool, error) {
	if err := uo.node.db.GetUserObjectsDB().RemoveUserObjectByID(uo.node.ctx, userObjectID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.WithMessage(err, "failed to remove user object by umid")
	}
	return true, nil
}

func (uo *userObjects) RemoveMany(userObjectIDs []entry.UserObjectID, updateDB bool) (bool, error) {
	if err := uo.node.db.GetUserObjectsDB().RemoveUserObjectsByIDs(uo.node.ctx, userObjectIDs); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.WithMessage(err, "failed to remove user objects by ids")
	}
	return true, nil
}
