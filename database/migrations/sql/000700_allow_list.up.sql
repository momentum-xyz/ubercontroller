BEGIN;

INSERT INTO node_attribute (plugin_id, attribute_name, value)
VALUES (
           'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
           'hosting_allow_list',
           '[]'::jsonb
       );

INSERT INTO object_type (object_type_id, asset_2d_id, asset_3d_id, object_type_name, category_name, description, options)
VALUES (
           'e31139ad-ff77-4124-825e-8c83f02b82f4',
           null,
           null,
           'Remote World',
           'Worlds',
           'Type for remotely hosted worlds',
           '{}'::jsonb
       );

COMMIT;