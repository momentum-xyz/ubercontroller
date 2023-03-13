-- auto-generated definition
create table blockchain
(
    blockchain_id               uuid      not null
        constraint blockchain_pk
            primary key,
    last_processed_block_number bigint    not null,
    blockchain_name             text      not null,
    rpc_url                     text      not null,
    updated_at                  timestamp not null
);

create unique index blockchain_blockchain_id_uindex
    on blockchain (blockchain_id);

-- auto-generated definition
create table wallet
(
    wallet_id     bytea not null,
    blockchain_id uuid  not null
        constraint wallet_blockchain_blockchain_id_fk
            references blockchain
            on update cascade on delete cascade,
    constraint wallet_id_pk
        primary key (blockchain_id, wallet_id)
);

-- auto-generated definition
create table contract
(
    contract_id bytea not null
        constraint contract_pk
            primary key,
    name        varchar(255)
);

create unique index contract_contract_id_uindex
    on contract (contract_id);

-- auto-generated definition
create table balance
(
    wallet_id                   bytea       not null,
    contract_id                 bytea       not null
        constraint balance_contract_contract_id_fk
            references contract
            on update cascade on delete cascade,
    blockchain_id               uuid        not null
        constraint balance_blockchain_blockchain_id_fk
            references blockchain
            on update cascade on delete cascade,
    balance                     numeric(78) not null,
    last_processed_block_number bigint      not null,
    constraint balance_pk
        primary key (wallet_id, contract_id, blockchain_id),
    FOREIGN KEY (blockchain_id, wallet_id) REFERENCES wallet (blockchain_id, wallet_id)

);


