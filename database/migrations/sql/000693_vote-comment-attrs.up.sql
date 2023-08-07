BEGIN;
insert into attribute_type
(
    plugin_id,
    attribute_name,
    description,
    options
)
values
    (
        '{{CORE_PLUGIN_ID}}',
        'vote',
        'Vote',
        '{
          "posbus_auto": {
            "scope": [
              "user", "object"
            ]
          }
        }'::jsonb
    ),
    (
        '{{CORE_PLUGIN_ID}}',
        'comments',
        'Comments',
        '{
          "posbus_auto": {
            "scope": [
              "user", "object"
            ]
          }
        }'::jsonb
    );

COMMIT;
