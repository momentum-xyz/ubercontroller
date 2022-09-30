package entry

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserID     uuid.UUID `db:"user_id"`
	UserTypeID uuid.UUID `db:"user_type_id"`
	Profile    any       `db:"profile"`
	Options    any       `db:"options"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
