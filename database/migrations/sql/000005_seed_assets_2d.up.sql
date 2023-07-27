insert into asset_2d
(
    asset_2d_id,
    meta,
    options
)
values
    (
        '7be0964f-df73-4880-91f5-22eef996aaaa',
        '{
          "name": "claimable",
          "pluginId": "{{CORE_PLUGIN_ID}}"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'a31722a6-26b7-46bc-97f9-435c380c3ca9',
        '{
          "name": "miro",
          "pluginId": "24071066-e8c6-4692-95b5-ae2dc3ed075c",
          "scopeName": "plugin_miro",
          "scriptUrl": "{{BASE_DOMAIN}}/plugins/miro/remoteEntry.js"
        }'::jsonb,
        '{
          "exact": true,
          "subPath": "miro",
          "iconName": "miro"
        }'::jsonb
    ),
    (
        'c601404b-61a2-47d5-a5c7-f3c704a8bf58',
        '{
          "name": "google Drive",
          "pluginId": "c3f89640-e0f0-4536-ae0d-8fc8a75ec0cd",
          "scopeName": "plugin_google_drive",
          "scriptUrl": "{{BASE_DOMAIN}}/plugins/google-drive/remoteEntry.js"
        }'::jsonb,
        '{
          "exact": true,
          "iconName": "drive"
        }'::jsonb
    ),
    (
        '7be0964f-df73-4880-91f5-22eef9967999',
        '{
          "name": "image",
          "pluginId": "ff40fbf0-8c22-437d-b27a-0258f99130fe"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'be0d0ca3-c50b-401a-89d9-0e59fc45c5c2',
        '{
          "name": "text",
          "pluginId": "fc9f2eb7-590a-4a1a-ac75-cd3bfeef28b2"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'bda25d5d-2aab-45b4-9e8a-23579514cec1',
        '{
          "name": "video",
          "pluginId": "308fdacc-8c2d-40dc-bd5f-d1549e3e03ba",
          "scopeName": "plugin_video",
          "scriptUrl": "{{BASE_DOMAIN}}/plugins/video/remoteEntry.js"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        '00000000-0000-0000-0000-000000000004',
        '{
          "name": ""
        }'::jsonb,
        '{}'
    ),
    (
        '00000000-0000-0000-0000-000000000008',
        '{
          "name": ""
        }'::jsonb,
        '{}'
    );