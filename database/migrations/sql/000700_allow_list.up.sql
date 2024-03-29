BEGIN;

INSERT INTO attribute_type
(
    plugin_id,
    attribute_name,
    description,
    options
)
VALUES
    (
        '{{CORE_PLUGIN_ID}}',
        'hosting_allow_list',
        'Hosting user whitelist',
        '{
          "permissions": {
            "read": "admin",
            "write": "admin"
          }
        }'::jsonb
    );

INSERT INTO attribute_type
(
    plugin_id,
    attribute_name,
    description,
    options
)
VALUES
    (
        '{{CORE_PLUGIN_ID}}',
        'node_private_key',
        'Node private key store',
        '{
          "permissions": {
            "read": "admin",
            "write": "admin"
          }
        }'::jsonb
    );

INSERT INTO attribute_type
(
    plugin_id,
    attribute_name,
    description,
    options
)
VALUES
    (
        '{{CORE_PLUGIN_ID}}',
        'node_public_key',
        'Node public key store',
        '{
          "permissions": {
            "read": "any",
            "write": "admin"
          }
        }'::jsonb
    );

INSERT INTO node_attribute (plugin_id, attribute_name, value)
VALUES (
           '{{CORE_PLUGIN_ID}}',
           'hosting_allow_list',
           '{"users": []}'::jsonb
       );

INSERT INTO node_attribute (plugin_id, attribute_name, value)
VALUES (
           '{{CORE_PLUGIN_ID}}',
           'node_private_key',
           '{{NODE_PRIVATE_KEY}}'::jsonb
       );

INSERT INTO node_attribute (plugin_id, attribute_name, value)
VALUES (
           '{{CORE_PLUGIN_ID}}',
           'node_public_key',
           '{{NODE_PUBLIC_KEY}}'::jsonb
       );

INSERT INTO object_type (object_type_id, asset_2d_id, asset_3d_id, object_type_name, category_name, description, options)
VALUES (
           'e31139ad-ff77-4124-825e-8c83f02b82f4',
           null,
           null,
           'Remote World',
           'Worlds',
           'Type for remotely hosted worlds',
           '{}'::jsonb
       );


COMMIT;