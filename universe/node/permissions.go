package node

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type OperationType string
type AttributeKind uint
type Permission uint

const (
	ReadOperation  OperationType = "read"
	WriteOperation OperationType = "write"
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
	attributeID := entry.NewAttributeID(pluginID, attributeName)
	attributeTypeID := entry.NewAttributeTypeID(pluginID, attributeName)
	attributeType, ok := n.GetAttributeTypes().GetAttributeType(attributeTypeID)
	if !ok {
		return false, errors.New("failed to get attributeType")
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		return false, errors.WithMessage(err, "failed to get userID from context")
	}

	permissions, err := n.GetPermissions(attributeID, attributeType, attributeKind, ownerID, userID)
	if err != nil {
		return false, errors.WithMessage(err, "failed to get permissions")
	}

	return n.AssessOperations(userID, ownerID, permissions, attributeKind, attributeID, operationType)
}

func (n *Node) GetPermissions(attributeID entry.AttributeID,
	attributeType universe.AttributeType,
	attributeKind AttributeKind, ownerID umid.UMID, userID umid.UMID,
) (map[string]string, error) {
	var attributeOptions *entry.AttributeOptions
	permissionsMap := make(map[string]string)
	defaultPermissions := map[string]string{
		"read":  "any",
		"write": "admin+user_owner",
	}

	switch attributeKind {
	case ObjectAttribute:
		options, ok := n.GetObjectAttributes().GetOptions(attributeID)
		if !ok {
			return nil, errors.New("failed to get objectAttribute options")
		}
		attributeOptions = options
	case ObjectUserAttribute:
		objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, ownerID, userID)
		options, ok := n.GetObjectUserAttributes().GetOptions(objectUserAttributeID)
		if !ok {
			return nil, errors.New("failed to get objectUserAttribute options")
		}
		attributeOptions = options
	case UserAttribute:
		userAttributeID := entry.NewUserAttributeID(attributeID, userID)
		options, ok := n.GetUserAttributes().GetOptions(userAttributeID)
		if !ok {
			return nil, errors.New("failed to get userAttribute options")
		}
		attributeOptions = options
	case UserUserAttribute:
		userUserAttributeID := entry.NewUserUserAttributeID(attributeID, userID, ownerID)
		options, ok := n.GetUserUserAttributes().GetOptions(userUserAttributeID)
		if !ok {
			return nil, errors.New("failed to get userUserAttribute options")
		}
		attributeOptions = options
	default:
		options := attributeType.GetOptions()
		attributeOptions = options
	}

	permissions := utils.GetFromAnyMap(*attributeOptions, "permissions", map[string]string(nil))
	if permissions != nil {
		permissionsMap = permissions
	} else {
		permissionsMap = defaultPermissions
	}

	return permissionsMap, nil
}

func (n *Node) GetUserPermissions(userID umid.UMID, permissions string) (universe.User, []string, []string, error) {
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
		return nil, nil, nil, errors.WithMessage(err, "failed to load user")
	}

	userType := user.GetUserType()
	guestUserTypeID, err := common.GetGuestUserTypeID()
	if err != nil {
		return nil, nil, nil, errors.WithMessage(err, "failed to get guest user type id")
	}

	if userType.GetID() != guestUserTypeID {
		userPermissions = append(userPermissions, User)
	}

	return user, userPermissions, attributeTypePermissions, nil
}

func (n *Node) AssessOperations(
	userID umid.UMID,
	ownerID umid.UMID, permissions map[string]string,
	attributeKind AttributeKind, attributeID entry.AttributeID, operationType OperationType,
) (bool, error) {
	user, userPermissions, attributeTypePermissions, err := n.GetUserPermissions(userID, permissions[string(operationType)])
	if err != nil {
		return false, errors.WithMessage(err, "failed to get user permissions")
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
	case UserAttribute:
		// any, users, user_owner
		userAttributeID := entry.NewUserAttributeID(attributeID, userID)
		userAttribute, err := n.db.GetUserAttributesDB().GetUserAttributeByID(n.ctx, userAttributeID)
		if err != nil {
			return false, errors.WithMessage(err, "failed to get user attribute")
		}
		if user.GetID() == userAttribute.UserID {
			userPermissions = append(userPermissions, UserOwner)
		}
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
	}

	permission := n.CompareReadPermissions(attributeTypePermissions, userPermissions)
	return permission, nil
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

func (n *Node) CompareWritePermissions(attributeTypePermissions []string, userPermissions []string) bool {
	for _, attributeTypePermission := range attributeTypePermissions {
		switch attributeTypePermission {
		case Owner:
		case UserOwner:
			if slices.Contains(userPermissions, Owner) {
				return true
			}
			return false
		case Admin:
			if slices.Contains(userPermissions, Admin) ||
				slices.Contains(userPermissions, Owner) {
				return true
			}
			return false
		case User:
			if slices.Contains(userPermissions, User) ||
				slices.Contains(userPermissions, Admin) ||
				slices.Contains(userPermissions, Owner) {
				return true
			}
			return false
		case Any:
			return true
		}
	}

	return false
}
