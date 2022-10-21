drop index if exists "idx_21929_index_plugin_name";

alter table "asset_2d" drop column "asset_2d_name";

alter table "asset_2d" add column "created_at" timestamp without time zone not null;

alter table "asset_2d" add column "meta" jsonb not null;

alter table "asset_2d" add column "updated_at" timestamp without time zone not null;

alter table "asset_3d" drop column "asset_3d_name";

alter table "asset_3d" add column "created_at" timestamp without time zone not null;

alter table "asset_3d" add column "meta" jsonb not null;

alter table "asset_3d" add column "updated_at" timestamp without time zone not null;

alter table "plugin" drop column "description";

alter table "plugin" drop column "plugin_name";

alter table "plugin" add column "meta" jsonb not null;


