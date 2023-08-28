BEGIN;

INSERT INTO asset_2d (asset_2d_id, meta)
	VALUES ('7be0964f-df73-4880-91f5-22eef996beef','{
  "name": "object_content",
  "pluginId": "f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0"
}')
ON CONFLICT (asset_2d_id) DO NOTHING;


INSERT INTO attribute_type (plugin_id, attribute_name, description, "options")
	VALUES ('f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0','texture','Object texture','{"cloneable": {}, "render_auto": {"slot_name": "object_texture", "slot_type": "texture", "content_type": "image"}}')
ON CONFLICT (plugin_id,attribute_name) DO NOTHING;

UPDATE attribute_type
	SET "options"='{"render_auto": {"slot_name": "object_image", "slot_type": "texture", "content_type": "image"}}'
	WHERE plugin_id='ff40fbf0-8c22-437d-b27a-0258f99130fe'::uuid AND attribute_name='state';


UPDATE object_attribute oa
SET    plugin_id = 'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
       attribute_name = 'texture'
WHERE  attribute_name = 'state'
       AND plugin_id = 'ff40fbf0-8c22-437d-b27a-0258f99130fe'
       AND oa.object_id IN (SELECT o.object_id
                            FROM   "object" o
                            WHERE  o.asset_3d_id IN (SELECT ad.asset_3d_id
                                                     FROM   asset_3d ad
                                                     WHERE
                                   ad.meta ->> 'category' =
                                   'basic'));

UPDATE object
  SET asset_2d_id = '7be0964f-df73-4880-91f5-22eef996beef'
  WHERE asset_2d_id IN (
    'be0d0ca3-c50b-401a-89d9-0e59fc45c5c2',
    'bda25d5d-2aab-45b4-9e8a-23579514cec1'
    );

UPDATE object
  SET asset_2d_id = NULL
  WHERE asset_2d_id = '7be0964f-df73-4880-91f5-22eef9967999';

UPDATE object
  SET asset_2d_id = '7be0964f-df73-4880-91f5-22eef996beef'
  WHERE object_id IN (
    SELECT object_id
    FROM object_attribute
    WHERE plugin_id = 'ff40fbf0-8c22-437d-b27a-0258f99130fe' AND attribute_name = 'state'
  );

COMMIT;