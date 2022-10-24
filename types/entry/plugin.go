package entry

import (
	"time"

	"github.com/google/uuid"
)

type Plugin struct {
	PluginID  uuid.UUID      `db:"plugin_id"`
	Meta      *PluginMeta    `db:"meta"`
	Options   *PluginOptions `db:"options"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt *time.Time     `db:"updated_at"`
}

type PluginMeta map[string]any

type PluginOptions struct {
	File string `db:"file" json:"file"`
}
