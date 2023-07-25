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
          "pluginId": "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'a31722a6-26b7-46bc-97f9-435c380c3ca9',
        '{
          "name": "miro",
          "pluginId": "24071066-e8c6-4692-95b5-ae2dc3ed075c",
          "scopeName": "plugin_miro",
          "scriptUrl": "https://dev.odyssey.ninja/plugins/miro/remoteEntry.js"
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
          "scriptUrl": "https://dev.odyssey.ninja/plugins/google-drive/remoteEntry.js"
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
          "scriptUrl": "https://dev.odyssey.ninja/plugins/video/remoteEntry.js"
        }'::jsonb,
        '{}'::jsonb
    );


insert into asset_3d
(
    asset_3d_id,
    meta,
    options
)
values
    (
        '839b21db-52ff-45ce-7484-fd1b59ebb087',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'a1f144de-b21a-d1e9-0635-6eb250927326',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        '313a97cc-fe1b-39bb-56e7-516d213cc23d',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'a55f9ca7-4b45-692e-204f-e37ed9dc3d78',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        '8a7e55f5-934d-8ebf-17bb-39e2d8d9bfa1',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        '46d923ad-21ff-276d-c3c4-ead2212bcb02',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        '97daa12f-9b2e-536d-7851-3b0837175e4c',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'dad4e8a4-cdcc-4174-9d77-f7e849bba352',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        '418c4963-623a-391c-795d-e6080be11899',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        '6e8fec1c-ff95-df66-1375-e312f6447b3d',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'eea924c0-6e33-393f-e06e-e6631e8860e9',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        '5b5bd872-0328-e38c-1b54-bf2bfa70fc85',
        '{
          "type": 2,
          "category": "basic"
        }'::jsonb,
        '{}'::jsonb
    );


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
          "scriptUrl": "https://dev.odyssey.ninja/plugins/video/remoteEntry.js"
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
          "scriptUrl": "https://dev.odyssey.ninja/plugins/google-drive/remoteEntry.js"
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
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
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
            "a31722a6-26b7-46bc-97f9-435c380c3ca9",
            "2a879830-b79e-4c35-accc-05607c51d504"
          ],
          "scopeName": "plugin_miro",
          "scriptUrl": "https://dev.odyssey.ninja/plugins/miro/remoteEntry.js"
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