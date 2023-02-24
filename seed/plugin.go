package seed

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func seedPlugins(ctx context.Context, node universe.Node) error {
	type pluginItem struct {
		ID   uuid.UUID
		Meta *entry.PluginMeta
	}

	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.New("failed to get config from context")
	}
	baseUrl := cfg.Settings.FrontendURL
	miroUrl, gdriveUrl, videoUrl, err := generatePluginUrls(baseUrl)
	if err != nil {
		return errors.Wrap(err, "Error generating plugin URLs")
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
				"scriptUrl": miroUrl,
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
				"scriptUrl": gdriveUrl,
			},
		},
		{
			ID: uuid.MustParse(videoPluginID),
			Meta: &entry.PluginMeta{
				"name": "Video",
				"assets2d": []string{
					videoPluginAsset2dID,
				},
				"scopeName": "plugin_video",
				"scriptUrl": videoUrl,
			},
		},
		{
			ID: uuid.MustParse(OdysseyHackatonPluginID),
			Meta: &entry.PluginMeta{
				"name": "Odyssey hackaton",
			},
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

func generatePluginUrls(baseUrl string) (miroUrl string, gdriveUrl string, videoUrl string, err error) {
	miroUrl, err = generateScriptUrl(baseUrl, "miro")
	if err != nil {
		err = errors.Wrap(err, "Could not generate plugin URL")
		return
	}
	gdriveUrl, err = generateScriptUrl(baseUrl, "google-drive")
	if err != nil {
		err = errors.Wrap(err, "Could not generate plugin URL")
		return
	}
	videoUrl, err = generateScriptUrl(baseUrl, "video")
	if err != nil {
		err = errors.Wrap(err, "Could not generate plugin URL")
		return
	}
	return
}

func generateScriptUrl(baseUrl string, pluginName string) (result string, err error) {
	result, err = url.JoinPath(baseUrl, "plugins", pluginName, "remoteEntry.js")
	return
}
