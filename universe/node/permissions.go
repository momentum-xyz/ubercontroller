package node

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type OperationType uint
type AttributeKind uint
type Permission uint

const (
	ReadOperation OperationType = iota
	WriteOperation
)

const (
	ObjectAttribute AttributeKind = iota
	ObjectUserAttribute
	UserAttribute
	UserUserAttribute
	NodeAttribute
)

const (
	Any        string = "any"
	User       string = "user"
	UserOwner  string = "user_owner"
	Admin      string = "admin"
	Owner      string = "owner"
	TargetUser string = "target_user"
	None       string = "none"
)

func (n *Node) AssessPermissions(
	c *gin.Context, pluginID umid.UMID, attributeName string, ownerID umid.UMID,
	operationType OperationType, attributeKind AttributeKind,
) (bool, error) {
	// If not available fall back to default
	// Handle exceptions

	defaultPermissions := map[string]string{
		"read":  "any",
		"write": "admin+user_owner",
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
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

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		return false, errors.WithMessage(err, "failed to get userID from context")
	}

	switch operationType {
	case ReadOperation:
		return n.AssessReadOperation(userID, ownerID, permissions["read"], attributeKind, attributeID)
	case WriteOperation:
		return n.AssessWriteOperation(userID, ownerID, permissions["write"], attributeKind, attributeID)
	}

	return false, nil
}

func (n *Node) AssessReadOperation(
	userID umid.UMID,
	ownerID umid.UMID, permissions string,
	attributeKind AttributeKind, attributeID entry.AttributeID,
) (bool, error) {
	userPermissions := make([]string, 0)
	attributeTypePermissions := make([]string, 0)
	if strings.Contains(permissions, "+") {
		attributeTypePermissions = strings.Split(permissions, "+")
	} else {
		attributeTypePermissions = append(attributeTypePermissions, permissions)
	}

	// Is the user a registered user or a guest?
	user, err := n.LoadUser(userID)
	if err != nil {
		return false, errors.WithMessage(err, "failed to load user")
	}

	userType := user.GetUserType()
	guestUserTypeID, err := common.GetGuestUserTypeID()
	if err != nil {
		return false, errors.WithMessage(err, "failed to get guest user type id")
	}

	if userType.GetID() != guestUserTypeID {
		userPermissions = append(userPermissions, User)
	}

	switch attributeKind {
	case ObjectAttribute:
		// any, users, admin
	case ObjectUserAttribute:
		// any, users, admin, user_owner, admin+user_owner
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

		isAdmin, err := n.db.GetUserObjectsDB().CheckIsIndirectAdminByID(n.ctx, userObjectID)
		if err != nil {
			return false, errors.WithMessage(err, "failed to check admin status")
		}
		if isAdmin {
			userPermissions = append(userPermissions, Admin)
		}
		// Then what is the user allowed to do?
		permission := n.CompareReadPermissions(attributeTypePermissions, userPermissions)
		return permission, nil
	case UserAttribute:
		// any, users, user_owner
		ownerUser, err := n.LoadUser(ownerID)
		if err != nil {
			return false, errors.WithMessage(err, "failed to load user")
		}
		if user.GetID() == ownerUser.GetID() {
			userPermissions = append(userPermissions, UserOwner)
		}

		permission := n.CompareReadPermissions(attributeTypePermissions, userPermissions)
		return permission, nil
	case UserUserAttribute:
		// user_owner == source_user
		// any, users, user_owner, target_user, user_owner+target_user
		userUserAttributeID := entry.NewUserUserAttributeID(attributeID, userID, ownerID)
		userUserAttribute, err := n.db.GetUserUserAttributesDB().GetUserUserAttributeByID(n.ctx, userUserAttributeID)
		if err != nil {
			return false, errors.WithMessage(err, "failed to get user user attribute")
		}
		if user.GetID() == userUserAttribute.SourceUserID {
			userPermissions = append(userPermissions, UserOwner)
		}
		if user.GetID() == userUserAttribute.TargetUserID {
			userPermissions = append(userPermissions, TargetUser)
		}

		permission := n.CompareReadPermissions(attributeTypePermissions, userPermissions)
		return permission, nil
	case NodeAttribute:
		permission := n.CompareReadPermissions(attributeTypePermissions, userPermissions)
		return permission, nil
	}

	return false, nil
}

func (n *Node) AssessWriteOperation(
	userID umid.UMID, ownerID umid.UMID, permissions string,
	attributeKind AttributeKind,
	attributeID entry.AttributeID,
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
		case UserOwner:
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
		case UserOwner:
			if slices.Contains(userPermissions, Owner) {
				return true
			}
			return false
		}
	}

	return false
}
