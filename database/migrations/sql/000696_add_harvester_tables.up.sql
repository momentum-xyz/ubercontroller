create table harvester_blockchain
(
    blockchain_id                   uuid             not null
        constraint harvester_blockchain_pk
            primary key,
    blockchain_name                 text             not null,
    last_processed_block_for_tokens bigint default 0 not null,
    last_processed_block_for_nfts   bigint default 0 not null,
    last_processed_block_for_ethers bigint default 0 not null,
    updated_at                      timestamp        not null
);

create table harvester_tokens
(
    wallet_id     bytea       not null,
    contract_id   bytea       not null,
    blockchain_id uuid        not null
        constraint harvester_tokens_blockchain_blockchain_id_fk
            references harvester_blockchain
            on update cascade on delete cascade,
    balance       numeric(78) not null,
    updated_at    timestamp   not null,
    PRIMARY KEY (blockchain_id, contract_id, wallet_id)
);

create table harvester_nfts
(
    wallet_id     bytea       not null,
    contract_id   bytea       not null,
    blockchain_id uuid        not null
        constraint harvester_nfts_blockchain_blockchain_id_fk
            references harvester_blockchain
            on update cascade on delete cascade,
    item_id       numeric(78) not null,
    updated_at    timestamp   not null,
    PRIMARY KEY (blockchain_id, contract_id, wallet_id, item_id)
);
