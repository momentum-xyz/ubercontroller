BEGIN;

UPDATE attribute_type
SET
options = jsonb_set(options, '{posbus_auto,scope}',
                   (options->'posbus_auto'->'scope' || '["object"]'::jsonb) - 'space')
WHERE options->'posbus_auto'->'scope' ? 'space';

COMMIT;
