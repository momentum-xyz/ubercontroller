BEGIN;

UPDATE attribute_type
SET options = options || '{"cloneable": {}}'::jsonb
WHERE attribute_name = 'object_color' AND NOT jsonb_exists(options, 'cloneable');

UPDATE attribute_type
SET options = options || '{"cloneable": {}}'::jsonb
WHERE attribute_name = 'state' AND NOT jsonb_exists(options, 'cloneable');

COMMIT;
