BEGIN;

UPDATE object_type
SET options = (options #- '{allowed_children}')::jsonb || '{"allowed_subobjects": []}'::jsonb
WHERE object_type_id = '00000000-0000-0000-0000-000000000001';

UPDATE object_type
SET options = (options #- '{allowed_children}')::jsonb || '{"allowed_subobjects": []}'::jsonb
WHERE object_type_id = 'a41ee21e-6c56-41b3-81a9-1c86578b6b3c';

COMMIT;
