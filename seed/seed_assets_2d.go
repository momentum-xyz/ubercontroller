package seed

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedAssets2d(node universe.Node) error {
	type item struct {
		asset2DID uuid.UUID
		options   *entry.Asset2dOptions
		meta      *entry.Asset2dMeta
	}

	items := []*item{
		{
			asset2DID: uuid.MustParse("a31722a6-26b7-46bc-97f9-435c380c3ca9"),
			options: &entry.Asset2dOptions{
				"exact":    true,
				"subPath":  "miro",
				"iconName": "miro",
			},
			meta: &entry.Asset2dMeta{
				"name":      "miro",
				"pluginId":  "24071066-e8c6-4692-95b5-ae2dc3ed075c",
				"scopeName": "plugin_miro",
				"scriptUrl": "http://localhost/plugins/miro/remoteEntry.js",
			},
		},
		{
			asset2DID: uuid.MustParse("c601404b-61a2-47d5-a5c7-f3c704a8bf58"),
			options: &entry.Asset2dOptions{
				"exact":    true,
				"iconName": "drive",
			},
			meta: &entry.Asset2dMeta{
				"name":      "google Drive",
				"pluginId":  "c3f89640-e0f0-4536-ae0d-8fc8a75ec0cd",
				"scopeName": "plugin_google_drive",
				"scriptUrl": "http://localhost/plugins/google-drive/remoteEntry.js",
			},
		},
		{
			asset2DID: uuid.MustParse("0d99e5aa-a627-4353-8bfa-1c0e7053db90"),
			options: &entry.Asset2dOptions{
				"exact":    true,
				"iconName": "drive",
			},
			meta: &entry.Asset2dMeta{
				"name":      "google drive Local",
				"pluginId":  "c3f89640-e0f0-4536-ae0d-8fc8a75ec0cd",
				"scopeName": "plugin_google_drive",
				"scriptUrl": "http://localhost:3002/remoteEntry.js",
			},
		},
		{
			asset2DID: uuid.MustParse("7be0964f-df73-4880-91f5-22eef9967999"),
			options:   &entry.Asset2dOptions{},
			meta: &entry.Asset2dMeta{
				"name":     "image",
				"pluginId": "ff40fbf0-8c22-437d-b27a-0258f99130fe",
			},
		},
		{
			asset2DID: uuid.MustParse("be0d0ca3-c50b-401a-89d9-0e59fc45c5c2"),
			options:   &entry.Asset2dOptions{},
			meta: &entry.Asset2dMeta{
				"name":     "text",
				"pluginId": "fc9f2eb7-590a-4a1a-ac75-cd3bfeef28b2",
			},
		},
		{
			asset2DID: uuid.MustParse("bda25d5d-2aab-45b4-9e8a-23579514cec1"),
			options:   &entry.Asset2dOptions{},
			meta: &entry.Asset2dMeta{
				"name":     "video",
				"pluginId": "308fdacc-8c2d-40dc-bd5f-d1549e3e03ba",
			},
		},
		{
			asset2DID: uuid.MustParse("2a879830-b79e-4c35-accc-05607c51d504"),
			options: &entry.Asset2dOptions{
				"exact":    true,
				"subPath":  "miro",
				"iconName": "miro",
			},
			meta: &entry.Asset2dMeta{
				"name":      "miro local",
				"pluginId":  "24071066-e8c6-4692-95b5-ae2dc3ed075c",
				"scopeName": "plugin_miro",
				"scriptUrl": "http://localhost:3001/remoteEntry.js",
			},
		},
		{
			asset2DID: uuid.MustParse("140c0f2e-2056-443f-b5a7-4a3c2e6b05da"),
			options:   &entry.Asset2dOptions{},
			meta: &entry.Asset2dMeta{
				"name": "Docking station",
			},
		},
	}

	for _, item := range items {
		asset, err := node.GetAssets2d().CreateAsset2d(item.asset2DID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create asset_2d: %s", item.asset2DID)
		}

		_, err = asset.SetOptions(modify.MergeWith(item.options), false)
		if err != nil {
			return errors.WithMessagef(err, "failed to set asset_2d options: %s", item.asset2DID)
		}

		if err = asset.SetMeta(item.meta, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset_2d meta: %s", item.asset2DID)
		}

	}

	return nil
}
