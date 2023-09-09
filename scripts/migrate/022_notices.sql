-- 2023-09-09: Add global notices

-- +migrate Up

create table notices (
    id serial primary key,
    notice text not null,
    start_time timestamptz not null default now(),
    end_time timestamptz not null
);

-- +migrate Down

drop table notices;
