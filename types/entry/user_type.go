package entry

import "github.com/google/uuid"

type UserTypeName string

type UserType struct {
	UserTypeID   *uuid.UUID    `db:"user_type_id"`
	UserTypeName *UserTypeName `db:"user_type_name"`
	Description  *string       `db:"description"`
	Options      any           `db:"options"`
}

const (
	USER           UserTypeName = "User"
	DEITY          UserTypeName = "Deity"
	TEMPORARY_USER UserTypeName = "Temporary User"
	TOKEN_GROUPS   UserTypeName = "Token Groups"
)
