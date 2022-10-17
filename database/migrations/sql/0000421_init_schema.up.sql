INSERT INTO attribute
(plugin_id, attribute_name, description)
VALUES ('F0F0F0F0-0F0F-4FF0-AF0F-F0F0F0F0F0F0', 'world_settings', '')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET description = excluded.description;

INSERT INTO attribute
(plugin_id, attribute_name, description)
VALUES ('F0F0F0F0-0F0F-4FF0-AF0F-F0F0F0F0F0F0', 'world_meta', '')
ON CONFLICT (plugin_id,attribute_name)
    DO UPDATE SET description = excluded.description;
