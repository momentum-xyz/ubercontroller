BEGIN;

UPDATE attribute_type
SET
options = jsonb_set(options, '{posbus_auto,scope}',
                   (options->'posbus_auto'->'scope' || '["space"]'::jsonb) - 'object')
WHERE options->'posbus_auto'->'scope' ? 'object';

COMMIT;
