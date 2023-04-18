create table stake
(
    wallet_id     bytea     not null,
    blockchain_id uuid      not null,
    object_id     uuid      not null,
    amount        integer   not null,
    last_comment  text      not null,
    updated_at    timestamp not null,
    created_at    timestamp not null,
    constraint stake_pk
        primary key (blockchain_id, wallet_id, object_id),
    constraint stake_wallet_fk
        foreign key (blockchain_id, wallet_id) references wallet (blockchain_id, wallet_id)
            on update cascade on delete cascade,
    constraint stake_object_fk
        foreign key (object_id) references object (object_id)
            on update cascade on delete cascade
);

create table pending_stake
(
    transaction_id bytea       not null,
    object_id      uuid        not null,
    wallet_id      bytea       not null,
    blockchain_id  uuid        not null,
    amount         numeric(78) not null,
    comment        text        not null,
    kind           integer     not null,
    updated_at     timestamp   not null,
    created_at     timestamp   not null,
    constraint pending_stake_pk
        primary key (blockchain_id, wallet_id, transaction_id, object_id),
    constraint pending_stake_wallet_fk
        foreign key (blockchain_id, wallet_id) references wallet (blockchain_id, wallet_id)
            on update cascade on delete cascade,
    constraint pending_stake_object_fk
        foreign key (object_id) references object (object_id)
            on update cascade on delete cascade
);