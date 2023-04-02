package node

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type OperationType uint
type AttributeKind uint

const (
	ReadOperation OperationType = iota
	WriteOperation
)

const (
	ObjectAttribute AttributeKind = iota
	ObjectUserAttribute
	UserAttribute
	UserUserAttribute
)

const (
	Any        string = "any"
	User       string = "user"
	Admin      string = "admin"
	Owner      string = "owner"
	TargetUser string = "target_user"
	None       string = "none"
)

func (n *Node) AssessPermissions(
	pluginID umid.UMID, attributeName string, userID umid.UMID, ownerID umid.UMID,
	operationType OperationType, attributeKind AttributeKind,
) (bool, error) {
	// If not available fall back to default
	// Handle exceptions

	defaultPermissions := map[string]string{
		"read":  "any",
		"write": "admin+user_owner",
	}

	attributeTypeID := entry.NewAttributeTypeID(pluginID, attributeName)
	attributeType, ok := n.GetAttributeTypes().GetAttributeType(attributeTypeID)
	if !ok {
		return false, errors.New("failed to get attributeType")
	}

	options := attributeType.GetOptions()
	permissions := make(map[string]string)
	if options != nil {
		// Todo: remove hardcode
		permissions = utils.GetFromAnyMap(*options, "permissions", map[string]string(nil))

		if permissions == nil {
			permissions = defaultPermissions
		}
	}

	switch operationType {
	case ReadOperation:
		return n.AssessReadOperation(userID, ownerID, permissions["read"], attributeKind)
	case WriteOperation:
		return n.AssessWriteOperation(userID, ownerID, permissions["write"], attributeKind)
	}

	return false, nil
}

func (n *Node) AssessReadOperation(userID umid.UMID, ownerID umid.UMID, permissions string, attributeKind AttributeKind) (bool, error) {
	userPermissions := make([]string, 0)
	attributeTypePermissions := make([]string, 0)
	if strings.Contains(permissions, "+") {
		attributeTypePermissions = strings.Split(permissions, "+")
	} else {
		attributeTypePermissions = append(attributeTypePermissions, permissions)
	}

	// Is the user a registered user or a guest?
	user, _ := n.GetUser(userID, false)
	userType := user.GetUserType()

	guestUserTypeID, _ := common.GetGuestUserTypeID()

	if userType.GetID() != guestUserTypeID {
		userPermissions = append(userPermissions, User)
	}

	switch attributeKind {
	case ObjectAttribute:
		userObjectID := entry.NewUserObjectID(userID, ownerID)
		// What rights does the user have?
		// Does the user own the object?
		object, ok := n.GetObjectFromAllObjects(ownerID)
		if !ok {
			return false, errors.New("failed to get object from all objects")
		}

		objectOwnerID := object.GetOwnerID()
		if objectOwnerID == userID {
			userPermissions = append(userPermissions, Owner)
		}

		isAdmin, _ := n.db.GetUserObjectsDB().CheckIsIndirectAdminByID(n.ctx, userObjectID)
		if isAdmin {
			userPermissions = append(userPermissions, Admin)
		}
		// Then what is the user allowed to do?
		permission := n.CompareReadPermissions(attributeTypePermissions, userPermissions)
		return permission, nil
		//case ObjectUserAttribute:
		//	dbInstance = n.db.GetObjectUserAttributesDB()
		//case UserAttribute:
		//	dbInstance = n.db.GetUserAttributesDB()
		//case UserUserAttribute:
		//	dbInstance = n.db.GetUserUserAttributesDB()
	}

	return false, nil
}

func (n *Node) AssessWriteOperation(
	userID umid.UMID, ownerID umid.UMID, permissions string,
	attributeKind AttributeKind,
) (bool, error) {
	sl := make([]string, 0)

	if strings.Contains(permissions, "+") {
		sl = strings.Split(permissions, "+")
	} else {
		sl = append(sl, permissions)
	}

	for _, permission := range sl {
		switch permission {
		case Owner:
			//
		case Admin:
			//
		case User:
			//
		case Any:
			//
		}
	}

	return false, nil
}

func (n *Node) CompareReadPermissions(attributeTypePermissions []string, userPermissions []string) bool {
	for _, attributeTypePermission := range attributeTypePermissions {
		switch attributeTypePermission {
		case Any:
			return true
		case User:
			if slices.Contains(userPermissions, User) ||
				slices.Contains(userPermissions, Admin) ||
				slices.Contains(userPermissions, Owner) {
				return true
			}
			return false
		case Admin:
			if slices.Contains(userPermissions, Admin) ||
				slices.Contains(userPermissions, Owner) {
				return true
			}
			return false
		case Owner:
			if slices.Contains(userPermissions, Owner) {
				return true
			}
			return false
		}
	}

	return false
}
