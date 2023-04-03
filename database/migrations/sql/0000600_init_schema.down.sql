BEGIN;

DROP TABLE object_attribute;
DROP TABLE object_user_attribute;
DROP TABLE user_object;
DROP TABLE node_attribute;
DROP TABLE object;
DROP TABLE user_attribute;
DROP TABLE user_membership;
DROP TABLE user_user_attribute;
DROP TABLE attribute_type;
DROP TABLE object_type;
DROP TABLE "user";
DROP TABLE asset_2d;
DROP TABLE asset_3d;
DROP TABLE plugin;
DROP TABLE user_type;

-- DROP INDEX FK_20 ON "user"
-- DROP UNIQUE INDEX idx_21018_ind_796 ON object_type USING btree;
-- DROP INDEX FK_17 ON object_type;
-- DROP INDEX FK_18 ON object_type;
-- DROP INDEX FK_3 ON attribute_type;
-- DROP INDEX FK_601 ON user_user_attribute;
-- DROP INDEX FK_603 ON user_user_attribute;
-- DROP INDEX FK_604 ON user_user_attribute;
-- DROP INDEX FK_202 ON user_membership;
-- DROP INDEX FK_203 ON user_membership;
-- DROP INDEX FK_101 ON user_attribute;
-- DROP INDEX FK_103 ON user_attribute;
-- DROP INDEX FK_10 ON object;
-- DROP INDEX FK_11 ON object;
-- DROP INDEX FK_12 ON object;
-- DROP INDEX FK_8 ON object;
-- DROP INDEX FK_9 ON object;
-- DROP INDEX FK_5 ON node_attribute;
-- DROP INDEX FK_1 ON user_object;
-- DROP INDEX FK_2 ON user_object;
-- DROP INDEX FK_401 ON object_user_attribute;
-- DROP INDEX FK_403 ON object_user_attribute;
-- DROP INDEX FK_404 ON object_user_attribute;
-- DROP INDEX FK_14 ON object_attribute;
-- DROP INDEX FK_15 ON object_attribute;


COMMIT;
