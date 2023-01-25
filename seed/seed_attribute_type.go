package seed

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedAttributeType(node universe.Node) error {
	type item struct {
		pluginID      uuid.UUID
		attributeName string
		description   string
		options       *entry.AttributeOptions
	}

	items := []*item{
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "node_settings",
			description:   "",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "world_settings",
			description:   "",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "world_meta",
			description:   "Holds world metadata and decorations",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "magic_links",
			description:   "",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "emojis",
			description:   "",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "emoji",
			description:   "",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "events",
			description:   "Space events",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "world_template",
			description:   "Basic template settings for any new world",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "news_feed",
			description:   "News feed storage",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "skybox_custom",
			description:   "Holds skybox data such as texture",
			options: &entry.AttributeOptions{
				"unity_auto": map[string]string{
					"slot_name":    "skybox_custom",
					"slot_type":    "texture",
					"content_type": "image",
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "skybox_list",
			description:   "Holds initial list of skyboxes",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "high_five",
			description:   "high fives",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "description",
			description:   "description",
			options: &entry.AttributeOptions{
				"permissions": "admin",
				"render_type": "texture",
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "screenshare",
			description:   "Odyssey screenshare state",
			options: &entry.AttributeOptions{
				"posbus_auto": map[string]any{
					"scope":   []string{"space"},
					"topic":   "screenshare-action",
					"send_to": 1,
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "name",
			description:   "Space name",
			options: &entry.AttributeOptions{
				"unity_auto": jsonStringToMap(
					`{
        "slot_name": "name",
        "slot_type": "texture",
        "value_field": "name",
        "content_type": "text",
        "text_render_template": {
            "x": 0,
            "y": 0,
            "text": {
                "padX": 0,
                "padY": 1,
                "wrap": false,
                "alignH": "center",
                "alignV": "center",
                "string": "%TEXT%",
                "fontfile": "",
                "fontsize": 0,
                "fontcolor": [
                    220,
                    220,
                    200,
                    255
                ]
            },
            "color": [
                0,
                255,
                0,
                0
            ],
            "width": 1024,
            "height": 64,
            "thickness": 0,
            "background": [
                0,
                0,
                0,
                0
            ]
        }
    }`),
			},
		},
	}

	for _, item := range items {
		attributeType, err := node.GetAttributeTypes().CreateAttributeType(entry.AttributeTypeID{
			PluginID: item.pluginID,
			Name:     item.attributeName,
		})
		if err != nil {
			return errors.WithMessagef(err, "failed to create attribute type: %s %s", item.pluginID, item.attributeName)
		}
		if err := attributeType.SetDescription(&item.description, false); err != nil {
			return errors.WithMessagef(err, "failed to set attribute type description: %s", item.description)
		}
		if _, err := attributeType.SetOptions(modify.MergeWith(item.options), false); err != nil {
			return errors.WithMessagef(err, "failed to set attribute type options: %s", item.options)
		}
	}

	return nil
}

func jsonStringToMap(rawJson string) map[string]any {
	var data map[string]any
	err := json.Unmarshal([]byte(rawJson), &data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}
