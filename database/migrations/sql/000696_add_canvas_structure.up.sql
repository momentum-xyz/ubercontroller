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

        INSERT INTO asset_3d (asset_3d_id, meta)
        VALUES
            (
                '2dc7df8e-a34a-829c-e3ca-b73bfe99faf0',
                '{"type": 1, "category": "explorer"}'::jsonb
            )
        ON CONFLICT (asset_3d_id) DO NOTHING;
        RAISE NOTICE 'Inserted or skipped existing asset_3d record with ID: 2dc7df8e-a34a-829c-e3ca-b73bfe99faf0';

        INSERT INTO asset_3d_user (asset_3d_id, user_id, meta, is_private)
        VALUES
            (
                '2dc7df8e-a34a-829c-e3ca-b73bfe99faf0',
                '00000000-0000-0000-0000-000000000003',
                '{"name": "orb"}'::jsonb,
                false
            )
        ON CONFLICT (asset_3d_id, user_id) DO NOTHING;
        RAISE NOTICE 'Inserted or skipped existing asset_3d_user record with ID: 2dc7df8e-a34a-829c-e3ca-b73bfe99faf0';

        IF NOT EXISTS (SELECT 1 FROM object_type WHERE object_type_id = '590028c4-2f9d-4c7e-abc3-791774fbe4c5') THEN
            INSERT INTO object_type (object_type_id, asset_2d_id, asset_3d_id, object_type_name, category_name, description, options)
            VALUES (
                       '590028c4-2f9d-4c7e-abc3-791774fbe4c5',
                       'd768aa3e-ca03-4f5e-b366-780a5361cc02',
                       '2dc7df8e-a34a-829c-e3ca-b73bfe99faf0',
                       'canvas',
                       'Canvas',
                       'Canvas object type',
                       '{"allowed_children": ["3eca3dd5-a2e1-4347-926a-19eab6da54b2"], "child_limit": 42, "child_placement": {
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

        IF NOT EXISTS (SELECT 1 FROM object_type WHERE object_type_id = '3eca3dd5-a2e1-4347-926a-19eab6da54b2') THEN
            INSERT INTO object_type (object_type_id, asset_2d_id, asset_3d_id, object_type_name, category_name, description, options)
            VALUES (
                       '3eca3dd5-a2e1-4347-926a-19eab6da54b2',
                       '8a0f9e8e-b32e-476a-8afe-e0c57260ff20',
                       '5b5bd872-0328-e38c-1b54-bf2bfa70fc85',
                       'canvas_child',
                       'Canvas child',
                       'Canvas child type',
                       '{}'::jsonb
                   );
            RAISE NOTICE 'Inserted new object_type with ID: 3eca3dd5-a2e1-4347-926a-19eab6da54b2';
        END IF;

        WITH updated_options AS (
            UPDATE object_type
                SET options = jsonb_set(options, '{allowed_children}', options->'allowed_children'||'["590028c4-2f9d-4c7e-abc3-791774fbe4c5"]'::jsonb)
                WHERE object_type_id = 'a41ee21e-6c56-41b3-81a9-1c86578b6b3c'
                    AND jsonb_exists(options, 'allowed_children')
                RETURNING object_type_id
        )

        UPDATE object_type
        SET options = options || '{"allowed_children": ["590028c4-2f9d-4c7e-abc3-791774fbe4c5"]}'::jsonb
        WHERE object_type_id = 'a41ee21e-6c56-41b3-81a9-1c86578b6b3c'
          AND NOT object_type_id IN (SELECT object_type_id FROM updated_options);

        /*
        -- Adds a canvas object to new worlds by default.
        WITH object_to_add AS (
            SELECT
                '{"object_id": "b3dbfce9-c635-4506-a823-09954a28dcd1",
                  "object_name": "Canvas",
                  "object_type_id": "590028c4-2f9d-4c7e-abc3-791774fbe4c5",
                  "asset_2d_id": "d768aa3e-ca03-4f5e-b366-780a5361cc02",
                  "asset_3d_id": "2dc7df8e-a34a-829c-e3ca-b73bfe99faf0",
                  "options": {
                    "spawn_point": {
                      "position": {"x": 0, "y": 0, "z": 50},
                      "rotation": {"x": 0, "y": 0, "z": 0}
                    }
                  }}'::jsonb AS obj
        )
        UPDATE public.node_attribute
        SET value = jsonb_set(value, '{objects}', value->'objects' || (SELECT obj FROM object_to_add))
        WHERE attribute_name = 'world_template'
          AND NOT EXISTS (
            SELECT 1
            FROM jsonb_array_elements(value->'objects') AS objects
            WHERE objects->>'object_id' = 'b3dbfce9-c635-4506-a823-09954a28dcd1'
        );
        */

        IF NOT EXISTS (SELECT 1 FROM attribute_type WHERE attribute_name = 'canvas') THEN
            INSERT INTO attribute_type (plugin_id, attribute_name, description, options)
            VALUES (
                       'c782385b-f518-4078-9988-24f356a37c72',
                       'canvas',
                       'Canvas creator configuration',
                       '{"permissions": {"read": "admin", "write": "admin"}}'::jsonb
                   );
            RAISE NOTICE 'Inserted new attribute_type with name: canvas';
        END IF;

        IF NOT EXISTS (SELECT 1 FROM attribute_type WHERE attribute_name = 'canvas_contribution') THEN
            INSERT INTO attribute_type (plugin_id, attribute_name, description, options)
            VALUES (
                       'c782385b-f518-4078-9988-24f356a37c72',
                       'canvas_contribution',
                       'Canvas user contributions',
                       '{"permissions": {"read": "any", "write": "admin"}, "render_auto": {
                         "slot_name": "object_texture",
                         "slot_type": "texture",
                         "value_field": "image_hash",
                         "content_type": "image"
                       }}'::jsonb
                   );
            RAISE NOTICE 'Inserted new attribute_type with name: canvas_contribution';
        END IF;

    END $$;
