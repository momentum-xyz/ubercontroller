package seed

import (
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
			pluginID:      uuid.MustParse(noname1PluginID),
			attributeName: "problem",
			description:   "Problem for space",
			options: &entry.AttributeOptions{
				"render_type":  "texture",
				"content_type": "text",
			},
		},
		{
			pluginID:      uuid.MustParse(noname1PluginID),
			attributeName: "solution",
			description:   "solution for space",
			options: &entry.AttributeOptions{
				"render_type":  "texture",
				"content_type": "text",
			},
		},
		{
			pluginID:      uuid.MustParse(noname1PluginID),
			attributeName: "tile",
			description:   "tile for space",
			options: &entry.AttributeOptions{
				"render_type": "texture",
			},
		},
		{
			pluginID:      uuid.MustParse(noname1PluginID),
			attributeName: "third",
			description:   "Third screen for space",
			options: &entry.AttributeOptions{
				"render_type": "texture",
			},
		},
		{
			pluginID:      uuid.MustParse(noname1PluginID),
			attributeName: "poster",
			description:   "Poster for space",
			options: &entry.AttributeOptions{
				"render_type": "texture",
			},
		},
		{
			pluginID:      uuid.MustParse(noname1PluginID),
			attributeName: "meme",
			description:   "Meme for space",
			options: &entry.AttributeOptions{
				"render_type": "texture",
			},
		},
		{
			pluginID:      uuid.MustParse(noname1PluginID),
			attributeName: "video",
			description:   "Video for space",
			options: &entry.AttributeOptions{
				"render_type":  "texture",
				"content_type": "video",
			},
		},
		//
		{
			pluginID:      uuid.MustParse(miroPluginID),
			attributeName: "state",
			description:   "Miro state",
			options: &entry.AttributeOptions{
				"posbus_auto": map[string]any{
					"scope":   []string{"space"},
					"send_to": 1,
				},
			},
		},
		{
			pluginID:      uuid.MustParse(miroPluginID),
			attributeName: "config",
			description:   "Miro configuration",
			options:       nil,
		},
		//
		{
			pluginID:      uuid.MustParse(videoPluginID),
			attributeName: "state",
			description:   "State of the video tile",
			options: &entry.AttributeOptions{
				"unity_auto": map[string]string{
					"slot_name":    "Block",
					"slot_type":    "texture",
					"value_field":  "value",
					"content_type": "video",
				},
			},
		},
		//
		{
			pluginID:      uuid.MustParse(high5PluginID),
			attributeName: "count",
			description:   "High5s given",
			options:       nil,
		},
		//
		{
			pluginID:      universe.GetKusamaPluginID(),
			attributeName: "challenge_store",
			description:   "auth challenge store",
			options:       nil,
		},
		{
			pluginID:      universe.GetKusamaPluginID(),
			attributeName: "wallet",
			description:   "Kusama/Substrate wallet",
			options:       nil,
		},
		//
		{
			pluginID:      uuid.MustParse(googleDrivePluginID),
			attributeName: "config",
			description:   "Google Drive configuration",
			options:       nil,
		},
		{
			pluginID:      uuid.MustParse(googleDrivePluginID),
			attributeName: "state",
			description:   "Google Drive state",
			options: &entry.AttributeOptions{
				"posbus_auto": map[string]any{
					"scope":   []string{"space"},
					"send_to": 1,
				},
			},
		},
		// CORE
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
				"unity_auto": map[string]any{
					"slot_name":    "name",
					"slot_type":    "texture",
					"value_field":  "name",
					"content_type": "text",
					"text_render_template": map[string]any{
						"x": 0,
						"y": 0,
						"text": map[string]any{
							"padX":     0,
							"padY":     1,
							"wrap":     false,
							"alignH":   "center",
							"alignV":   "center",
							"string":   "%TEXT%",
							"fontfile": "",
							"fontsize": 0,
							"fontcolor": []any{
								220,
								220,
								200,
								255,
							},
						},
						"color": []any{
							0,
							255,
							0,
							0,
						},
						"width":     1024,
						"height":    64,
						"thickness": 0,
						"background": []any{
							0,
							0,
							0,
							0,
						},
					},
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "jwt_key",
			description:   "",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "VoiceChatAction",
			description:   "Voice chat user actions",
			options: &entry.AttributeOptions{
				"posbus_auto": map[string]any{
					"scope":   []string{"space"},
					"topic":   "voice-chat-action",
					"send_to": 1,
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "VoiceChatUser",
			description:   "Voice chat users",
			options: &entry.AttributeOptions{
				"posbus_auto": map[string]any{
					"scope":   []string{"space"},
					"topic":   "voice-chat-user",
					"send_to": 1,
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "staking",
			description:   "Odyssey staking information",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "teleport",
			description:   "Target World ID to teleport",
			options:       nil,
		},
		//
		{
			pluginID:      uuid.MustParse(textPluginID),
			attributeName: "state",
			description:   "State of the text tile",
			options: &entry.AttributeOptions{
				"unity_auto": map[string]any{
					"slot_name":    "description",
					"slot_type":    "Block",
					"value_field":  "value",
					"content_type": "text",
					"text_render_template": map[string]any{
						"x": 0,
						"y": 0,
						"sub": []any{
							map[string]any{
								"x": 4,
								"y": 4,
								"text": map[string]any{
									"padX":     4,
									"padY":     4,
									"wrap":     true,
									"align":    0,
									"alignH":   "left",
									"alignV":   "top",
									"string":   "%TEXT%",
									"fontfile": "IBMPlexSans-SemiBold",
									"fontsize": 17,
									"fontcolor": []any{
										0,
										255,
										255,
									},
								},
								"color": []any{
									0,
									255,
									255,
								},
								"width":     1012,
								"height":    504,
								"thickness": 1,
								"background": []any{
									0,
									40,
									0,
								},
							},
						},
						"color": []any{
							20,
							20,
							20,
						},
						"width":     1024,
						"height":    512,
						"thickness": 4,
						"background": []any{
							10,
							10,
							10,
						},
					},
				},
			},
		},
		//
		//
		{
			pluginID:      uuid.MustParse(imagePluginID),
			attributeName: "state",
			description:   "State of the image tile",
			options: &entry.AttributeOptions{
				"unity_auto": map[string]any{
					"slot_name":    "Block",
					"slot_type":    "texture",
					"content_type": "image",
				},
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