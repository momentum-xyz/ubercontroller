--
-- PostgreSQL database dump
--

-- Dumped from database version 15.2 (Debian 15.2-1.pgdg110+1)
-- Dumped by pg_dump version 15.2 (Debian 15.2-1.pgdg110+1)

--
-- Name: object;
--

CREATE TABLE object
(
    object_id      uuid                                                  NOT NULL,
    object_type_id uuid                                                  NOT NULL,
    owner_id       uuid                                                  NOT NULL,
    parent_id      uuid                                                  NOT NULL,
    asset_2d_id    uuid,
    asset_3d_id    uuid,
    options        jsonb,
    transform      jsonb,
    updated_at     timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at     timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

--
-- Name: activity;
--

CREATE TABLE activity
(
    activity_id uuid                                                  NOT NULL,
    user_id     uuid                                                  NOT NULL,
    object_id   uuid                                                  NOT NULL,
    type        character varying(255)                                NOT NULL,
    data        jsonb                                                 NOT NULL,
    created_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: asset_2d;
--

CREATE TABLE asset_2d
(
    asset_2d_id uuid                                                  NOT NULL,
    meta        jsonb                       DEFAULT '{}'::jsonb       NOT NULL,
    options     jsonb                       DEFAULT '{}'::jsonb,
    created_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: asset_3d;
--

CREATE TABLE asset_3d
(
    asset_3d_id uuid                                                  NOT NULL,
    meta        jsonb                       DEFAULT '{}'::jsonb       NOT NULL,
    options     jsonb                       DEFAULT '{}'::jsonb,
    created_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: asset_3d_user;
--

CREATE TABLE asset_3d_user
(
    asset_3d_id uuid                                                  NOT NULL,
    user_id     uuid                                                  NOT NULL,
    meta        jsonb                                                 NOT NULL,
    is_private  boolean                     DEFAULT true              NOT NULL,
    updated_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: attribute_type;
--

CREATE TABLE attribute_type
(
    plugin_id      uuid                   NOT NULL,
    attribute_name character varying(255) NOT NULL,
    description    text DEFAULT ''::text  NOT NULL,
    options        jsonb
);


--
-- Name: balance;
--

CREATE TABLE balance
(
    wallet_id                   bytea          NOT NULL,
    contract_id                 bytea          NOT NULL,
    blockchain_id               uuid           NOT NULL,
    balance                     numeric(78, 0) NOT NULL,
    last_processed_block_number bigint         NOT NULL
);


--
-- Name: blockchain;
--

CREATE TABLE blockchain
(
    blockchain_id               uuid                        NOT NULL,
    last_processed_block_number bigint                      NOT NULL,
    blockchain_name             text                        NOT NULL,
    rpc_url                     text                        NOT NULL,
    updated_at                  timestamp without time zone NOT NULL
);


--
-- Name: contract;
--

CREATE TABLE contract
(
    contract_id bytea NOT NULL,
    name        character varying(255)
);


--
-- Name: nft;
--

CREATE TABLE nft
(
    wallet_id     bytea NOT NULL,
    blockchain_id uuid  NOT NULL,
    object_id     uuid  NOT NULL,
    contract_id   bytea NOT NULL,
    created_at    timestamp without time zone,
    updated_at    timestamp without time zone
);


--
-- Name: node_attribute;
--

CREATE TABLE node_attribute
(
    plugin_id      uuid                   NOT NULL,
    attribute_name character varying(255) NOT NULL,
    value          jsonb                  NOT NULL,
    options        jsonb
);


--
-- Name: object_activity;
--

CREATE TABLE object_activity
(
    object_id   uuid                                                  NOT NULL,
    activity_id uuid                                                  NOT NULL,
    created_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: object_attribute;
--

CREATE TABLE object_attribute
(
    plugin_id      uuid                   NOT NULL,
    attribute_name character varying(255) NOT NULL,
    object_id      uuid                   NOT NULL,
    value          jsonb,
    options        jsonb
);


--
-- Name: object_type;
--

CREATE TABLE object_type
(
    object_type_id   uuid                                                  NOT NULL,
    asset_2d_id      uuid,
    asset_3d_id      uuid,
    object_type_name character varying(255)                                NOT NULL,
    category_name    character varying(255)                                NOT NULL,
    description      text                        DEFAULT ''::text,
    options          jsonb                                                 NOT NULL,
    created_at       timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at       timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: object_user_attribute;
--

CREATE TABLE object_user_attribute
(
    user_id        uuid                   NOT NULL,
    object_id      uuid                   NOT NULL,
    attribute_name character varying(255) NOT NULL,
    plugin_id      uuid                   NOT NULL,
    value          jsonb,
    options        jsonb
);


--
-- Name: pending_stake;
--

CREATE TABLE pending_stake
(
    transaction_id bytea                       NOT NULL,
    object_id      uuid                        NOT NULL,
    wallet_id      bytea                       NOT NULL,
    blockchain_id  uuid                        NOT NULL,
    amount         numeric(78, 0)              NOT NULL,
    comment        text                        NOT NULL,
    kind           integer                     NOT NULL,
    updated_at     timestamp without time zone NOT NULL,
    created_at     timestamp without time zone NOT NULL
);


--
-- Name: plugin;
--

CREATE TABLE plugin
(
    plugin_id  uuid                                                  NOT NULL,
    meta       jsonb                       DEFAULT '{}'::jsonb       NOT NULL,
    options    jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: stake;
--

CREATE TABLE stake
(
    wallet_id     bytea                       NOT NULL,
    blockchain_id uuid                        NOT NULL,
    object_id     uuid                        NOT NULL,
    amount        numeric(78, 0)              NOT NULL,
    last_comment  text                        NOT NULL,
    updated_at    timestamp without time zone NOT NULL,
    created_at    timestamp without time zone NOT NULL,
    kind          integer                     NOT NULL
);


--
-- Name: user;
--

CREATE TABLE "user"
(
    user_id      uuid                                                  NOT NULL,
    user_type_id uuid                                                  NOT NULL,
    profile      jsonb                       DEFAULT '{}'::jsonb       NOT NULL,
    updated_at   timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at   timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    options      jsonb
);


--
-- Name: user_activity;
--

CREATE TABLE user_activity
(
    user_id     uuid                                                  NOT NULL,
    activity_id uuid                                                  NOT NULL,
    created_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: user_attribute;
--

CREATE TABLE user_attribute
(
    user_id        uuid                   NOT NULL,
    plugin_id      uuid                   NOT NULL,
    attribute_name character varying(255) NOT NULL,
    value          jsonb                  NOT NULL,
    options        jsonb
);


--
-- Name: user_membership;
--

CREATE TABLE user_membership
(
    member_of  uuid                                                  NOT NULL,
    user_id    uuid                                                  NOT NULL,
    value      json,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_at  timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: user_object;
--

CREATE TABLE user_object
(
    object_id  uuid                                                  NOT NULL,
    user_id    uuid                                                  NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    value      jsonb
);


--
-- Name: user_type;
--

CREATE TABLE user_type
(
    user_type_id   uuid                   NOT NULL,
    user_type_name character varying(255) NOT NULL,
    description    text                   NOT NULL,
    options        jsonb
);


--
-- Name: user_user_attribute;
--

CREATE TABLE user_user_attribute
(
    plugin_id      uuid                   NOT NULL,
    attribute_name character varying(255) NOT NULL,
    source_user_id uuid                   NOT NULL,
    target_user_id uuid                   NOT NULL,
    value          jsonb                  NOT NULL,
    options        jsonb
);


--
-- Name: wallet;
--

CREATE TABLE wallet
(
    wallet_id     bytea NOT NULL,
    blockchain_id uuid  NOT NULL
);


--
-- Name: activity activity_pkey;
--

ALTER TABLE ONLY activity
    ADD CONSTRAINT activity_pkey PRIMARY KEY (activity_id);


--
-- Name: balance balance_pk;
--

ALTER TABLE ONLY balance
    ADD CONSTRAINT balance_pk PRIMARY KEY (wallet_id, contract_id, blockchain_id);


--
-- Name: blockchain blockchain_pk;
--

ALTER TABLE ONLY blockchain
    ADD CONSTRAINT blockchain_pk PRIMARY KEY (blockchain_id);


--
-- Name: contract contract_pk;
--

ALTER TABLE ONLY contract
    ADD CONSTRAINT contract_pk PRIMARY KEY (contract_id);


--
-- Name: nft nft_pk;
--

ALTER TABLE ONLY nft
    ADD CONSTRAINT nft_pk PRIMARY KEY (wallet_id, contract_id, blockchain_id, object_id);


--
-- Name: object_activity object_activity_pkey;
--

ALTER TABLE ONLY object_activity
    ADD CONSTRAINT object_activity_pkey PRIMARY KEY (object_id, activity_id);


--
-- Name: pending_stake pending_stake_pk;
--

ALTER TABLE ONLY pending_stake
    ADD CONSTRAINT pending_stake_pk PRIMARY KEY (blockchain_id, wallet_id, transaction_id, object_id);


--
-- Name: asset_2d pk_1;
--

ALTER TABLE ONLY asset_2d
    ADD CONSTRAINT pk_1 PRIMARY KEY (asset_2d_id);


--
-- Name: user_attribute pk_102;
--

ALTER TABLE ONLY user_attribute
    ADD CONSTRAINT pk_102 PRIMARY KEY (user_id, plugin_id, attribute_name);


--
-- Name: asset_3d_user pk_103;
--

ALTER TABLE ONLY asset_3d_user
    ADD CONSTRAINT pk_103 PRIMARY KEY (asset_3d_id, user_id);


--
-- Name: object pk_13;
--

ALTER TABLE ONLY object
    ADD CONSTRAINT pk_13 PRIMARY KEY (object_id);


--
-- Name: object_attribute pk_16;
--

ALTER TABLE ONLY object_attribute
    ADD CONSTRAINT pk_16 PRIMARY KEY (plugin_id, attribute_name, object_id);


--
-- Name: object_type pk_19;
--

ALTER TABLE ONLY object_type
    ADD CONSTRAINT pk_19 PRIMARY KEY (object_type_id);


--
-- Name: asset_3d pk_2;
--

ALTER TABLE ONLY asset_3d
    ADD CONSTRAINT pk_2 PRIMARY KEY (asset_3d_id);


--
-- Name: user_membership pk_203;
--

ALTER TABLE ONLY user_membership
    ADD CONSTRAINT pk_203 PRIMARY KEY (member_of, user_id);


--
-- Name: user pk_21;
--

ALTER TABLE ONLY "user"
    ADD CONSTRAINT pk_21 PRIMARY KEY (user_id);


--
-- Name: user_object pk_3;
--

ALTER TABLE ONLY user_object
    ADD CONSTRAINT pk_3 PRIMARY KEY (object_id, user_id);


--
-- Name: attribute_type pk_4;
--

ALTER TABLE ONLY attribute_type
    ADD CONSTRAINT pk_4 PRIMARY KEY (plugin_id, attribute_name);


--
-- Name: object_user_attribute pk_402;
--

ALTER TABLE ONLY object_user_attribute
    ADD CONSTRAINT pk_402 PRIMARY KEY (user_id, object_id, attribute_name, plugin_id);


--
-- Name: user_type pk_501;
--

ALTER TABLE ONLY user_type
    ADD CONSTRAINT pk_501 PRIMARY KEY (user_type_id);


--
-- Name: node_attribute pk_6;
--

ALTER TABLE ONLY node_attribute
    ADD CONSTRAINT pk_6 PRIMARY KEY (plugin_id, attribute_name);


--
-- Name: user_user_attribute pk_602;
--

ALTER TABLE ONLY user_user_attribute
    ADD CONSTRAINT pk_602 PRIMARY KEY (plugin_id, attribute_name, source_user_id, target_user_id);


--
-- Name: plugin pk_7;
--

ALTER TABLE ONLY plugin
    ADD CONSTRAINT pk_7 PRIMARY KEY (plugin_id);


--
-- Name: stake stake_pk;
--

ALTER TABLE ONLY stake
    ADD CONSTRAINT stake_pk PRIMARY KEY (blockchain_id, wallet_id, object_id, kind);


--
-- Name: user_activity user_activity_pkey;
--

ALTER TABLE ONLY user_activity
    ADD CONSTRAINT user_activity_pkey PRIMARY KEY (user_id, activity_id);


--
-- Name: wallet wallet_id_pk;
--

ALTER TABLE ONLY wallet
    ADD CONSTRAINT wallet_id_pk PRIMARY KEY (blockchain_id, wallet_id);


--
-- Name: activity_object_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX activity_object_idx ON activity USING btree (object_id);


--
-- Name: activity_user_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX activity_user_idx ON activity USING btree (user_id);


--
-- Name: blockchain_blockchain_id_uindex; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX blockchain_blockchain_id_uindex ON blockchain USING btree (blockchain_id);


--
-- Name: contract_contract_id_uindex; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX contract_contract_id_uindex ON contract USING btree (contract_id);


--
-- Name: fk_1; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_1 ON user_object USING btree (user_id);


--
-- Name: fk_10; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_10 ON object USING btree (object_type_id);


--
-- Name: fk_101; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_101 ON user_attribute USING btree (user_id);


--
-- Name: fk_103; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_103 ON user_attribute USING btree (plugin_id, attribute_name);


--
-- Name: fk_11; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_11 ON object USING btree (parent_id);


--
-- Name: fk_110; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_110 ON asset_3d_user USING btree (asset_3d_id);


--
-- Name: fk_111; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_111 ON asset_3d_user USING btree (user_id);


--
-- Name: fk_12; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_12 ON object USING btree (asset_3d_id);


--
-- Name: fk_14; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_14 ON object_attribute USING btree (plugin_id, attribute_name);


--
-- Name: fk_15; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_15 ON object_attribute USING btree (object_id);


--
-- Name: fk_17; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_17 ON object_type USING btree (asset_3d_id);


--
-- Name: fk_18; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_18 ON object_type USING btree (asset_2d_id);


--
-- Name: fk_2; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_2 ON user_object USING btree (object_id);


--
-- Name: fk_20; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_20 ON "user" USING btree (user_type_id);


--
-- Name: fk_202; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_202 ON user_membership USING btree (user_id);


--
-- Name: fk_203; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_203 ON user_membership USING btree (member_of);


--
-- Name: fk_3; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_3 ON attribute_type USING btree (plugin_id);


--
-- Name: fk_401; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_401 ON object_user_attribute USING btree (user_id);


--
-- Name: fk_403; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_403 ON object_user_attribute USING btree (plugin_id, attribute_name);


--
-- Name: fk_404; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_404 ON object_user_attribute USING btree (object_id);


--
-- Name: fk_5; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_5 ON node_attribute USING btree (plugin_id, attribute_name);


--
-- Name: fk_601; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_601 ON user_user_attribute USING btree (plugin_id, attribute_name);


--
-- Name: fk_603; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_603 ON user_user_attribute USING btree (source_user_id);


--
-- Name: fk_604; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_604 ON user_user_attribute USING btree (target_user_id);


--
-- Name: fk_8; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_8 ON object USING btree (owner_id);


--
-- Name: fk_9; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX fk_9 ON object USING btree (asset_2d_id);


--
-- Name: idx_21018_ind_796; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_21018_ind_796 ON object_type USING btree (object_type_name);


--
-- Name: oa_activity_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX oa_activity_idx ON object_activity USING btree (activity_id);


--
-- Name: oa_object_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX oa_object_idx ON object_activity USING btree (object_id);


--
-- Name: ua_activity_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX ua_activity_idx ON user_activity USING btree (activity_id);


--
-- Name: ua_user_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX ua_user_idx ON user_activity USING btree (user_id);


--
-- Name: balance balance_blockchain_blockchain_id_fk;
--

ALTER TABLE ONLY balance
    ADD CONSTRAINT balance_blockchain_blockchain_id_fk FOREIGN KEY (blockchain_id) REFERENCES blockchain (blockchain_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: balance balance_blockchain_id_wallet_id_fkey;
--

ALTER TABLE ONLY balance
    ADD CONSTRAINT balance_blockchain_id_wallet_id_fkey FOREIGN KEY (blockchain_id, wallet_id) REFERENCES wallet (blockchain_id, wallet_id);


--
-- Name: balance balance_contract_contract_id_fk;
--

ALTER TABLE ONLY balance
    ADD CONSTRAINT balance_contract_contract_id_fk FOREIGN KEY (contract_id) REFERENCES contract (contract_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: attribute_type fk_1;
--

ALTER TABLE ONLY attribute_type
    ADD CONSTRAINT fk_1 FOREIGN KEY (plugin_id) REFERENCES plugin (plugin_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object_user_attribute fk_10;
--

ALTER TABLE ONLY object_user_attribute
    ADD CONSTRAINT fk_10 FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute_type (plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object_user_attribute fk_11;
--

ALTER TABLE ONLY object_user_attribute
    ADD CONSTRAINT fk_11 FOREIGN KEY (object_id) REFERENCES object (object_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_attribute fk_12;
--

ALTER TABLE ONLY user_attribute
    ADD CONSTRAINT fk_12 FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_attribute fk_13;
--

ALTER TABLE ONLY user_attribute
    ADD CONSTRAINT fk_13 FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute_type (plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user fk_14;
--

ALTER TABLE ONLY "user"
    ADD CONSTRAINT fk_14 FOREIGN KEY (user_type_id) REFERENCES user_type (user_type_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_membership fk_16;
--

ALTER TABLE ONLY user_membership
    ADD CONSTRAINT fk_16 FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_user_attribute fk_19;
--

ALTER TABLE ONLY user_user_attribute
    ADD CONSTRAINT fk_19 FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute_type (plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: node_attribute fk_2;
--

ALTER TABLE ONLY node_attribute
    ADD CONSTRAINT fk_2 FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute_type (plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_user_attribute fk_20;
--

ALTER TABLE ONLY user_user_attribute
    ADD CONSTRAINT fk_20 FOREIGN KEY (source_user_id) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_user_attribute fk_21;
--

ALTER TABLE ONLY user_user_attribute
    ADD CONSTRAINT fk_21 FOREIGN KEY (target_user_id) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_membership fk_23;
--

ALTER TABLE ONLY user_membership
    ADD CONSTRAINT fk_23 FOREIGN KEY (member_of) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_object fk_23_1;
--

ALTER TABLE ONLY user_object
    ADD CONSTRAINT fk_23_1 FOREIGN KEY (user_id) REFERENCES "user" (user_id);


--
-- Name: object_type fk_24;
--

ALTER TABLE ONLY object_type
    ADD CONSTRAINT fk_24 FOREIGN KEY (asset_3d_id) REFERENCES asset_3d (asset_3d_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object_type fk_24_1;
--

ALTER TABLE ONLY object_type
    ADD CONSTRAINT fk_24_1 FOREIGN KEY (asset_2d_id) REFERENCES asset_2d (asset_2d_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_object fk_24_2;
--

ALTER TABLE ONLY user_object
    ADD CONSTRAINT fk_24_2 FOREIGN KEY (object_id) REFERENCES object (object_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object fk_25;
--

ALTER TABLE ONLY object
    ADD CONSTRAINT fk_25 FOREIGN KEY (asset_2d_id) REFERENCES asset_2d (asset_2d_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object fk_26;
--

ALTER TABLE ONLY object
    ADD CONSTRAINT fk_26 FOREIGN KEY (asset_3d_id) REFERENCES asset_3d (asset_3d_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object_attribute fk_3;
--

ALTER TABLE ONLY object_attribute
    ADD CONSTRAINT fk_3 FOREIGN KEY (plugin_id, attribute_name) REFERENCES attribute_type (plugin_id, attribute_name) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: asset_3d_user fk_32;
--

ALTER TABLE ONLY asset_3d_user
    ADD CONSTRAINT fk_32 FOREIGN KEY (asset_3d_id) REFERENCES asset_3d (asset_3d_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: asset_3d_user fk_33;
--

ALTER TABLE ONLY asset_3d_user
    ADD CONSTRAINT fk_33 FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object_attribute fk_4;
--

ALTER TABLE ONLY object_attribute
    ADD CONSTRAINT fk_4 FOREIGN KEY (object_id) REFERENCES object (object_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object fk_6;
--

ALTER TABLE ONLY object
    ADD CONSTRAINT fk_6 FOREIGN KEY (object_type_id) REFERENCES object_type (object_type_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object fk_7;
--

ALTER TABLE ONLY object
    ADD CONSTRAINT fk_7 FOREIGN KEY (parent_id) REFERENCES object (object_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object fk_8a;
--

ALTER TABLE ONLY object
    ADD CONSTRAINT fk_8a FOREIGN KEY (owner_id) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: activity fk_activity_object;
--

ALTER TABLE ONLY activity
    ADD CONSTRAINT fk_activity_object FOREIGN KEY (object_id) REFERENCES object (object_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: activity fk_activity_user;
--

ALTER TABLE ONLY activity
    ADD CONSTRAINT fk_activity_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object_activity fk_oa_activity;
--

ALTER TABLE ONLY object_activity
    ADD CONSTRAINT fk_oa_activity FOREIGN KEY (activity_id) REFERENCES activity (activity_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object_activity fk_oa_object;
--

ALTER TABLE ONLY object_activity
    ADD CONSTRAINT fk_oa_object FOREIGN KEY (object_id) REFERENCES object (object_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: object_user_attribute fk_oua_user;
--

ALTER TABLE ONLY object_user_attribute
    ADD CONSTRAINT fk_oua_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_activity fk_ua_activity;
--

ALTER TABLE ONLY user_activity
    ADD CONSTRAINT fk_ua_activity FOREIGN KEY (activity_id) REFERENCES activity (activity_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_activity fk_ua_user;
--

ALTER TABLE ONLY user_activity
    ADD CONSTRAINT fk_ua_user FOREIGN KEY (user_id) REFERENCES "user" (user_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: nft nft_blockchain_blockchain_id_fk;
--

ALTER TABLE ONLY nft
    ADD CONSTRAINT nft_blockchain_blockchain_id_fk FOREIGN KEY (blockchain_id) REFERENCES blockchain (blockchain_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: pending_stake pending_stake_object_fk;
--

ALTER TABLE ONLY pending_stake
    ADD CONSTRAINT pending_stake_object_fk FOREIGN KEY (object_id) REFERENCES object (object_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: pending_stake pending_stake_wallet_fk;
--

ALTER TABLE ONLY pending_stake
    ADD CONSTRAINT pending_stake_wallet_fk FOREIGN KEY (blockchain_id, wallet_id) REFERENCES wallet (blockchain_id, wallet_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: stake stake_wallet_fk;
--

ALTER TABLE ONLY stake
    ADD CONSTRAINT stake_wallet_fk FOREIGN KEY (blockchain_id, wallet_id) REFERENCES wallet (blockchain_id, wallet_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: wallet wallet_blockchain_blockchain_id_fk;
--

ALTER TABLE ONLY wallet
    ADD CONSTRAINT wallet_blockchain_blockchain_id_fk FOREIGN KEY (blockchain_id) REFERENCES blockchain (blockchain_id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

