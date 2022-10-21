alter table "asset_2d" drop column "created_at";

alter table "asset_2d" drop column "meta";

alter table "asset_2d" drop column "updated_at";

alter table "asset_2d" add column "asset_2d_name" character varying(255) not null;

alter table "asset_3d" drop column "created_at";

alter table "asset_3d" drop column "meta";

alter table "asset_3d" drop column "updated_at";

alter table "asset_3d" add column "asset_3d_name" character varying(255) not null;

alter table "plugin" drop column "meta";

alter table "plugin" add column "description" text not null default ''::text;

alter table "plugin" add column "plugin_name" character varying(255) not null;

CREATE UNIQUE INDEX idx_21929_index_plugin_name ON public.plugin USING btree (plugin_name);


