BEGIN;

UPDATE object_type
SET options = (options #- '{allowed_subobjects}')::jsonb || '{"allowed_children": ["a41ee21e-6c56-41b3-81a9-1c86578b6b3c"]}'::jsonb
WHERE object_type_id = '00000000-0000-0000-0000-000000000001';

UPDATE object_type
SET options = (options #- '{allowed_subobjects}')::jsonb || '{"allowed_children": ["4ed3a5bb-53f8-4511-941b-079029111111", "4ed3a5bb-53f8-4511-941b-07902982c31c"]}'::jsonb
WHERE object_type_id = 'a41ee21e-6c56-41b3-81a9-1c86578b6b3c';

COMMIT;

