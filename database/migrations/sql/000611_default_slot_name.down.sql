BEGIN;

UPDATE attribute_type
SET
options = jsonb_set(options, '{render_auto,slot_name}', '"Block"')
WHERE options->'render_auto'->>'slot_name' = 'object_texture';


COMMIT;
