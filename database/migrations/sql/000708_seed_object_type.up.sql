insert into object_type
(
    object_type_id,
    asset_2d_id,
    asset_3d_id,
    object_type_name,
    category_name,
    description,
    options
)
values
    (
        '4ed3a5bb-53f8-4511-941b-079029111111',
        '7be0964f-df73-4880-91f5-22eef996aaaa',
        null,
        'Custom claimable objects',
        'Custom claimable',
        'Custom placed objects which can be claimed',
        '{
          "visible": 3
        }'::jsonb
    ),
    (
        '75b56447-c4f1-4020-b8fc-d68704a11d65',
        null,
        null,
        'Generic Space',
        'Generic Spaces',
        '',
        '{
          "subs": null,
          "minimap": true,
          "private": false,
          "visible": 3,
          "editable": true
        }'::jsonb
    ),
    (
        'f9607e55-63e8-4cb1-ae47-66395199975d',
        null,
        null,
        'morgue',
        'Morgues',
        'morgue',
        '{
          "subs": null,
          "minimap": true,
          "private": false,
          "visible": 0,
          "editable": true
        }'::jsonb
    ),
    (
        '00000000-0000-0000-0000-000000000001',
        null,
        null,
        'Node',
        'Nodes',
        'Root of it all',
        '{
          "subs": {
            "asset2d_plugins": [
              "24071066-e8c6-4692-95b5-ae2dc3ed075c"
            ]
          },
          "minimap": true,
          "private": false,
          "visible": 1,
          "editable": true,
          "allowed_children": [
            "a41ee21e-6c56-41b3-81a9-1c86578b6b3c"
          ]
        }'::jsonb
    ),
    (
        '88415343-90db-4d23-a9e7-79a11aaaaf04',
        null,
        null,
        'anchor',
        'Anchors',
        '',
        '{
          "subs": null,
          "minimap": true,
          "private": false,
          "visible": 3,
          "editable": true,
          "child_placement": {
            "00000000-0000-0000-0000-000000000000": {
              "algo": "circular",
              "options": {
                "R": 55,
                "angle": 0
              }
            }
          }
        }'::jsonb
    ),
    (
        'a41ee21e-6c56-41b3-81a9-1c86578b6b3c',
        null,
        null,
        'World',
        'Worlds',
        'World Type',
        '{
          "subs": null,
          "minimap": true,
          "private": false,
          "visible": 3,
          "editable": true,
          "allowed_children": [
            "4ed3a5bb-53f8-4511-941b-079029111111",
            "4ed3a5bb-53f8-4511-941b-07902982c31c"
          ]
        }'::jsonb
    ),
    (
        '4ed3a5bb-53f8-4511-941b-07902982c31c',
        null,
        null,
        'Custom objects',
        'Custom',
        'Custom placed objects',
        '{
          "subs": null,
          "minimap": true,
          "private": false,
          "visible": 3,
          "editable": true
        }'::jsonb
    ),
    (
        '69d8ae40-df9b-4fc8-af95-32b736d2bbcd',
        null,
        null,
        'Service Space',
        'Service Spaces',
        '',
        '{
          "subs": null,
          "minimap": true,
          "private": false,
          "visible": 0,
          "editable": true
        }'::jsonb
    );