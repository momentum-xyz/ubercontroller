package seed

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func Node(ctx context.Context, node universe.Node) error {
	group, _ := errgroup.WithContext(ctx)

	group.Go(func() error {
		return seedPlugins(node)
	})

	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to seed plugins")
	}

	return node.Save()
}

func seedPlugins(node universe.Node) error {
	type pluginItem struct {
		ID   uuid.UUID
		Meta *entry.PluginMeta
	}

	data := []*pluginItem{
		{
			ID: universe.GetSystemPluginID(),
			Meta: &entry.PluginMeta{
				"name": "Core",
			},
		},
	}

	for _, p := range data {
		plugin, err := node.GetPlugins().CreatePlugin(p.ID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create plugin: %s", p.ID)
		}
		if err := plugin.SetMeta(p.Meta, false); err != nil {
			return errors.WithMessagef(err, "failed to set meta: %s", p.Meta)
		}
	}

	return nil
}
