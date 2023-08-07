BEGIN;
DELETE FROM attribute_type
WHERE plugin_id='{{CORE_PLUGIN_ID}}'
  AND attribute_name IN ('vote', 'comments')
;
COMMIT;
