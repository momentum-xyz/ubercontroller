BEGIN;
INSERT INTO attribute_type(
	plugin_id, attribute_name, description, options)
	VALUES (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,
        'timeline_last_seen',
        'Last recorded activity of user viewing timeline',
        NULL
    )
    ON CONFLICT (plugin_id, attribute_name) DO NOTHING
;
COMMIT;
