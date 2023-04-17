BEGIN;

DROP FUNCTION IF EXISTS getindirectobjectadmins(t_pid uuid);
DROP FUNCTION IF EXISTS getobjectancestorsids(t_pid uuid, t_lev integer);

DROP FUNCTION IF EXISTS getobjectancestorstableadminrole(t_pid uuid, t_uid uuid);
DROP VIEW IF EXISTS objects_plus_level_isadmin;

COMMIT;
