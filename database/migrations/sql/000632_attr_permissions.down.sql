
BEGIN;

UPDATE attribute_type
	SET options=NULL
	WHERE plugin_id='86dc3ae7-9f3d-42cb-85a3-a71abc3c3cb8'::uuid 
        AND attribute_name='wallet';

COMMIT;