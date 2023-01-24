package entry

import (
	"time"

	"github.com/google/uuid"
)

type Plugin struct {
	PluginID  uuid.UUID      `db:"plugin_id" json:"plugin_id"`
	Meta      *PluginMeta    `db:"meta" json:"meta"`
	Options   *PluginOptions `db:"options" json:"options"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time     `db:"updated_at" json:"updated_at"`
}

type PluginMeta map[string]any

type PluginOptions map[string]any
