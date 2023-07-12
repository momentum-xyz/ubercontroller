BEGIN;

INSERT INTO attribute_type (plugin_id, attribute_name, description, options)
VALUES ('f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0',
        'user_customisable_data',
        'Data for user customisable objects',
        '{
            "render_auto": {
                "slot_type": "texture",
                "content_type": "image",
				"slot_name":    "object_texture",
				"value_field":  "image_hash"
            }
        }')
    ON CONFLICT (plugin_id, attribute_name) DO UPDATE
    SET description=EXCLUDED.description, options=EXCLUDED.options
;
COMMIT;