package entry

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID     *uuid.UUID   `db:"user_id"`
	UserTypeID *uuid.UUID   `db:"user_type_id"`
	Profile    *UserProfile `db:"profile"`
	Options    *UserOptions `db:"options"`
	CreatedAt  *time.Time   `db:"created_at"`
	UpdatedAt  *time.Time   `db:"updated_at"`
}

type UserOptions struct {
	IsGuest *bool `db:"is_guest" json:"is_guest"`
}

type UserProfile struct {
	Name *string `db:"name" json:"name"`
}

type UserAttribute struct {
	PluginID uuid.UUID         `db:"plugin_id"`
	UserID   uuid.UUID         `db:"user_id"`
	Name     string            `db:"attribute_name"`
	Value    *AttributeValue   `db:"value"`
	Options  *AttributeOptions `db:"options"`
}

type UserUserAttribute struct {
	PluginID     uuid.UUID         `db:"plugin_id"`
	SourceUserID uuid.UUID         `db:"source_user_id"`
	TargetUserID uuid.UUID         `db:"target_user_id"`
	Name         string            `db:"attribute_name"`
	Value        *AttributeValue   `db:"value"`
	Options      *AttributeOptions `db:"options"`
}
