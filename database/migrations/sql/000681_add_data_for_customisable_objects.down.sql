BEGIN;

DELETE
FROM attribute_type
WHERE attribute_name = 'user_customisable_data'
  AND plugin_id = 'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0';

DELETE
FROM object_type
WHERE object_type_id = '4ed3a5bb-53f8-4511-941b-079029111111';

DELETE
FROM asset_2d
WHERE asset_2d_id = '7be0964f-df73-4880-91f5-22eef996aaaa';

COMMIT;
