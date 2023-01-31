package seed

import (
	"crypto/rand"
	"math/big"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedNodeAttributes(node universe.Node) error {
	type item struct {
		pluginID      uuid.UUID
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
				"id":                      "0e347f9e-bba9-48c1-a3f9-4258ca230481",
				"name":                    "dev2-node",
				"user_id_salt":            "63a81341-39c3-4d23-8261-3d22d70839bb",
				"entrance_world":          "d83670c7-a120-47a4-892d-f9ec75604f74",
				"guest_user_type":         "76802331-37b3-44fa-9010-35008b0cbaec",
				"normal_user_type":        "00000000-0000-0000-0000-000000000006",
				"docking_hub_space_type":  "b59abd4d-f54d-4a97-8b6d-16a2037ddd8f",
				"dock_station_space_type": "27456794-f9aa-44b5-90a5-e307fc21bc3d",
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
				"spaces": []any{
					map[string]any{
						"label": "effects_emitter",
						"options": map[string]any{
							"editable": false,
						},
						"space_name":    "Effects emitter",
						"space_type_id": "4fe4ed05-9024-461a-97d6-22666e8a4f46",
					},
					map[string]any{
						"label": "skybox",
						"options": map[string]any{
							"editable": false,
						},
						"space_name":    "Skybox",
						"space_type_id": "d7a41cbd-5cfe-454b-b522-76f22fa55026",
					},
					map[string]any{
						"label": "docking_station",
						"options": map[string]any{
							"editable": true,
						},
						"space_name":    "Docking station",
						"space_type_id": "27456794-f9aa-44b5-90a5-e307fc21bc3d",
					},
				},
				"space_type_id": "a41ee21e-6c56-41b3-81a9-1c86578b6b3c",
				"space_attributes": []any{
					map[string]any{
						"value": map[string]any{
							"lod": []any{
								6400,
								40000,
								160000,
							},
							"decorations":       []any{},
							"avatar_controller": "2eeb56e2-5ac4-42c2-a348-d167862260f6",
							"skybox_controller": "658611b8-a86a-4bf0-a956-12129b06dbfd",
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
