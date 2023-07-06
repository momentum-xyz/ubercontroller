BEGIN;

alter table stake
    drop constraint stake_pk;

alter table stake
    add constraint stake_pk
        primary key (blockchain_id, wallet_id, object_id);

alter table stake
    drop column kind;

COMMIT;
