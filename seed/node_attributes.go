package seed

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedNodeAttributes(ctx context.Context, node universe.Node) error {
	type item struct {
		pluginID      umid.UMID
		attributeName string
		value         *entry.AttributeValue
		options       *entry.AttributeOptions
	}

	secret, err := generateRandomString(40)
	if err != nil {
		return errors.WithMessage(err, "failed to generate secret")
	}

	signature, err := generateRandomString(30)
	if err != nil {
		return errors.WithMessage(err, "failed to generate signature")
	}

	items := []*item{
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "node_settings",
			value: &entry.AttributeValue{
				"umid":                     umid.New(),
				"name":                     "dev2-node",
				"user_id_salt":             umid.New(),
				"entrance_world":           "d83670c7-a120-47a4-892d-f9ec75604f74",
				"guest_user_type":          guestUserTypeID,
				"normal_user_type":         normalUserTypeID,
				"docking_hub_object_type":  "b59abd4d-f54d-4a97-8b6d-16a2037ddd8f",
				"dock_station_object_type": "27456794-f9aa-44b5-90a5-e307fc21bc3d",
			},
			options: nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "jwt_key",
			value: &entry.AttributeValue{
				"secret":    secret,
				"signature": signature,
			},
			options: nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "world_template",
			value: &entry.AttributeValue{
				"objects": []any{
					map[string]any{
						"label": "effects_emitter",
						"options": map[string]any{
							"editable": false,
						},
						"object_name":    "Effects emitter",
						"object_type_id": "4fe4ed05-9024-461a-97d6-22666e8a4f46",
					},
					map[string]any{
						"label": "docking_station",
						"options": map[string]any{
							"editable": true,
						},
						"object_name":    "Docking station",
						"object_type_id": "27456794-f9aa-44b5-90a5-e307fc21bc3d",
					},
					map[string]any{
						"label": "skybox",
						"options": map[string]any{
							"editable": false,
						},
						"object_name":    "Skybox",
						"asset_3d_id":    "313a597a-8b9a-47a7-9908-52bdc7a21a3e",
						"object_type_id": "d7a41cbd-5cfe-454b-b522-76f22fa55026",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "Skybox",
									"auto_render_hash": "7e7aab88b51141f184de7a102810c3e0",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
							map[string]any{
								"value": map[string]any{
									"render_hash": "26485e74acb29223ba7a9fa600d36c7f",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "active_skybox",
							},
						},
					},
					map[string]any{
						"label": "big_boiler",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 1,
								"y": 1,
								"z": 1,
							},
							"position": map[string]any{
								"x": 15.887797,
								"y": 51.115723,
								"z": 144.49608,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 91.23572,
								"z": 0,
							},
						},
						"object_name":    "Big Boiler",
						"asset_3d_id":    "3aa837cf-0d47-44d1-ad3e-19828ea9d0b6",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "boiler",
									"auto_render_hash": "597af7a55a21fa1276c4a844816341ec",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
						},
					},
					map[string]any{
						"label": "innerverse_core",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 1,
								"y": 1,
								"z": 1,
							},
							"position": map[string]any{
								"x": 86.6892,
								"y": 59.070778,
								"z": 156.44916,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 0,
							},
						},
						"object_name":    "Innerverse Core",
						"asset_3d_id":    "ca276094-caa2-420f-892e-c715315d1b75",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "InverseCore",
									"auto_render_hash": "aa590822e031b9580f78472b5825c7c3",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
						},
					},
					map[string]any{
						"label": "scatter_slabs",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 2.1864054,
								"y": 1.1553961,
								"z": 2.0322535,
							},
							"position": map[string]any{
								"x": 51.27168,
								"y": 47.102135,
								"z": 147.76534,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 0,
							},
						},
						"object_name":    "Scatter Slabs",
						"asset_3d_id":    "036736dd-165d-4315-b4f0-2b78d9c28adf",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "scatterslab",
									"auto_render_hash": "28ae0dd94fa5ca49c5e167bcaa4c2f47",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
						},
					},
					map[string]any{
						"label": "single_slab",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 2.4160378,
								"y": 1.0045278,
								"z": 2.4337356,
							},
							"position": map[string]any{
								"x": 15.918898,
								"y": 46.477386,
								"z": 144.60236,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 0,
							},
						},
						"object_name":    "Single Slab 1",
						"asset_3d_id":    "ce89b37b-d4e7-4e83-9351-4009d4c0f14e",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "slab",
									"auto_render_hash": "04ea486bf0e78574ad589a275a4caf49",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
						},
					},
					map[string]any{
						"label": "single_slab",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 2.0890555,
								"y": 1,
								"z": 2.2992187,
							},
							"position": map[string]any{
								"x": 30.522566,
								"y": 45.030613,
								"z": 164.6621,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 0,
							},
						},
						"object_name":    "Single Slab 2",
						"asset_3d_id":    "036736dd-165d-4315-b4f0-2b78d9c28adf",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "slab",
									"auto_render_hash": "04ea486bf0e78574ad589a275a4caf49",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
						},
					},
					map[string]any{
						"label": "portal",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 1.9816073,
								"y": 2.3933778,
								"z": 1,
							},
							"position": map[string]any{
								"x": 30.702652,
								"y": 51.685955,
								"z": 163.36012,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 180.06071,
							},
						},
						"object_name":    "Portal",
						"asset_3d_id":    "de240de6-d911-4d84-9406-8b81550dfea8",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "Portal to Odyssey's Odyssey",
									"auto_render_hash": "3261815daa038b46adc0dabc096f47cf",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
						},
					},
					map[string]any{
						"label": "portal",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 1.9816073,
								"y": 2.3933778,
								"z": 1,
							},
							"position": map[string]any{
								"x": 30.702652,
								"y": 51.685955,
								"z": 163.36012,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 180.06071,
							},
						},
						"object_name":    "Portal",
						"asset_3d_id":    "de240de6-d911-4d84-9406-8b81550dfea8",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "Portal to Odyssey's Odyssey",
									"auto_render_hash": "3261815daa038b46adc0dabc096f47cf",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
						},
					},
					map[string]any{
						"label": "video",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 4.390896,
								"y": 1,
								"z": 3.0823991,
							},
							"position": map[string]any{
								"x": 56.173256,
								"y": 50.4917,
								"z": 143.80066,
							},
							"rotation": map[string]any{
								"x": 89.11309,
								"y": 142.04346,
								"z": 233.6865,
							},
						},
						"object_name":    "Video",
						"asset_2d_id":    "bda25d5d-2aab-45b4-9e8a-23579514cec1",
						"asset_3d_id":    "3aa77816-345c-4f63-8b0d-3c1ec5585b23",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "Tutorial Poster",
									"auto_render_hash": "9f82e480a6c4cadcae915e5789b023d5",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
							map[string]any{
								"value": map[string]any{
									"render_hash": "4e15e728443dc4e8304a60c2625dddc6",
								},
								"options":        nil,
								"plugin_id":      "ff40fbf0-8c22-437d-b27a-0258f99130fe",
								"attribute_name": "state",
							},
							map[string]any{
								"value": map[string]any{
									"youtube_url": "https://www.youtube.com/watch?v=UmSJIEZQAyQ",
								},
								"options":        nil,
								"plugin_id":      "308fdacc-8c2d-40dc-bd5f-d1549e3e03ba",
								"attribute_name": "state",
							},
						},
					},
					map[string]any{
						"label": "video",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 1.6822052,
								"y": 1,
								"z": 2.8605652,
							},
							"position": map[string]any{
								"x": 63.490574,
								"y": 51.481842,
								"z": 160.97379,
							},
							"rotation": map[string]any{
								"x": 87.640175,
								"y": 180.00003,
								"z": 306.21182,
							},
						},
						"object_name":    "Video 2",
						"asset_2d_id":    "bda25d5d-2aab-45b4-9e8a-23579514cec1",
						"asset_3d_id":    "3aa77816-345c-4f63-8b0d-3c1ec5585b23",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
						"object_attributes": []any{
							map[string]any{
								"value": map[string]any{
									"name":             "Staking Poster",
									"auto_render_hash": "af99f07c4b86c8f20e95cdf52c2b7e5f",
								},
								"options":        nil,
								"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
								"attribute_name": "name",
							},
							map[string]any{
								"value": map[string]any{
									"render_hash": "c7b276ad8f3882d3002b310cdf04141a",
								},
								"options":        nil,
								"plugin_id":      "ff40fbf0-8c22-437d-b27a-0258f99130fe",
								"attribute_name": "state",
							},
							map[string]any{
								"value": map[string]any{
									"youtube_url": "https://www.youtube.com/watch?v=UmSJIEZQAyQ",
								},
								"options":        nil,
								"plugin_id":      "308fdacc-8c2d-40dc-bd5f-d1549e3e03ba",
								"attribute_name": "state",
							},
						},
					},
				},
				"random_spaces": []any{
					map[string]any{
						"label": "kiddo_jelly_1",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 1,
								"y": 1,
								"z": 1,
							},
							"position": map[string]any{
								"x": 62.438713,
								"y": 54.17204,
								"z": 152.16934,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 0,
							},
						},
						"object_name":    "Kiddo Jelly 1",
						"asset_3d_id":    "a0ef8db2-dfa0-4d17-83bc-e662a2f25e74",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
					},
					map[string]any{
						"label": "kiddo_jelly_2",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 1,
								"y": 1,
								"z": 1,
							},
							"position": map[string]any{
								"x": 62.438713,
								"y": 54.17204,
								"z": 152.16934,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 0,
							},
						},
						"object_name":    "Kiddo Jelly 2",
						"asset_3d_id":    "f30cd02f-a51e-4cab-ad66-c25cc6c76c66",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
					},
					map[string]any{
						"label": "mega_jelly",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 1,
								"y": 1,
								"z": 1,
							},
							"position": map[string]any{
								"x": 62.438713,
								"y": 54.17204,
								"z": 152.16934,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 0,
							},
						},
						"object_name":    "Mega Jelly",
						"asset_3d_id":    "a76d85b3-7831-429e-ab5a-63d8e69cd380",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
					},
					map[string]any{
						"label": "object_whale",
						"options": map[string]any{
							"subs": map[string]any{
								"visible": 1,
							},
						},
						"position": map[string]any{
							"scale": map[string]any{
								"x": 1,
								"y": 1,
								"z": 1,
							},
							"position": map[string]any{
								"x": 62.438713,
								"y": 54.17204,
								"z": 152.16934,
							},
							"rotation": map[string]any{
								"x": 0,
								"y": 0,
								"z": 0,
							},
						},
						"object_name":    "Object Whale",
						"asset_3d_id":    "355657f1-95af-49a7-817f-aecb431cf4dd",
						"object_type_id": "4ed3a5bb-53f8-4511-941b-07902982c31c",
					},
				},
				"object_type_id": "a41ee21e-6c56-41b3-81a9-1c86578b6b3c",
				"object_attributes": []any{
					map[string]any{
						"value": map[string]any{
							"lod": []any{
								6400,
								40000,
								160000,
							},
							"decorations":       []any{},
							"avatar_controller": "2eeb56e2-5ac4-42c2-a348-d167862260f6",
							"skybox_controller": skyboxArrivalAsset3dID,
						},
						"options":        nil,
						"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
						"attribute_name": "world_meta",
					},
				},
			},
			options: nil,
		},
		{
			pluginID:      universe.GetKusamaPluginID(),
			attributeName: "challenge_store",
			value:         &entry.AttributeValue{},
			options:       nil,
		},
	}

	for _, item := range items {
		payload := entry.AttributePayload{
			Value:   item.value,
			Options: item.options,
		}
		_, err := node.GetNodeAttributes().Upsert(
			entry.NewAttributeID(item.pluginID, item.attributeName),
			modify.MergeWith(&payload),
			false,
		)
		if err != nil {
			return errors.WithMessagef(err, "failed to upsert node attribute: %s %s", item.pluginID, item.attributeName)
		}
	}

	return nil
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
