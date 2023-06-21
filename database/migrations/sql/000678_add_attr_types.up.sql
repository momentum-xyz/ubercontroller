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
        'blockadelabs',
        'Blockadelabs API key storage attribute',
        '{"permissions": {"read": "admin", "write": "admin"}}'
    ), (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,
        'skybox_ai',
        'Generated skybox storage attribute',
        '{"permissions": {"read": "admin+user_owner", "write": "admin+user_owner"}, "posbus_auto": {"scope": ["user"]}}'
    )
    ON CONFLICT (plugin_id, attribute_name) DO NOTHING
;

INSERT INTO node_attribute(
    plugin_id, attribute_name, value, options)
VALUES (
           'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,
           'blockadelabs',
           '{}',
           null
       )
    ON CONFLICT (plugin_id, attribute_name) DO NOTHING
;

COMMIT;