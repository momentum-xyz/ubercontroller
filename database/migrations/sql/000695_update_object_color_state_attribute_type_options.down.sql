BEGIN;

UPDATE attribute_type
SET options = options #- '{cloneable}'
WHERE attribute_name = 'object_color' AND jsonb_exists(options, 'cloneable');

UPDATE attribute_type
SET options = options #- '{cloneable}'
WHERE attribute_name = 'state' AND jsonb_exists(options, 'cloneable');

COMMIT;
