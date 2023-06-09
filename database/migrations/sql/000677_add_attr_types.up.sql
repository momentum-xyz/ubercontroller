BEGIN;

-- Conflicting with seeder here, when running from scratch. So have to add this.
INSERT INTO plugin(
	plugin_id, meta, options, created_at, updated_at)
	VALUES (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,
        '{"name": "Core"}'::jsonb,
        NULL,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
        )
    ON CONFLICT (plugin_id) DO NOTHING
;

INSERT INTO attribute_type(
	plugin_id, attribute_name, description, options)
	VALUES (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,
        'spatial_audio',
        'Spatial audio',
        '{"posbus_auto": {"scope": ["world"]}}'
    ), (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,
        'soundtrack',
        'Playlist',
        '{"posbus_auto": {"scope": ["object"]}}'
    ), (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,
        'shaders',
        '3D shader FX',
        '{"posbus_auto": {"scope": ["world"]}}'
    )
    ON CONFLICT (plugin_id, attribute_name) DO NOTHING
;

COMMIT;