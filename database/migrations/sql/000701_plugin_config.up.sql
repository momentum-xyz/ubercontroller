BEGIN;

INSERT INTO attribute_type
(
    plugin_id,
    attribute_name,
    description,
    options
)
VALUES
    (
        '{{CORE_PLUGIN_ID}}',
        'frontend_plugins',
        'Plugin configuration',
        '{
          "permissions": {
            "read": "any",
            "write": "admin"
          }
        }'::jsonb
    );

COMMIT;