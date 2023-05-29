package seed

import (
	"context"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedAssets3d(ctx context.Context, node universe.Node) error {
	/*
		select asset_3d_id, meta
		from asset_3d
		where options != 'null'::jsonb
		order by created_at, meta
	*/

	items := []*entry.Asset3d{
		{
			Asset3dID: umid.MustParse("a55f9ca74b45692e204fe37ed9dc3d78"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Bucky with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("eea924c06e33393fe06ee6631e8860e9"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Capsule with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("8a7e55f5934d8ebf17bb39e2d8d9bfa1"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Cone with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("5b5bd8720328e38c1b54bf2bfa70fc85"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Cube with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("46d923ad21ff276dc3c4ead2212bcb02"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Cylinder with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("97daa12f9b2e536d78513b0837175e4c"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Disc with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("839b21db52ff45ce7484fd1b59ebb087"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Dodeca with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("dad4e8a4cdcc41749d77f7e849bba352"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Icosa with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("a1f144deb21ad1e906356eb250927326"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Octo with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("418c4963623a391c795de6080be11899"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Quad with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("313a97ccfe1b39bb56e7516d213cc23d"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Sphere with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
		{
			Asset3dID: umid.MustParse("6e8fec1cff95df661375e312f6447b3d"),
			Options:   &entry.Asset3dOptions{},
			Meta: &entry.Asset3dMeta{
				"name":     "Torus with image",
				"type":     dto.BasicAsset3dType,
				"category": "basic",
			},
		},
	}

	for _, item := range items {
		assetUserId := universe.AssetUserIDPair{
			AssetID: item.Asset3dID,
			UserID:  umid.MustParse("00000000-0000-0000-0000-000000000003"), // Odin
		}

		asset, err := node.GetAssets3d().CreateAsset3d(assetUserId)
		if err != nil {
			return errors.WithMessagef(err, "failed to create asset_3d: %s", item.Asset3dID)
		}

		_, err = asset.SetOptions(modify.MergeWith(item.Options), false)
		if err != nil {
			return errors.WithMessagef(err, "failed to set asset_3d options: %s", item.Asset3dID)
		}

		if err = asset.SetMeta(item.Meta, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset_3d meta: %s", item.Asset3dID)
		}
	}

	return nil
}
