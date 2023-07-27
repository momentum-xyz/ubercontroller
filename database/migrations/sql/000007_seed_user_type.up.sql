insert into user_type
(
    user_type_id,
    user_type_name,
    description,
    options
)
values
    (
        '00000000-0000-0000-0000-000000000002',
        'Deity',
        'They rule the world',
        '{
          "is_guest": false
        }'::jsonb
    ),
    (
        '00000000-0000-0000-0000-000000000006',
        'User',
        'Momentum user',
        '{
          "is_guest": false
        }'::jsonb
    ),
    (
        '76802331-37b3-44fa-9010-35008b0cbaec',
        'Temporary User',
        'Temporary Momentum user',
        '{
          "is_guest": true
        }'::jsonb
    );