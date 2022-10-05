package space

import (
	"github.com/google/uuid"
)

type AttributeIndex struct {
	PluginId      uuid.UUID
	AttributeName string
}

type UserAttributeIndex struct {
	PluginId      uuid.UUID
	UserId        uuid.UUID
	AttributeName string
}
