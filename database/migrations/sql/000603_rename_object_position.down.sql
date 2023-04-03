BEGIN;

ALTER TABLE "object"
  RENAME COLUMN "transform" TO "position";

DROP FUNCTION getobjectancestorstableadminrole;
DROP VIEW objects_plus_level_isadmin;

CREATE OR REPLACE VIEW objects_plus_level_isadmin
            (object_id, object_type_id, owner_id, parent_id, asset_2d_id, asset_3d_id,
            options, position, updated_at, created_at, level, is_admin)
AS
SELECT object.object_id,
       object.object_type_id,
       object.owner_id,
       object.parent_id,
       object.asset_2d_id,
       object.asset_3d_id,
       object.options,
       object."position",
       object.updated_at,
       object.created_at,
       NULL::integer AS level,
       NULL::boolean AS is_admin
FROM object
WHERE false;

create function getobjectancestorstableadminrole(t_pid uuid, t_uid uuid) returns SETOF objects_plus_level_isadmin
    language sql
as
$$
with recursive cte as (
    select *,
           0 as level
    from object
    where object_id = t_pid
    union all
    select o.*,
           level + 1
    from object o
             inner join cte on cte.parent_id = o.object_id
    where cte.object_id != cte.parent_id
)
select cte.*,
       bool_or(COALESCE((user_object.value->>'role')::varchar = 'admin', false)) OVER (
           order by level DESC
           ) as is_admin
from cte
         left join user_object on cte.object_id = user_object.object_id
    and user_object.user_id = t_uid;
$$;


COMMIT;

COMMIT;
