package seed

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

func seedObjectTypes(node universe.Node) error {

	type item struct {
		id             umid.UMID
		asset2dID      *umid.UMID
		asset3dID      *umid.UMID
		objectTypeName string
		categoryName   string
		description    *string
		options        *entry.ObjectOptions
	}

	items := []*item{
		{
			id:             umid.MustParse(NodeObjectTypeID),
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
				Private:        utils.GetPTR(false),
				Visible:        utils.GetPTR(entry.UI2DObjectVisibleType),
				DefaultTiles:   []any{},
				FrameTemplates: map[string]any{},
				AllowedChildren: []umid.UMID{
					umid.MustParse("a41ee21e-6c56-41b3-81a9-1c86578b6b3c"),
				},
			},
		},
		{
			id:             umid.MustParse("88415343-90db-4d23-a9e7-79a11aaaaf04"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "anchor",
			categoryName:   "Anchors",
			description:    utils.GetPTR(""),
			options: &entry.ObjectOptions{
				Private: utils.GetPTR(false),
				Visible: utils.GetPTR(entry.AllObjectVisibleType),
				ChildPlacements: map[umid.UMID]*entry.ObjectChildPlacement{
					umid.MustParse("00000000-0000-0000-0000-000000000000"): &entry.ObjectChildPlacement{
						Algo: utils.GetPTR("circular"),
						Options: map[string]any{
							"R":     55,
							"angle": 0,
						},
					},
				},
				FrameTemplates:  map[string]any{},
				AllowedChildren: []umid.UMID{},
			},
		},
		{
			id:             umid.MustParse("a41ee21e-6c56-41b3-81a9-1c86578b6b3c"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "World",
			categoryName:   "Worlds",
			description:    utils.GetPTR("World Type"),
			options: &entry.ObjectOptions{
				Private:        utils.GetPTR(false),
				Visible:        utils.GetPTR(entry.AllObjectVisibleType),
				DefaultTiles:   []any{},
				FrameTemplates: map[string]any{},
				AllowedChildren: []umid.UMID{
					umid.MustParse("4ed3a5bb-53f8-4511-941b-079029111111"), // Custom claimable
					umid.MustParse("4ed3a5bb-53f8-4511-941b-07902982c31c"), // Custom objects
				},
			},
		},
		{
			id:             umid.MustParse("4ed3a5bb-53f8-4511-941b-07902982c31c"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "Custom objects",
			categoryName:   "Custom",
			description:    utils.GetPTR("Custom placed objects"),
			options: &entry.ObjectOptions{
				Visible: utils.GetPTR(entry.AllObjectVisibleType),
			},
		},
		{
			id:             umid.MustParse("69d8ae40-df9b-4fc8-af95-32b736d2bbcd"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "Service Space",
			categoryName:   "Service Spaces",
			description:    utils.GetPTR(""),
			options: &entry.ObjectOptions{
				Private:         utils.GetPTR(false),
				Visible:         utils.GetPTR(entry.InvisibleObjectVisibleType),
				DefaultTiles:    []any{},
				FrameTemplates:  map[string]any{},
				AllowedChildren: []umid.UMID{},
			},
		},
		{
			id:             umid.MustParse("75b56447-c4f1-4020-b8fc-d68704a11d65"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "Generic Space",
			categoryName:   "Generic Spaces",
			description:    utils.GetPTR(""),
			options: &entry.ObjectOptions{
				Private:         utils.GetPTR(false),
				Visible:         utils.GetPTR(entry.AllObjectVisibleType),
				DefaultTiles:    []any{},
				FrameTemplates:  map[string]any{},
				AllowedChildren: []umid.UMID{},
			},
		},
		{
			id:             umid.MustParse("f9607e55-63e8-4cb1-ae47-66395199975d"),
			asset2dID:      nil,
			asset3dID:      nil,
			objectTypeName: "morgue",
			categoryName:   "Morgues",
			description:    utils.GetPTR("morgue"),
			options: &entry.ObjectOptions{
				Private:         utils.GetPTR(false),
				Visible:         utils.GetPTR(entry.InvisibleObjectVisibleType),
				DefaultTiles:    []any{},
				FrameTemplates:  map[string]any{},
				AllowedChildren: []umid.UMID{},
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
				return fmt.Errorf("failed to create object type: asset_2d not found: %s", item.asset2dID)
			}
			if err := objectType.SetAsset2d(asset2d, false); err != nil {
				return errors.WithMessagef(
					err, "failed to set asset_2d: %s for object_type: %s", item.asset2dID, item.id,
				)
			}
		}

		if item.asset3dID != nil {
			asset3d, ok := node.GetAssets3d().GetAsset3d(*item.asset3dID)
			if !ok {
				return fmt.Errorf("failed to create object type: asset_3d not found: %s", item.asset3dID)
			}
			if err := objectType.SetAsset3d(asset3d, false); err != nil {
				return errors.WithMessagef(
					err, "failed to set asset_3d: %s for object_type: %s", item.asset3dID, item.id,
				)
			}
		}

		if err := objectType.SetName(item.objectTypeName, false); err != nil {
			return errors.WithMessagef(err, "failed to set name: %s for object_type: %s", item.objectTypeName, item.id)
		}

		if err := objectType.SetCategoryName(item.categoryName, false); err != nil {
			return errors.WithMessagef(
				err, "failed to set category name: %s for object_type: %s", item.categoryName, item.id,
			)
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
