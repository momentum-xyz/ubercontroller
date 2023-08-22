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
        'open_ai',
        'Chat GPT settings',
        '{
          "permissions": {
            "read": "admin",
            "write": "admin"
          }
        }'::jsonb
    );

INSERT INTO node_attribute
(
    plugin_id,
    attribute_name,
    value,
    options
)
VALUES
    (
        '{{CORE_PLUGIN_ID}}',
        'open_ai',
        '{
          "api_key": "set_your_api_key_here",
          "temperature": 0.7,
          "max_tokens": 30
        }'::jsonb,
        null
    );