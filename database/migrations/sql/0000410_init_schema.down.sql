alter table "attribute" drop constraint "fk_1";

alter table "node_attributes" drop constraint "fk_2";

alter table "space" drop constraint "fk_25";

alter table "space" drop constraint "fk_26";

alter table "space" drop constraint "fk_6";

alter table "space" drop constraint "fk_7";

alter table "space" drop constraint "fk_8a";

alter table "space_attribute" drop constraint "fk_3";

alter table "space_attribute" drop constraint "fk_4";

alter table "space_type" drop constraint "fk_24";

alter table "space_type" drop constraint "fk_24_1";

alter table "space_user_attribute" drop constraint "fk_10";

alter table "space_user_attribute" drop constraint "fk_11";

alter table "space_user_attribute" drop constraint "fk_9";

alter table "user" drop constraint "fk_14";

alter table "user_attribute" drop constraint "fk_12";

alter table "user_attribute" drop constraint "fk_13";

alter table "user_membership" drop constraint "fk_16";

alter table "user_membership" drop constraint "fk_23";

alter table "user_space" drop constraint "fk_23_1";

alter table "user_space" drop constraint "fk_24_2";

alter table "user_user_attribute" drop constraint "fk_19";

alter table "user_user_attribute" drop constraint "fk_20";

alter table "user_user_attribute" drop constraint "fk_21";

alter table "asset_2d" drop constraint "pk_1";

alter table "asset_3d" drop constraint "pk_2";

alter table "attribute" drop constraint "pk_4";

alter table "node_attributes" drop constraint "pk_6";

alter table "plugin" drop constraint "pk_7";

alter table "space" drop constraint "pk_13";

alter table "space_attribute" drop constraint "pk_16";

alter table "space_type" drop constraint "pk_19";

alter table "space_user_attribute" drop constraint "pk_402";

alter table "user" drop constraint "pk_21";

alter table "user_attribute" drop constraint "pk_102";

alter table "user_membership" drop constraint "pk_203";

alter table "user_space" drop constraint "pk_3";

alter table "user_type" drop constraint "pk_501";

alter table "user_user_attribute" drop constraint "pk_602";

drop index if exists "fk_1";

drop index if exists "fk_10";

drop index if exists "fk_101";

drop index if exists "fk_103";

drop index if exists "fk_11";

drop index if exists "fk_12";

drop index if exists "fk_14";

drop index if exists "fk_15";

drop index if exists "fk_17";

drop index if exists "fk_18";

drop index if exists "fk_2";

drop index if exists "fk_20";

drop index if exists "fk_202";

drop index if exists "fk_203";

drop index if exists "fk_3";

drop index if exists "fk_401";

drop index if exists "fk_403";

drop index if exists "fk_404";

drop index if exists "fk_5";

drop index if exists "fk_601";

drop index if exists "fk_603";

drop index if exists "fk_604";

drop index if exists "fk_8";

drop index if exists "fk_9";

drop index if exists "idx_20929_index_attributes_name";

drop index if exists "idx_21018_ind_796";

drop index if exists "idx_21929_index_plugin_name";

drop index if exists "pk_1";

drop index if exists "pk_102";

drop index if exists "pk_13";

drop index if exists "pk_16";

drop index if exists "pk_19";

drop index if exists "pk_2";

drop index if exists "pk_203";

drop index if exists "pk_21";

drop index if exists "pk_3";

drop index if exists "pk_4";

drop index if exists "pk_402";

drop index if exists "pk_501";

drop index if exists "pk_6";

drop index if exists "pk_602";

drop index if exists "pk_7";

drop table "asset_2d";

drop table "asset_3d";

drop table "attribute";

drop table "node_attributes";

drop table "plugin";

drop table "space";

drop table "space_attribute";

drop table "space_type";

drop table "space_user_attribute";

drop table "user";

drop table "user_attribute";

drop table "user_membership";

drop table "user_space";

drop table "user_type";

drop table "user_user_attribute";


