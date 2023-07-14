BEGIN;

-- Conflicting with seeder here, when running from scratch. So have to add this.
INSERT INTO plugin(
	plugin_id, meta, options, created_at, updated_at)
	VALUES (
        '3023d874-c186-45bf-a7a8-60e2f57b8877'::uuid,
        '{"name": "SDK"}'::jsonb,
        NULL,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
        )
    ON CONFLICT (plugin_id) DO NOTHING
;

INSERT INTO attribute_type(
	plugin_id, attribute_name, description, options)
	VALUES (
        '3023d874-c186-45bf-a7a8-60e2f57b8877'::uuid,
        'bot',
        'Bot state storage',
        ''
    ), (
        '3023d874-c186-45bf-a7a8-60e2f57b8877'::uuid,
        'bot_world',
        'Bot state broadcasted to the world',
        '{"posbus_auto": {"scope": ["world"]}}'
    )
    ON CONFLICT (plugin_id, attribute_name) DO NOTHING
;

COMMIT;

