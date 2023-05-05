package seed

import (
	"context"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedAssets2d(ctx context.Context, node universe.Node) error {
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.New("failed to get config from context")
	}
	baseUrl := cfg.Settings.FrontendURL
	miroUrl, gdriveUrl, videoUrl, err := generatePluginUrls(baseUrl)
	if err != nil {
		return errors.Wrap(err, "Error generating plugin URLs")
	}

	items := []*entry.Asset2d{
		{
			Asset2dID: umid.MustParse(miroPluginAsset2dID),
			Options: &entry.Asset2dOptions{
				"exact":    true,
				"subPath":  "miro",
				"iconName": "miro",
			},
			Meta: entry.Asset2dMeta{
				"name":      "miro",
				"pluginId":  miroPluginID,
				"scopeName": "plugin_miro",
				"scriptUrl": miroUrl,
			},
		},
		{
			Asset2dID: umid.MustParse(googleDrivePluginAsset2dID),
			Options: &entry.Asset2dOptions{
				"exact":    true,
				"iconName": "drive",
			},
			Meta: entry.Asset2dMeta{
				"name":      "google Drive",
				"pluginId":  googleDrivePluginID,
				"scopeName": "plugin_google_drive",
				"scriptUrl": gdriveUrl,
			},
		},
		{
			Asset2dID: umid.MustParse("7be0964f-df73-4880-91f5-22eef9967999"),
			Options:   &entry.Asset2dOptions{},
			Meta: entry.Asset2dMeta{
				"name":     "image",
				"pluginId": imagePluginID,
			},
		},
		{
			Asset2dID: umid.MustParse("be0d0ca3-c50b-401a-89d9-0e59fc45c5c2"),
			Options:   &entry.Asset2dOptions{},
			Meta: entry.Asset2dMeta{
				"name":     "text",
				"pluginId": textPluginID,
			},
		},
		{
			Asset2dID: umid.MustParse(videoPluginAsset2dID),
			Options:   &entry.Asset2dOptions{},
			Meta: entry.Asset2dMeta{
				"name":      "video",
				"pluginId":  videoPluginID,
				"scopeName": "plugin_video",
				"scriptUrl": videoUrl,
			},
		},
	}

	for _, item := range items {
		asset, err := node.GetAssets2d().CreateAsset2d(item.Asset2dID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create asset_2d: %s", item.Asset2dID)
		}

		_, err = asset.SetOptions(modify.MergeWith(item.Options), false)
		if err != nil {
			return errors.WithMessagef(err, "failed to set asset_2d options: %s", item.Asset2dID)
		}

		if err = asset.SetMeta(item.Meta, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset_2d meta: %s", item.Asset2dID)
		}

	}

	return nil
}
