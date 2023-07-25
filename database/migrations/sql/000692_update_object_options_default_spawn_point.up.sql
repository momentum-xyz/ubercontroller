DO $$
DECLARE
    object_type_uuid UUID := 'a41ee21e-6c56-41b3-81a9-1c86578b6b3c';
    spawn_point JSONB := '{
        "position": { "x": 50, "y": 50, "z": 150 },
        "rotation": {"x": 0, "y": 0, "z": 0}
    }'::jsonb;
BEGIN
    UPDATE "object"
    SET "options" = jsonb_set(coalesce("options", '{}'::jsonb), '{spawn_point}', spawn_point, true)
    WHERE "object_type_id" = object_type_uuid;
END $$;