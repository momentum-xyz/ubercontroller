BEGIN;
DELETE FROM attribute_type
	WHERE plugin_id='3023d874-c186-45bf-a7a8-60e2f57b8877'
      AND attribute_name IN ('bot', 'bot_world')
    ;

DELETE FROM plugin
    WHERE plugin_id='3023d874-c186-45bf-a7a8-60e2f57b8877'
;
COMMIT;
