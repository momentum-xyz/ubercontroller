BEGIN;

UPDATE attribute_type
SET options = options #- '{cloneable}'
WHERE attribute_name = 'object_color' AND jsonb_exists(options, 'cloneable');

UPDATE attribute_type
SET options = options #- '{cloneable}'
WHERE attribute_name = 'state' AND jsonb_exists(options, 'cloneable') AND plugin_id = 'ff40fbf0-8c22-437d-b27a-0258f99130fe';

COMMIT;
