BEGIN;

INSERT INTO node_attribute (plugin_id, attribute_name, value)
VALUES (
           'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
           'hosting_allow_list',
           '[]'::jsonb
       );

COMMIT;