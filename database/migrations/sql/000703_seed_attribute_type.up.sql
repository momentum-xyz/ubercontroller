insert into attribute_type
(
    plugin_id,
    attribute_name,
    description,
    options
)
values
    (
        '{{CORE_PLUGIN_ID}}',
        'magic_links',
        '',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'events',
        'Object events',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'active_skybox',
        'Holds skybox data such as texture',
        '{
          "render_auto": {
            "slot_name": "skybox_custom",
            "slot_type": "texture",
            "content_type": "image"
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'timeline_last_seen',
        'Last recorded activity of user viewing timeline',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'world_meta',
        'Holds world metadata and decorations',
        null
    ),
    (
        '220578c8-5fec-42c8-ade1-14d970e714bd',
        'solution',
        'solution for object',
        '{
          "render_type": "texture",
          "content_type": "text"
        }'::jsonb
    ),
    (
        '220578c8-5fec-42c8-ade1-14d970e714bd',
        'poster',
        'Poster for object',
        '{
          "render_type": "texture"
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'world_settings',
        '',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'last_known_position',
        'Last known position for user in the world',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'skybox_list',
        'Holds initial list of skyboxes',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'VoiceChatAction',
        'Voice chat user actions',
        '{
          "posbus_auto": {
            "scope": [
              "object"
            ],
            "topic": "voice-chat-action",
            "send_to": 1
          }
        }'::jsonb
    ),
    (
        '3023d874-c186-45bf-a7a8-60e2f57b8877',
        'bot',
        'Bot state storage',
        '{}'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'high_five',
        'high fives',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'description',
        'description',
        '{
          "render_type": "texture"
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'name',
        'Object name',
        '{
          "render_auto": {
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
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'blockadelabs',
        'Blockadelabs API key storage attribute',
        '{
          "permissions": {
            "read": "admin",
            "write": "admin"
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'staking',
        'Odyssey staking information',
        null
    ),
    (
        'fc9f2eb7-590a-4a1a-ac75-cd3bfeef28b2',
        'state',
        'State of the text tile',
        '{
          "render_auto": {
            "slot_name": "description",
            "slot_type": "object_texture",
            "value_field": "value",
            "content_type": "text",
            "text_render_template": {
              "x": 0,
              "y": 0,
              "sub": [
                {
                  "x": 4,
                  "y": 4,
                  "text": {
                    "padX": 4,
                    "padY": 4,
                    "wrap": true,
                    "align": 0,
                    "alignH": "left",
                    "alignV": "top",
                    "string": "%TEXT%",
                    "fontfile": "IBMPlexSans-SemiBold",
                    "fontsize": 17,
                    "fontcolor": [
                      0,
                      255,
                      255
                    ]
                  },
                  "color": [
                    0,
                    255,
                    255
                  ],
                  "width": 1012,
                  "height": 504,
                  "thickness": 1,
                  "background": [
                    0,
                    40,
                    0
                  ]
                }
              ],
              "color": [
                20,
                20,
                20
              ],
              "width": 1024,
              "height": 512,
              "thickness": 4,
              "background": [
                10,
                10,
                10
              ]
            }
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'spatial_audio',
        'Spatial audio',
        '{
          "posbus_auto": {
            "scope": [
              "world"
            ],
            "topic": ""
          },
          "render_auto": {
            "slot_name": "spatial",
            "slot_type": "audio",
            "value_field": "",
            "content_type": "audio",
            "text_render_template": ""
          }
        }'::jsonb
    ),
    (
        '220578c8-5fec-42c8-ade1-14d970e714bd',
        'video',
        'Video for object',
        '{
          "render_type": "texture",
          "content_type": "video"
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'emoji',
        '',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'soundtrack',
        'Playlist',
        '{
          "posbus_auto": {
            "scope": [
              "object"
            ],
            "topic": ""
          }
        }'::jsonb
    ),
    (
        'c3f89640-e0f0-4536-ae0d-8fc8a75ec0cd',
        'config',
        'Google Drive configuration',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'object_effect',
        'Visual 3D effect for object',
        '{
          "render_auto": {
            "slot_type": "string",
            "content_type": "string"
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'world_template',
        'Basic template settings for any new world',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'shaders',
        '3D shader FX',
        '{
          "posbus_auto": {
            "scope": [
              "world"
            ],
            "topic": ""
          }
        }'::jsonb
    ),
    (
        'c3f89640-e0f0-4536-ae0d-8fc8a75ec0cd',
        'state',
        'Google Drive state',
        '{
          "posbus_auto": {
            "scope": [
              "object"
            ],
            "send_to": 1
          }
        }'::jsonb
    ),
    (
        '24071066-e8c6-4692-95b5-ae2dc3ed075c',
        'state',
        'Miro state',
        '{
          "posbus_auto": {
            "scope": [
              "object"
            ],
            "send_to": 1
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'VoiceChatUser',
        'Voice chat users',
        '{
          "posbus_auto": {
            "scope": [
              "object"
            ],
            "topic": "voice-chat-user",
            "send_to": 1
          }
        }'::jsonb
    ),
    (
        '220578c8-5fec-42c8-ade1-14d970e714bd',
        'tile',
        'tile for object',
        '{
          "render_type": "texture"
        }'::jsonb
    ),
    (
        '3253d616-215f-47a9-ba9d-93185eb3e6b5',
        'count',
        'High5s given',
        null
    ),
    (
        '86dc3ae7-9f3d-42cb-85a3-a71abc3c3cb8',
        'challenge_store',
        'auth challenge store',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'emojis',
        '',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'user_customisable_data',
        'Data for user customisable objects',
        '{
          "render_auto": {
            "slot_name": "object_texture",
            "slot_type": "texture",
            "value_field": "image_hash",
            "content_type": "image"
          }
        }'::jsonb
    ),
    (
        '308fdacc-8c2d-40dc-bd5f-d1549e3e03ba',
        'state',
        'State of the video tile',
        '{
          "render_auto": {
            "slot_name": "object_texture",
            "slot_type": "texture",
            "value_field": "value",
            "content_type": "video"
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'screenshare',
        'Odyssey screenshare state',
        '{
          "posbus_auto": {
            "scope": [
              "object"
            ],
            "topic": "screenshare-action",
            "send_to": 1
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'leonardo',
        'Leonardo API key storage attribute',
        '{
          "permissions": {
            "read": "admin",
            "write": "admin"
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'skybox_ai',
        'Generated skybox storage attribute',
        '{
          "permissions": {
            "read": "admin+user_owner",
            "write": "admin+user_owner"
          },
          "posbus_auto": {
            "scope": [
              "user"
            ],
            "topic": ""
          }
        }'::jsonb
    ),
    (
        '220578c8-5fec-42c8-ade1-14d970e714bd',
        'problem',
        'Problem for object',
        '{
          "render_type": "texture",
          "content_type": "text"
        }'::jsonb
    ),
    (
        '220578c8-5fec-42c8-ade1-14d970e714bd',
        'third',
        'Third screen for object',
        '{
          "render_type": "texture"
        }'::jsonb
    ),
    (
        '220578c8-5fec-42c8-ade1-14d970e714bd',
        'meme',
        'Meme for object',
        '{
          "render_type": "texture"
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'world_avatar',
        '',
        '{
          "render_auto": {
            "slot_type": "texture",
            "content_type": "image"
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'website_link',
        '',
        null
    ),
    (
        'ff40fbf0-8c22-437d-b27a-0258f99130fe',
        'state',
        'State of the image tile',
        '{
          "render_auto": {
            "slot_name": "object_texture",
            "slot_type": "texture",
            "content_type": "image"
          }
        }'::jsonb
    ),
    (
        '24071066-e8c6-4692-95b5-ae2dc3ed075c',
        'config',
        'Miro configuration',
        null
    ),
    (
        '86dc3ae7-9f3d-42cb-85a3-a71abc3c3cb8',
        'wallet',
        'Kusama/Substrate wallet',
        '{
          "permissions": {
            "read": "admin+user_owner",
            "write": "admin"
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'node_settings',
        '',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'jwt_key',
        '',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'tracker_ai_usage',
        'Track AI usages',
        null
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'object_color',
        'Holds the object color',
        '{
          "render_auto": {
            "slot_type": "string",
            "content_type": "string"
          }
        }'::jsonb
    ),
    (
        '3023d874-c186-45bf-a7a8-60e2f57b8877',
        'bot_world',
        'Bot state broadcasted to the world',
        '{
          "posbus_auto": {
            "scope": [
              "world"
            ],
            "topic": ""
          }
        }'::jsonb
    );