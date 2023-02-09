package seed

import (
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/pkg/errors"
)

func seedObjectTypes(node universe.Node) error {
	type item struct {
		id             uuid.UUID
		asset2dID      *uuid.UUID
		asset3dID      *uuid.UUID
		objectTypeName string
		categoryName   string
		description    *string
		options        *entry.ObjectOptions
	}

	items := []*item{
		{
			id:             uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			asset2dID:      utils.GetPTR(uuid.MustParse("00000000-0000-0000-0000-000000000004")),
			asset3dID:      nil,
			objectTypeName: "Node",
			categoryName:   "Nodes",
			description:    utils.GetPTR("Root of it all"),
			options: &entry.ObjectOptions{
				Subs: map[string]any{
					"asset2d_plugins": []any{
						miroPluginID,
					},
				},
				Private:           utils.GetPTR(false),
				Visible:           utils.GetPTR(entry.ReactObjectVisibleType),
				DefaultTiles:      []any{},
				FrameTemplates:    map[string]any{},
				AllowedSubObjects: []uuid.UUID{},
			},
		},
		{
			id:             uuid.MustParse("88415343-90db-4d23-a9e7-79a11aaaaf04"),
			asset2dID:      utils.GetPTR(uuid.MustParse("a31722a6-26b7-46bc-97f9-435c380c3ca9")),
			asset3dID:      utils.GetPTR(uuid.MustParse(noname1Asset3dID)),
			objectTypeName: "anchor",
			categoryName:   "Anchors",
			description:    utils.GetPTR(""),
			options: &entry.ObjectOptions{
				Subs: map[string]any{
					"asset2d_plugins": []any{
						miroPluginID,
					},
				},
				Private: utils.GetPTR(false),
				Visible: utils.GetPTR(entry.ReactUnityObjectVisibleType),
				DefaultTiles: []any{
					map[string]any{
						"row":           0,
						"hash":          "53e9a2811a7a6cd93011a6df7c23edc7",
						"type":          "tile_type_media",
						"column":        0,
						"edited":        1,
						"render":        1,
						"content":       map[string]any{},
						"permanentType": "poster",
					},
					map[string]any{
						"row":           1,
						"hash":          "69e2b342788fe70273c15b62f618ef22",
						"type":          "tile_type_media",
						"column":        0,
						"edited":        1,
						"render":        1,
						"content":       map[string]any{},
						"permanentType": "meme",
					},
					map[string]any{
						"row":    1,
						"hash":   "9ae1db04e863bb9d1a572d8a6727c665",
						"type":   "tile_type_text",
						"column": 1,
						"edited": 1,
						"render": 1,
						"content": map[string]any{
							"text":  "Description goes here",
							"title": "Description",
						},
						"permanentType": "description",
					},
					map[string]any{
						"row":    0,
						"hash":   "fc4d116880460bf808b4487954823e80",
						"type":   "tile_type_video",
						"column": 2,
						"edited": 1,
						"render": 1,
						"content": map[string]any{
							"url": "https://www.youtube.com/watch?v=mwpj70Gcatg",
						},
						"permanentType": "video",
					},
					map[string]any{
						"row":           0,
						"hash":          "3b441ce7b41c54693fbc798ca896c88e",
						"type":          "tile_type_media",
						"column":        2,
						"edited":        1,
						"render":        1,
						"content":       map[string]any{},
						"permanentType": "third",
					},
				},
				ChildPlacements: map[uuid.UUID]*entry.ObjectChildPlacement{
					uuid.MustParse("00000000-0000-0000-0000-000000000000"): &entry.ObjectChildPlacement{
						Algo: utils.GetPTR("circular"),
						Options: map[string]any{
							"R":     55,
							"angle": 0,
						},
					},
				},
				FrameTemplates:    map[string]any{},
				AllowedSubObjects: []uuid.UUID{},
			},
		},
		{
			id:             uuid.MustParse("27456794-f9aa-44b5-90a5-e307fc21bc3d"),
			asset2dID:      nil,
			asset3dID:      utils.GetPTR(uuid.MustParse(dockingStationAsset3dID)),
			objectTypeName: "Docking station",
			categoryName:   "Docking station",
			description:    utils.GetPTR("Odyssey docking hub"),
			options: &entry.ObjectOptions{
				Visible: utils.GetPTR(entry.ReactUnityObjectVisibleType),
				ChildPlacements: map[uuid.UUID]*entry.ObjectChildPlacement{
					uuid.MustParse("00000000-0000-0000-0000-000000000000"): &entry.ObjectChildPlacement{
						Algo: utils.GetPTR("circular"),
						Options: map[string]any{
							"R":     42,
							"angle": 0,
						},
					},
				},
			},
		},
		//
		{
			id:             uuid.MustParse("a41ee21e-6c56-41b3-81a9-1c86578b6b3c"),
			asset2dID:      utils.GetPTR(uuid.MustParse("00000000-0000-0000-0000-000000000008")),
			asset3dID:      utils.GetPTR(uuid.MustParse("b2ef3600-9595-2743-ac9d-0a86c1a327a2")),
			objectTypeName: "World",
			categoryName:   "Worlds",
			description:    utils.GetPTR("World Type"),
			options: &entry.ObjectOptions{
				Subs: map[string]any{
					"asset2d_plugins": []any{
						"24071066-e8c6-4692-95b5-ae2dc3ed075c",
					},
				},
				Private:           utils.GetPTR(false),
				Visible:           utils.GetPTR(entry.ReactUnityObjectVisibleType),
				DefaultTiles:      []any{},
				FrameTemplates:    map[string]any{},
				AllowedSubObjects: []uuid.UUID{},
			},
		},

		{
			id:             uuid.MustParse("b59abd4d-f54d-4a97-8b6d-16a2037ddd8f"),
			asset2dID:      utils.GetPTR(uuid.MustParse("140c0f2e-2056-443f-b5a7-4a3c2e6b05da")),
			asset3dID:      utils.GetPTR(uuid.MustParse("a6862b31-8f80-497d-b9d6-8234e6a71773")),
			objectTypeName: "Docking bulb",
			categoryName:   "Docking bulbs",
			description:    utils.GetPTR("Odyssey docking bulb"),
			options: &entry.ObjectOptions{
				Visible:  utils.GetPTR(entry.ReactUnityObjectVisibleType),
				Editable: utils.GetPTR(false),
			},
		},

		{
			id:             uuid.MustParse("4ed3a5bb-53f8-4511-941b-07902982c31c"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "Custom objects",
			categoryName:   "Custom",
			description:    utils.GetPTR("Custom placed objects"),
			options: &entry.ObjectOptions{
				Visible: utils.GetPTR(entry.UnityObjectVisibleType), //TODO should be 0
			},
		},
	}

	for _, item := range items {
		objectType, err := node.GetObjectTypes().CreateObjectType(item.id)
		if err != nil {
			return errors.WithMessagef(err, "failed to create object type: %s", item.id)
		}

		if item.asset2dID != nil {
			asset2d, ok := node.GetAssets2d().GetAsset2d(*item.asset2dID)
			if !ok {
				return errors.WithMessagef(err, "failed to create object type: asset_2d not found: %s", item.asset2dID)
			}
			if err := objectType.SetAsset2d(asset2d, false); err != nil {
				return errors.WithMessagef(err, "failed to set asset_2d: %s for object_type: %s", item.asset2dID, item.id)
			}
		}

		if item.asset3dID != nil {
			asset3d, ok := node.GetAssets3d().GetAsset3d(*item.asset3dID)
			if !ok {
				return errors.WithMessagef(err, "failed to create object type: asset_3d not found: %s", item.asset3dID)
			}
			if err := objectType.SetAsset3d(asset3d, false); err != nil {
				return errors.WithMessagef(err, "failed to set asset_3d: %s for object_type: %s", item.asset3dID, item.id)
			}
		}

		if err := objectType.SetName(item.objectTypeName, false); err != nil {
			return errors.WithMessagef(err, "failed to set name: %s for object_type: %s", item.objectTypeName, item.id)
		}

		if err := objectType.SetCategoryName(item.categoryName, false); err != nil {
			return errors.WithMessagef(err, "failed to set category name: %s for object_type: %s", item.categoryName, item.id)
		}

		if err := objectType.SetDescription(item.description, false); err != nil {
			return errors.WithMessagef(err, "failed to set description for object_type: %s", item.id)
		}

		_, err = objectType.SetOptions(modify.MergeWith(item.options), false)
		if err != nil {
			return errors.WithMessagef(err, "failed to set options for object_type: %s", item.id)
		}
	}

	return nil
}
