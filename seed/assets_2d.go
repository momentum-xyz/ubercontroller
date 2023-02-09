package seed

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedAssets2d(node universe.Node) error {
	items := []*entry.Asset2d{
		{
			Asset2dID: uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Options:   &entry.Asset2dOptions{},
			Meta: &entry.Asset2dMeta{
				"name": "",
			},
		},
		{
			Asset2dID: uuid.MustParse("00000000-0000-0000-0000-000000000008"),
			Options:   &entry.Asset2dOptions{},
			Meta: &entry.Asset2dMeta{
				"name": "",
			},
		},
		{
			Asset2dID: uuid.MustParse("a31722a6-26b7-46bc-97f9-435c380c3ca9"),
			Options: &entry.Asset2dOptions{
				"exact":    true,
				"subPath":  "miro",
				"iconName": "miro",
			},
			Meta: &entry.Asset2dMeta{
				"name":      "miro",
				"pluginId":  miroPluginID,
				"scopeName": "plugin_miro",
				"scriptUrl": "http://localhost/plugins/miro/remoteEntry.js",
			},
		},
		{
			Asset2dID: uuid.MustParse("c601404b-61a2-47d5-a5c7-f3c704a8bf58"),
			Options: &entry.Asset2dOptions{
				"exact":    true,
				"iconName": "drive",
			},
			Meta: &entry.Asset2dMeta{
				"name":      "google Drive",
				"pluginId":  googleDrivePluginID,
				"scopeName": "plugin_google_drive",
				"scriptUrl": "http://localhost/plugins/google-drive/remoteEntry.js",
			},
		},
		{
			Asset2dID: uuid.MustParse("0d99e5aa-a627-4353-8bfa-1c0e7053db90"),
			Options: &entry.Asset2dOptions{
				"exact":    true,
				"iconName": "drive",
			},
			Meta: &entry.Asset2dMeta{
				"name":      "google drive Local",
				"pluginId":  googleDrivePluginID,
				"scopeName": "plugin_google_drive",
				"scriptUrl": "http://localhost:3002/remoteEntry.js",
			},
		},
		{
			Asset2dID: uuid.MustParse("7be0964f-df73-4880-91f5-22eef9967999"),
			Options:   &entry.Asset2dOptions{},
			Meta: &entry.Asset2dMeta{
				"name":     "image",
				"pluginId": imagePluginID,
			},
		},
		{
			Asset2dID: uuid.MustParse("be0d0ca3-c50b-401a-89d9-0e59fc45c5c2"),
			Options:   &entry.Asset2dOptions{},
			Meta: &entry.Asset2dMeta{
				"name":     "text",
				"pluginId": textPluginID,
			},
		},
		{
			Asset2dID: uuid.MustParse("bda25d5d-2aab-45b4-9e8a-23579514cec1"),
			Options:   &entry.Asset2dOptions{},
			Meta: &entry.Asset2dMeta{
				"name":     "video",
				"pluginId": videoPluginID,
			},
		},
		{
			Asset2dID: uuid.MustParse("2a879830-b79e-4c35-accc-05607c51d504"),
			Options: &entry.Asset2dOptions{
				"exact":    true,
				"subPath":  "miro",
				"iconName": "miro",
			},
			Meta: &entry.Asset2dMeta{
				"name":      "miro local",
				"pluginId":  miroPluginID,
				"scopeName": "plugin_miro",
				"scriptUrl": "http://localhost:3001/remoteEntry.js",
			},
		},
		{
			Asset2dID: uuid.MustParse("140c0f2e-2056-443f-b5a7-4a3c2e6b05da"),
			Options:   &entry.Asset2dOptions{},
			Meta: &entry.Asset2dMeta{
				"name": "Docking station",
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
