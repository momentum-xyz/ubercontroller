package user

import (
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

type UserUserAttributeIndex struct {
	entry.AttributeID

	SourceUserId uuid.UUID
	TargetUserId uuid.UUID
}

type UserAttributeIndex struct {
	entry.AttributeID
}
