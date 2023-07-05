BEGIN;


INSERT INTO asset_2d (asset_2d_id, options, created_at, meta, updated_at)
VALUES ('7be0964f-df73-4880-91f5-22eef996aaaa',
        '{}',
        NOW(),
        '{
            "name": "claimable",
            "pluginId": "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0"
        }',
        NOW());

INSERT INTO object_type (object_type_id, asset_2d_id, asset_3d_id, object_type_name, category_name, description,
                         options, created_at, updated_at)
VALUES ('4ed3a5bb-53f8-4511-941b-079029111111',
        '7be0964f-df73-4880-91f5-22eef996aaaa',
        null,
        'Custom claimable objects',
        'Custom claimable',
        'Custom placed objects which can be claimed',
        '{
            "visible": 3
        }',
        NOW(),
        NOW());


INSERT INTO attribute_type (plugin_id, attribute_name, description, options)
VALUES ('f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
        'user_customisable_data',
        'Data for user customisable objects',
        '{
            "render_auto": {
                "slot_type": "texture",
                "content_type": "image"
            }
        }');

COMMIT;
