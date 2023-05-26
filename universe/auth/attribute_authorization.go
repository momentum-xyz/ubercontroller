package auth

// Authorization logic for plugin attributes.

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type operationType string

const (
	ReadOperation  operationType = "read"
	WriteOperation operationType = "write"
)

func defaultPermissions() *entry.PermissionsAttributeOption {
	return &entry.PermissionsAttributeOption{
		Read:  string(entry.PermissionAny),
		Write: string(entry.PermissionAdmin) + "+" + string(entry.PermissionUserOwner), //TODO:again, resolve on input so this is not a string here
	}
}

// Interface plugin attributes need to implement for authorization.
type AttributePermissionsAuthorizer[T comparable] interface {
	universe.AttributeUserRoleGetter[T]
	universe.AttributeOptionsGetter[T]
}

// Check authorization for an operation on an plugin attribute.
func CheckAttributePermissions[ID comparable](
	ctx context.Context,
	attrType entry.AttributeType,
	attrStore AttributePermissionsAuthorizer[ID],
	targetID ID,
	userID umid.UMID, // The user executing
	opType operationType,

) (bool, error) {
	permissions, err := getPermissions[ID](ctx, attrType, attrStore, targetID)
	if err != nil {
		return false, fmt.Errorf("get attribute permissions: %w", err)
	}
	if permissions == nil {
		permissions = defaultPermissions()
	}
	allowedRoles := getOperationRoles(*permissions, opType)
	// short-circuit the 'any' permission handling
	if slices.Contains(allowedRoles, entry.PermissionAny) {
		return true, nil
	}

	roles, err := attrStore.GetUserRoles(ctx, attrType, targetID, userID)
	if err != nil {
		return false, fmt.Errorf("get user roles: %w", err)
	}
	return hasRole(roles, allowedRoles), nil
}

// Check if user is authorized to read all attributes.
func CheckReadAllPermissions[ID comparable](
	ctx context.Context,
	attrType entry.AttributeType,
	attrStore AttributePermissionsAuthorizer[ID],
	userID umid.UMID, // The user executing
) (bool, error) {
	// This is an edge case, that is a big TODO to do properly.
	// Would need to check each individual attribute, since options can be overridden.
	// And return a filtered list??
	// For now, to keep stuff working: you need read: "user" permission on the attribute type
	// can not be overriden.
	var permissions *entry.PermissionsAttributeOption
	options := attrType.Options
	if options != nil {
		var err error
		permissions, err = permissionsFromOptions(*options)
		if err != nil {
			return false, err
		}
	}
	if permissions == nil {
		permissions = defaultPermissions()
	}
	allowedRoles := getOperationRoles(*permissions, ReadOperation)
	if slices.Contains(allowedRoles, entry.PermissionAny) {
		return true, nil
	}
	// assume the use always has the 'user' role (no other user types implemented yet)
	return hasRole([]entry.PermissionsRoleType{entry.PermissionUser}, allowedRoles), nil
}

func getPermissions[ID comparable](
	ctx context.Context,
	attrType entry.AttributeType,
	attrStore universe.AttributeOptionsGetter[ID],
	targetID ID,
) (*entry.PermissionsAttributeOption, error) {
	options, ok := attrStore.GetEffectiveOptions(targetID)
	if !ok {
		options = attrType.Options
	}
	if options != nil {
		return permissionsFromOptions(*options)
	}
	return nil, nil
}

func getOperationRoles(
	option entry.PermissionsAttributeOption,
	opType operationType) []entry.PermissionsRoleType {
	// TODO: move decoding this to the input (where we retrieve it from database)
	// not here, somewhere in the middle
	var perms string
	if opType == ReadOperation {
		perms = option.Read
	} else {
		perms = option.Write
	}
	var roles []entry.PermissionsRoleType
	for _, v := range strings.Split(perms, "+") {
		roles = append(roles, entry.PermissionsRoleType(v))
	}
	return roles
}

func hasRole(
	userRoles []entry.PermissionsRoleType,
	allowedRoles []entry.PermissionsRoleType) bool {
	for _, r := range allowedRoles {
		if slices.Contains(userRoles, r) {
			return true
		}
	}
	return false
}

// Get permissions from attribute options
func permissionsFromOptions(options entry.AttributeOptions) (*entry.PermissionsAttributeOption, error) {
	// TODO: move decoding this to the input (where we retrieve it from database)
	// so we don't need this function.
	permissions, ok := options["permissions"]
	if ok && permissions != nil {
		result := &entry.PermissionsAttributeOption{}
		err := utils.MapDecode(permissions, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, nil //todo?

}
