insert into "user"
(
    user_id,
    user_type_id,
    profile,
    options
)
values
    (
        '00000000-0000-0000-0000-000000000003',
        '00000000-0000-0000-0000-000000000002',
        '{
          "bio": null,
          "name": null,
          "location": null,
          "onboarded": null,
          "avatar_hash": null,
          "profile_link": null
        }'::jsonb,
        null
    );