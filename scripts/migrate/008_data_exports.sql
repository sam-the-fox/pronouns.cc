-- 2023-03-15: Add data export

-- +migrate Up

create table data_exports (
    id         serial primary key,
    user_id    text not null references users (id) on delete cascade,
    filename   text not null,
    created_at timestamptz not null default now()
);

-- +migrate Down

drop table data_exports;
