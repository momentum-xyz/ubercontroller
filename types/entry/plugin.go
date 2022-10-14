package entry

import (
	"time"

	"github.com/google/uuid"
)

type Plugin struct {
	PluginID    uuid.UUID      `db:"plugin_id"`
	PluginName  string         `db:"plugin_name"`
	Description *string        `db:"description"`
	Options     *PluginOptions `db:"options"`
	CreatedAt   *time.Time     `db:"created_at"`
	UpdatedAt   *time.Time     `db:"updated_at"`
}

type PluginOptions struct {
	File string `db:"file"`
}
