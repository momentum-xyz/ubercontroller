insert into asset_3d_user (asset_3d_id, meta, user_id, is_private)
select asset_3d_id,
       json_build_object('name', asset_3d.meta -> 'name', 'preview_hash', asset_3d.meta -> 'preview_hash'),
       '00000000-0000-0000-0000-000000000003',
       false
from asset_3d;

update asset_3d_user
set meta = meta - ARRAY ['preview_hash']
where meta -> 'preview_hash' = 'null';

UPDATE asset_3d
SET meta = meta - ARRAY ['name', 'preview_hash'];
