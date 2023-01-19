package entry

import (
	"time"

	"github.com/google/uuid"
)

type UserType struct {
	UserTypeID   uuid.UUID    `db:"user_type_id" json:"user_type_id"`
	UserTypeName string       `db:"user_type_name" json:"user_type_name"`
	Description  *string      `db:"description" json:"description"`
	Options      *UserOptions `db:"options" json:"options"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time   `db:"updated_at" json:"updated_at"`
}
