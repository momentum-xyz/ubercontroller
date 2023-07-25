--
-- Name: getindirectobjectadmins(uuid); Type: FUNCTION;
--

CREATE FUNCTION getindirectobjectadmins(t_pid uuid) RETURNS TABLE(userid uuid)
    LANGUAGE sql
AS $$
WITH RECURSIVE
    objects_cte as (
        select object_id, parent_id,owner_id
        from object
        where object_id = t_pid
        union all
        select o.object_id,
               o.parent_id,o.owner_id
        from object o
                 inner join objects_cte on o.object_id = objects_cte.parent_id
        where  objects_cte.object_id != objects_cte.parent_id
    ),
    object_users_cte AS (
        SELECT user_id, objects_cte.object_id AS object_id
        FROM user_object
                 INNER JOIN objects_cte ON user_object.object_id = objects_cte.object_id
        WHERE user_object.value->>'role' = 'admin'
        union all
        SELECT owner_id as user_id, object_id from objects_cte
    ),
    users_cte as (
        select user_id
        from object_users_cte
        union all
        select u.user_id as user_id
        from user_membership u
                 inner join users_cte on u.member_of = users_cte.user_id
    )
SELECT distinct(users_cte.user_id)
FROM users_cte;
$$;


--
-- Name: getobjectancestorsids(uuid, integer); Type: FUNCTION;
--

CREATE FUNCTION getobjectancestorsids(t_pid uuid, t_lev integer) RETURNS TABLE(id uuid, parent_id uuid, level integer)
    LANGUAGE sql
AS $$
with recursive cte as (
    select object_id,
           parent_id,
           0 as level
    from object
    where object_id = t_pid
    union all
    select o.object_id,
           o.parent_id,
           level + 1
    from object o
             inner join cte on cte.parent_id = o.object_id
    where level < t_lev and cte.parent_id != cte.object_id
)
select object_id,
       parent_id,
       level
from cte;
$$;

--
-- Name: objects_plus_level_isadmin; Type: VIEW;
--

CREATE VIEW objects_plus_level_isadmin AS
SELECT object.object_id,
       object.object_type_id,
       object.owner_id,
       object.parent_id,
       object.asset_2d_id,
       object.asset_3d_id,
       object.options,
       object.transform,
       object.updated_at,
       object.created_at,
       NULL::integer AS level,
       NULL::boolean AS is_admin
FROM object
WHERE false;


--
-- Name: getobjectancestorstableadminrole(uuid, uuid); Type: FUNCTION;
--

CREATE FUNCTION getobjectancestorstableadminrole(t_pid uuid, t_uid uuid) RETURNS SETOF objects_plus_level_isadmin
    LANGUAGE sql
AS $$
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