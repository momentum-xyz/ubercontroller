BEGIN;

UPDATE attribute_type
	SET "options"='{"cloneable": {}, "render_auto": {"slot_name": "object_texture", "slot_type": "texture", "content_type": "image"}}'
	WHERE plugin_id='ff40fbf0-8c22-437d-b27a-0258f99130fe'::uuid AND attribute_name='state';

UPDATE object_attribute oa
SET    plugin_id = 'ff40fbf0-8c22-437d-b27a-0258f99130fe',
       attribute_name = 'state'
WHERE  attribute_name = 'texture'
AND NOT EXISTS (
	SELECT 1
	FROM object_attribute
	WHERE object_id = oa.object_id AND plugin_id = 'ff40fbf0-8c22-437d-b27a-0258f99130fe' AND attribute_name = 'state'
);

UPDATE object
	SET asset_2d_id = '7be0964f-df73-4880-91f5-22eef9967999'
	WHERE object_id IN
	(
		SELECT object_id
		FROM object_attribute
		WHERE plugin_id = 'ff40fbf0-8c22-437d-b27a-0258f99130fe' AND attribute_name = 'state'
	);

UPDATE object
	SET asset_2d_id = 'be0d0ca3-c50b-401a-89d9-0e59fc45c5c2'
	WHERE object_id IN
	(
		SELECT object_id
		FROM object_attribute
		WHERE plugin_id = 'fc9f2eb7-590a-4a1a-ac75-cd3bfeef28b2' AND attribute_name = 'state'
	);

UPDATE object
	SET asset_2d_id = 'bda25d5d-2aab-45b4-9e8a-23579514cec1'
	WHERE object_id IN
	(
		SELECT object_id
		FROM object_attribute
		WHERE plugin_id = '308fdacc-8c2d-40dc-bd5f-d1549e3e03ba' AND attribute_name = 'state'
	);
	

DELETE FROM attribute_type 
WHERE plugin_id = 'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0' AND attribute_name = 'texture';

DELETE FROM asset_2d
WHERE asset_2d_id = '7be0964f-df73-4880-91f5-22eef996beef';

COMMIT;