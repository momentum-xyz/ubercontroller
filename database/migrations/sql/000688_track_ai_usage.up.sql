BEGIN;

INSERT INTO attribute_type(
	plugin_id, attribute_name, description, options)
	VALUES (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,
        'tracker_ai_usage',
        'Track AI usage',
        '{}'
    )
    ON CONFLICT (plugin_id, attribute_name) DO NOTHING
;

INSERT INTO node_attribute(plugin_id, attribute_name, value)
VALUES ('f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,'tracker_ai_usage','{}'::jsonb)
    ON CONFLICT (plugin_id, attribute_name) DO NOTHING
;
COMMIT;