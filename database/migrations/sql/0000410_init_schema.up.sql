create table "asset_2d" (
    "asset_2d_id" uuid not null,
    "asset_2d_name" character varying(255) not null,
    "options" jsonb default '{}'::jsonb
);


create table "asset_3d" (
    "asset_3d_id" uuid not null,
    "asset_3d_name" character varying(255) not null,
    "options" jsonb default '{}'::jsonb
);


create table "attribute" (
    "plugin_id" uuid not null,
    "attribute_name" character varying(255) not null,
    "description" text not null default ''::text,
    "options" jsonb
);


create table "node_attributes" (
    "plugin_id" uuid not null,
    "attribute_name" character varying(255) not null,
    "value" jsonb not null,
    "options" jsonb
);


create table "plugin" (
    "plugin_id" uuid not null,
    "plugin_name" character varying(255) not null,
    "description" text not null default ''::text,
    "options" jsonb,
    "created_at" timestamp without time zone not null default CURRENT_TIMESTAMP,
    "updated_at" timestamp without time zone not null default CURRENT_TIMESTAMP
);


create table "space" (
    "space_id" uuid not null,
    "space_type_id" uuid not null,
    "owner_id" uuid not null,
    "parent_id" uuid not null,
    "asset_2d_id" uuid,
    "asset_3d_id" uuid,
    "options" jsonb,
    "position" jsonb,
    "updated_at" timestamp without time zone not null default CURRENT_TIMESTAMP,
    "created_at" timestamp without time zone not null default CURRENT_TIMESTAMP
);


create table "space_attribute" (
    "plugin_id" uuid not null,
    "attribute_name" character varying(255) not null,
    "space_id" uuid not null,
    "value" jsonb,
    "options" jsonb
);


create table "space_type" (
    "space_type_id" uuid not null,
    "asset_2d_id" uuid,
    "asset_3d_id" uuid,
    "space_type_name" character varying(255) not null,
    "category_name" character varying(255) not null,
    "description" text default ''::text,
    "options" jsonb not null,
    "created_at" timestamp without time zone not null default CURRENT_TIMESTAMP,
    "updated_at" timestamp without time zone not null default CURRENT_TIMESTAMP
);


create table "space_user_attribute" (
    "user_id" uuid not null,
    "space_id" uuid not null,
    "attribute_name" character varying(255) not null,
    "plugin_id" uuid not null,
    "value" jsonb,
    "options" jsonb
);


create table "user" (
    "user_id" uuid not null,
    "user_type_id" uuid not null,
    "profile" jsonb not null default '{}'::jsonb,
    "updated_at" timestamp without time zone not null default CURRENT_TIMESTAMP,
    "created_at" timestamp without time zone not null default CURRENT_TIMESTAMP,
    "options" jsonb
);


create table "user_attribute" (
    "user_id" uuid not null,
    "plugin_id" uuid not null,
    "attribute_name" character varying(255) not null,
    "value" jsonb not null,
    "options" jsonb
);


create table "user_membership" (
    "member_of" uuid not null,
    "user_id" uuid not null,
    "value" json
);


create table "user_space" (
    "space_id" uuid not null,
    "user_id" uuid not null,
    "created_at" timestamp without time zone not null default CURRENT_TIMESTAMP,
    "updated_at" timestamp without time zone not null default CURRENT_TIMESTAMP,
    "value" jsonb
);


create table "user_type" (
    "user_type_id" uuid not null,
    "user_type_name" character varying(255) not null,
    "description" text not null,
    "options" jsonb
);


create table "user_user_attribute" (
    "plugin_id" uuid not null,
    "attribute_name" character varying(255) not null,
    "source_user_id" uuid not null,
    "target_user_id" uuid not null,
    "value" jsonb not null,
    "options" jsonb
);


CREATE INDEX fk_1 ON public.user_space USING btree (user_id);

CREATE INDEX fk_10 ON public.space USING btree (space_type_id);

CREATE INDEX fk_101 ON public.user_attribute USING btree (user_id);

CREATE INDEX fk_103 ON public.user_attribute USING btree (plugin_id, attribute_name);

CREATE INDEX fk_11 ON public.space USING btree (parent_id);

CREATE INDEX fk_12 ON public.space USING btree (asset_3d_id);

CREATE INDEX fk_14 ON public.space_attribute USING btree (plugin_id, attribute_name);

CREATE INDEX fk_15 ON public.space_attribute USING btree (space_id);

CREATE INDEX fk_17 ON public.space_type USING btree (asset_3d_id);

CREATE INDEX fk_18 ON public.space_type USING btree (asset_2d_id);

CREATE INDEX fk_2 ON public.user_space USING btree (space_id);

CREATE INDEX fk_20 ON public."user" USING btree (user_type_id);

CREATE INDEX fk_202 ON public.user_membership USING btree (user_id);

CREATE INDEX fk_203 ON public.user_membership USING btree (member_of);

CREATE INDEX fk_3 ON public.attribute USING btree (plugin_id);

CREATE INDEX fk_401 ON public.space_user_attribute USING btree (user_id);

CREATE INDEX fk_403 ON public.space_user_attribute USING btree (plugin_id, attribute_name);

CREATE INDEX fk_404 ON public.space_user_attribute USING btree (space_id);

CREATE INDEX fk_5 ON public.node_attributes USING btree (plugin_id, attribute_name);

CREATE INDEX fk_601 ON public.user_user_attribute USING btree (plugin_id, attribute_name);

CREATE INDEX fk_603 ON public.user_user_attribute USING btree (source_user_id);

CREATE INDEX fk_604 ON public.user_user_attribute USING btree (target_user_id);

CREATE INDEX fk_8 ON public.space USING btree (owner_id);

CREATE INDEX fk_9 ON public.space USING btree (asset_2d_id);

CREATE UNIQUE INDEX idx_20929_index_attributes_name ON public.attribute USING btree (attribute_name);

CREATE UNIQUE INDEX idx_21018_ind_796 ON public.space_type USING btree (space_type_name);

CREATE UNIQUE INDEX idx_21929_index_plugin_name ON public.plugin USING btree (plugin_name);

CREATE UNIQUE INDEX pk_1 ON public.asset_2d USING btree (asset_2d_id);

CREATE UNIQUE INDEX pk_102 ON public.user_attribute USING btree (user_id, plugin_id, attribute_name);

CREATE UNIQUE INDEX pk_13 ON public.space USING btree (space_id);

CREATE UNIQUE INDEX pk_16 ON public.space_attribute USING btree (plugin_id, attribute_name, space_id);

CREATE UNIQUE INDEX pk_19 ON public.space_type USING btree (space_type_id);

CREATE UNIQUE INDEX pk_2 ON public.asset_3d USING btree (asset_3d_id);

CREATE UNIQUE INDEX pk_203 ON public.user_membership USING btree (member_of, user_id);

CREATE UNIQUE INDEX pk_21 ON public."user" USING btree (user_id);

CREATE UNIQUE INDEX pk_3 ON public.user_space USING btree (space_id, user_id);

CREATE UNIQUE INDEX pk_4 ON public.attribute USING btree (plugin_id, attribute_name);

CREATE UNIQUE INDEX pk_402 ON public.space_user_attribute USING btree (user_id, space_id, attribute_name, plugin_id);

CREATE UNIQUE INDEX pk_501 ON public.user_type USING btree (user_type_id);

CREATE UNIQUE INDEX pk_6 ON public.node_attributes USING btree (plugin_id, attribute_name);

CREATE UNIQUE INDEX pk_602 ON public.user_user_attribute USING btree (plugin_id, attribute_name, source_user_id, target_user_id);

CREATE UNIQUE INDEX pk_7 ON public.plugin USING btree (plugin_id);

alter table "asset_2d" add constraint "pk_1" PRIMARY KEY using index "pk_1";

alter table "asset_3d" add constraint "pk_2" PRIMARY KEY using index "pk_2";

alter table "attribute" add constraint "pk_4" PRIMARY KEY using index "pk_4";

alter table "node_attributes" add constraint "pk_6" PRIMARY KEY using index "pk_6";

alter table "plugin" add constraint "pk_7" PRIMARY KEY using index "pk_7";

alter table "space" add constraint "pk_13" PRIMARY KEY using index "pk_13";

alter table "space_attribute" add constraint "pk_16" PRIMARY KEY using index "pk_16";

alter table "space_type" add constraint "pk_19" PRIMARY KEY using index "pk_19";

alter table "space_user_attribute" add constraint "pk_402" PRIMARY KEY using index "pk_402";

alter table "user" add constraint "pk_21" PRIMARY KEY using index "pk_21";

alter table "user_attribute" add constraint "pk_102" PRIMARY KEY using index "pk_102";

alter table "user_membership" add constraint "pk_203" PRIMARY KEY using index "pk_203";

alter table "user_space" add constraint "pk_3" PRIMARY KEY using index "pk_3";

alter table "user_type" add constraint "pk_501" PRIMARY KEY using index "pk_501";

alter table "user_user_attribute" add constraint "pk_602" PRIMARY KEY using index "pk_602";

alter table "attribute" add constraint "fk_1" FOREIGN KEY (plugin_id) REFERENCES plugin(plugin_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "attribute" validate constraint "fk_1";

alter table "node_attributes" add constraint "fk_2" FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute(plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "node_attributes" validate constraint "fk_2";

alter table "space" add constraint "fk_25" FOREIGN KEY (asset_2d_id) REFERENCES asset_2d(asset_2d_id) not valid;

alter table "space" validate constraint "fk_25";

alter table "space" add constraint "fk_26" FOREIGN KEY (asset_3d_id) REFERENCES asset_3d(asset_3d_id) not valid;

alter table "space" validate constraint "fk_26";

alter table "space" add constraint "fk_6" FOREIGN KEY (space_type_id) REFERENCES space_type(space_type_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "space" validate constraint "fk_6";

alter table "space" add constraint "fk_7" FOREIGN KEY (parent_id) REFERENCES space(space_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "space" validate constraint "fk_7";

alter table "space" add constraint "fk_8a" FOREIGN KEY (owner_id) REFERENCES "user"(user_id) not valid;

alter table "space" validate constraint "fk_8a";

alter table "space_attribute" add constraint "fk_3" FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute(plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "space_attribute" validate constraint "fk_3";

alter table "space_attribute" add constraint "fk_4" FOREIGN KEY (space_id) REFERENCES space(space_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "space_attribute" validate constraint "fk_4";

alter table "space_type" add constraint "fk_24" FOREIGN KEY (asset_3d_id) REFERENCES asset_3d(asset_3d_id) not valid;

alter table "space_type" validate constraint "fk_24";

alter table "space_type" add constraint "fk_24_1" FOREIGN KEY (asset_2d_id) REFERENCES asset_2d(asset_2d_id) not valid;

alter table "space_type" validate constraint "fk_24_1";

alter table "space_user_attribute" add constraint "fk_10" FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute(plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "space_user_attribute" validate constraint "fk_10";

alter table "space_user_attribute" add constraint "fk_11" FOREIGN KEY (space_id) REFERENCES space(space_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "space_user_attribute" validate constraint "fk_11";

alter table "space_user_attribute" add constraint "fk_9" FOREIGN KEY (user_id) REFERENCES "user"(user_id) not valid;

alter table "space_user_attribute" validate constraint "fk_9";

alter table "user" add constraint "fk_14" FOREIGN KEY (user_type_id) REFERENCES user_type(user_type_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "user" validate constraint "fk_14";

alter table "user_attribute" add constraint "fk_12" FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "user_attribute" validate constraint "fk_12";

alter table "user_attribute" add constraint "fk_13" FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute(plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "user_attribute" validate constraint "fk_13";

alter table "user_membership" add constraint "fk_16" FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "user_membership" validate constraint "fk_16";

alter table "user_membership" add constraint "fk_23" FOREIGN KEY (member_of) REFERENCES "user"(user_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "user_membership" validate constraint "fk_23";

alter table "user_space" add constraint "fk_23_1" FOREIGN KEY (user_id) REFERENCES "user"(user_id) not valid;

alter table "user_space" validate constraint "fk_23_1";

alter table "user_space" add constraint "fk_24_2" FOREIGN KEY (space_id) REFERENCES space(space_id) not valid;

alter table "user_space" validate constraint "fk_24_2";

alter table "user_user_attribute" add constraint "fk_19" FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute(plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "user_user_attribute" validate constraint "fk_19";

alter table "user_user_attribute" add constraint "fk_20" FOREIGN KEY (source_user_id) REFERENCES "user"(user_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "user_user_attribute" validate constraint "fk_20";

alter table "user_user_attribute" add constraint "fk_21" FOREIGN KEY (target_user_id) REFERENCES "user"(user_id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "user_user_attribute" validate constraint "fk_21";


