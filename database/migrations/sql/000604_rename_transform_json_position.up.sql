BEGIN;

-- Rename location to position inside the transform jsonb column.
UPDATE "object"
SET "transform" = jsonb_set(
    "transform" #- '{location}', 
    '{position}', "transform" -> 'location'
)
WHERE "transform" ? 'location';

COMMIT;