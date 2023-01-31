package seed

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedObjectTypes(node universe.Node) error {
	type item struct {
		id             uuid.UUID
		asset2dID      uuid.UUID
		asset3dID      *uuid.UUID
		objectTypeName string
		categoryName   string
		description    *string
		options        *entry.ObjectOptions
	}

	items := []*item{
		{
			id:             uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			asset2dID:      uuid.MustParse("00000000-0000-0000-0000-000000000004"),
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
			asset2dID:      uuid.MustParse("a31722a6-26b7-46bc-97f9-435c380c3ca9"),
			asset3dID:      utils.GetPTR(uuid.MustParse("e5609914-e673-7f48-a666-574ed2baff92")),
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
	}

	for _, item := range items {
		objectType, err := node.GetObjectTypes().CreateObjectType(item.id)
		if err != nil {
			return errors.WithMessagef(err, "failed to create object type: %s", item.id)
		}

		asset2d, ok := node.GetAssets2d().GetAsset2d(item.asset2dID)
		if !ok {
			return errors.WithMessagef(err, "failed to create object type: asset_2d not found: %s", item.asset2dID)
		}
		if err := objectType.SetAsset2d(asset2d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset_2d: %s for object_type: %s", item.asset2dID, item.id)
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
