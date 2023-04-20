alter table stake
    alter column amount type integer using amount::integer;