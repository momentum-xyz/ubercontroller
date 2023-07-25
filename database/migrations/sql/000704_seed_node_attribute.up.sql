insert into node_attribute
(
    plugin_id,
    attribute_name,
    value,
    options
)
values
    (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
        'leonardo',
        '{}'::jsonb,
        null
    ),
    (
        '86dc3ae7-9f3d-42cb-85a3-a71abc3c3cb8',
        'challenge_store',
        '{}'::jsonb,
        null
    ),
    (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
        'tracker_ai_usage',
        '{}'::jsonb,
        null
    ),
    (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
        'node_settings',
        '{
          "name": "dev2-node",
          "umid": "{{NODE_SETTINGS_ID}}",
          "user_id_salt": "{{NODE_SETTINGS_USER_ID_SALT}}",
          "entrance_world": "d83670c7-a120-47a4-892d-f9ec75604f74",
          "guest_user_type": "76802331-37b3-44fa-9010-35008b0cbaec",
          "normal_user_type": "00000000-0000-0000-0000-000000000006"
        }'::jsonb,
        null
    ),
    (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
        'jwt_key',
        '{
          "secret": "{{JWT_KEY_SECRET}}",
          "signature": "{{JWT_KEY_SIGNATURE}}"
        }'::jsonb,
        null
    ),
    (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
        'blockadelabs',
        '{}'::jsonb,
        null
    ),
    (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
        'world_template',
        '{
          "objects": [],
          "random_spaces": [],
          "object_type_id": "a41ee21e-6c56-41b3-81a9-1c86578b6b3c",
          "object_attributes": [
            {
              "value": {
                "lod": [
                  6400,
                  40000,
                  160000
                ]
              },
              "options": null,
              "plugin_id": "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
              "attribute_name": "world_meta"
            },
            {
              "value": {
                "render_hash": "26485e74acb29223ba7a9fa600d36c7f"
              },
              "plugin_id": "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0",
              "attribute_name": "active_skybox"
            }
          ]
        }'::jsonb,
        null
    );