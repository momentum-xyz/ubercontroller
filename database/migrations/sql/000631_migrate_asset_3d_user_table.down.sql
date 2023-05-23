UPDATE asset_3d
set meta = jsonb_set(asset_3d.meta, '{name}', asset_3d_user.meta->'name')
from asset_3d_user
where asset_3d_user.asset_3d_id = asset_3d.asset_3d_id;


UPDATE asset_3d
set meta = jsonb_set(asset_3d.meta, '{preview_hash}', asset_3d_user.meta->'preview_hash')
from asset_3d_user
where asset_3d_user.asset_3d_id = asset_3d.asset_3d_id and asset_3d_user.meta->'preview_hash' != 'null';
