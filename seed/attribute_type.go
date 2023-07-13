package seed

import (
	"context"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func seedAttributeType(ctx context.Context, node universe.Node) error {
	type item struct {
		pluginID      umid.UMID
		attributeName string
		description   string
		options       *entry.AttributeOptions
	}

	items := []*item{
		{
			pluginID:      umid.MustParse(OdysseyHackatonPluginID),
			attributeName: "problem",
			description:   "Problem for space",
			options: &entry.AttributeOptions{
				"render_type":  "texture",
				"content_type": "text",
			},
		},
		{
			pluginID:      umid.MustParse(OdysseyHackatonPluginID),
			attributeName: "solution",
			description:   "solution for space",
			options: &entry.AttributeOptions{
				"render_type":  "texture",
				"content_type": "text",
			},
		},
		{
			pluginID:      umid.MustParse(OdysseyHackatonPluginID),
			attributeName: "tile",
			description:   "tile for space",
			options: &entry.AttributeOptions{
				"render_type": "texture",
			},
		},
		{
			pluginID:      umid.MustParse(OdysseyHackatonPluginID),
			attributeName: "third",
			description:   "Third screen for space",
			options: &entry.AttributeOptions{
				"render_type": "texture",
			},
		},
		{
			pluginID:      umid.MustParse(OdysseyHackatonPluginID),
			attributeName: "poster",
			description:   "Poster for space",
			options: &entry.AttributeOptions{
				"render_type": "texture",
			},
		},
		{
			pluginID:      umid.MustParse(OdysseyHackatonPluginID),
			attributeName: "meme",
			description:   "Meme for space",
			options: &entry.AttributeOptions{
				"render_type": "texture",
			},
		},
		{
			pluginID:      umid.MustParse(OdysseyHackatonPluginID),
			attributeName: "video",
			description:   "Video for space",
			options: &entry.AttributeOptions{
				"render_type":  "texture",
				"content_type": "video",
			},
		},
		//
		{
			pluginID:      umid.MustParse(miroPluginID),
			attributeName: "state",
			description:   "Miro state",
			options: &entry.AttributeOptions{
				"posbus_auto": map[string]any{
					"scope":   []string{"object"},
					"send_to": 1,
				},
			},
		},
		{
			pluginID:      umid.MustParse(miroPluginID),
			attributeName: "config",
			description:   "Miro configuration",
			options:       nil,
		},
		//
		{
			pluginID:      umid.MustParse(videoPluginID),
			attributeName: "state",
			description:   "State of the video tile",
			options: &entry.AttributeOptions{
				"render_auto": map[string]string{
					"slot_name":    "object_texture",
					"slot_type":    "texture",
					"value_field":  "value",
					"content_type": "video",
				},
			},
		},
		//
		{
			pluginID:      umid.MustParse(high5PluginID),
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
			options: &entry.AttributeOptions{
				"permissions": entry.PermissionsAttributeOption{
					Read:  string(entry.PermissionAdmin) + "+" + string(entry.PermissionUserOwner),
					Write: string(entry.PermissionAdmin),
				},
			},
		},
		//
		{
			pluginID:      umid.MustParse(googleDrivePluginID),
			attributeName: "config",
			description:   "Google Drive configuration",
			options:       nil,
		},
		{
			pluginID:      umid.MustParse(googleDrivePluginID),
			attributeName: "state",
			description:   "Google Drive state",
			options: &entry.AttributeOptions{
				"posbus_auto": map[string]any{
					"scope":   []string{"object"},
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
			attributeName: "world_avatar",
			description:   "",
			options: &entry.AttributeOptions{
				"render_auto": map[string]string{
					"slot_type":    "texture",
					"content_type": "image",
				},
			},
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
			attributeName: "website_link",
			description:   "",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "last_known_position",
			description:   "Holds users last known position",
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
			attributeName: "active_skybox",
			description:   "Holds skybox data such as texture",
			options: &entry.AttributeOptions{
				"render_auto": map[string]string{
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
				"render_type": "texture",
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "screenshare",
			description:   "Odyssey screenshare state",
			options: &entry.AttributeOptions{
				"posbus_auto": map[string]any{
					"scope":   []string{"object"},
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
				"render_auto": map[string]any{
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
					"scope":   []string{"object"},
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
					"scope":   []string{"object"},
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
			description:   "Target World UMID to teleport",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "timeline_last_seen",
			description:   "Last recorded activity of user viewing timeline",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "blockadelabs",
			description:   "Blockadelabs API key storage attribute",
			options: &entry.AttributeOptions{
				"permissions": entry.PermissionsAttributeOption{
					Read:  string(entry.PermissionAdmin),
					Write: string(entry.PermissionAdmin),
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "leonardo",
			description:   "Leonardo API key storage attribute",
			options: &entry.AttributeOptions{
				"permissions": entry.PermissionsAttributeOption{
					Read:  string(entry.PermissionAdmin),
					Write: string(entry.PermissionAdmin),
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "skybox_ai",
			description:   "Generated skybox storage attribute",
			options: &entry.AttributeOptions{
				"permissions": entry.PermissionsAttributeOption{
					Read:  string(entry.PermissionAdmin) + "+" + string(entry.PermissionUserOwner),
					Write: string(entry.PermissionAdmin) + "+" + string(entry.PermissionUserOwner),
				},
				"posbus_auto": entry.PosBusAutoAttributeOption{
					Scope: []entry.PosBusAutoScopeAttributeOption{entry.UserPosBusAutoScopeAttributeOption},
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "tracker_ai_usage",
			description:   "Track AI usages",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "last_known_position",
			description:   "Last known position for user in the world",
			options:       nil,
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "object_color",
			description:   "Holds the object color",
			options: &entry.AttributeOptions{
				"render_auto": map[string]any{
					"slot_type":    "string",
					"content_type": "string",
				},
			},
		},
		//
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "spatial_audio",
			description:   "Spatial audio",
			options: &entry.AttributeOptions{
				"posbus_auto": entry.PosBusAutoAttributeOption{
					Scope: []entry.PosBusAutoScopeAttributeOption{entry.WorldPosBusAutoScopeAttributeOption},
				},
				"render_auto": entry.RenderAutoAttributeOption{
					SlotType:    entry.SlotTypeAudio,
					ContentType: entry.SlotContentTypeAudio,
					SlotName:    "spatial",
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "soundtrack",
			description:   "Playlist",
			options: &entry.AttributeOptions{
				"posbus_auto": entry.PosBusAutoAttributeOption{
					Scope: []entry.PosBusAutoScopeAttributeOption{entry.ObjectPosBusAutoScopeAttributeOption},
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "object_effect",
			description:   "Visual 3D effect for object",
			options: &entry.AttributeOptions{
				"render_auto": map[string]any{
					"slot_type":    "string",
					"content_type": "string",
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "user_customisable_data",
			description:   "Data for user customisable objects",
			options: &entry.AttributeOptions{
				"render_auto": map[string]string{
					"slot_type":    "texture",
					"content_type": "image",
					"slot_name":    "object_texture",
					"value_field":  "image_hash",
				},
			},
		},
		{
			pluginID:      universe.GetSystemPluginID(),
			attributeName: "shaders",
			description:   "3D shader FX",
			options: &entry.AttributeOptions{
				"posbus_auto": entry.PosBusAutoAttributeOption{
					Scope: []entry.PosBusAutoScopeAttributeOption{entry.WorldPosBusAutoScopeAttributeOption},
				},
			},
		},
		{
			pluginID:      umid.MustParse(textPluginID),
			attributeName: "state",
			description:   "State of the text tile",
			options: &entry.AttributeOptions{
				"render_auto": map[string]any{
					"slot_name":    "description",
					"slot_type":    "object_texture",
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
			pluginID:      umid.MustParse(imagePluginID),
			attributeName: "state",
			description:   "State of the image tile",
			options: &entry.AttributeOptions{
				"render_auto": map[string]any{
					"slot_name":    "object_texture",
					"slot_type":    "texture",
					"content_type": "image",
				},
			},
		},
	}

	for _, item := range items {
		attributeType, err := node.GetAttributeTypes().CreateAttributeType(
			entry.AttributeTypeID{
				PluginID: item.pluginID,
				Name:     item.attributeName,
			},
		)
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
