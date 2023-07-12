BEGIN;

INSERT INTO attribute_type (plugin_id, attribute_name, description, options)
VALUES ('f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
        'object_effect',
        'Visual 3D effect for object',
        '{
            "render_auto": {
                "slot_type": "string",
                "content_type": "string"
            }
        }')
    ON CONFLICT (plugin_id, attribute_name) DO NOTHING
;
COMMIT;