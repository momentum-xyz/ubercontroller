BEGIN;

DELETE FROM object_type WHERE object_type_id = 'e31139ad-ff77-4124-825e-8c83f02b82f4';

DELETE FROM node_attribute WHERE plugin_id = '{{CORE_PLUGIN_ID}}' AND attribute_name = 'hosting_allow_list';

DELETE FROM node_attribute WHERE plugin_id = '{{CORE_PLUGIN_ID}}' AND attribute_name = 'node_key';

DELETE FROM attribute_type WHERE plugin_id = '{{CORE_PLUGIN_ID}}' AND attribute_name = 'hosting_allow_list';

DELETE FROM attribute_type WHERE plugin_id = '{{CORE_PLUGIN_ID}}' AND attribute_name = 'node_key';

COMMIT;