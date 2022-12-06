package entry

import (
	"time"

	"github.com/google/uuid"
)

type UserSpace struct {
	SpaceID   uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Value     UserSpaceValue
}

type UserSpaceValue map[string]any
