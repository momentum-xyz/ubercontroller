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
		{
			ID: uuid.MustParse("fd94324b-b2bb-4d79-8328-7f33362385b2"),
			Meta: &entry.PluginMeta{
				"name":      "Template",
				"scopeName": "momentum_plugin_template",
				"scriptUrl": "http://localhost:3002/remoteEntry.js",
			},
		},
		{
			ID: uuid.MustParse("2b92edbc-5ef5-4028-89a6-d510f8583887"),
			Meta: &entry.PluginMeta{
				"name":        "Event Calendar",
				"description": "Event calendar plugin",
			},
		},
		{
			ID: uuid.MustParse("308fdacc-8c2d-40dc-bd5f-d1549e3e03ba"),
			Meta: &entry.PluginMeta{
				"name": "Video",
			},
		},
		{
			ID: uuid.MustParse("ff40fbf0-8c22-437d-b27a-0258f99130fe"),
			Meta: &entry.PluginMeta{
				"name": "Image",
			},
		},
		{
			ID: uuid.MustParse("fc9f2eb7-590a-4a1a-ac75-cd3bfeef28b2"),
			Meta: &entry.PluginMeta{
				"name": "Text",
			},
		},
		{
			ID: uuid.MustParse("24071066-e8c6-4692-95b5-ae2dc3ed075c"),
			Meta: &entry.PluginMeta{
				"name": "Miro",
				"assets2d": []string{
					"a31722a6-26b7-46bc-97f9-435c380c3ca9",
					"2a879830-b79e-4c35-accc-05607c51d504",
				},
				"scopeName": "plugin_miro",
				"scriptUrl": "http://localhost/plugins/miro/remoteEntry.js",
			},
		},
		{
			ID: uuid.MustParse("c3f89640-e0f0-4536-ae0d-8fc8a75ec0cd"),
			Meta: &entry.PluginMeta{
				"name": "Google Drive",
				"assets2d": []string{
					"c601404b-61a2-47d5-a5c7-f3c704a8bf58",
					"0d99e5aa-a627-4353-8bfa-1c0e7053db90",
				},
				"scopeName": "plugin_google_drive",
				"scriptUrl": "http://localhost/plugins/google-drive/remoteEntry.js",
			},
		},
		{
			ID: uuid.MustParse("3159dfc4-1ba2-4a76-bd81-dc08846d8557"),
			Meta: &entry.PluginMeta{
				"name": "Twitch",
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
