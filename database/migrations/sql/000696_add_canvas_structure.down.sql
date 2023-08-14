DO $$
BEGIN

    IF EXISTS (SELECT 1
               FROM object_type
               WHERE object_type_id = '590028c4-2f9d-4c7e-abc3-791774fbe4c5') THEN

DELETE FROM object_type
WHERE object_type_id = '590028c4-2f9d-4c7e-abc3-791774fbe4c5';

END IF;

    IF EXISTS (SELECT 1
               FROM asset_2d
               WHERE asset_2d_id IN ('d768aa3e-ca03-4f5e-b366-780a5361cc02',
                                     '8a0f9e8e-b32e-476a-8afe-e0c57260ff20')) THEN

DELETE FROM asset_2d
WHERE asset_2d_id IN ('d768aa3e-ca03-4f5e-b366-780a5361cc02',
                      '8a0f9e8e-b32e-476a-8afe-e0c57260ff20');

END IF;

    IF EXISTS (SELECT 1
               FROM plugin
               WHERE plugin_id = 'c782385b-f518-4078-9988-24f356a37c72') THEN

DELETE FROM plugin
WHERE plugin_id = 'c782385b-f518-4078-9988-24f356a37c72';

END IF;

EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE 'An error occurred: %', SQLERRM;
ROLLBACK;

END $$;