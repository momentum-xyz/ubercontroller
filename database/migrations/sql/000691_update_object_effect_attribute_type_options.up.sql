BEGIN;

UPDATE attribute_type
SET options = options || '{"cloneable": {"use_default": {"value": "transparent"}}}'::jsonb
WHERE attribute_name = 'object_effect' AND NOT jsonb_exists(options, 'cloneable');

COMMIT;
