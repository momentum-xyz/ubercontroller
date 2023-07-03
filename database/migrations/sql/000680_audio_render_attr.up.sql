BEGIN;


INSERT INTO attribute_type(
	plugin_id, attribute_name, description, options)
	VALUES (
        'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0'::uuid,
        'spatial_audio',
        'Spatial audio',
        '{"posbus_auto": {"scope": ["world"]},
          "render_auto": {"slot_type": "audio", "content_type": "audio", "slot_name": "spatial"}
        }'
    )
    ON CONFLICT (plugin_id, attribute_name) DO UPDATE
    SET description=EXCLUDED.description, options=EXCLUDED.options
;

COMMIT;
