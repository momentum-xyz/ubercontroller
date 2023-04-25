alter table stake
    alter column amount type numeric(78) using amount::numeric(78);