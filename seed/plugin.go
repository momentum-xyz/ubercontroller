package seed

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

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
			ID: uuid.MustParse("2b92edbc-5ef5-4028-89a6-d510f8583887"),
			Meta: &entry.PluginMeta{
				"name":        "Event Calendar",
				"description": "Event calendar plugin",
			},
		},
		{
			ID: uuid.MustParse(videoPluginID),
			Meta: &entry.PluginMeta{
				"name": "Video",
			},
		},
		{
			ID: uuid.MustParse(imagePluginID),
			Meta: &entry.PluginMeta{
				"name": "Image",
			},
		},
		{
			ID: uuid.MustParse(textPluginID),
			Meta: &entry.PluginMeta{
				"name": "Text",
			},
		},
		{
			ID: uuid.MustParse(miroPluginID),
			Meta: &entry.PluginMeta{
				"name": "Miro",
				"assets2d": []string{
					miroPluginAsset2dID,
					"2a879830-b79e-4c35-accc-05607c51d504",
				},
				"scopeName": "plugin_miro",
				"scriptUrl": "http://localhost/plugins/miro/remoteEntry.js",
			},
		},
		{
			ID: uuid.MustParse(googleDrivePluginID),
			Meta: &entry.PluginMeta{
				"name": "Google Drive",
				"assets2d": []string{
					googleDrivePluginAsset2dID,
					noname4Asset2dID,
				},
				"scopeName": "plugin_google_drive",
				"scriptUrl": "http://localhost/plugins/google-drive/remoteEntry.js",
			},
		},
		{
			ID:   uuid.MustParse(noname1PluginID),
			Meta: &entry.PluginMeta{},
		},
		{
			ID:   universe.GetKusamaPluginID(),
			Meta: &entry.PluginMeta{},
		},
		{
			ID:   uuid.MustParse(high5PluginID),
			Meta: &entry.PluginMeta{},
		},
	}

	for _, p := range data {
		plugin, err := node.GetPlugins().CreatePlugin(p.ID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create plugin: %s", p.ID)
		}
		if err := plugin.SetMeta(*p.Meta, false); err != nil {
			return errors.WithMessagef(err, "failed to set meta: %s", p.Meta)
		}
	}

	return nil
}
