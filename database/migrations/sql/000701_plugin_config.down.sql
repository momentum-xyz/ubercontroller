BEGIN;

DELETE FROM node_attribute WHERE plugin_id = '{{CORE_PLUGIN_ID}}' AND attribute_name = 'frontend_plugins';

COMMIT;