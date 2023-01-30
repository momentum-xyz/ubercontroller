package seed

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedAssets3d(node universe.Node) error {
	/*
		select asset_3d_id, meta
		from asset_3d
		where options != 'null'::jsonb
		order by created_at, meta
	*/

	type item struct {
		asset3DID uuid.UUID
		options   *entry.Asset3dOptions
		meta      *entry.Asset3dMeta
	}

	items := []*item{
		{
			asset3DID: uuid.MustParse("658611b8-a86a-4bf0-a956-12129b06dbfd"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":         "Arrival",
				"type":         0,
				"category":     "skybox",
				"preview_hash": "bbe7629c15c6f96559def3ace4b203ca",
			},
		},
		{
			asset3DID: uuid.MustParse("b4e10a57-de25-4793-b9f1-35fcb6daf404"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Abyss",
				"type":     0,
				"category": "skybox",
			},
		},
		{
			asset3DID: uuid.MustParse("a0b93937-fbab-4b9f-b5c1-93c47ffef624"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":         "QuantumFlux",
				"type":         0,
				"category":     "skybox",
				"preview_hash": "3805f85c58e0bc817420efabc946517c",
			},
		},
		{
			asset3DID: uuid.MustParse("581a8563-1329-4de6-8c8e-869770526ff0"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Docking station",
				"type":     0,
				"category": "odyssey",
			},
		},
		{
			asset3DID: uuid.MustParse("a6862b31-8f80-497d-b9d6-8234e6a71773"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Innerverse Bulb",
				"type":     0,
				"category": "odyssey",
			},
		},
		{
			asset3DID: uuid.MustParse("de240de6-d911-4d84-9406-8b81550dfea8"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Portal",
				"type":     0,
				"category": "odyssey",
			},
		},
		{
			asset3DID: uuid.MustParse("7e20a110-149b-4c6e-b1ab-a25cbdc066e6"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Bucky with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("bda945cc-8fb6-4e4d-94e3-0d0480c78893"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Capsule with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("9bfd83c1-7dad-4cc9-a97b-69c7b9ad931d"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Cone with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("ad49552f-67c8-47f4-bcad-fc6f6deac1fc"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Cube with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("97a8bd60-bbdb-4c28-964c-280322f84d4a"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Cylinder with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("021a6576-25c2-4245-a48e-73f1e9c4c25a"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Disc with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("b238c592-ba69-4721-a275-30f9738db31e"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Dodeca with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("414cfe78-a3b1-4d48-a473-5b1cf163ea3e"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Icosa with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("e50d9cef-4588-4032-80ed-3bb2fb133835"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Octo with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("3aa77816-345c-4f63-8b0d-3c1ec5585b23"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Quad with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("e369c559-a1ca-4c5b-9e16-d1c942bb86b8"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Sphere with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("c57d792d-ee61-4b2d-9ea3-b49c6ce9991a"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Torus with image",
				"type":     0,
				"category": "basic",
			},
		},
		{
			asset3DID: uuid.MustParse("9ac3b215-f1fd-4d23-bb8b-7849f4e13659"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Bucky",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("fcb944f4-a952-4d72-bf68-8d7bf249fda9"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Capsule",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("01f475b3-8922-4acf-8bb5-1c4e870aab7a"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Cone",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("008472fd-6033-4ecb-81ca-fe345334791e"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Cube",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("1b093918-8e7d-4ee3-9f5d-af5f209ae84a"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Cylinder",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("aeda3d26-d5dd-455f-b162-014d2c2e36ab"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Disc",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("0576b67f-3214-4862-8973-c984a30dfda9"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Dodeca",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("5044b89c-1d5c-457d-a06f-36e05455f0d0"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Icosa",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("5b447fea-b639-4895-ba3a-4ac8487252c6"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Octa",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("bece5db8-1ae3-4839-8e46-63afb947c96d"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Quad",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("c4338f9f-9f5b-4ca0-9939-3644bbddbc9e"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Sphere",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("ee7961ea-e01f-4d1d-9ad8-673c2fb49fb2"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":     "Torus",
				"type":     0,
				"category": "noclick",
			},
		},
		{
			asset3DID: uuid.MustParse("eb0fe08b-155d-4783-a6b4-a49bd5be6a8e"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":         "Gogogo",
				"type":         0,
				"category":     "skybox",
				"preview_hash": "0543ba59db159b6c6a2b60395460f7dd",
			},
		},
		{
			asset3DID: uuid.MustParse("f7be7dac-f103-458f-9aea-8d937e6e493c"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":         "Honey",
				"type":         0,
				"category":     "skybox",
				"preview_hash": "8a14a7f55c1c419db7595f9dbc59dd78",
			},
		},
		{
			asset3DID: uuid.MustParse("f70dceda-98cc-4fce-8a0d-0b2ce864e7bd"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":         "PlasmaSummer",
				"type":         0,
				"category":     "skybox",
				"preview_hash": "995b64c7c7efb2795a8ceade7ba75995",
			},
		},
		{
			asset3DID: uuid.MustParse("5079b26d-3653-419c-97fd-6aa6d0361a56"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":         "ShowTime",
				"type":         0,
				"category":     "skybox",
				"preview_hash": "8268b7490370f60d715540a8f6ff68f2",
			},
		},
		{
			asset3DID: uuid.MustParse("67f3e7e9-8dea-4458-8e54-26e05246296c"),
			options:   &entry.Asset3dOptions{},
			meta: &entry.Asset3dMeta{
				"name":         "Temptations",
				"type":         0,
				"category":     "skybox",
				"preview_hash": "7c8524c0304d8bc68af0093f2d6ff472",
			},
		},
	}

	for _, item := range items {
		asset, err := node.GetAssets3d().CreateAsset3d(item.asset3DID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create asset_3d: %s", item.asset3DID)
		}

		_, err = asset.SetOptions(modify.MergeWith(item.options), false)
		if err != nil {
			return errors.WithMessagef(err, "failed to set asset_3d options: %s", item.asset3DID)
		}

		if err = asset.SetMeta(item.meta, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset_3d meta: %s", item.asset3DID)
		}

	}

	return nil
}
