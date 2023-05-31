BEGIN;

UPDATE attribute_type
	SET options='{
    "permissions": {"read": "admin", "write": "admin"}
  }'::jsonb
	WHERE plugin_id='86dc3ae7-9f3d-42cb-85a3-a71abc3c3cb8'::uuid 
        AND attribute_name='wallet';

COMMIT;