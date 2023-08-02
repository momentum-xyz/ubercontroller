BEGIN;

UPDATE attribute_type
SET options = options || '{"cloneable": {}}'::jsonb
WHERE attribute_name = 'object_color' AND NOT jsonb_exists(options, 'cloneable');

UPDATE attribute_type
SET options = options || '{"cloneable": {}}'::jsonb
WHERE attribute_name = 'state' AND plugin_id = 'ff40fbf0-8c22-437d-b27a-0258f99130fe' AND NOT jsonb_exists(options, 'cloneable');

COMMIT;
