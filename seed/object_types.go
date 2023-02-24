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
			id:             uuid.MustParse(NodeObjectTypeID),
			asset2dID:      nil,
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
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "anchor",
			categoryName:   "Anchors",
			description:    utils.GetPTR(""),
			options: &entry.ObjectOptions{
				Private: utils.GetPTR(false),
				Visible: utils.GetPTR(entry.ReactUnityObjectVisibleType),
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
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "World",
			categoryName:   "Worlds",
			description:    utils.GetPTR("World Type"),
			options: &entry.ObjectOptions{
				Private:           utils.GetPTR(false),
				Visible:           utils.GetPTR(entry.ReactUnityObjectVisibleType),
				DefaultTiles:      []any{},
				FrameTemplates:    map[string]any{},
				AllowedSubObjects: []uuid.UUID{},
			},
		},

		{
			id:             uuid.MustParse("b59abd4d-f54d-4a97-8b6d-16a2037ddd8f"),
			asset2dID:      utils.GetPTR(uuid.MustParse(dockingStationAsset2dID)),
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
				Visible: utils.GetPTR(entry.ReactUnityObjectVisibleType),
			},
		},

		{
			id:             uuid.MustParse("4fe4ed05-9024-461a-97d6-22666e8a4f46"),
			asset2dID:      nil,
			asset3dID:      utils.GetPTR(uuid.MustParse("6846dba3-38b1-4540-a80d-4ba04af4111e")),
			objectTypeName: "Effects emitter",
			categoryName:   "Effects",
			description:    utils.GetPTR("Effects emitter"),
			options: &entry.ObjectOptions{
				Visible: utils.GetPTR(entry.UnityObjectVisibleType),
			},
		},

		{
			id:             uuid.MustParse("69d8ae40-df9b-4fc8-af95-32b736d2bbcd"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "Service Space",
			categoryName:   "Service Spaces",
			description:    utils.GetPTR(""),
			options: &entry.ObjectOptions{
				Private:           utils.GetPTR(false),
				Visible:           utils.GetPTR(entry.InvisibleObjectVisibleType),
				DefaultTiles:      []any{},
				FrameTemplates:    map[string]any{},
				AllowedSubObjects: []uuid.UUID{},
			},
		},

		{
			id:             uuid.MustParse("d7a41cbd-5cfe-454b-b522-76f22fa55026"),
			asset2dID:      nil,
			asset3dID:      utils.GetPTR(uuid.MustParse("313a597a-8b9a-47a7-9908-52bdc7a21a3e")),
			objectTypeName: "Skybox",
			categoryName:   "Skybox",
			description:    utils.GetPTR("Skybox"),
			options: &entry.ObjectOptions{
				Editable: utils.GetPTR(false),
				Visible:  utils.GetPTR(entry.ReactUnityObjectVisibleType),
			},
		},

		{
			id:             uuid.MustParse("75b56447-c4f1-4020-b8fc-d68704a11d65"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "Generic Space",
			categoryName:   "Generic Spaces",
			description:    utils.GetPTR(""),
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
			id:             uuid.MustParse("f9607e55-63e8-4cb1-ae47-66395199975d"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "morgue",
			categoryName:   "Morgues",
			description:    utils.GetPTR("morgue"),
			options: &entry.ObjectOptions{
				Subs: map[string]any{
					"asset2d_plugins": []any{
						"24071066-e8c6-4692-95b5-ae2dc3ed075c",
					},
				},
				Private:           utils.GetPTR(false),
				Visible:           utils.GetPTR(entry.InvisibleObjectVisibleType),
				DefaultTiles:      []any{},
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
