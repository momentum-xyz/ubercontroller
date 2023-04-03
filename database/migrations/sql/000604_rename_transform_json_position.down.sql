BEGIN;

-- Rename position to location inside the transform jsonb column.
UPDATE "object"
SET "transform" = jsonb_set(
    "transform" #- '{position}', 
    '{location}', "transform" -> 'position'
)
WHERE "transform" ? 'position';

COMMIT;