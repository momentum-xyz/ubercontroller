BEGIN;

DELETE FROM object_type WHERE object_type_id = 'e31139ad-ff77-4124-825e-8c83f02b82f4';

DELETE FROM node_attribute WHERE plugin_id = 'f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0' AND attribute_name = 'hosting_allow_list';

COMMIT;