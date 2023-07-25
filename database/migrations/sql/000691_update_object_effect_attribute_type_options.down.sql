BEGIN;

UPDATE attribute_type
SET options = options #- '{cloneable}'
WHERE attribute_name = 'object_effect' AND jsonb_exists(options, 'cloneable');

COMMIT;
