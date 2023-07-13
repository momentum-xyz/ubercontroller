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
				"umid":             umid.New(),
				"name":             "dev2-node",
				"user_id_salt":     umid.New(),
				"entrance_world":   "d83670c7-a120-47a4-892d-f9ec75604f74",
				"guest_user_type":  guestUserTypeID,
				"normal_user_type": normalUserTypeID,
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
			attributeName: "blockadelabs",
			value:         &entry.AttributeValue{},
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "world_template",
			value: &entry.AttributeValue{
				"objects":        []any{},
				"random_spaces":  []any{},
				"object_type_id": "a41ee21e-6c56-41b3-81a9-1c86578b6b3c",
				"object_attributes": []any{
					map[string]any{
						"value": map[string]any{
							"lod": []any{
								6400,
								40000,
								160000,
							},
						},
						"options":        nil,
						"plugin_id":      "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
						"attribute_name": "world_meta",
					},
					map[string]any{
						"plugin_id":      universe.GetSystemPluginID(),
						"attribute_name": "active_skybox",
						"value": map[string]any{
							"render_hash": "26485e74acb29223ba7a9fa600d36c7f", // TODO: default skybox hash
						},
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
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "tracker_ai_usage",
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
