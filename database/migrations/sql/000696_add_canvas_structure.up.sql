DO $$
    BEGIN

        IF NOT EXISTS (SELECT 1 FROM plugin WHERE plugin_id = 'c782385b-f518-4078-9988-24f356a37c72') THEN
            INSERT INTO plugin (plugin_id, meta)
            VALUES ('c782385b-f518-4078-9988-24f356a37c72', '{"name": "Canvas creator"}'::jsonb);
            RAISE NOTICE 'Inserted new plugin with ID: c782385b-f518-4078-9988-24f356a37c72';
        END IF;

        INSERT INTO asset_2d (asset_2d_id, meta)
        VALUES
            (
                'd768aa3e-ca03-4f5e-b366-780a5361cc02',
                '{"name": "Canvas object UI", "pluginId": "c782385b-f518-4078-9988-24f356a37c72"}'::jsonb
            )
        ON CONFLICT (asset_2d_id) DO NOTHING;
        RAISE NOTICE 'Inserted or skipped existing asset_2d record with ID: d768aa3e-ca03-4f5e-b366-780a5361cc02';

        INSERT INTO asset_2d (asset_2d_id, meta)
        VALUES
            (
                '8a0f9e8e-b32e-476a-8afe-e0c57260ff20',
                '{"name": "Canvas child objects", "pluginId": "c782385b-f518-4078-9988-24f356a37c72"}'::jsonb
            )
        ON CONFLICT (asset_2d_id) DO NOTHING;
        RAISE NOTICE 'Inserted or skipped existing asset_2d record with ID: 8a0f9e8e-b32e-476a-8afe-e0c57260ff20';

        IF NOT EXISTS (SELECT 1 FROM object_type WHERE object_type_id = '590028c4-2f9d-4c7e-abc3-791774fbe4c5') THEN
            INSERT INTO object_type (object_type_id, asset_2d_id, asset_3d_id, object_type_name, category_name, description, options)
            VALUES (
                       '590028c4-2f9d-4c7e-abc3-791774fbe4c5',
                       'd768aa3e-ca03-4f5e-b366-780a5361cc02',
                       '2dc7df8e-a34a-829c-e3ca-b73bfe99faf0',
                       'canvas',
                       'Canvas',
                       'Canvas object type',
                       '{"child_limit": 42, "child_placement": {
                         "00000000-0000-0000-0000-000000000000": {
                           "algo": "helix",
                           "options": {
                             "angle": 7.2,
                             "helixVshift": 15,
                             "spiralScale": 50
                           }
                         }
                       }}'::jsonb
                   );
            RAISE NOTICE 'Inserted new object_type with ID: 590028c4-2f9d-4c7e-abc3-791774fbe4c5';
        END IF;
END $$;
