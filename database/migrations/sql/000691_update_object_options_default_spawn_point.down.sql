DO $$
DECLARE
    object_type_uuid UUID := 'a41ee21e-6c56-41b3-81a9-1c86578b6b3c';
BEGIN
    UPDATE "object"
    SET "options" = "options" - 'spawn_point'
    WHERE "object_type_id" = object_type_uuid AND "options" ? 'spawn_point';
END $$;