create table harvester_blockchain
(
    blockchain_id                   uuid      not null
        constraint harvester_blockchain_pk
            primary key,
    blockchain_name                 text      not null,
    last_processed_block_for_tokens bigint    not null,
    last_processed_block_for_nfts   bigint    not null,
    last_processed_block_for_ethers bigint    not null,
    updated_at                      timestamp not null
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
