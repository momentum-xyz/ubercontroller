insert into object
(
    object_id,
    object_type_id,
    owner_id,
    parent_id,
    asset_2d_id,
    asset_3d_id,
    options,
    transform
)
values
    (
        '{{NODE_ID}}',
        '00000000-0000-0000-0000-000000000001',
        '00000000-0000-0000-0000-000000000003',
        '{{NODE_ID}}',
        null,
        null,
        null,
        null
    ),
    (
        'd83670c7-a120-47a4-892d-f9ec75604f74',
        'a41ee21e-6c56-41b3-81a9-1c86578b6b3c',
        '00000000-0000-0000-0000-000000000003',
        '{{NODE_ID}}',
        '00000000-0000-0000-0000-000000000008',
        null,
        '{
          "visible": 3,
          "spawn_point": {
            "position": {
              "x": 50,
              "y": 50,
              "z": 150
            },
            "rotation": {
              "x": 0,
              "y": 0,
              "z": 0
            }
          },
          "allowed_subobjects": [
            "0fcd7e48-8d88-41b8-aa74-1630d9dfbe72"
          ],
          "child_placement_backup": {
            "00000000-0000-0000-0000-000000000000": {
              "algo": "circular",
              "options": {
                "R": 330,
                "angle": 0
              }
            },
            "86229140-93a5-4206-ab3b-75713c38f6a6": {
              "algo": "circular",
              "options": {
                "R": 600,
                "angle": 0,
                "Vshift": 300
              }
            }
          }
        }'::jsonb,
        null
    );