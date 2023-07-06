BEGIN;

alter table stake
    drop constraint stake_pk;

alter table stake
    add kind int not null;

alter table stake
    add constraint stake_pk
        primary key (blockchain_id, wallet_id, object_id, kind);

COMMIT;
