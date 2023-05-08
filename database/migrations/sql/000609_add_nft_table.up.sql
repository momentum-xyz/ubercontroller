create table nft
(
    wallet_id     bytea not null,
    blockchain_id uuid  not null
        constraint nft_blockchain_blockchain_id_fk
            references blockchain
            on update cascade on delete cascade,
    object_id     uuid  not null,
    contract_id   bytea not null,
    created_at    timestamp,
    updated_at    timestamp,
    constraint nft_pk
        primary key (wallet_id, contract_id, blockchain_id, object_id)
);