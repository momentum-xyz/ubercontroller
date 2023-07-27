insert into plugin
(
    plugin_id,
    meta,
    options
)
values
    (
        '308fdacc-8c2d-40dc-bd5f-d1549e3e03ba',
        '{
          "name": "Video",
          "assets2d": [
            "bda25d5d-2aab-45b4-9e8a-23579514cec1"
          ],
          "scopeName": "plugin_video",
          "scriptUrl": "{{BASE_DOMAIN}}/plugins/video/remoteEntry.js"
        }'::jsonb,
        'null'
    ),
    (
        'ff40fbf0-8c22-437d-b27a-0258f99130fe',
        '{
          "name": "Image"
        }'::jsonb,
        'null'
    ),
    (
        'fc9f2eb7-590a-4a1a-ac75-cd3bfeef28b2',
        '{
          "name": "Text"
        }'::jsonb,
        'null'
    ),
    (
        'c3f89640-e0f0-4536-ae0d-8fc8a75ec0cd',
        '{
          "name": "Google Drive",
          "assets2d": [
            "c601404b-61a2-47d5-a5c7-f3c704a8bf58"
          ],
          "scopeName": "plugin_google_drive",
          "scriptUrl": "{{BASE_DOMAIN}}/plugins/google-drive/remoteEntry.js"
        }'::jsonb,
        'null'
    ),
    (
        '220578c8-5fec-42c8-ade1-14d970e714bd',
        '{
          "name": "Odyssey hackaton"
        }'::jsonb,
        'null'
    ),
    (
        '3253d616-215f-47a9-ba9d-93185eb3e6b5',
        '{
          "name": "High five"
        }'::jsonb,
        'null'
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        '{
          "name": "Core"
        }'::jsonb,
        'null'
    ),
    (
        '2b92edbc-5ef5-4028-89a6-d510f8583887',
        '{
          "name": "Event Calendar",
          "description": "Event calendar plugin"
        }'::jsonb,
        'null'
    ),
    (
        '24071066-e8c6-4692-95b5-ae2dc3ed075c',
        '{
          "name": "Miro",
          "assets2d": [
            "a31722a6-26b7-46bc-97f9-435c380c3ca9"
          ],
          "scopeName": "plugin_miro",
          "scriptUrl": "{{BASE_DOMAIN}}/plugins/miro/remoteEntry.js"
        }'::jsonb,
        'null'
    ),
    (
        '86dc3ae7-9f3d-42cb-85a3-a71abc3c3cb8',
        '{
          "name": "Kusama"
        }'::jsonb,
        'null'
    ),
    (
        '3023d874-c186-45bf-a7a8-60e2f57b8877',
        '{
          "name": "SDK"
        }'::jsonb,
        'null'
    );